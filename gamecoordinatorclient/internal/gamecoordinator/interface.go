package gamecoordinator

import (
	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
	shared "github.com/Cludch/csgo-microservices/shared/proto"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
)

type GamecoordinatorUseCase interface {
	Connect(username, password, twoFactorSecret string)
	RequestMatchDetails(*share_code.ShareCodeData) chan *shared.MatchDetails
	HandleGCPacket(packet *gamecoordinator.GCPacket)
}
