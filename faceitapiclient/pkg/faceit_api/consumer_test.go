package faceit_api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/Cludch/csgo-microservices/faceitapiclient/pkg/faceit_api"
	shared_mocks "github.com/Cludch/csgo-microservices/shared/mocks"
	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	apiKey    = "api"
	playerID  = entity.NewID()
	matchID   = "1-07079673-6318-44b5-9aca-2f10d7cfb440"
	startTime = int64(1)
)

type ConsumerTestSuite struct {
	suite.Suite
	httpClientMock *shared_mocks.HttpClient
	consumer       *faceit_api.FaceitApiConsumerService
}

func (suite *ConsumerTestSuite) SetupTest() {
	suite.httpClientMock = new(shared_mocks.HttpClient)
	suite.consumer = faceit_api.New(suite.httpClientMock)
}

func TestConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ConsumerTestSuite))
}

func (suite *ConsumerTestSuite) TestGetMatchHistory() {
	response := &faceit_api.PlayerMatchHistoryResponse{
		Result: []*faceit_api.PlayerMatchHistoryEntry{
			{
				MatchId: matchID,
			},
		},
	}
	b, _ := json.Marshal(response)
	r := io.NopCloser(bytes.NewReader(b))
	suite.httpClientMock.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK, Body: r}, nil)

	res, err := suite.consumer.GetPlayerMatchHistory(apiKey, playerID)
	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal(res.Result[0].MatchId, matchID)
}

func (suite *ConsumerTestSuite) TestGetMatchDetail() {
	response := &faceit_api.MatchDetailsResponse{
		DemoUrl: []string{
			"url",
		},
		StartTime: startTime,
		Status:    "FINISHED",
	}
	b, _ := json.Marshal(response)
	r := io.NopCloser(bytes.NewReader(b))
	suite.httpClientMock.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK, Body: r}, nil)

	res, err := suite.consumer.GetMatchDetails(apiKey, matchID)
	suite.Nil(err)
	suite.NotNil(res)
	suite.Equal("url", res.DemoUrl[0])
	suite.Equal(startTime, res.StartTime)
	suite.Equal("FINISHED", res.Status)
}
