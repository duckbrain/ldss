package lib

import (
	"fmt"
	"strings"
)

func AutoLookup(lang *Language, q string) <-chan Message {
	return AutoDownload(func() (interface{}, error) {
		return Lookup(lang, q)
	})
}

func AutoLookupPath(lang *Language, q string) <-chan Message {
	return AutoDownload(func() (interface{}, error) {
		return LookupPath(lang, q)
	})
}

// Look up a user given string passing it through the reference parser.
func Lookup(lang *Language, q string) (Reference, error) {
	ref, err := lang.ref()
	if err == nil {
		q, err = ref.lookup(q)
		if err != nil {
			return nil, err
		}
	}
	item, err := LookupPath(lang, q)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Finds an Item by it's path. Expects a fully qualified path. An empty string
// or "/" will return this catalog. Will return an error if there is an error
// loading the item or it is not downloaded.
func LookupPath(lang *Language, path string) (Reference, error) {
	c, err := lang.Catalog()
	if err != nil {
		return nil, err
	}
	path = strings.TrimSpace(path)
	if path == "" || path == "/" {
		return c, nil
	}
	path = strings.TrimRight(path, "/ ")
	if folder, ok := c.foldersByPath[path]; ok {
		return folder, nil
	}
	sections := strings.Split(path, "/")
	if sections[0] != "" {
		return nil, fmt.Errorf("Invalid path \"%v\", must start with '/'", path)
	}
	for i := 2; i <= len(sections); i++ {
		temppath := strings.Join(sections[0:i], "/")
		if book, ok := c.booksByPath[temppath]; ok {
			if path == book.Path() {
				return book, nil
			}
			node := &Node{Book: b}
			db, err := b.db()
			if err != nil {
				return nil, err
			}
			err = db.stmtUri.QueryRow(uri).Scan(&node.id, &node.name, &node.path, &node.parentId, &node.hasContent, &node.childCount)
			if err != nil {
				return nil, fmt.Errorf("Path %v not found", uri)
			}
			return node, err

			return book.lookupPath(path)
		}
	}
	return nil, fmt.Errorf("Path \"%v\" not found", path)
}
