package user

import (
	"context"

	pb "github.com/Cludch/csgo-microservices/faceitapiclient/proto"
	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
)

type UserHandler struct {
	userService UserUseCase
}

func NewUserHandler(u UserUseCase) *UserHandler {
	return &UserHandler{
		userService: u,
	}
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	id, err := entity.StringToID(req.GetId())
	if err != nil {
		return nil, err
	}

	u, err := h.userService.GetUser(id)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:         u.ID.String(),
		ApiEnabled: u.ApiEnabled,
	}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	id, err := entity.StringToID(req.GetId())
	if err != nil {
		return nil, err
	}
	u, err := h.userService.CreateUser(id)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id: u.ID.String(),
	}, nil
}
