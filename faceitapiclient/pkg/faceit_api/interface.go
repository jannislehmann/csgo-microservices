package faceit_api

import (
	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
)

type FaceitApiConsumerUseCase interface {
	GetPlayerMatchHistory(faceitAPIKey string, playerID entity.ID) (*PlayerMatchHistoryResponse, error)
	GetMatchDetails(faceitAPIKey string, matchID string) (*MatchDetailsResponse, error)
}
