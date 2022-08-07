package download

import (
	"context"
	"time"

	"github.com/Cludch/csgo-microservices/demodownloader/pkg/downloader"
	pb "github.com/Cludch/csgo-microservices/demodownloader/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DownloadHandler struct {
	downloaderService downloader.DownloaderUseCase
	demoDir           string
}

func NewDownloadHandler(d downloader.DownloaderUseCase, demoDir string) *DownloadHandler {
	return &DownloadHandler{
		downloaderService: d,
		demoDir:           demoDir,
	}
}

func (h *DownloadHandler) DownloadDemo(ctx context.Context, req *pb.DownloadDemoRequest) (*emptypb.Empty, error) {
	if err := h.downloaderService.DownloadDemo(req.GetDemoUrl(), h.demoDir, time.Now()); err != nil {
		return nil, err
	}

	return nil, nil
}
