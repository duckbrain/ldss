package lib

import (
	"errors"
)

var ErrPathNotFound error

func init() {
	ErrPathNotFound = errors.New("Path not found")
}

// Lookup finds an Item by it's path. Expects a fully qualified path. "/" will
// return the catalog. Will return an error if there is an error
// loading the item or it is not downloaded.
func (r Reference) Lookup() (Item, error) {
	for srcName, src := range srcs {
		lang := LanguageFromSource(*r.Language, srcName)
		item, err := src.Lookup(*lang, r.Path)
		if err != ErrPathNotFound {
			return nil, err
		}
		if err == nil {
			return item, nil
		}
	}
	return nil, ErrPathNotFound
}
