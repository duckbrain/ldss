package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var _ strconv.NumError
var _ strings.Reader

var ErrNotFound error

func init() {
	ErrNotFound = errors.New("Item not found")
}

type Catalog struct {
	folderBase
	language      *Language
	foldersById   map[int]*Folder
	foldersByPath map[string]*Folder
	booksById     map[int]*Book
	booksByPath   map[string]*Book
}

func (c *Catalog) String() string {
	return fmt.Sprintf("%v {folders[%v] books[%v]}", c.base.Name, len(c.base.Folders), len(c.base.Books))
}

func (c *Catalog) Path() string {
	return "/"
}

func (c *Catalog) Parent() Item {
	return nil
}

func (c *Catalog) Next() Item {
	return nil
}

func (c *Catalog) Previous() Item {
	return nil
}

func (c *Catalog) Language() *Language {
	return c.language
}

func newCatalog(lang *Language) (*Catalog, error) {
	var desc jsonCatalogBase
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

type glCatalogDescrpition struct {
	Catalog         *Catalog `json:"catalog"`
	CoverArtBaseUrl string   `json:"cover_art_base_url"`
}

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

func (l *Catalog) BookByUnknown(id string) (*Book, error) {
	for _, book := range l.booksById {
		if book.Name() == id || fmt.Sprintf("%v", book.ID) == id || book.URL() == id || book.Path() == id {
			return book, nil
		}
	}
	return nil, errors.New("Book not found")
}

func (c *Catalog) Lookup(id string) (Item, error) {
	if id[0] == '/' {
		return c.LookupPath(id)
	} else {
		return nil, errors.New("Non-path lookup not implemented")
	}
}

func (c *Catalog) LookupBook(q string) (*Book, error) {
	i, err := c.Lookup(q)
	if err != nil {
		return nil, err
	}
	book, ok := i.(*Book)
	if !ok {
		return nil, fmt.Errorf("Result \"%v\" is not a book", i)
	}
	return book, nil
}

func (c *Catalog) LookupPath(path string) (Item, error) {
	if path == "" {
		return nil, fmt.Errorf("Cannot use empty string as a path")
	}
	if path == "/" {
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
