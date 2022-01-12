package user

import (
	"context"

	pb "github.com/Cludch/csgo-microservices/valveapiclient/proto"
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
	u, err := h.userService.GetUser(req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:            u.ID,
		ApiEnabled:    u.ApiEnabled,
		ApiKey:        u.ApiKey,
		AuthCode:      u.AuthCode,
		LastShareCode: u.LastShareCode,
	}, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	u, err := h.userService.CreateUser(req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id: u.ID,
	}, nil
}

func (h *UserHandler) UpdateUserApiCredentials(ctx context.Context, req *pb.UpdateUserApiCredentialsRequest) (*pb.StatusResponse, error) {
	u, err := h.userService.GetUser(req.GetId())
	if err != nil {
		// Do not leak information about non-existing users for now.
		return &pb.StatusResponse{}, nil
	}

	errAdd := h.userService.AddSteamMatchHistoryAuthenticationCode(u, req.ApiKey, req.AuthCode, req.LastShareCode)
	if errAdd != nil {
		return &pb.StatusResponse{}, errAdd
	}

	return &pb.StatusResponse{
		Success: true,
	}, nil
}
