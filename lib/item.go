// This package provides functions to download content from the https://lds.org
// and lookup scriptures, magazines, General Conference talks, and much more from
// the downloaded content. It exports a variety of functions to access all the
// content.
package lib

import (
	"context"
	"errors"
	"html/template"
)

type Index struct {
	Path string
	Lang Lang
}

func (i Index) Valid() bool {
	return i.Lang != "" && i.Path != ""
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

type Store interface {
	Item(ctx context.Context, index Index) (Item, error)
	Store(ctx context.Context, item Item) error
	Header(ctx context.Context, index Index) (Header, error)
	Metadata(ctx context.Context, index Index, data interface{}) error
	SetMetadata(ctx context.Context, data interface{}) error
	Search(ctx context.Context, query string, results chan<- Result) error
	// Remove(ctx context.Context, item Item) error
}

type Source interface {
	Load(ctx context.Context, store Store, index Index) error
}

var ErrNotFound = errors.New("not found")
