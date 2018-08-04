package ldsorg

import (
	"github.com/duckbrain/ldss/lib"
	"github.com/duckbrain/ldss/lib/dl"
)

// Represents a book in the catalog or one of it's folders and
// provides a way to access the nodes in it's database if
// downloaded.
// Used for parsing books in the catalog's JSON file
type book struct {
	ID          int    `json:"id"`
	JSONName    string `json:"name"`
	DownloadURL string `json:"url"`
	GlURI       string `json:"gl_uri"`
	catalog     *catalog
	parent      lib.Item
	dl.Template
	children []lib.Item
}

func (b *book) init(catalog *catalog, parent lib.Item) {
	b.catalog = catalog
	b.parent = parent
}

func (b *book) Open() error {
	path := bookPath(b)
	l, err := opendb(path)
	if err != nil {
		return err
	}
	defer l.Close()

	// Populate Children
	nodes, err := l.childrenByParentID(0, b, b)
	if err != nil {
		return err
	}
	b.children = nodes

	langCode := b.Lang().Code()
	for _, node := range nodes {
		itemsByLangAndPath[ref{langCode, node.Path()}] = node
	}

	return nil
}

// The name of this book.
func (b *book) Name() string {
	return b.JSONName
}

// The Gospel Library Path of this book, unique within it's language
func (b *book) Path() string {
	return b.GlURI
}

// The language this book is in
func (b *book) Lang() lib.Lang {
	return b.catalog.lang
}

// Children in this book.
func (b *book) Children() []lib.Item {
	return b.children
}

// Parent Folder or Catalog of this book
func (b *book) Parent() lib.Item {
	return b.parent
}

// Next book in the Folder
func (b *book) Next() lib.Item {
	return lib.GenericNextPrevious(b, 1)
}

// Previous book in the Folder
func (b *book) Prev() lib.Item {
	return lib.GenericNextPrevious(b, -1)
}
