package valve_api

import (
	"context"

	pb "github.com/Cludch/csgo-microservices/valveapiclient/proto"
)

// ValveMatchApiHandlerUseCase defines the external interface for the gRPC server.
type ValveMatchApiHandlerUseCase interface {
	GetNextShareCode(ctx context.Context, req *pb.ShareCodeRequest) (*pb.ShareCode, error)
}
