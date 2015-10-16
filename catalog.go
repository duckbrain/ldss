package main

import (
	"encoding/json"
	"fmt"
)

type CatalogLoader struct {
	content      Content
	foldersById  map[int]*Folder
	booksById    map[int]*Book
	booksByGlURI map[string]*Book
	catalog      *Catalog
	language     *Language
}

func NewCatalogLoader(lang *Language, content Content) *CatalogLoader {
	c := new(CatalogLoader)
	c.language = lang
	c.content = content
	return c
}

type CatalogItem interface {
	GetName() string
}

type glCatalogDescrpition struct {
	Catalog         *Catalog `json:"catalog"`
	Success         bool     `json:"success"`
	CoverArtBaseUrl string   `json:"cover_art_base_url"`
}

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

type Node struct {
	Name string
}

func (l *CatalogLoader) populateIfNeeded() {
	if l.catalog != nil {
		return
	}

	var description glCatalogDescrpition
	file := l.content.OpenRead(l.content.GetCatalogPath(l.language))
	dec := json.NewDecoder(file)
	err := dec.Decode(&description)
	if err != nil {
		panic(err)
	}

	l.foldersById = make(map[int]*Folder)
	l.booksById = make(map[int]*Book)
	l.booksByGlURI = make(map[string]*Book)

	l.catalog = description.Catalog
	l.addFolders(description.Catalog.Folders)
	l.addBooks(description.Catalog.Books)
}

func (l *CatalogLoader) addFolders(folders []*Folder) {
	for _, f := range folders {
		l.foldersById[f.ID] = f
		l.addFolders(f.Folders)
		l.addBooks(f.Books)
	}
}

func (l *CatalogLoader) addBooks(books []*Book) {
	for _, b := range books {
		l.booksById[b.ID] = b
		l.booksByGlURI[b.GlURI] = b
		b.Language = l.language
	}
}

func (l *CatalogLoader) GetCatalog() *Catalog {
	l.populateIfNeeded()
	return l.catalog
}

func (l *CatalogLoader) GetFolderById(id int) *Folder {
	l.populateIfNeeded()
	return l.foldersById[id]
}

func (l *CatalogLoader) GetBookById(id int) *Book {
	l.populateIfNeeded()
	return l.booksById[id]
}

func (l *CatalogLoader) GetBookByGlURI(glUri string) *Book {
	l.populateIfNeeded()
	return l.booksByGlURI[glUri]
}

func (l *CatalogLoader) GetBookByUnknown(id string) *Book {
	gl := ParseForBook(id)

	if gl != "" {
		return l.GetBookByGlURI(gl)
	}

	return l.GetBookByGlURI(id)
}
