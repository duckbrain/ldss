// This package provides functions to download content from the https://lds.org
// and lookup scriptures, magazines, General Conference talks, and much more from
// the downloaded content. It exports a variety of functions to access all the
// content.
package lib

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const DefaultLang Lang = "en"

type Lang string
type langDesc struct {
	Name   string
	GLCode string
}

func ParseLang(s string) (lang Lang, err error) {
	var ok bool
	lang, ok = languageAliases[strings.ToLower(s)]
	if !ok {
		err = ErrNotFound
	}
	return
}

func (l Lang) Name() string {
	if desc, ok := languageDescs[l]; ok {
		return desc.Name
	}
	return fmt.Sprintf("%v (unknown language", l)
}
func (l Lang) GLCode() string {
	if desc, ok := languageDescs[l]; ok {
		return desc.GLCode
	}
	return ""
}

func (l Lang) String() string {
	return string(l)
}

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
	// IsLoaded is set to true when the item has been fully downloaded. Loaders
	// can store non-loaded items to reference partially loaded Items that should
	// be loaded fully/directly before display.
	IsLoaded bool
	Index
}

type Item struct {
	Header

	Breadcrumbs []Header
	Children    []Header
	Parent      Header
	Next        Header
	Prev        Header
	Media       []Media

	Content   Content
	Footnotes []Footnote

	// Metadata may be used by sources for storing information.
	Metadata map[string]string
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
