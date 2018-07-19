package download

import (
	"sync"
)

type Stage byte

const (
	// Waiting indicates the download is queued, but not yet started.
	Waiting Stage = iota

	// Started indiciates the download has been started, but has not given
	// any status updates, so the Progress should be ignored.
	Started

	// Progress indicates that the download is progressing and Progress
	// is accurate.
	Progress

	// Complete indicates that the download is complete, and will receive
	// no further updates. Err should be checked to see if the download
	// was a success.
	Complete
)

type Status struct {
	Err      error
	Progress float64
}

type DownloadStatus struct {
	Downloader
	Status
	Stage Stage
}

var lock sync.Mutex
var allDownloads map[Downloader]DownloadStatus

var newDownloads chan Downloader
var statusUpdates chan DownloadStatus

var subscribers map[Downloader]chan<- Status
var queueSubscribers []chan<- DownloadStatus

func init() {
	allDownloads = make(map[Downloader]DownloadStatus)
	newDownloads = make(chan Downloader, 50)
	statusUpdates = make(chan DownloadStatus, 100)
	subscribers = make(map[Downloader]chan<- Status)
	queueSubscribers = make([]chan<- DownloadStatus, 0)
	go messageWorker()
	go worker()
}

// Listen subscribes to events in download status. All uncompleted downloads
// are immediately passed in so the subscriber will know of their existance
// and current state.
func Listen(updates chan<- DownloadStatus) {
	lock.Lock()

	queueSubscribers = append(queueSubscribers, updates)
	for _, ds := range allDownloads {
		updates <- ds
	}

	lock.Unlock()
}

// Enqueue queues a Downloader to be downloaded. If ret is passed it will
// receive updates for the download. Can be nil to ignore.
//
// NOTE: If multiple subscribers to a single download are wanted. That could be
// implemented in the future by calling the function again with the same
// Downloader, but for now, DO NOT PASS THE SAME DOWNLOADER MORE THAN ONCE.
func Enqueue(d Downloader, ret chan<- Status) {
	lock.Lock()

	allDownloads[d] = DownloadStatus{
		Downloader: d,
		Stage:      Waiting,
	}

	if ret != nil {
		subscribers[d] = ret
	}

	lock.Unlock()

	newDownloads <- d
}

// EnqueueAndWait enques a download, waits for completion, and returns its error if any.
func EnqueueAndWait(d Downloader) (err error) {
	c := make(chan Status, 20)
	Enqueue(d, c)
	for s := range c {
		err = s.Err
	}
	return
}

func messageWorker() {
	for s := range statusUpdates {
		lock.Lock()
		//TODO: If the channels are closed, the write will panic. This
		// is a good interface for unsubscribing, so we'll need to
		// handle that.
		for _, sub := range queueSubscribers {
			sub <- s
		}
		if sub, ok := subscribers[s.Downloader]; ok {
			sub <- s.Status
		}
		allDownloads[s.Downloader] = s
		if s.Stage == Complete {
			delete(allDownloads, s.Downloader)
			close(subscribers[s.Downloader])
			delete(subscribers, s.Downloader)
		}
		lock.Unlock()
	}
}

func worker() {
	for dl := range newDownloads {
		dls := make(chan Status)
		go func() {
			for s := range dls {
				statusUpdates <- DownloadStatus{
					Status:     s,
					Downloader: dl,
					Stage:      Progress,
				}
			}
		}()
		statusUpdates <- DownloadStatus{
			Downloader: dl,
			Stage:      Started,
		}
		err := dl.Download(dls)
		close(dls)
		statusUpdates <- DownloadStatus{
			Status: Status{
				Err:      err,
				Progress: 1,
			},
			Downloader: dl,
			Stage:      Complete,
		}
	}
}
