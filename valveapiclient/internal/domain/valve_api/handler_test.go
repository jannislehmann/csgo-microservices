package valve_api_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
	"github.com/Cludch/csgo-microservices/valveapiclient/internal/domain/valve_api"
	"github.com/Cludch/csgo-microservices/valveapiclient/mocks"
	pb "github.com/Cludch/csgo-microservices/valveapiclient/proto"
	"github.com/stretchr/testify/suite"
)

var (
	id            = uint64(1)
	apiKey        = "api"
	authCode      = "auth"
	lastShareCode = "CSGO-Y4DVh-amkvh-OyBrh-SyMHN-2SvPB"
	nextShareCode = "CSGO-kYpP8-BDS3P-9EBqy-CzoGX-RQCmQ"
)

type HandlerTestSuite struct {
	suite.Suite
	consumerMock *mocks.ValveMatchApiConsumerUseCase
	handler      *valve_api.ValveMatchApiHandler
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.consumerMock = new(mocks.ValveMatchApiConsumerUseCase)
	suite.handler = valve_api.NewValveMatchApiHandler(suite.consumerMock)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestGetNextShareCode() {
	decoded, _ := share_code.Decode(nextShareCode)
	suite.consumerMock.On("RequestNextShareCode", apiKey, id, authCode, lastShareCode).Return(nextShareCode, nil)

	sc, err := suite.handler.GetNextShareCode(context.TODO(), &pb.ShareCodeRequest{
		SteamId:              id,
		ApiKey:               apiKey,
		MatchHistoryAuthCode: authCode,
		PreviousShareCode:    lastShareCode,
	})
	suite.Nil(err)
	suite.NotNil(sc)
	suite.Equal(sc.Encoded, decoded.Encoded)
	suite.Equal(sc.MatchId, decoded.MatchID)
	suite.Equal(sc.OutcomeId, decoded.OutcomeID)
	suite.Equal(sc.Token, decoded.Token)
}

func (suite *HandlerTestSuite) TestGetNextShareCode_InvalidShareCodeResponse() {
	suite.consumerMock.On("RequestNextShareCode", apiKey, id, authCode, lastShareCode).Return("invalid share code", nil)

	sc, err := suite.handler.GetNextShareCode(context.TODO(), &pb.ShareCodeRequest{
		SteamId:              id,
		ApiKey:               apiKey,
		MatchHistoryAuthCode: authCode,
		PreviousShareCode:    lastShareCode,
	})
	suite.NotNil(err)
	suite.Nil(sc)
}

func (suite *HandlerTestSuite) TestGetNextShareCode_ConsumerError() {
	suite.consumerMock.On("RequestNextShareCode", apiKey, id, authCode, lastShareCode).Return("", errors.New(""))

	sc, err := suite.handler.GetNextShareCode(context.TODO(), &pb.ShareCodeRequest{
		SteamId:              id,
		ApiKey:               apiKey,
		MatchHistoryAuthCode: authCode,
		PreviousShareCode:    lastShareCode,
	})
	suite.NotNil(err)
	suite.Nil(sc)
}

func (suite *HandlerTestSuite) TestGetNextShareCode_NoNewShareCode() {
	emptyShareCode := ""
	suite.consumerMock.On("RequestNextShareCode", apiKey, id, authCode, lastShareCode).Return(emptyShareCode, nil)

	sc, err := suite.handler.GetNextShareCode(context.TODO(), &pb.ShareCodeRequest{
		SteamId:              id,
		ApiKey:               apiKey,
		MatchHistoryAuthCode: authCode,
		PreviousShareCode:    lastShareCode,
	})
	suite.NotNil(err)
	suite.Nil(sc)
}
