package valve_match_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Cludch/csgo-microservices/valveapiclient/pkg/api_client"
)

const apiUrl = "https://api.steampowered.com/ICSGOPlayers_730/GetNextMatchSharingCode/v1"

type ValveMatchApiConsumerService struct {
	apiCient api_client.ApiClient
}

func New(apiClient api_client.ApiClient) *ValveMatchApiConsumerService {
	return &ValveMatchApiConsumerService{
		apiCient: apiClient,
	}
}

// MatchResponse contains information about the latest match.
type MatchResponse struct {
	Result struct {
		Nextcode string `json:"nextcode"`
	} `json:"result"`
}

// InvalidSteamID is used to notify when the supplied credentials are not valid / cannot be used with the api.
type InvalidSteamID struct {
	SteamID string
}

func (e *InvalidSteamID) Error() string {
	const msg = "Invalid steam id %v."
	return fmt.Sprintf(msg, e.SteamID)
}

// InvalidMatchHistoryCredentials is used to notify when the supplied credentials are not valid / cannot be used with the api.
type InvalidApiKeyOrAuthCode struct {
	SteamID string
}

func (e *InvalidApiKeyOrAuthCode) Error() string {
	const msg = "Invalid api key or auth code for steam id %v."
	return fmt.Sprintf(msg, e.SteamID)
}

// RequestNextShareCode returns the next match's share code.
// It uses the saved share codes as the current one.
func (s *ValveMatchApiConsumerService) RequestNextShareCode(steamAPIKey string, steamID uint64, historyAuthenticationCode string, lastShareCode string) (string, error) {
	// Get latest match
	u, err := url.Parse(apiUrl)
	if err != nil {
		return "", errors.New("valveapi: unable to parse url")
	}

	steamIDString := strconv.FormatUint(steamID, 10)

	// Build query
	q := u.Query()
	q.Set("key", steamAPIKey)
	q.Set("steamid", steamIDString)
	q.Set("steamidkey", historyAuthenticationCode)
	q.Set("knowncode", lastShareCode)
	u.RawQuery = q.Encode()

	matchResponse := &MatchResponse{}

	// Request match code.
	r, err := s.apiCient.Get(u.String())
	if err != nil {
		return "", err
	}

	// Forbidden = wrong api keys.
	if r.StatusCode == http.StatusForbidden {
		r.Body.Close()
		return "", &InvalidApiKeyOrAuthCode{SteamID: steamIDString}
	}

	// Precondition Failed = Steam id wrong.
	if r.StatusCode == http.StatusPreconditionFailed {
		r.Body.Close()
		return "", &InvalidSteamID{SteamID: steamIDString}
	}

	// 500 or 504 indicate an error with the API from Valve and not an authentication problem.
	// We should retry the same request later again.
	if r.StatusCode == http.StatusInternalServerError || r.StatusCode == http.StatusGatewayTimeout {
		return "", nil
	}

	// Accepted means that there is no recent match code available.
	if r.StatusCode == http.StatusAccepted {
		r.Body.Close()
		return "", nil
	}

	if err = json.NewDecoder(r.Body).Decode(matchResponse); err != nil {
		r.Body.Close()
		return "", err
	}

	defer r.Body.Close()

	return matchResponse.Result.Nextcode, nil
}
