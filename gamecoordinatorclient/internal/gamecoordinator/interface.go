package gamecoordinator

import (
	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
)

type GamecoordinatorUseCase interface {
	Connect(username, password, twoFactorSecret string)
	RequestMatchDetails(*share_code.ShareCodeData) chan *MatchDetails
	HandleGCPacket(packet *gamecoordinator.GCPacket)
}
