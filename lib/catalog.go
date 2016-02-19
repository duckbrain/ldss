package lib

import (
	"encoding/json"
	"fmt"
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
	file, err := source.Open(source.CatalogPath(lang))
	if err != nil {
		dlErr := NotDownloadedCatalogErr{lang: lang}
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
		catalog.booksByPath[base.GlURI] = b
	}
	return books
}

// Finds an Item by it's path. Expects a fully qualified path. An empty string
// or "/" will return this catalog. Will return an error if there is an error
// loading the item or it is not downloaded.
func (c *Catalog) LookupPath(path string) (Item, error) {
	path = strings.TrimSpace(path)
	if path == "" || path == "/" {
		return c, nil
	}
	path = strings.TrimRight(path, "/ ")
	if folder, ok := c.foldersByPath[path]; ok {
		return folder, nil
	}
	sections := strings.Split(path, "/")
	if sections[0] != "" {
		return nil, fmt.Errorf("Invalid path \"%v\", must start with '/'", path)
	}
	for i := 2; i <= len(sections); i++ {
		temppath := strings.Join(sections[0:i], "/")
		if book, ok := c.booksByPath[temppath]; ok {
			if path == book.Path() {
				return book, nil
			}
			return book.lookupPath(path)
		}
	}
	return nil, fmt.Errorf("Path \"%v\" not found", path)
}
