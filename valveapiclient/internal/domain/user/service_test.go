package user_test

import (
	"errors"
	"testing"

	"github.com/Cludch/csgo-microservices/shared/pkg/share_code"
	"github.com/Cludch/csgo-microservices/valveapiclient/internal/domain/user"
	"github.com/Cludch/csgo-microservices/valveapiclient/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	repositoryMock  *mocks.UserRepository
	apiConsumerMock *mocks.ValveMatchApiConsumerUseCase
	service         *user.UserService
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.repositoryMock = new(mocks.UserRepository)
	suite.apiConsumerMock = new(mocks.ValveMatchApiConsumerUseCase)
	suite.service = user.NewService(suite.repositoryMock, suite.apiConsumerMock)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) TestGetUser() {
	suite.repositoryMock.On("Find", TestId).Return(&user.User{ID: TestId}, nil)

	u, err := suite.service.GetUser(TestId)
	suite.Nil(err)
	suite.NotNil(u)
	suite.Equal(u.ID, TestId)
}

func (suite *ServiceTestSuite) TestGetUser_NotFound() {
	suite.repositoryMock.On("Find", TestId).Return(nil, errors.New(""))

	u, err := suite.service.GetUser(TestId)
	suite.Nil(u)
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestCreateUser() {
	suite.repositoryMock.On("Find", TestId).Return(nil, nil)
	suite.repositoryMock.On("Create", &user.User{ID: TestId}).Return(nil)

	u, err := suite.service.CreateUser(TestId)
	suite.NotNil(u)
	suite.Equal(u.ID, TestId)
	suite.Nil(err)
}

func (suite *ServiceTestSuite) TestCreateUser_Exists() {
	suite.repositoryMock.On("Find", TestId).Return(&user.User{ID: TestId}, nil)

	u, err := suite.service.CreateUser(TestId)
	suite.Nil(u)
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestGetUsersWithApiEnabled() {
	testUser := &user.User{ID: TestId}
	suite.repositoryMock.On("ListAllWithApiEnabled").Return([]*user.User{testUser}, nil)

	u, err := suite.service.GetUsersWithApiEnabled()
	suite.NotNil(u)
	suite.Equal(len(u), 1)
	suite.Equal(u[0], testUser)
	suite.Nil(err)
}

func (suite *ServiceTestSuite) TestUpdateLatestShareCode() {
	testUser := &user.User{ID: TestId}
	sc, _ := share_code.Decode(TestShareCode)
	suite.repositoryMock.On("UpdateLatestShareCode", testUser).Return(nil)

	err := suite.service.UpdateLatestShareCode(testUser, sc)
	suite.Nil(err)
	suite.Equal(testUser.LastShareCode, sc.Encoded)
}

func (suite *ServiceTestSuite) TestUpdateLatestShareCode_Error() {
	testUser := &user.User{ID: TestId}
	sc, _ := share_code.Decode(TestShareCode)
	suite.repositoryMock.On("UpdateLatestShareCode", testUser).Return(errors.New(""))

	err := suite.service.UpdateLatestShareCode(testUser, sc)
	suite.NotNil(err)
	suite.NotEqual(testUser.LastShareCode, sc.Encoded)
}

func (suite *ServiceTestSuite) TestUpdateSteamApiUsage_Enable() {
	testUser := &user.User{ID: TestId, AuthCode: TestAuthCode}
	suite.repositoryMock.On("UpdateSteamApiUsage", testUser).Return(nil)

	err := suite.service.UpdateSteamApiUsage(testUser, true)
	suite.Nil(err)
	suite.Equal(testUser.ApiEnabled, true)
}

func (suite *ServiceTestSuite) TestUpdateSteamApiUsage_Disable() {
	testUser := &user.User{ID: TestId, AuthCode: TestAuthCode}
	suite.repositoryMock.On("UpdateSteamApiUsage", testUser).Return(nil)

	err := suite.service.UpdateSteamApiUsage(testUser, false)
	suite.Nil(err)
	suite.Equal(testUser.ApiEnabled, false)
}

func (suite *ServiceTestSuite) TestUpdateSteamApiUsage_MissingAuthCode() {
	testUser := &user.User{ID: TestId, ApiEnabled: false}
	suite.repositoryMock.On("UpdateSteamApiUsage", testUser).Return(nil)

	err := suite.service.UpdateSteamApiUsage(testUser, true)
	suite.NotNil(err)
	suite.Equal(testUser.ApiEnabled, false)
}

func (suite *ServiceTestSuite) TestUpdateSteamApiUsage_DisableWithMissingAuthCode() {
	testUser := &user.User{ID: TestId, ApiEnabled: true}
	suite.repositoryMock.On("UpdateSteamApiUsage", testUser).Return(nil)

	err := suite.service.UpdateSteamApiUsage(testUser, false)
	suite.Nil(err)
	suite.Equal(testUser.ApiEnabled, false)
}

func (suite *ServiceTestSuite) TestUpdateSteamApiUsage_Error() {
	testUser := &user.User{ID: TestId, ApiEnabled: true}
	suite.repositoryMock.On("UpdateSteamApiUsage", testUser).Return(errors.New(""))

	err := suite.service.UpdateSteamApiUsage(testUser, false)
	suite.NotNil(err)
	suite.NotEqual(testUser.ApiEnabled, false)
}

func (suite *ServiceTestSuite) TestAddSteamMatchHistoryAuthenticationCode() {
	testUser := &user.User{ID: TestId}
	suite.apiConsumerMock.On("RequestNextShareCode", TestApiKey, testUser.ID, TestAuthCode, TestShareCode).Return("", nil)

	suite.repositoryMock.On("UpdateApiKey", testUser).Return(nil)
	suite.repositoryMock.On("UpdateMatchAuthCode", testUser).Return(nil)
	suite.repositoryMock.On("UpdateLatestShareCode", testUser).Return(nil)
	suite.repositoryMock.On("UpdateSteamApiUsage", testUser).Return(nil)

	err := suite.service.AddSteamMatchHistoryAuthenticationCode(testUser, TestApiKey, TestAuthCode, TestShareCode)
	suite.Nil(err)
	suite.Equal(TestAuthCode, testUser.AuthCode)
	suite.Equal(TestApiKey, testUser.ApiKey)
	suite.Equal(TestShareCode, testUser.LastShareCode)
	suite.Equal(true, testUser.ApiEnabled)
}

func (suite *ServiceTestSuite) TestAddSteamMatchHistoryAuthenticationCode_InvalidShareCode() {
	testUser := &user.User{ID: TestId}
	err := suite.service.AddSteamMatchHistoryAuthenticationCode(testUser, TestApiKey, TestAuthCode, "invalid")
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestAddSteamMatchHistoryAuthenticationCode_ConsumerError() {
	testUser := &user.User{ID: TestId}
	suite.apiConsumerMock.On("RequestNextShareCode", TestApiKey, testUser.ID, TestAuthCode, TestShareCode).Return("", errors.New(""))

	err := suite.service.AddSteamMatchHistoryAuthenticationCode(testUser, TestApiKey, TestAuthCode, TestShareCode)
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestAddSteamMatchHistoryAuthenticationCode_UpdateApiKeyError() {
	testUser := &user.User{ID: TestId}
	suite.apiConsumerMock.On("RequestNextShareCode", TestApiKey, testUser.ID, TestAuthCode, TestShareCode).Return("", errors.New(""))
	suite.repositoryMock.On("UpdateApiKey", testUser).Return(errors.New(""))

	err := suite.service.AddSteamMatchHistoryAuthenticationCode(testUser, TestApiKey, TestAuthCode, TestShareCode)
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestAddSteamMatchHistoryAuthenticationCode_UpdateMatchAuthCodeError() {
	testUser := &user.User{ID: TestId}
	suite.apiConsumerMock.On("RequestNextShareCode", TestApiKey, testUser.ID, TestAuthCode, TestShareCode).Return("", errors.New(""))
	suite.repositoryMock.On("UpdateApiKey", testUser).Return(nil)
	suite.repositoryMock.On("UpdateMatchAuthCode", testUser).Return(errors.New(""))

	err := suite.service.AddSteamMatchHistoryAuthenticationCode(testUser, TestApiKey, TestAuthCode, TestShareCode)
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestAddSteamMatchHistoryAuthenticationCode_UpdatLatestShareCodeError() {
	testUser := &user.User{ID: TestId}
	suite.apiConsumerMock.On("RequestNextShareCode", TestApiKey, testUser.ID, TestAuthCode, TestShareCode).Return("", errors.New(""))
	suite.repositoryMock.On("UpdateApiKey", testUser).Return(nil)
	suite.repositoryMock.On("UpdateMatchAuthCode", testUser).Return(nil)
	suite.repositoryMock.On("UpdateLatestShareCode", testUser).Return(errors.New(""))

	err := suite.service.AddSteamMatchHistoryAuthenticationCode(testUser, TestApiKey, TestAuthCode, TestShareCode)
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestAddSteamMatchHistoryAuthenticationCode_UpdatSteamApiUsageError() {
	testUser := &user.User{ID: TestId}
	suite.apiConsumerMock.On("RequestNextShareCode", TestApiKey, testUser.ID, TestAuthCode, TestShareCode).Return("", errors.New(""))
	suite.repositoryMock.On("UpdateApiKey", testUser).Return(nil)
	suite.repositoryMock.On("UpdateMatchAuthCode", testUser).Return(nil)
	suite.repositoryMock.On("UpdateLatestShareCode", testUser).Return(nil)
	suite.repositoryMock.On("UpdateSteamApiUsage", testUser).Return(errors.New(""))

	err := suite.service.AddSteamMatchHistoryAuthenticationCode(testUser, TestApiKey, TestAuthCode, TestShareCode)
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestQueryLatestShareCode() {
	testUser := &user.User{ID: TestId, ApiEnabled: true, AuthCode: TestAuthCode, LastShareCode: TestShareCode, ApiKey: TestApiKey}
	suite.apiConsumerMock.On("RequestNextShareCode", testUser.ApiKey, testUser.ID, testUser.AuthCode, testUser.LastShareCode).Return(TestShareCode, nil)

	sc, err := suite.service.QueryLatestShareCode(testUser)
	suite.NotNil(sc)
	suite.Equal(sc.Encoded, TestShareCode)
	suite.Nil(err)
}

func (suite *ServiceTestSuite) TestQueryLatestShareCode_ApiDisabled() {
	testUser := &user.User{ID: TestId, ApiEnabled: false, AuthCode: TestAuthCode, LastShareCode: TestShareCode, ApiKey: TestApiKey}

	sc, err := suite.service.QueryLatestShareCode(testUser)
	suite.Nil(sc)
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestQueryLatestShareCode_ConsumerError() {
	testUser := &user.User{ID: TestId, ApiEnabled: true, AuthCode: TestAuthCode, LastShareCode: TestShareCode, ApiKey: TestApiKey}
	suite.apiConsumerMock.On("RequestNextShareCode", testUser.ApiKey, testUser.ID, testUser.AuthCode, testUser.LastShareCode).Return("", errors.New(""))
	suite.repositoryMock.On("UpdateSteamApiUsage", testUser).Return(nil)

	sc, err := suite.service.QueryLatestShareCode(testUser)
	suite.repositoryMock.AssertNumberOfCalls(suite.T(), "UpdateSteamApiUsage", 1)
	suite.Nil(sc)
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestQueryLatestShareCode_NoNewMatch() {
	testUser := &user.User{ID: TestId, ApiEnabled: true, AuthCode: TestAuthCode, LastShareCode: TestShareCode, ApiKey: TestApiKey}
	suite.apiConsumerMock.On("RequestNextShareCode", testUser.ApiKey, testUser.ID, testUser.AuthCode, testUser.LastShareCode).Return("", nil)

	sc, err := suite.service.QueryLatestShareCode(testUser)
	suite.Nil(sc)
	suite.Nil(err)
}
