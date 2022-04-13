package faceit_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Cludch/csgo-microservices/shared/pkg/api_client"
	"github.com/google/uuid"
)

const apiUrl = "https://open.faceit.com/data/v4"

type FaceitApiConsumerService struct {
	httpClient api_client.HttpClient
}

func New(httpClient api_client.HttpClient) *FaceitApiConsumerService {
	return &FaceitApiConsumerService{
		httpClient: httpClient,
	}
}

// RequestNextShareCode returns the next match's share code.
// It uses the saved share codes as the current one.
func (s *FaceitApiConsumerService) GetMatchDetails(faceitAPIKey string, matchID string) (*MatchDetailsResponse, error) {
	playerResponse := &MatchDetailsResponse{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/matches/%s", apiUrl, matchID), nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", faceitAPIKey))

	r, err := s.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	if r.StatusCode == http.StatusTooManyRequests || r.StatusCode == http.StatusServiceUnavailable {
		r.Body.Close()
		return nil, &FaceitApiConnectionIssues{}
	}

	if r.StatusCode == http.StatusUnauthorized || r.StatusCode == http.StatusForbidden {
		r.Body.Close()
		return nil, &InvalidFaceitApiCredentials{}
	}

	if err = json.NewDecoder(r.Body).Decode(playerResponse); err != nil {
		r.Body.Close()
		return nil, err
	}

	defer r.Body.Close()

	return playerResponse, nil
}

// GetPlayerMatchHistory returns the match history for a given player.
func (s *FaceitApiConsumerService) GetPlayerMatchHistory(faceitAPIKey string, playerId uuid.UUID) (*PlayerMatchHistoryResponse, error) {
	u, err := url.Parse(fmt.Sprintf("https://open.faceit.com/data/v4/players/%s/history", playerId))
	if err != nil {
		return nil, errors.New("faceitapi: unable to parse url")
	}

	// Build query
	q := u.Query()
	q.Set("game", "csgo")
	q.Set("offset", "0")
	q.Set("limit", "20")
	// Query the last 24 hours only
	q.Set("from", fmt.Sprint(time.Now().AddDate(0, -1, 0).Unix()))
	u.RawQuery = q.Encode()

	matchResponse := &PlayerMatchHistoryResponse{}

	// Request player match history.
	request, _ := http.NewRequest("GET", u.String(), nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", faceitAPIKey))
	r, rErr := s.httpClient.Do(request)
	if rErr != nil {
		return nil, rErr
	}

	if r.StatusCode == http.StatusTooManyRequests || r.StatusCode == http.StatusServiceUnavailable {
		r.Body.Close()
		return nil, &FaceitApiConnectionIssues{}
	}

	if r.StatusCode == http.StatusUnauthorized || r.StatusCode == http.StatusForbidden {
		r.Body.Close()
		return nil, &InvalidFaceitApiCredentials{}
	}

	if err = json.NewDecoder(r.Body).Decode(matchResponse); err != nil {
		r.Body.Close()
		return nil, err
	}

	defer r.Body.Close()

	return matchResponse, nil
}
