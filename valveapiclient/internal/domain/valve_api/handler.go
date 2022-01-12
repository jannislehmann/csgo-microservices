package valve_api

import (
	"context"

	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
	"github.com/Cludch/csgo-microservices/valveapiclient/pkg/valve_match_api"
	pb "github.com/Cludch/csgo-microservices/valveapiclient/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValveMatchApiHandler struct {
	consumer valve_match_api.ValveMatchApiConsumerUseCase
}

func NewValveMatchApiHandler(c valve_match_api.ValveMatchApiConsumerUseCase) *ValveMatchApiHandler {
	return &ValveMatchApiHandler{
		consumer: c,
	}
}

func (h *ValveMatchApiHandler) GetNextShareCode(ctx context.Context, req *pb.ShareCodeRequest) (*pb.ShareCode, error) {
	// Consume valve's api to get the next sharecode.
	encoded, err := h.consumer.RequestNextShareCode(req.ApiKey, req.GetSteamId(), req.MatchHistoryAuthCode, req.PreviousShareCode)
	if err != nil {
		return nil, err
	}

	if encoded == "" {
		return nil, status.Error(codes.NotFound, "no new share code was found.")
	}

	// Decode the encoded share code to extract further information.
	sc, err := share_code.Decode(encoded)
	if err != nil {
		return nil, err
	}

	return &pb.ShareCode{
		Encoded:   sc.Encoded,
		MatchId:   sc.MatchID,
		OutcomeId: sc.OutcomeID,
		Token:     sc.Token,
	}, nil
}
