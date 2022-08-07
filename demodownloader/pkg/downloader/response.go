package downloader

import (
	"errors"
	"fmt"
)

// ErrInvalidDownloadURL is returned when the url to be download is invalid or malicious.
var ErrInvalidDownloadURL = errors.New("invalid download url")

// DemoNotFoundError is used when a valid matchid / demo is not found or can no longer be downloaded.
type errDemoNotFound struct {
	URL string
}

func (e errDemoNotFound) Error() string {
	const msg = "demo no longer downloadable: %s"
	return fmt.Sprintf(msg, e.URL)
}

func IsDemoNotFoundError(err error) bool {
	_, ok := err.(errDemoNotFound)
	return ok
}
