package ldsorg

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/duckbrain/ldss/lib"
	"github.com/duckbrain/ldss/lib/dl"
)

var languages []lib.Lang

func (s source) Langs() ([]lib.Lang, error) {
	if languages != nil {
		return languages, nil
	}

	var root struct {
		Languages []*lang `json:"languages"`
	}

	file, err := os.Open(languagesPath())
	if err != nil {
		return nil, dl.ErrNotDownloaded(s)
	}
	if err = json.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	languages = make([]lib.Lang, len(root.Languages))
	for i, l := range root.Languages {
		languages[i] = l
	}

	return languages, nil
}

// Lang defines a language as from the server. The fields should not be modified.
type lang struct {
	// The Gospel Library ID for the language. Used for downloads.
	ID int `json:"id"`
	// Native representation of the language in the language observed
	JsonName string `json:"name"`

	// English representation of the language
	JsonEnglishName string `json:"eng_name"`

	// The internationalization (i18n) code used in most programs
	JsonCode string `json:"code"`

	// Gospel Library language code, seen in the urls of https://lds.org
	GlCode string `json:"code_three"`
}

func (l lang) Name() string {
	return l.JsonName
}

func (l lang) EnglishName() string {
	return l.JsonEnglishName
}

func (l lang) Code() string {
	return l.JsonCode
}

func (l lang) Matches(s string) bool {
	s = strings.ToLower(s)
	return s == strings.ToLower(l.JsonName)
}
