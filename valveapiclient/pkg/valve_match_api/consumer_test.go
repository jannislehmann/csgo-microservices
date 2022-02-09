package valve_match_api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Cludch/csgo-microservices/valveapiclient/mocks"
	"github.com/Cludch/csgo-microservices/valveapiclient/pkg/valve_match_api"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	id            = uint64(1)
	apiKey        = "api"
	authCode      = "auth"
	lastShareCode = "CSGO-Y4DVh-amkvh-OyBrh-SyMHN-2SvPB"
	nextShareCode = "CSGO-kYpP8-BDS3P-9EBqy-CzoGX-RQCmQ"
)

type ConsumerTestSuite struct {
	suite.Suite
	httpClientMock *mocks.ApiClient
	consumer       *valve_match_api.ValveMatchApiConsumerService
}

func (suite *ConsumerTestSuite) SetupTest() {
	suite.httpClientMock = new(mocks.ApiClient)
	suite.consumer = valve_match_api.New(suite.httpClientMock)
}

func TestConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ConsumerTestSuite))
}

func (suite *ConsumerTestSuite) TestRequestNextShareCode() {
	response := &valve_match_api.MatchResponse{}
	response.Result.Nextcode = nextShareCode
	b, _ := json.Marshal(response)
	r := io.NopCloser(bytes.NewReader(b))
	suite.httpClientMock.On("Get", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK, Body: r}, nil)

	sc, err := suite.consumer.RequestNextShareCode(apiKey, id, authCode, lastShareCode)
	suite.Nil(err)
	suite.NotNil(sc)
	suite.Equal(nextShareCode, sc)
}

func (suite *ConsumerTestSuite) TestRequestNextShareCode_Accepted() {
	suite.httpClientMock.On("Get", mock.Anything).Return(&http.Response{StatusCode: http.StatusAccepted, Body: io.NopCloser(strings.NewReader(""))}, nil)

	sc, err := suite.consumer.RequestNextShareCode(apiKey, id, authCode, lastShareCode)
	suite.Nil(err)
	suite.NotNil(sc)
	suite.Equal("", sc)
}

func (suite *ConsumerTestSuite) TestRequestNextShareCode_GetError() {
	suite.httpClientMock.On("Get", mock.Anything).Return(nil, errors.New(""))

	sc, err := suite.consumer.RequestNextShareCode(apiKey, id, authCode, lastShareCode)
	suite.NotNil(err)
	suite.NotNil(sc)
	suite.Equal("", sc)
}

func (suite *ConsumerTestSuite) TestRequestNextShareCode_Forbidden() {
	suite.httpClientMock.On("Get", mock.Anything).Return(&http.Response{StatusCode: http.StatusForbidden, Body: io.NopCloser(strings.NewReader(""))}, nil)

	sc, err := suite.consumer.RequestNextShareCode(apiKey, id, authCode, lastShareCode)
	suite.NotNil(err)
	var invalidApiKeyOrAuthCode *valve_match_api.InvalidApiKeyOrAuthCode
	suite.ErrorAs(err, &invalidApiKeyOrAuthCode)
	suite.NotNil(sc)
	suite.Equal("", sc)
}

func (suite *ConsumerTestSuite) TestRequestNextShareCode_NoNewData() {
	suite.httpClientMock.On("Get", mock.Anything).Return(&http.Response{StatusCode: http.StatusPreconditionFailed, Body: io.NopCloser(strings.NewReader(""))}, nil)

	sc, err := suite.consumer.RequestNextShareCode(apiKey, id, authCode, lastShareCode)
	suite.NotNil(err)
	var invalidSteamIDErr *valve_match_api.InvalidSteamID
	suite.ErrorAs(err, &invalidSteamIDErr)
	suite.NotNil(sc)
	suite.Equal("", sc)
}

func (suite *ConsumerTestSuite) TestRequestNextShareCode_InternalServerError() {
	suite.httpClientMock.On("Get", mock.Anything).Return(&http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))}, nil)

	sc, err := suite.consumer.RequestNextShareCode(apiKey, id, authCode, lastShareCode)
	suite.Nil(err)
	suite.NotNil(sc)
	suite.Equal("", sc)
}

func (suite *ConsumerTestSuite) TestRequestNextShareCode_GatewayTimeout() {
	suite.httpClientMock.On("Get", mock.Anything).Return(&http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))}, nil)

	sc, err := suite.consumer.RequestNextShareCode(apiKey, id, authCode, lastShareCode)
	suite.Nil(err)
	suite.NotNil(sc)
	suite.Equal("", sc)
}
