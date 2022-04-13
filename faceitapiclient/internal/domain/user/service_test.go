package user_test

import (
	"errors"
	"testing"

	"github.com/Cludch/csgo-microservices/faceitapiclient/internal/domain/user"
	"github.com/Cludch/csgo-microservices/faceitapiclient/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	repositoryMock  *mocks.UserRepository
	apiConsumerMock *mocks.FaceitApiConsumerUseCase
	service         *user.UserService
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.repositoryMock = new(mocks.UserRepository)
	suite.apiConsumerMock = new(mocks.FaceitApiConsumerUseCase)
	suite.service = user.NewService(suite.repositoryMock, suite.apiConsumerMock)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) TestGetUser() {
	suite.repositoryMock.On("Find", TestID).Return(&user.User{ID: TestID}, nil)

	u, err := suite.service.GetUser(TestID)
	suite.Nil(err)
	suite.NotNil(u)
	suite.Equal(u.ID, TestID)
}

func (suite *ServiceTestSuite) TestGetUser_NotFound() {
	suite.repositoryMock.On("Find", TestID).Return(nil, errors.New(""))

	u, err := suite.service.GetUser(TestID)
	suite.Nil(u)
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestCreateUser() {
	suite.repositoryMock.On("Find", TestID).Return(nil, nil)
	suite.repositoryMock.On("Create", &user.User{ID: TestID}).Return(nil)

	u, err := suite.service.CreateUser(TestID)
	suite.NotNil(u)
	suite.Equal(u.ID, TestID)
	suite.Nil(err)
}

func (suite *ServiceTestSuite) TestCreateUser_Exists() {
	suite.repositoryMock.On("Find", TestID).Return(&user.User{ID: TestID}, nil)

	u, err := suite.service.CreateUser(TestID)
	suite.Nil(u)
	suite.NotNil(err)
}

func (suite *ServiceTestSuite) TestGetUsersWithApiEnabled() {
	testUser := &user.User{ID: TestID}
	suite.repositoryMock.On("ListAllWithApiEnabled").Return([]*user.User{testUser}, nil)

	u, err := suite.service.GetUsersWithApiEnabled()
	suite.NotNil(u)
	suite.Equal(len(u), 1)
	suite.Equal(u[0], testUser)
	suite.Nil(err)
}

func (suite *ServiceTestSuite) TestUpdateFaceitApiUsage_Enable() {
	testUser := &user.User{ID: TestID}
	suite.repositoryMock.On("UpdateFaceitApiUsage", testUser).Return(nil)

	err := suite.service.UpdateFaceitApiUsage(testUser, true)
	suite.Nil(err)
	suite.Equal(testUser.ApiEnabled, true)
}

func (suite *ServiceTestSuite) TestUpdateFaceitApiUsage_Disable() {
	testUser := &user.User{ID: TestID}
	suite.repositoryMock.On("UpdateFaceitApiUsage", testUser).Return(nil)

	err := suite.service.UpdateFaceitApiUsage(testUser, false)
	suite.Nil(err)
	suite.Equal(testUser.ApiEnabled, false)
}

func (suite *ServiceTestSuite) TestUpdateFaceitApiUsage_Error() {
	testUser := &user.User{ID: TestID, ApiEnabled: true}
	suite.repositoryMock.On("UpdateFaceitApiUsage", testUser).Return(errors.New(""))

	err := suite.service.UpdateFaceitApiUsage(testUser, false)
	suite.NotNil(err)
	suite.NotEqual(testUser.ApiEnabled, false)
}
