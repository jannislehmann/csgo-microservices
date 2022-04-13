package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Cludch/csgo-microservices/faceitapiclient/internal/domain/user"
	"github.com/Cludch/csgo-microservices/faceitapiclient/mocks"
	pb "github.com/Cludch/csgo-microservices/faceitapiclient/proto"
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
		ID:         TestID,
		ApiEnabled: true,
	}, nil)

	u, err := suite.handler.GetUser(context.TODO(), &pb.GetUserRequest{Id: TestID.String()})
	suite.Nil(err)
	suite.NotNil(u)
	suite.Equal(true, u.ApiEnabled)
}

func (suite *HandlerTestSuite) TestGetUser_NotFound() {
	suite.serviceMock.On("GetUser", TestID).Return(nil, errors.New(""))

	u, err := suite.handler.GetUser(context.TODO(), &pb.GetUserRequest{Id: TestID.String()})
	suite.Nil(u)
	suite.NotNil(err)
}

func (suite *HandlerTestSuite) TestGetUser_InvalidIdFormat() {
	u, err := suite.handler.GetUser(context.TODO(), &pb.GetUserRequest{Id: "h-i"})
	suite.Nil(u)
	suite.NotNil(err)
}

func (suite *HandlerTestSuite) TestCreateUser() {
	suite.serviceMock.On("CreateUser", TestID).Return(&user.User{
		ID:         TestID,
		ApiEnabled: false,
	}, nil)

	u, err := suite.handler.CreateUser(context.TODO(), &pb.CreateUserRequest{Id: TestID.String()})
	suite.Nil(err)
	suite.NotNil(u)
	suite.Equal(false, u.ApiEnabled)
}

func (suite *HandlerTestSuite) TestCreateUser_Error() {
	suite.serviceMock.On("CreateUser", TestID).Return(nil, errors.New(""))

	u, err := suite.handler.CreateUser(context.TODO(), &pb.CreateUserRequest{Id: TestID.String()})
	suite.Nil(u)
	suite.NotNil(err)
}

func (suite *HandlerTestSuite) TestCreateUser_InvalidIdFormat() {
	u, err := suite.handler.CreateUser(context.TODO(), &pb.CreateUserRequest{Id: "h-i"})
	suite.Nil(u)
	suite.NotNil(err)
}
