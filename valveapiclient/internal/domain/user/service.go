package user

import (
	"errors"
	"fmt"
	"os"

	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
	"github.com/Cludch/csgo-microservices/valveapiclient/pkg/valve_match_api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	userCreationRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "valveapiclient_api_user_creation",
		Help: "The total number of requested share codes",
	})
)

type UserService struct {
	repo     UserRepository
	consumer valve_match_api.ValveMatchApiConsumerUseCase
}

func NewService(r UserRepository, c valve_match_api.ValveMatchApiConsumerUseCase) *UserService {
	return &UserService{
		repo:     r,
		consumer: c,
	}
}

func (s *UserService) GetUser(id uint64) (*User, error) {
	return s.repo.Find(id)
}

func (s *UserService) GetUsersWithApiEnabled() ([]*User, error) {
	return s.repo.ListAllWithApiEnabled()
}

func (s *UserService) CreateUser(id uint64) (*User, error) {
	u := NewUser(id)

	dbUser, _ := s.GetUser(u.ID)
	if dbUser != nil {
		return nil, fmt.Errorf("user with steam id %d already exists", u.ID)
	}

	// increment prometheus metric.
	userCreationRequests.Inc()

	return u, s.repo.Create(u)
}

func (s *UserService) AddSteamMatchHistoryAuthenticationCode(user *User, apiKey string, authCode string, sc string) error {
	shareCode, errCode := share_code.Decode(sc)
	if errCode != nil {
		return errors.New("invalid share code")
	}

	// Test credentials
	_, errTest := s.consumer.RequestNextShareCode(apiKey, user.ID, authCode, sc)
	if errTest != nil {
		return errors.New("invalid authentication code or last share code")
	}

	user.AuthCode = authCode
	user.ApiKey = apiKey

	errApiKey := s.repo.UpdateApiKey(user)
	if errApiKey != nil {
		return errApiKey
	}

	errAdd := s.repo.UpdateMatchAuthCode(user)
	if errAdd != nil {
		return errAdd
	}

	errSc := s.UpdateLatestShareCode(user, shareCode)
	if errSc != nil {
		return errSc
	}

	errApi := s.UpdateSteamApiUsage(user, true)
	if errApi != nil {
		return errApi
	}

	return nil
}

func (s *UserService) UpdateSteamApiUsage(u *User, active bool) error {
	if active && u.AuthCode == "" {
		return errors.New("missing steam api auth code")
	}

	err := s.repo.UpdateSteamApiUsage(u)
	if err == nil {
		u.ApiEnabled = active
	}
	return err
}

func (s *UserService) UpdateLatestShareCode(u *User, sc *share_code.ShareCodeData) error {
	err := s.repo.UpdateLatestShareCode(u)
	if err == nil {
		u.LastShareCode = sc.Encoded
	}
	return err
}

func (s *UserService) QueryLatestShareCode(u *User) (*share_code.ShareCodeData, error) {
	if !u.ApiEnabled {
		return nil, errors.New("user: api usage is disabled")
	}

	steamID := u.ID
	shareCode, err := s.consumer.RequestNextShareCode(u.ApiKey, steamID, u.AuthCode, u.LastShareCode)

	// Disable user on error.
	if err != nil {
		if os.IsTimeout(err) {
			return nil, errors.New("user: lost connection while querying the steam api for the latest sharecode")
		}

		updateErr := s.UpdateSteamApiUsage(u, false)
		if updateErr != nil {
			const msg = "disabled csgo user %d due to an error (%t) in fetching the share code"
			log.Warnf(msg, steamID, err)
		}
		return nil, err
	}

	// No new match.
	if shareCode == "" {
		const msg = "no new match found for %d"
		log.Trace(fmt.Sprintf(msg, steamID))
		return nil, nil
	}

	const msg = "found match share code %v for %d"
	log.Info(fmt.Sprintf(msg, shareCode, steamID))

	sc, err := share_code.Decode(shareCode)
	if err != nil {
		const msg = "invalid share code %s"
		return nil, fmt.Errorf(msg, sc.Encoded)
	}

	return sc, nil
}
