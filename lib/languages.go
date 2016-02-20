package lib

import (
	"encoding/json"
	"errors"
	"fmt"
)

var languages cache

type glLanguageDescription struct {
	Languages []*Language `json:"languages"`
	Success   bool        `json:"success"`
}

type Language struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	EnglishName  string `json:"eng_name"`
	Code         string `json:"code"`
	GlCode       string `json:"code_three"`
	catalogCache cache
	reference    cache
}

func (l *Language) String() string {
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

func (l *Language) Catalog() (*Catalog, error) {
	l.catalogCache.construct = func() (interface{}, error) {
		return newCatalog(l)
	}
	c, err := l.catalogCache.get()
	if err != nil {
		return nil, err
	}
	return c.(*Catalog), err
}

func (l *Language) ref() (*refParser, error) {
	ref, err := l.reference.get()
	if err != nil {
		return nil, err
	}
	return ref.(*refParser), nil
}

func (l *Language) Reference(q string) (string, error) {
	ref, err := l.ref()
	if err != nil {
		return "", err
	}
	return ref.lookup(q)
}

func init() {
	languages.construct = func() (interface{}, error) {
		var description glLanguageDescription
		file, err := source.Open(source.LanguagesPath())
		if err != nil {
			return nil, &NotDownloadedLanguageErr{err}
		}
		err = json.NewDecoder(file).Decode(&description)
		return description.Languages, err
	}
}

func Languages() ([]*Language, error) {
	if !source.Exist(source.LanguagesPath()) {
		if err := DownloadLanguages(); err != nil {
			return nil, err
		}
	}
	langs, err := languages.get()
	if err != nil {
		return nil, err
	}
	return langs.([]*Language), err
}

func LookupLanguage(id string) (*Language, error) {
	langs, err := Languages()
	if err != nil {
		return nil, err
	}
	for _, lang := range langs {
		if lang.Name == id || fmt.Sprintf("%v", lang.ID) == id || lang.EnglishName == id || lang.Code == id || lang.GlCode == id {
			return lang, nil
		}
	}
	return nil, errors.New("Language not found")
}

func DefaultLanguage() (*Language, error) {
	return LookupLanguage(Config().Get("Language").(string))
}
