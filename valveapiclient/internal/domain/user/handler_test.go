package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Cludch/csgo-microservices/valveapiclient/internal/domain/user"
	"github.com/Cludch/csgo-microservices/valveapiclient/mocks"
	pb "github.com/Cludch/csgo-microservices/valveapiclient/proto"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	serviceMock *mocks.UserUseCase
	handler     *user.UserHandler
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.serviceMock = new(mocks.UserUseCase)
	suite.handler = user.NewUserHandler(suite.serviceMock)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestGetUser() {
	suite.serviceMock.On("GetUser", TestID).Return(&user.User{
		ID:            TestID,
		ApiEnabled:    true,
		ApiKey:        TestApiKey,
		AuthCode:      TestAuthCode,
		LastShareCode: TestShareCode,
	}, nil)

	u, err := suite.handler.GetUser(context.TODO(), &pb.GetUserRequest{Id: TestID})
	suite.Nil(err)
	suite.NotNil(u)
	suite.Equal(true, u.ApiEnabled)
	suite.Equal(TestApiKey, u.ApiKey)
	suite.Equal(TestAuthCode, u.AuthCode)
	suite.Equal(TestShareCode, u.LastShareCode)
}

func (suite *HandlerTestSuite) TestGetUser_NotFound() {
	suite.serviceMock.On("GetUser", TestID).Return(nil, errors.New(""))

	u, err := suite.handler.GetUser(context.TODO(), &pb.GetUserRequest{Id: TestID})
	suite.Nil(u)
	suite.NotNil(err)
}

func (suite *HandlerTestSuite) TestCreateUser() {
	suite.serviceMock.On("CreateUser", TestID).Return(&user.User{
		ID:         TestID,
		ApiEnabled: false,
	}, nil)

	u, err := suite.handler.CreateUser(context.TODO(), &pb.CreateUserRequest{Id: TestID})
	suite.Nil(err)
	suite.NotNil(u)
	suite.Equal(false, u.ApiEnabled)
}

func (suite *HandlerTestSuite) TestCreateUser_Error() {
	suite.serviceMock.On("CreateUser", TestID).Return(nil, errors.New(""))

	u, err := suite.handler.CreateUser(context.TODO(), &pb.CreateUserRequest{Id: TestID})
	suite.Nil(u)
	suite.NotNil(err)
}

func (suite *HandlerTestSuite) TestUpdateUserApiCredentials() {
	u := &user.User{
		ID:            TestID,
		ApiEnabled:    true,
		ApiKey:        TestApiKey,
		AuthCode:      TestAuthCode,
		LastShareCode: TestShareCode,
	}

	suite.serviceMock.On("GetUser", TestID).Return(u, nil)
	suite.serviceMock.On("AddSteamMatchHistoryAuthenticationCode", u, u.ApiKey, u.AuthCode, u.LastShareCode).Return(nil)

	status, err := suite.handler.UpdateUserApiCredentials(context.TODO(), &pb.UpdateUserApiCredentialsRequest{
		Id:            TestID,
		ApiKey:        TestApiKey,
		AuthCode:      TestAuthCode,
		LastShareCode: TestShareCode,
	})
	suite.NotNil(status)
	suite.Equal(true, status.Success)
	suite.Nil(err)
}

func (suite *HandlerTestSuite) TestUpdateUserApiCredentials_NotFound() {
	suite.serviceMock.On("GetUser", TestID).Return(nil, errors.New(""))

	status, err := suite.handler.UpdateUserApiCredentials(context.TODO(), &pb.UpdateUserApiCredentialsRequest{
		Id:            TestID,
		ApiKey:        TestApiKey,
		AuthCode:      TestAuthCode,
		LastShareCode: TestShareCode,
	})
	suite.NotNil(status)
	suite.Equal(false, status.Success)
	suite.Nil(err)
}

func (suite *HandlerTestSuite) TestUpdateUserApiCredentials_Error() {
	u := &user.User{
		ID:            TestID,
		ApiEnabled:    true,
		ApiKey:        TestApiKey,
		AuthCode:      TestAuthCode,
		LastShareCode: TestShareCode,
	}

	suite.serviceMock.On("GetUser", TestID).Return(u, nil)
	suite.serviceMock.On("AddSteamMatchHistoryAuthenticationCode", u, u.ApiKey, u.AuthCode, u.LastShareCode).Return(errors.New(""))

	status, err := suite.handler.UpdateUserApiCredentials(context.TODO(), &pb.UpdateUserApiCredentialsRequest{
		Id:            TestID,
		ApiKey:        TestApiKey,
		AuthCode:      TestAuthCode,
		LastShareCode: TestShareCode,
	})
	suite.NotNil(status)
	suite.Equal(false, status.Success)
	suite.NotNil(err)
}
