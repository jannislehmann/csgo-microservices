package user_test

import (
	"context"
	"net"
	"testing"

	"github.com/Cludch/csgo-microservices/faceitapiclient/mocks"
	pb "github.com/Cludch/csgo-microservices/faceitapiclient/proto"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type MainTestSuite struct {
	suite.Suite
	userHandlerMock *mocks.UserHandlerUseCase
}

func (suite *MainTestSuite) SetupTest() {
	suite.userHandlerMock = new(mocks.UserHandlerUseCase)

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, suite.userHandlerMock)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func (suite *MainTestSuite) TestGetUser() {
	suite.userHandlerMock.On("GetUser", mock.Anything, mock.Anything).Return(&pb.User{}, nil)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		suite.T().Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	resp, err := client.GetUser(ctx, &pb.GetUserRequest{})
	if err != nil {
		suite.T().Fatalf("GetUser failed: %v", err)
	}
	log.Printf("Response: %+v", resp)

	suite.Nil(err)
	suite.NotNil(resp)
}

func (suite *MainTestSuite) TestCreateUser() {
	suite.userHandlerMock.On("CreateUser", mock.Anything, mock.Anything).Return(&pb.User{}, nil)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		suite.T().Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)
	resp, err := client.CreateUser(ctx, &pb.CreateUserRequest{})
	if err != nil {
		suite.T().Fatalf("GetUser failed: %v", err)
	}
	log.Printf("Response: %+v", resp)

	suite.Nil(err)
	suite.NotNil(resp)
}
