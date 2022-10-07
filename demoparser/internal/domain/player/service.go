package player

import (
	"errors"

	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
)

type PlayerService struct {
	repo Repository
}

func NewService(r Repository) *PlayerService {
	return &PlayerService{
		repo: r,
	}
}

func (s *PlayerService) CreatePlayer(id uint64) (*Player, error) {
	p, _ := NewPlayer(id)
	return p, s.repo.Create(p)
}

func (s *PlayerService) GetAll() ([]*Player, error) {
	return s.repo.List()
}

func (s *PlayerService) GetPlayer(id uint64) (*Player, error) {
	p, err := s.repo.Find(id)
	if p == nil {
		return s.CreatePlayer(id)
	}

	return p, err
}

func (s *PlayerService) GetResult(p *Player, matchId entity.ID) (*PlayerResult, error) {
	for _, result := range p.Results {
		if result.MatchID == matchId {
			return result, nil
		}
	}

	return nil, entity.ErrNotFound
}

func (s *PlayerService) AddResult(p *Player, r *PlayerResult) error {
	matchId := r.MatchID

	// Delete old result.
	dbResult, err := s.GetResult(p, matchId)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return err
	}

	if dbResult != nil {
		err = s.DeleteResult(p, matchId)
		if err != nil {
			return err
		}
	}

	p.Results = append(p.Results, r)
	return s.repo.AddResult(p, r)
}

func (s *PlayerService) DeleteResult(p *Player, matchId entity.ID) error {
	return s.repo.DeleteResult(p, matchId)
}
