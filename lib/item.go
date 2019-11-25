// This package provides functions to download content from the https://lds.org
// and lookup scriptures, magazines, General Conference talks, and much more from
// the downloaded content. It exports a variety of functions to access all the
// content.
package lib

import (
	"context"
	"errors"
	"fmt"
)

const DefaultLang Lang = "en"

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

	Breadcrumbs []Index
	Children    []Index
	Parent      Index
	Next        Index
	Prev        Index
	Media       []Media

	Content   Content
	Footnotes []Footnote
}
type ItemDetails struct {
	Header

	Breadcrumbs []Header
	Children    []Header
	Parent      Header
	Next        Header
	Prev        Header

	Content   Content
	Footnotes []Footnote
}
type Media struct {
	Type string
	Desc string
	URL  string
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
	Clear(ctx context.Context) error
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

func (lib *Library) Register(l Loader) {
	lib.Sources = append(lib.Sources, l)
	ctx := lib.ctx(context.Background())
	if x, ok := l.(interface {
		LoadParser(context.Context, *ReferenceParser)
	}); ok {
		x.LoadParser(ctx, Default.Parser)
	}
}

var Default = &Library{
	Parser:  NewReferenceParser(),
	Sources: []Loader{},
}
