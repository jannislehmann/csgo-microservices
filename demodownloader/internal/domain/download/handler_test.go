package download_test

import (
	"context"
	"testing"

	"github.com/Cludch/csgo-microservices/demodownloader/internal/domain/download"
	"github.com/Cludch/csgo-microservices/demodownloader/mocks"
	pb "github.com/Cludch/csgo-microservices/demodownloader/proto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	demoUrl = "demoUrl.bz2"
)

type HandlerTestSuite struct {
	suite.Suite
	serviceMock *mocks.DownloaderUseCase
	handler     *download.DownloadHandler
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.serviceMock = new(mocks.DownloaderUseCase)
	suite.handler = download.NewDownloadHandler(suite.serviceMock, "/tmp")
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestDownloadDemoRequest() {
	suite.serviceMock.On("DownloadDemo", demoUrl, mock.Anything, mock.Anything).Return(nil, nil)
	_, err := suite.handler.DownloadDemo(context.TODO(), &pb.DownloadDemoRequest{DemoUrl: demoUrl})
	suite.Nil(err)
}
