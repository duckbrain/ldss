package ldslib

import "fmt"

type CatalogItem interface {
	Name() string
	Children() []CatalogItem
	String() string
}

/*
 * Catalog
 */

type Catalog struct {
	Name    string    `json:"name"`
	Folders []*Folder `json:"folders"`
	Books   []*Book   `json:"books"`
}

func (c *Catalog) String() string {
	if c == nil {
		return "Catalog: nil"
	} else {
		return fmt.Sprintf("Catalog: %v {folders[%v] books[%v]}", c.Name, len(c.Folders), len(c.Books))
	}
}

/*
 * Folder
 */

type Folder struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Folders []*Folder `json:"folders"`
	Books   []*Book   `json:"books"`
}

func (f *Folder) String() string {
	return fmt.Sprintf("Folder: %v {folders[%v] books[%v]}", f.Name, len(f.Folders), len(f.Books))
}

func (f *Folder) GetName() string {
	return f.Name
}

/*
 * Book
 */

type Book struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	GlURI    string `json:"gl_uri"`
	Language *Language
}

func (b *Book) String() string {
	return fmt.Sprintf("Book: %v {%v}", b.Name, b.GlURI)
}

func (b *Book) GetName() string {
	return b.Name
}

/*
 * Node
 */

type Node struct {
	ID       int
	Name     string
	GlURI    string
	Language *Language
}

func (n *Node) GetName() string {
	return n.Name
}

func (n *Node) String() string {
	return n.Name
}
