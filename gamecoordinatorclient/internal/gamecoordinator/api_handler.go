package gamecoordinator

import (
	"context"

	pb "github.com/Cludch/csgo-microservices/gamecoordinatorclient/proto"
	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GamecoordinatorApiHandler struct {
	service GamecoordinatorUseCase
}

func NewGamecoordinatorApiHandler(s GamecoordinatorUseCase) *GamecoordinatorApiHandler {
	return &GamecoordinatorApiHandler{
		service: s,
	}
}

func (h *GamecoordinatorApiHandler) GetMatchDetails(ctx context.Context, req *pb.MatchDetailsRequest) (*pb.MatchDetailsResponse, error) {
	sc := &share_code.ShareCodeData{MatchID: req.MatchId, OutcomeID: req.OutcomeId, Token: req.Token}
	matchDetails := <-h.service.RequestMatchDetails(sc)
	if matchDetails == nil {
		return nil, status.Error(codes.NotFound, "no match details could be received. Maybe try again later.")
	}

	return &pb.MatchDetailsResponse{
		MatchId:     matchDetails.MatchId,
		MatchTime:   timestamppb.New(matchDetails.MatchTime),
		DownloadUrl: matchDetails.DownloadUrl,
	}, nil
}
