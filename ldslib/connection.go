package ldslib

import (
	"io"
)

type Reader interface {
	Language(id string) *Language
	Languages() []Language
	Index(lang *Language) []CatalogItem
	Children(item CatalogItem) []CatalogItem
	Lookup(path string) CatalogItem
}

type Writer interface {
	SetLanguages(languages []Language)
	SetCatalog(catalog Catalog)
	SetBook(stream io.Reader)
}