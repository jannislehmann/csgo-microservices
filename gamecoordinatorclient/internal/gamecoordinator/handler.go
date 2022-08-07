package gamecoordinator

import (
	"time"

	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
	shared "github.com/Cludch/csgo-microservices/shared/proto"
	csgo "github.com/Philipp15b/go-steam/v3/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v3/protocol/gamecoordinator"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Channel is used to request demos one after another.
var matchResponse chan *shared.MatchDetails

// RequestMatchDetails sends a protobuf message to the gc to request details for the requested match.
func (s *GamecoordinatorService) RequestMatchDetails(sc *share_code.ShareCodeData) chan *shared.MatchDetails {
	matchResponse = make(chan *shared.MatchDetails)
	const msg = "requesting match details for %v %d"
	log.Debugf(msg, sc.Encoded, sc.MatchID)
	go s.write(uint32(csgo.ECsgoGCMsg_k_EMsgGCCStrike15_v2_MatchListRequestFullGameInfo), &csgo.CMsgGCCStrike15V2_MatchListRequestFullGameInfo{
		Matchid:   &sc.MatchID,
		Outcomeid: &sc.OutcomeID,
		Token:     &sc.Token,
	})

	return matchResponse
}

// HandleGCPacket takes incoming packets from the GC and coordinates them to the handler funcs.
func (s *GamecoordinatorService) HandleGCPacket(packet *gamecoordinator.GCPacket) {
	if packet.AppId != AppID {
		log.Debug("wrong app id")
		return
	}

	if packet.MsgType == uint32(csgo.EGCBaseClientMsg_k_EMsgGCClientWelcome) {
		s.handleClientWelcome(packet)
	} else if packet.MsgType == uint32(csgo.ECsgoGCMsg_k_EMsgGCCStrike15_v2_MatchList) {
		s.handleMatchList(packet)
	}
}

// handleClientWelcome logs that the client connected successfully.
func (s *GamecoordinatorService) handleClientWelcome(packet *gamecoordinator.GCPacket) {
	log.Info("connected to csgo gc")
}

// handleMatchList handles a gc message containing matches and tries to download those.
func (s *GamecoordinatorService) handleMatchList(packet *gamecoordinator.GCPacket) {
	matchList := new(csgo.CMsgGCCStrike15V2_MatchList)
	packet.ReadProtoMsg(matchList)

	for _, matchEntry := range matchList.GetMatches() {
		for _, round := range matchEntry.GetRoundstatsall() {
			// Demo link is only linked in the last round and in this case the reserveration id is set.
			// Reservation id is the outcome id.
			if round.GetReservationid() == 0 {
				continue
			}

			matchDetails := &shared.MatchDetails{
				MatchId:     round.GetReservationid(),
				MatchTime:   timestamppb.New(time.Unix(int64(*matchEntry.Matchtime), 0)),
				DownloadUrl: round.GetMap(),
			}

			const msg = "received match details for %d"
			log.Infof(msg, matchDetails.MatchId)
			matchResponse <- matchDetails
		}
	}
}
