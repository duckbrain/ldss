package ldslib

import "fmt"

type CatalogItem interface {
	DisplayName() string
	//Children() []CatalogItem
	String() string
	Path() string
	Language() *Language
}

/*
 * Catalog
 */

type Catalog struct {
	Name     string    `json:"name"`
	Folders  []*Folder `json:"folders"`
	Books    []*Book   `json:"books"`
	language *Language
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

func (c Catalog) Language() *Language {
	return c.language
}

/*
 * Folder
 */

type Folder struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Folders  []*Folder `json:"folders"`
	Books    []*Book   `json:"books"`
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

/*
 * Book
 */

type Book struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	GlURI    string `json:"gl_uri"`
	Catalog  *Catalog
}

func (b Book) String() string {
	return fmt.Sprintf("Book: %v {%v}", b.Name, b.GlURI)
}

func (b Book) DisplayName() string {
	return b.Name
}

func (b Book) Path() string {
	return b.GlURI
}

func (b Book) Language() *Language {
	return b.Catalog.language
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