package downloader

import "time"

type DownloaderUseCase interface {
	DownloadDemo(url string, demoDir string, lastModified time.Time) error
}
