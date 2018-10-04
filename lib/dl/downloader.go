package dl

// The directory to read and store the gospel library cache and configurations
var DataDirectory = ".ldss"

type Downloader interface {
	Downloaded() bool
	Download(chan<- Status) error
	Name() string
	Hash() string
}
