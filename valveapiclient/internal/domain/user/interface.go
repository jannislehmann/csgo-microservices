package user

import (
	"context"

	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
	pb "github.com/Cludch/csgo-microservices/valveapiclient/proto"
)

// UserRepository defines repository functions for user entities.
type UserRepository interface {
	Create(*User) error

	Find(uint64) (*User, error)

	ListAllWithApiEnabled() ([]*User, error)

	UpdateLatestShareCode(*User) error
	UpdateMatchAuthCode(*User) error
	UpdateSteamApiUsage(*User) error
	UpdateApiKey(*User) error
}

// UserUseCase defines the user service functions.
type UserUseCase interface {
	CreateUser(id uint64) (*User, error)

	GetUser(uint64) (*User, error)
	GetUsersWithApiEnabled() ([]*User, error)

	AddSteamMatchHistoryAuthenticationCode(user *User, apiKey string, authCode string, sc string) error
	UpdateSteamApiUsage(*User, bool) error
	UpdateLatestShareCode(*User, *share_code.ShareCodeData) error

	QueryLatestShareCode(*User) (*share_code.ShareCodeData, error)
}

// UserHandlerUseCase defines the external interface for the gRPC server.
type UserHandlerUseCase interface {
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error)
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error)
	UpdateUserApiCredentials(ctx context.Context, req *pb.UpdateUserApiCredentialsRequest) (*pb.StatusResponse, error)
}
