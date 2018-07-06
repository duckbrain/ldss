package lib

type Downloader interface {
	Downloaded() bool
	Download() error
}
