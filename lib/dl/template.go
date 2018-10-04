package dl

import (
	"io"
	"net/http"
	"os"
)

type Template struct {
	Src       string
	Dest      string
	Transform func(io.Reader) (io.ReadCloser, error)
}

func (t Template) Downloaded() bool {
	_, err := os.Stat(t.Dest)
	return !os.IsNotExist(err)
}

// https://gist.github.com/albulescu/e61979cc852e4ee8f49c to download with progress
// The best way to do it is to create a reader that can eavesdrop on the transform
// reading, so compression, etc. won't matter. The gist has a good example of
// getting the content length. If the length cannot be gotten, we just won't send
// status updates. (This is good because renderers should display this correctly).
func (t Template) Download(chan<- Status) error {
	var input io.Reader

	response, err := http.Get(t.Src)
	if err != nil {
		return err
	}
	body := response.Body
	defer response.Body.Close()

	if t.Transform != nil {
		in, err := t.Transform(body)
		if err != nil {
			return err
		}
		defer in.Close()
		input = in
	} else {
		input = body
	}

	file, err := os.Create(t.Dest)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, input)
	return err
}
