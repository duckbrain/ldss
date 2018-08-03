package ldsorg

import (
	"compress/zlib"
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
var _ lib.Item = &catalog{}

// Represents a catalog, exports ways to lookup children by path and ID
type catalog struct {
	jsonFolder
	dl.Template
	lang          lib.Lang
	itemsByPath   map[string]lib.Item
	foldersById   map[int]*folder
	foldersByPath map[string]*folder
	booksById     map[int]*book
	booksByPath   map[string]*book
}

// Creates a catalog object and populates it with it's Folders and Books
func newCatalog(lang *lang) *catalog {
	c := &catalog{}
	c.lang = lang
	c.JsonName = fmt.Sprintf("All %v Content", lang.EnglishName())
	c.Template.Src = getServerAction(fmt.Sprintf("catalog.query&languageid=%v&platformid=%v", lang.ID, platformID))
	c.Template.Dest = catalogPath(lang)

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
		return dl.ErrNotDownloaded(c)
	}
	if err = json.NewDecoder(file).Decode(&desc); err != nil {
		return err
	}

	c.itemsByPath = make(map[string]lib.Item)
	c.foldersById = make(map[int]*folder)
	c.foldersByPath = make(map[string]*folder)
	c.booksById = make(map[int]*book)
	c.booksByPath = make(map[string]*book)

	c.traverseFolders(c.Folders, c)
	c.traverseBooks(c.Books, c)

	return nil
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
	return c.lang
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
	for i, f := range f.Folders {
		items[i] = f
	}
	for i, f := range f.Books {
		items[folderLen+i] = f
	}
	return items
}

// Recursively converts jsonFolders into Folders
func (catalog *catalog) traverseFolders(folders []*folder, parent lib.Item) {
	langCode := catalog.Lang().Code()
	for _, f := range folders {
		catalog.traverseFolders(f.Folders, f)
		catalog.traverseBooks(f.Books, f)
		catalog.foldersById[f.ID] = f
		catalog.foldersByPath[f.Path()] = f
		catalog.foldersByPath[fmt.Sprintf("/%v", f.ID)] = f
		f.catalog = catalog
		f.path = f.computePath()
		itemsByLangAndPath[ref{langCode, f.path}] = f
	}
}

// Converts books into Books and sets their parent item
func (catalog *catalog) traverseBooks(books []*book, parent lib.Item) {
	langCode := catalog.Lang().Code()
	for _, b := range books {
		catalog.booksById[b.ID] = b
		if _, ok := catalog.booksByPath[b.GlURI]; !ok {
			catalog.booksByPath[b.GlURI] = b
			itemsByLangAndPath[ref{langCode, b.Path()}] = b
		}
		b.catalog = catalog
		b.Template.Src = b.DownloadURL
		b.Template.Dest = bookPath(b)
		b.Template.Transform = zlib.NewReader
	}
}
