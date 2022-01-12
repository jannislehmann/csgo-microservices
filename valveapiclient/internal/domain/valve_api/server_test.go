package valve_api_test

import (
	"context"
	"net"
	"testing"

	"github.com/Cludch/csgo-microservices/valveapiclient/mocks"
	pb "github.com/Cludch/csgo-microservices/valveapiclient/proto"
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
	valveMatchApiHandlerMock *mocks.ValveMatchApiHandlerUseCase
}

func (suite *MainTestSuite) SetupTest() {
	suite.valveMatchApiHandlerMock = new(mocks.ValveMatchApiHandlerUseCase)

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterValveMatchApiServiceServer(s, suite.valveMatchApiHandlerMock)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func (suite *MainTestSuite) TestGetNextShareCode() {
	suite.valveMatchApiHandlerMock.On("GetNextShareCode", mock.Anything, mock.Anything).Return(&pb.ShareCode{Encoded: ""}, nil)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		suite.T().Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewValveMatchApiServiceClient(conn)
	resp, err := client.GetNextShareCode(ctx, &pb.ShareCodeRequest{})
	if err != nil {
		suite.T().Fatalf("GetUser failed: %v", err)
	}
	log.Printf("Response: %+v", resp)

	suite.Nil(err)
	suite.NotNil(resp)
}
