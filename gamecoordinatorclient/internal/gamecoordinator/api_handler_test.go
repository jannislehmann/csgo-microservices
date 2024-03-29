package gamecoordinator_test

import (
	"context"
	"testing"
	"time"

	"github.com/Cludch/csgo-microservices/gamecoordinatorclient/internal/gamecoordinator"
	"github.com/Cludch/csgo-microservices/gamecoordinatorclient/mocks"
	pb "github.com/Cludch/csgo-microservices/gamecoordinatorclient/proto"
	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
	shared "github.com/Cludch/csgo-microservices/shared/proto"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var testShareCode, _ = share_code.Decode("CSGO-Z6wMP-JeoHt-C23L7-HTJ7B-feQ3A")

type HandlerTestSuite struct {
	suite.Suite
	serviceMock *mocks.GamecoordinatorUseCase
	handler     *gamecoordinator.GamecoordinatorApiHandler
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.serviceMock = new(mocks.GamecoordinatorUseCase)
	suite.handler = gamecoordinator.NewGamecoordinatorApiHandler(suite.serviceMock)
	testShareCode.Encoded = ""
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestGetMatchDetails() {
	channel := make(chan *shared.MatchDetails)
	suite.serviceMock.On("RequestMatchDetails", testShareCode).Return(channel)

	details := &shared.MatchDetails{MatchId: 1, MatchTime: timestamppb.New(time.Now()), DownloadUrl: "download"}

	go func() {
		time.Sleep(100 * time.Microsecond)
		channel <- details
	}()

	res, err := suite.handler.GetMatchDetails(context.TODO(),
		&pb.MatchDetailsRequest{MatchId: testShareCode.MatchID, OutcomeId: testShareCode.OutcomeID, Token: testShareCode.Token})

	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal(details.MatchId, res.MatchId)
	suite.Equal(details.MatchTime, res.MatchTime)
	suite.Equal(details.DownloadUrl, res.DownloadUrl)
}

func (suite *HandlerTestSuite) TestGetMatchDetails_Nil() {
	channel := make(chan *shared.MatchDetails)
	suite.serviceMock.On("RequestMatchDetails", testShareCode).Return(channel)

	go func() {
		time.Sleep(100 * time.Microsecond)
		channel <- nil
	}()

	res, err := suite.handler.GetMatchDetails(context.TODO(),
		&pb.MatchDetailsRequest{MatchId: testShareCode.MatchID, OutcomeId: testShareCode.OutcomeID, Token: testShareCode.Token})

	suite.NotNil(err)
	suite.Nil(res)
}
