package download

type Info struct {
	err      error
	Item     Item
	Language Lang
	Started  bool
	Progress float64
}

type Status struct {
	Err      error
	Progress float64
}

type Receiver func([]Info)

var downloadUpdate chan Info
var downloadSubscribers []Receiver

func init() {
	go func() {
		dl := make([]Info, 0)
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

// Returns a slice of Downloaders that are currently queued and a map of Downloaders
func Queue() ([]Downloader, map[Download]<-chan Status) {
}

// Enqueue queues a download and waits for a result
func Enqueue(d Downloader) (ret <-chan Status) {
	return nil
}
