package user

import (
	"context"

	pb "github.com/Cludch/csgo-microservices/faceitapiclient/proto"
	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
)

// UserRepository defines repository functions for user entities.
type UserRepository interface {
	Create(*User) error

	Find(entity.ID) (*User, error)

	ListAllWithApiEnabled() ([]*User, error)

	UpdateFaceitApiUsage(*User) error
}

// UserUseCase defines the user service functions.
type UserUseCase interface {
	CreateUser(id entity.ID) (*User, error)

	GetUser(entity.ID) (*User, error)
	GetUsersWithApiEnabled() ([]*User, error)
}

// UserHandlerUseCase defines the external interface for the gRPC server.
type UserHandlerUseCase interface {
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error)
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error)
}
