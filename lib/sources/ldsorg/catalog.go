package ldsorg

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/duckbrain/ldss/lib"
	"github.com/duckbrain/ldss/lib/dl"
)

var _ strconv.NumError
var _ strings.Reader
var _ lib.Item = catalog{}

// Represents a catalog, exports ways to lookup children by path and ID
type catalog struct {
	jsonFolder
	dl.Template
	lang          Lang
	itemsByPath   map[string]lib.Item
	foldersById   map[int]*folder
	foldersByPath map[string]*folder
	booksById     map[int]*book
	booksByPath   map[string]*book
}

// Creates a catalog object and populates it with it's Folders and Books
func newCatalog(lang Lang) *catalog {
	c := &catalog{}
	c.lang = lang
	c.JsonName = fmt.Sprintf("All %v Content", lang.EnglishName())
	c.Template.Src = getServerAction(fmt.Sprintf("catalog.query&languageid=%v&platformid=%v", lang.ID, platformID))
	c.Template.Dest = catalogPath(language)

	return c
}

func (c *catalog) Open() error {
	if c.itemsByPath != nil {
		return nil
	}
	var desc = struct {
		Catalog         *catalog `json:"catalog"`
		CoverArtBaseUrl string   `json:"cover_art_base_url"`
	}{
		Catalog: c,
	}
	file, err := os.Open(c.Template.Dest)
	if err != nil {
		return nil, dl.ErrNotDownloaded(c)
	}
	if err = json.NewDecoder(file).Decode(&desc); err != nil {
		return nil, err
	}

	c.base = desc.Catalog
	c.itemsByPath = make(map[int]lib.Item)
	c.foldersById = make(map[int]*folder)
	c.foldersByPath = make(map[string]*folder)
	c.booksById = make(map[int]*book)
	c.booksByPath = make(map[string]*book)
	c.language = lang

	c.folders = c.addFolders(c.base.Folders, c)
	c.books = c.addBooks(c.base.Books, c)

	return c, nil
}

// The Gospel Library Path of this catalog. Every catalog's path is "/"
func (c *catalog) Path() string {
	return "/"
}

// Parent of this catalog. Always nil
func (c *catalog) Parent() lib.Item {
	return nil
}

// Next sibling of this catalog. Always nil
func (c *catalog) Next() lib.Item {
	return nil
}

// Previous sibling of this catalog. Always nil
func (c *catalog) Prev() lib.Item {
	return nil
}

// The language this catalog represents
func (c *catalog) Lang() lib.Lang {
	return c.language
}

func (c *catalog) Children() []lib.Item {
}

// Used for parsing folders in the catalog's JSON file
type jsonFolder struct {
	JsonName string    `json:"name"`
	Folders  []*folder `json:"folders"`
	Books    []*book   `json:"books"`
	ID       int       `json:"id"`
}

func (f *jsonFolder) Name() string {
	return f.JsonName
}

// Combined Folders and Books
func (f *jsonFolder) Children() []lib.Item {
	folderLen := len(f.Folders)
	items := make([]lib.Item, folderLen+len(f.Books))
	for i, f := range f.Folders() {
		items[i] = f
	}
	for i, f := range f.Books() {
		items[folderLen+i] = f
	}
	return items
}

// Recursively converts jsonFolders into Folders
func (catalog *catalog) addFolders(jFolders []*jsonFolder, parent Item) []*Folder {
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
func (catalog *catalog) addBooks(jBooks []*jsonBook, parent Item) []*Book {
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
