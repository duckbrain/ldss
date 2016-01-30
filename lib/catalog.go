package lib

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrNotFound error

func init() {
	ErrNotFound = errors.New("Item not found")
}

type Catalog struct {
	base         *jsonCatalog
	language     *Language
	source       Source
	foldersById  map[int]*Folder
	booksById    map[int]*Book
	booksByGlURI map[string]*Book
}

func (c Catalog) Name() string {
	return c.base.Name
}

func (c Catalog) String() string {
	return fmt.Sprintf("%v {folders[%v] books[%v]}", c.base.Name, len(c.base.Folders), len(c.base.Books))
}

func (c Catalog) Path() string {
	return "/"
}

func (c Catalog) Parent() Item {
	return nil
}

func (c Catalog) Language() *Language {
	return c.language
}

func (f *Catalog) Children() ([]Item, error) {
	folderLen := len(f.base.Folders)
	items := make([]Item, folderLen+len(f.base.Books))
	for i, f := range f.Folders() {
		items[i] = f
	}
	for i, f := range f.Books() {
		items[folderLen+i] = f
	}
	return items, nil
}

func (f *Catalog) Folders() []*Folder {
	return f.base.Folders
}

func (f *Catalog) Books() []*Book {
	return f.base.Books
}

func newCatalog(lang *Language, source Source) (*Catalog, error) {
	var desc jsonCatalogBase
	file, err := source.Open(source.CatalogPath(lang))
	if err != nil {
		return nil, &NotDownloadedBookErr{err, lang}
	}
	if err = json.NewDecoder(file).Decode(&desc); err != nil {
		return nil, err
	}

	c := &Catalog{base: desc.Catalog}
	c.foldersById = make(map[int]*Folder)
	c.booksById = make(map[int]*Book)
	c.booksByGlURI = make(map[string]*Book)
	c.language = lang
	c.source = source

	c.addFolders(c.Folders(), c)
	c.addBooks(c.Books(), c)

	return c, nil
}

type glCatalogDescrpition struct {
	Catalog         *Catalog `json:"catalog"`
	CoverArtBaseUrl string   `json:"cover_art_base_url"`
}

func (l *Catalog) addFolders(folders []*Folder, parent Item) {
	for _, f := range folders {
		f.catalog = l
		f.parent = parent
		l.foldersById[f.ID()] = f
		l.addFolders(f.Folders(), f)
		l.addBooks(f.Books(), f)
	}
}

func (l *Catalog) addBooks(books []*Book, parent Item) {
	for _, b := range books {
		b.catalog = l
		b.parent = parent
		l.booksById[b.base.ID] = b
		l.booksByGlURI[b.Path()] = b
	}
}

func (l *Catalog) Folder(id int) (*Folder, error) {
	c, ok := l.foldersById[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (l *Catalog) Book(id int) (*Book, error) {
	c, ok := l.booksById[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (l *Catalog) BookByGlURI(glUri string) (*Book, error) {
	c, ok := l.booksByGlURI[glUri]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (l *Catalog) BookByUnknown(id string) (*Book, error) {
	for _, book := range l.booksById {
		if book.Name() == id || fmt.Sprintf("%v", book.ID) == id || book.URL() == id || book.Path() == id {
			return book, nil
		}
	}
	return nil, errors.New("Book not found")
}
