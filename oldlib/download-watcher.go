package lib

type DownloadInfo struct {
	err      error
	Item     Item
	Language Lang
	Started  bool
	Progress float64
}

type DownloadReceiver func([]DownloadInfo)

var downloadUpdate chan DownloadInfo
var downloadSubscribers []DownloadReceiver

func init() {
	go func() {
		dl := make([]DownloadInfo, 0)
		for uInfo := range downloadUpdate {
			var found bool
			for index, oInfo := range dl {
				if uInfo.Item == oInfo.Item && uInfo.Language == oInfo.Language {
					dl[index] = uInfo
					found = true
				} else if oInfo.Progress == 1 || oInfo.err != nil {
					dl = append(dl[:index], dl[index+1:]...)
				}
			}
			if !found {
				dl = append(dl, uInfo)
			}

			//TODO Send a copy instead of the original slice
			for _, receiver := range downloadSubscribers {
				receiver(dl)
			}
		}

	}()
}
