// This package provides functions to download content from the https://lds.org
// and lookup scriptures, magazines, General Conference talks, and much more from
// the downloaded content. It exports a variety of functions to access all the
// content.
package lib

import (
	"context"
	"errors"
	"fmt"
	"html/template"
)

type Lang string

type Index struct {
	Path string
	Lang Lang
}

func (i Index) Valid() bool {
	return i.Lang != "" && i.Path != ""
}
func (i Index) String() string {
	return fmt.Sprintf("%v?lang=%v", i.Path, i.Lang)
}
func (i Index) Hash() []byte {
	return []byte(i.String())
}

type Header struct {
	Name        string
	Subtitle    string
	SectionName string
	ShortTitle  string
	Index
}

type Item struct {
	Header

	Children []Index
	Parent   Index
	Next     Index
	Prev     Index

	Content   Content
	Footnotes []Footnote
}

type Footnote struct {
	Name     string        `json:"name"`
	LinkName string        `json:"linkName"`
	Content  template.HTML `json:"content"`
}

type Result struct {
	Item
	Rank int64
}
type Results []Result

func (r Results) Len() int {
	return len(r)
}
func (r Results) Less(a, b int) bool {
	return r[a].Rank < r[b].Rank
}
func (r Results) Swap(a, b int) {
	r[a], r[b] = r[b], r[a]
}

type Storer interface {
	Item(ctx context.Context, index Index) (Item, error)
	Store(ctx context.Context, item Item) error
	Header(ctx context.Context, index Index) (Header, error)
	Metadata(ctx context.Context, index Index, data interface{}) error
	SetMetadata(ctx context.Context, index Index, data interface{}) error
	// Remove(ctx context.Context, item Item) error
}
type BulkStorer interface {
	BulkRead(func(Storer) error) error
	BulkEdit(func(Storer) error) error
}
type Indexer interface {
	Index(ctx context.Context, item Item) error
	Search(ctx context.Context, ref Reference, results chan<- Result) error
}

type Loader interface {
	Load(ctx context.Context, store Storer, index Index) error
}

var ErrNotFound = errors.New("not found")
