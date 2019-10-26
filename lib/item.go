// This package provides functions to download content from the https://lds.org
// and lookup scriptures, magazines, General Conference talks, and much more from
// the downloaded content. It exports a variety of functions to access all the
// content.
package lib

import (
	"context"
	"errors"

	"github.com/duckbrain/ldss/lib"
)

type Index struct {
	Path string
	Lang Lang
}

type Header struct {
	Name        string
	Subtitle    string
	SectionName string
	Index
}

type Item struct {
	Header

	Children []Header
	Parent   *Header
	Next     *Header
	Prev     *Header

	Content   string
	Footnotes []Footnote
}

type Result struct {
	Item
	Rank
}

type Store interface {
	Item(ctx context.Context, index Index) (Item, error)
	Store(ctx context.Context, item Item) error
	Metadata(ctx context.Context, index Index, data interface{}) error
	SetMetadata(ctx context.Context, data interface{}) error
	Search(ctx context.Context, query string, results chan<- Result) error
	// Remove(ctx context.Context, item Item) error
}

type Source interface {
	Load(ctx context.Context, store lib.Store, index lib.Index) error
}

var ErrNotFound = errors.New("not found")
