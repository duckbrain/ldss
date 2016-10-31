package lib

var downloads []DownloadInfo
var downloadSubscribers []DownloadNotifier

type DownloadInfo NotDownloadedErr

type DownloadNotifier interface {
	Notify([]DownloadInfo)
}

func init() {
	downloads = make([]DownloadInfo, 0)
	go func() {

	}()
}

func downloadNotifyAll() {
}

func downloadNotifyStart(d DownloadInfo) {
}

func downloadNotifyEnd(d DownloadInfo) {
}
