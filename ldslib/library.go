package ldslib

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"strconv"
)

type Library struct {
	source               Source
	languages            []Language
	catalogsByLanguageId map[int]*catalogParser
	booksByLangBookId    map[langBookID]*bookParser
}

type langBookID struct {
	langID int
	bookID int
}

func NewLibrary(source Source) *Library {
	p := &Library{}
	p.source = source
	p.catalogsByLanguageId = make(map[int]*catalogParser)
	p.booksByLangBookId = make(map[langBookID]*bookParser)
	return p
}

func (l *Library) populateLanguages() error {
	if l.languages != nil {
		return nil
	}

	var description glLanguageDescription
	file, err := l.source.Open(l.source.LanguagesPath())
	if err != nil {
		return err
	}
	dec := json.NewDecoder(file)
	err = dec.Decode(&description)
	if err != nil {
		return err
	}

	l.languages = description.Languages
	return nil
}

func (l *Library) populateCatalog(lang *Language) *catalogParser {
	c, ok := l.catalogsByLanguageId[lang.ID]
	if !ok {
		c = newCatalogLoader(lang, l.source)
		l.catalogsByLanguageId[lang.ID] = c
	}
	return c
}

func (l *Library) populateBook(book *Book) *bookParser {
	id := langBookID{book.Catalog.language.ID, book.ID}
	b, ok := l.booksByLangBookId[id]
	if !ok {
		b = newBookParser(book, l.source)
		l.booksByLangBookId[id] = b
		book.parser = b
	}
	return b
}

func (l *Library) Language(id string) (*Language, error) {
	if err := l.populateLanguages(); err != nil {
		return nil, err
	}
	for _, lang := range l.languages {
		if lang.Name == id || fmt.Sprintf("%v", lang.ID) == id || lang.EnglishName == id || lang.Code == id || lang.GlCode == id {
			return &lang, nil
		}
	}
	return nil, errors.New("Language not found")
}

func (l *Library) Languages() ([]Language, error) {
	return l.languages, l.populateLanguages()
}

func (l *Library) Catalog(lang *Language) (*Catalog, error) {
	return l.populateCatalog(lang).Catalog()
}
func (l *Library) Book(path string, catalog *Catalog) (*Book, error) {
	return l.populateCatalog(catalog.Language()).BookByUnknown(path)
}

func (l *Library) lookupGlURI(path string, catalog *Catalog) (CatalogItem, error) {
	c := l.populateCatalog(catalog.Language())
	if path == "/" {
		return c.Catalog()
	}
	if folderId, err := strconv.Atoi(path[1:]); err == nil {
		return c.Folder(folderId)
	}
	sections := strings.Split(path, "/")
	if sections[0] != "" {
		return nil, fmt.Errorf("Invalid path \"%v\", must start with '/'", path)
	}
	for i := 1; i <= len(sections); i++ {
		temppath := strings.Join(sections[0:i], "/")
		book, err := c.BookByGlURI(temppath)
		if err == nil {
			if path == book.GlURI {
				return book, nil
			}
			// Look for a node
			b := l.populateBook(book)
			return b.GlUri(path)
		}
	}
	return nil, fmt.Errorf("Path \"%v\" not found", path)
}

func (l *Library) Lookup(id string, catalog *Catalog) (CatalogItem, error) {
	if id[0] == '/' {
		return l.lookupGlURI(id, catalog)
	}
	p := NewRefParser(l, catalog)
	p.Load(id)
	return p.Item()
}

func (l *Library) Children(item CatalogItem) ([]CatalogItem, error) {
	switch item.(type) {
	case *Book:
		l.populateBook(item.(*Book))
		return item.Children()
	default:
		return item.Children()
	}
}

func (l *Library) Content(node Node) (string, error) {
	return l.populateBook(node.Book).Content(node)
}

//	Index(lang *Language) []CatalogItem
//	Children(item CatalogItem) []CatalogItem
