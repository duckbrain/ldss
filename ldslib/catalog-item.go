package ldslib

import (
	"fmt"
	"html/template"
)

type CatalogItem interface {
	DisplayName() string
	Children() ([]CatalogItem, error)
	String() string
	Path() string
	Language() *Language
	Parent() CatalogItem
}

/*
 * Catalog
 */

type folder struct {
	Name     string    `json:"name"`
	Folders  []*Folder `json:"folders"`
	Books    []*Book   `json:"books"`
}

func (f *folder) Children() ([]CatalogItem, error) {
	folderLen := len(f.Folders)
	items := make([]CatalogItem, folderLen + len(f.Books))
	for i, f := range f.Folders {
		items[i] = f
	}
	for i, f := range f.Books {
		items[folderLen + i] = f
	}
	return items, nil
}

type Catalog struct {
	folder
	language *Language
	parser   *catalogParser
}

func (c Catalog) DisplayName() string {
	return c.Name
}

func (c Catalog) String() string {
	return fmt.Sprintf("Catalog: %v {folders[%v] books[%v]}", c.Name, len(c.Folders), len(c.Books))
}

func (c Catalog) Path() string {
	return "/"
}

func (c Catalog) Parent() CatalogItem {
	return nil
}

func (c Catalog) Language() *Language {
	return c.language
}

/*
 * Folder
 */

type Folder struct {
	ID       int       `json:"id"`
	folder
	parent   CatalogItem
	Catalog  *Catalog
}

func (f Folder) String() string {
	return fmt.Sprintf("Folder: %v {folders[%v] books[%v]}", f.Name, len(f.Folders), len(f.Books))
}

func (f Folder) DisplayName() string {
	return f.Name
}

func (f Folder) Path() string {
	return fmt.Sprintf("/%v", f.ID)
}

func (f Folder) Language() *Language {
	return f.Catalog.language
}

func (f Folder) Parent() CatalogItem {
	return f.parent
}

/*
 * Book
 */

type Book struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	GlURI    string `json:"gl_uri"`
	Catalog  *Catalog
	parser   *bookParser
	parent   CatalogItem
}

func (b *Book) String() string {
	return fmt.Sprintf("Book: %v {%v}", b.Name, b.GlURI)
}

func (b *Book) DisplayName() string {
	return b.Name
}

func (b *Book) Path() string {
	return b.GlURI
}

func (b *Book) Language() *Language {
	return b.Catalog.language
}

func (b *Book) Children() ([]CatalogItem, error) {
	if b.parser == nil {
		b.parser = newBookParser(b, b.Catalog.parser.source)
	}
	nodes, err := b.parser.Index()
	if err != nil {
		return nil, err
	}
	items := make([]CatalogItem, len(nodes))
	for i, n := range nodes {
		items[i] = n
	}
	return items, nil
}

func (b *Book) Parent() CatalogItem {
	return b.parent
}

/*
 * Node
 */

type Node struct {
	ID         int
	Name       string
	GlURI      string
	Book      *Book
	HasContent bool
	ChildCount int
	parent	   CatalogItem
}

func (n Node) DisplayName() string {
	return n.Name
}

func (n Node) String() string {
	return n.Name
}

func (n Node) Path() string {
	return n.GlURI
}

func (n Node) Language() *Language {
	return n.Book.Language()
}

func (n Node) Children() ([]CatalogItem, error) {
	nodes, err := n.Book.parser.Children(n)
	if err != nil {
		return nil, err
	}
	items := make([]CatalogItem, len(nodes))
	for i, n := range nodes {
		items[i] = n
	}
	return items, nil
}

func (n Node) Content() (template.HTML, error) {
	html, err := n.Book.parser.Content(n)
	return template.HTML(html), err
}

func (n Node) Parent() CatalogItem {
	return n.parent
}
