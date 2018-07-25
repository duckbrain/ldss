package dl

import (
	"fmt"
)

func ErrNotDownloaded(d Downloader) error {
	return errNotDownloaded{dl: d}
}

func IsNotDownloaded(err error) (Downloader, bool) {
	dlErr, ok := err.(errNotDownloaded)
	return dlErr.dl, ok
}

type errNotDownloaded struct {
	dl Downloader
}

func (err errNotDownloaded) Error() string {
	return fmt.Sprintf("%v is not downloaded", err.dl.Name())
}
