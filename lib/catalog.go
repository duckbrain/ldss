package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var _ strconv.NumError
var _ strings.Reader

// Represents a catalog, exports ways to lookup children by path and ID
type Catalog struct {
	folderBase
	language      *Language
	foldersById   map[int]*Folder
	foldersByPath map[string]*Folder
	booksById     map[int]*Book
	booksByPath   map[string]*Book
}

// A short human-readable representation of the catalog, mostly useful for debugging.
func (c *Catalog) String() string {
	return fmt.Sprintf("%v {folders[%v] books[%v]}", c.base.Name, len(c.base.Folders), len(c.base.Books))
}

// The Gospel Library Path of this catalog. Every catalog's path is "/"
func (c *Catalog) Path() string {
	return "/"
}

// Parent of this catalog. Always nil
func (c *Catalog) Parent() Item {
	return nil
}

// Next sibling of this catalog. Always nil
func (c *Catalog) Next() Item {
	return nil
}

// Previous sibling of this catalog. Always nil
func (c *Catalog) Previous() Item {
	return nil
}

// The language this catalog represents
func (c *Catalog) Language() *Language {
	return c.language
}

// Creates a catalog object and populates it with it's Folders and Books
func newCatalog(lang *Language) (*Catalog, error) {
	var desc struct {
		Catalog         *jsonFolder `json:"catalog"`
		CoverArtBaseUrl string      `json:"cover_art_base_url"`
	}
	file, err := os.Open(catalogPath(lang))
	if err != nil {
		dlErr := notDownloadedCatalogErr{lang: lang}
		dlErr.err = err
		return nil, dlErr
	}
	if err = json.NewDecoder(file).Decode(&desc); err != nil {
		return nil, err
	}

	c := &Catalog{}
	c.base = desc.Catalog
	c.foldersById = make(map[int]*Folder)
	c.foldersByPath = make(map[string]*Folder)
	c.booksById = make(map[int]*Book)
	c.booksByPath = make(map[string]*Book)
	c.language = lang

	c.folders = c.addFolders(c.base.Folders, c)
	c.books = c.addBooks(c.base.Books, c)

	return c, nil
}

// Used for parsing folders in the catalog's JSON file
type jsonFolder struct {
	Name    string        `json:"name"`
	Folders []*jsonFolder `json:"folders"`
	Books   []*jsonBook   `json:"books"`
	ID      int           `json:"id"`
}

// Used for parsing books in the catalog's JSON file
type jsonBook struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	GlURI string `json:"gl_uri"`
}

// Recursively converts jsonFolders into Folders
func (catalog *Catalog) addFolders(jFolders []*jsonFolder, parent Item) []*Folder {
	folders := make([]*Folder, len(jFolders))
	for i, base := range jFolders {
		f := &Folder{
			parent:  parent,
			catalog: catalog,
		}
		f.base = base
		f.folders = catalog.addFolders(base.Folders, f)
		f.books = catalog.addBooks(base.Books, f)
		folders[i] = f
		catalog.foldersById[base.ID] = f
		catalog.foldersByPath[f.Path()] = f
		catalog.foldersByPath[fmt.Sprintf("/%v", f.ID())] = f
	}
	return folders
}

// Converts jsonBooks into Books and sets their parent item
func (catalog *Catalog) addBooks(jBooks []*jsonBook, parent Item) []*Book {
	books := make([]*Book, len(jBooks))
	for i, base := range jBooks {
		b := newBook(base, catalog, parent)
		books[i] = b
		catalog.booksById[base.ID] = b
		if _, ok := catalog.booksByPath[base.GlURI]; !ok {
			catalog.booksByPath[base.GlURI] = b
		}
	}
	return books
}
