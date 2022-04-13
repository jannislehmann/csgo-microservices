package user

import (
	"fmt"

	"github.com/Cludch/csgo-microservices/faceitapiclient/pkg/faceit_api"
	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	userCreationRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "faceitapiclient_api_user_creation",
		Help: "The total number of created faceit api client users",
	})
)

type UserService struct {
	repo     UserRepository
	consumer faceit_api.FaceitApiConsumerUseCase
}

func NewService(r UserRepository, c faceit_api.FaceitApiConsumerUseCase) *UserService {
	return &UserService{
		repo:     r,
		consumer: c,
	}
}

func (s *UserService) GetUser(id entity.ID) (*User, error) {
	return s.repo.Find(id)
}

func (s *UserService) GetUsersWithApiEnabled() ([]*User, error) {
	return s.repo.ListAllWithApiEnabled()
}

func (s *UserService) CreateUser(id entity.ID) (*User, error) {
	u := NewUser(id)

	dbUser, _ := s.GetUser(u.ID)
	if dbUser != nil {
		return nil, fmt.Errorf("user with steam id %d already exists", u.ID)
	}

	// increment prometheus metric.
	userCreationRequests.Inc()

	return u, s.repo.Create(u)
}

func (s *UserService) UpdateFaceitApiUsage(u *User, active bool) error {
	u.ApiEnabled = active
	err := s.repo.UpdateFaceitApiUsage(u)
	if err != nil {
		u.ApiEnabled = !active
	}
	return err
}
