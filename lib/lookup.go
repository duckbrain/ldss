package lib

import (
	"errors"

	"github.com/duckbrain/ldss/lib/dl"
)

var ErrNotFound error

func init() {
	ErrNotFound = errors.New("Path not found")
}

// Lookup finds an Item by it's path. Expects a fully qualified path. "/" will
// return the catalog. Will return an error if there is an error
// loading the item or it is not downloaded.
func (r Reference) Lookup() (Item, error) {
	// TODO: Special case for / where it returns a combined catalog

	var item Item

	for srcName, src := range srcs {
		lang := languageFromSource(r.Lang, srcName)
		i, err := src.Lookup(lang, r.Path)
		if err == nil {
			item = i
		} else if err != ErrNotFound {
			return nil, err
		}
	}

	if item == nil {
		return nil, ErrNotFound
	}

	if x, ok := item.(dl.Downloader); ok && !x.Downloaded() {
		if err := dl.EnqueueAndWait(x); err != nil {
			return item, err
		}
	}

	if x, ok := item.(Opener); ok {
		if err := x.Open(); err != nil {
			return item, err
		}
	}

	return item, nil
}
