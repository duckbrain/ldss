package ldslib

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrNotFound error

func init() {
	ErrNotFound = errors.New("Item not found")
}

type catalogParser struct {
	source       Source
	foldersById  map[int]*Folder
	booksById    map[int]*Book
	booksByGlURI map[string]*Book
	catalog      *Catalog
	language     *Language
}

func newCatalogLoader(lang *Language, source Source) *catalogParser {
	c := new(catalogParser)
	c.language = lang
	c.source = source
	return c
}

type glCatalogDescrpition struct {
	Catalog         *Catalog `json:"catalog"`
	Success         bool     `json:"success"`
	CoverArtBaseUrl string   `json:"cover_art_base_url"`
}

func (l *catalogParser) populateIfNeeded() error {
	if l.catalog != nil {
		return nil
	}

	var description glCatalogDescrpition
	file, err := l.source.Open(l.source.CatalogPath(l.language))
	if err != nil {
		return err
	}
	dec := json.NewDecoder(file)
	err = dec.Decode(&description)
	if err != nil {
		return err
	}

	l.foldersById = make(map[int]*Folder)
	l.booksById = make(map[int]*Book)
	l.booksByGlURI = make(map[string]*Book)
	l.catalog = description.Catalog
	l.catalog.Language = l.language
	l.addFolders(description.Catalog.Folders)
	l.addBooks(description.Catalog.Books)

	return nil
}

func (l *catalogParser) addFolders(folders []*Folder) {
	for _, f := range folders {
		f.Language = l.language
		f.Catalog = l.catalog
		l.foldersById[f.ID] = f
		l.addFolders(f.Folders)
		l.addBooks(f.Books)
	}
}

func (l *catalogParser) addBooks(books []*Book) {
	for _, b := range books {
		b.Catalog = l.catalog
		l.booksById[b.ID] = b
		l.booksByGlURI[b.GlURI] = b
	}
}

func (l *catalogParser) Catalog() (*Catalog, error) {
	if err := l.populateIfNeeded(); err != nil {
		return nil, err
	}
	return l.catalog, nil
}

func (l *catalogParser) Folder(id int) (*Folder, error) {
	if err := l.populateIfNeeded(); err != nil {
		return nil, err
	}
	c, ok := l.foldersById[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (l *catalogParser) Book(id int) (*Book, error) {
	if err := l.populateIfNeeded(); err != nil {
		return nil, err
	}
	c, ok := l.booksById[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (l *catalogParser) BookByGlURI(glUri string) (*Book, error) {
	if err := l.populateIfNeeded(); err != nil {
		return nil, err
	}
	c, ok := l.booksByGlURI[glUri]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (l *catalogParser) BookByUnknown(id string) (*Book, error) {
	if err := l.populateIfNeeded(); err != nil {
		return nil, err
	}
	for _, book := range l.booksById {
		if book.Name == id || fmt.Sprintf("%v", book.ID) == id || book.URL == id || book.GlURI == id {
			return book, nil
		}
	}
	return nil, errors.New("Book not found")
}
