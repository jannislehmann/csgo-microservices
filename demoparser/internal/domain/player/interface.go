package player

import (
	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
)

type Repository interface {
	Create(*Player) error

	Find(uint64) (*Player, error)

	List() ([]*Player, error)

	AddResult(*Player, *PlayerResult) error

	DeleteResult(*Player, entity.ID) error
}

type UseCase interface {
	CreatePlayer(steamId uint64) (*Player, error)

	GetAll() ([]*Player, error)
	GetPlayer(uint64) (*Player, error)
	GetResult(p *Player, matchId entity.ID) (*PlayerResult, error)

	AddResult(*Player, *PlayerResult) error

	DeleteResult(p *Player, matchId entity.ID) error
}
