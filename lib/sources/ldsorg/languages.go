package ldsorg

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var languages cache

type glLangDesc struct {
	Languages []*Lang `json:"languages"`
	Success   bool    `json:"success"`
}

// Lang defines a language as from the server. The fields should not be modified.
type Lang struct {
	// The Gospel Library ID for the language. Used for downloads.
	ID int `json:"id"`

	// Native representation of the language in the language observed
	Name string `json:"name"`

	// English representation of the language
	EnglishName string `json:"eng_name"`

	// The internationalization (i18n) code used in most programs
	Code string `json:"code"`

	// Gospel Library language code, seen in the urls of https://lds.org
	GlCode string `json:"code_three"`

	catalogCache cache
}

// Returns a human readable version of the language that is appropriate to
// show an end user. It will format the language to contain it's native
// representation as well as the English representation. It will also show
// the standard internationalization code as well as the Gospel Library
// language code.
func (l *Lang) String() string {
	var id, name, code string

	id = fmt.Sprintf("%v: ", l.ID)
	if l.Name == l.EnglishName {
		name = l.Name
	} else {
		name = fmt.Sprintf("%v (%v)", l.Name, l.EnglishName)
	}
	if l.Code == l.GlCode {
		code = fmt.Sprintf(" [%v]", l.Code)
	} else {
		code = fmt.Sprintf(" [%v/%v]", l.Code, l.GlCode)
	}

	return id + name + code
}

// Gets the catalog for this language. If cached, it will return the cached version.
func (l *Lang) Catalog() (*Catalog, error) {
	l.catalogCache.construct = func() (interface{}, error) {
		return newCatalog(l)
	}
	c, err := l.catalogCache.get()
	if err != nil {
		return nil, err
	}
	return c.(*Catalog), err
}

func init() {
	languages.construct = func() (interface{}, error) {
		var description glLangDesc
		file, err := os.Open(languagesPath())
		if err != nil {
			return nil, &notDownloadedLanguageErr{err}
		}
		err = json.NewDecoder(file).Decode(&description)
		return description.Languages, err
	}
}

// Returns a list of all languages available. Downloads the languages if not already downloaded first.
func Languages() ([]*Lang, error) {
	if !fileExist(languagesPath()) {
		if err := DownloadLanguages(); err != nil {
			return nil, err
		}
	}
	langs, err := languages.get()
	if err != nil {
		return nil, err
	}
	return langs.([]*Lang), err
}

// Finds a language by any of the accepted methods, compares ID, Code, and GlCode
func LookupLanguage(id string) (*Lang, error) {
	langs, err := Languages()
	if err != nil {
		return nil, err
	}
	for _, lang := range langs {
		if lang.Name == id || fmt.Sprintf("%v", lang.ID) == id || lang.Code == id || lang.GlCode == id {
			return lang, nil
		}
	}
	return nil, errors.New("Language not found")
}
