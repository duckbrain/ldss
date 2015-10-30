package ldslib

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Library struct {
	source               Source
	languages            []Language
	catalogsByLanguageId map[int]*catalogParser
}

func NewLibrary(source Source) *Library {
	p := &Library{}
	p.source = source
	p.catalogsByLanguageId = make(map[int]*catalogParser)
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
	if err := l.populateLanguages(); err != nil {
		return nil, err
	}
	return l.languages, nil
}

func (l *Library) Catalog(lang *Language) *Catalog {
	c, ok := l.catalogsByLanguageId[lang.ID]
	if !ok {
		c = newCatalogLoader(lang, l.source)
		l.catalogsByLanguageId[lang.ID] = c
	}
	return c.catalog
}
func (r *Library) Book(path string, catalog *Catalog) *Book {
	return nil
}

//	Index(lang *Language) []CatalogItem
//	Children(item CatalogItem) []CatalogItem
//	Lookup(path string) CatalogItem
