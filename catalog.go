package main

import (
	"encoding/json"
	"fmt"
)

type CatalogLoader struct {
	content      Content
	foldersById  map[int]*Folder
	booksByGlURI map[string]*Book
	booksById    map[int]*Book
	catalog      *Catalog
	language     *Language
}

func NewCatalogLoader(lang *Language, content Content) *CatalogLoader {
	c := new(CatalogLoader)
	c.language = lang
	c.content = content
	return c
}

type glCatalogDescrpition struct {
	Catalog         *Catalog `json:"catalog"`
	Success         bool     `json:"success"`
	CoverArtBaseUrl string   `json:"cover_art_base_url"`
}

type Catalog struct {
	Name           string        `json:"name"`
	Folders        []*Folder     `json:"folders"`
	Books          []*Book       `json:"books"`
	DateChanged    string        `json:"date_changed"`
	DisplayOrder   int           `json:"display_order"`
	Sections       []interface{} `json:"sections"`
	FolderContents []interface{} `json:"folder_contents"`
}

func (c *Catalog) String() string {
	if c == nil {
		return "Catalog: nil"
	} else {
		return fmt.Sprintf("Catalog: %v {folders[%v] books[%v]}", c.Name, c.Folders, c.Books)
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

type Book struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	GlURI string `json:"gl_uri"`
}

func (b *Book) String() string {
	return fmt.Sprintf("Book: %v {%v}", b.Name, b.GlURI)
}

type Node struct {
	Name string
}

func (l *CatalogLoader) populateIfNeeded() {
	if l.catalog != nil {
		return
	}

	var description glCatalogDescrpition
	file := l.content.OpenRead(l.content.GetLanguagesPath())
	dec := json.NewDecoder(file)
	err := dec.Decode(&description)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Success: %v, Cover art: %v\n", description.Success, description.CoverArtBaseUrl)

	l.catalog = description.Catalog
	//l.addFolders(description.Catalog.Folders)
	//l.addBooks(description.Catalog.Books)
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
	}
}

func (l *CatalogLoader) GetCatalog() *Catalog {
	l.populateIfNeeded()
	return l.catalog
}
