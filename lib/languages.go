package lib

import (
	"errors"
	"fmt"
)

var langsBySrc map[string][]*Lang
var langs map[string]langWithSrcs

func init() {
	langsBySrc = make(map[string][]*Lang)
	langs = make(map[string]langWithSrcs)
}

type langWithSrcs struct {
	Lang
	srcs map[string]*Lang
}

// Lang defines a language as from the server. The fields should not be modified.
type Lang struct {
	// Source internal ID for the language. May be different between sources
	ID int
	// Source internal name for the language. May be different between sources.
	InternalCode string

	// Native representation of the language in the language observed (optional)
	Name string
	// English representation of the language (optional)
	EnglishName string
	// The internationalization (i18n) code used in most programs. Used as key to join languages from other sources
	Code string
}

// Returns a human readable version of the language that is appropriate to
// show an end user. It will format the language to contain it's native
// representation as well as the English representation. It will also show
// the standard internationalization code as well as the Gospel Library
// language code.
func (l *Lang) String() string {
	var name, code string

	if l.Name == l.EnglishName {
		name = l.Name
	} else {
		name = fmt.Sprintf("%v (%v)", l.Name, l.EnglishName)
	}
	if l.Code == l.InternalCode {
		code = fmt.Sprintf(" [%v]", l.Code)
	} else {
		code = fmt.Sprintf(" [%v/%v]", l.Code, l.InternalCode)
	}

	return name + code
}

func Languages() []Lang {
	res := []Lang{}
	for _, l := range langs {
		res = append(res, l.Lang)
	}
	return res
}

func LanguageFromSource(lang Lang, srcName string) *Lang {
	return langs[lang.Code].srcs[srcName]
}

func registerLanguage(srcName string, srcLangs []*Lang) {
	langsBySrc[srcName] = srcLangs
	for _, srcLang := range srcLangs {
		if lang, ok := langs[srcLang.Code]; ok {
			lang.srcs[srcName] = srcLang
			// TODO Merge other fields to fill in the blanks
		} else {
			lang := langWithSrcs{
				Lang: *srcLang,
				srcs: map[string]*Lang{
					srcName: srcLang,
				},
			}
			langs[lang.Code] = lang
		}
	}
}

// LookupLanguage finds a language by any of the accepted methods, compares ID, Code, and InternalCode
func LookupLanguage(id string) (Lang, error) {
	langs := Languages()
	for _, lang := range langs {
		if lang.Name == id || fmt.Sprintf("%v", lang.ID) == id || lang.Code == id || lang.InternalCode == id {
			return lang, nil
		}
	}
	return Lang{}, errors.New("Language not found")
}
