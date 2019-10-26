package ldsorg

import (
	"fmt"

	"github.com/duckbrain/ldss/lib"
)

type ItemType string

const (
	TypeCatalog ItemType = "Catalog"
	TypeFolder  ItemType = "Folder"
	TypeBook    ItemType = "Book"
	TypeNode    ItemType = "Node"
)

type Metadata struct {
	Type ItemType

	// Book only
	DownloadURL string

	// Catalog only
	Languages map[string]Lang
}

type Catalog struct {
	Folder
}

type Folder struct {
	Name    string   `json:"name"`
	Folders []Folder `json:"folders"`
	Books   []Book   `json:"books"`
	ID      int      `json:"id"`
}

func (f Folder) Header(lang lib.Lang) lib.Header {
	return lib.Header{
		Index: lib.Index{
			Lang: lang,
			Path: f.Path(),
		},
		Name: f.Name,
	}
}
func (f Folder) Path() string {
	return fmt.Sprintf("folder-%v", f.ID)
}

type Book struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DownloadURL string `json:"url"`
	Path        string `json:"gl_uri"`
}

type Node struct {
	ID          int64
	ParentID    int64
	Path        string
	Name        string
	Subtitle    string
	SectionName string
	ShortTitle  string

	Content lib.Content
}

func (n Node) Header(lang lib.Lang) lib.Header {
	return lib.Header{
		Index: lib.Index{
			Lang: lang,
			Path: n.Path,
		},
		Name:        n.Name,
		Subtitle:    n.Subtitle,
		SectionName: n.SectionName,
		ShortTitle:  n.ShortTitle,
	}
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
}
