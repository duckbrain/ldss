package ldsorg

import (
	"fmt"
	"strings"

	"github.com/duckbrain/ldss/lib"
)

type source struct{}

// Finds an Item by it's path. Expects a fully qualified path. "/" will
// return the catalog. Will return an error if there is an error
// loading the item or it is not downloaded.
func (s source) Lookup(lang lib.Lang, path string) (Item, error) {
	if err := r.Check(); err != nil {
		return nil, err
	}

	c, err := r.Language.Catalog()
	if err != nil {
		return nil, err
	}
	if r.Path == "/" {
		return c, nil
	}
	if folder, ok := c.foldersByPath[r.Path]; ok {
		return folder, nil
	}
	sections := strings.Split(r.Path, "/")
	if sections[0] != "" {
		return nil, fmt.Errorf("Invalid path \"%v\", must start with '/'", r.Path)
	}
	for i := 2; i <= len(sections); i++ {
		temppath := strings.Join(sections[0:i], "/")
		if book, ok := c.booksByPath[temppath]; ok {
			if r.Path == book.Path() {
				return book, nil
			}
			node, err := book.lookupPath(r.Path)
			if err != nil {
				return nil, fmt.Errorf("Path %v not found - %v", r.Path, err.Error())
			}
			return node, err
		}
	}
	return nil, fmt.Errorf("Path \"%v\" not found", r.Path)
}
