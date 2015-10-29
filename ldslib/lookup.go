package ldslib

import ()

type LookupLoader struct {
	language *Language
	catalog  *CatalogLoader
	books    map[string]*Book
}

func NewLookupLoader(lang *Language, content Content) *LookupLoader {
	l := new(LookupLoader)
	l.language = lang
	l.catalog = NewCatalogLoader(lang, content)
	l.books = make(map[string]*Book)
	return l
}
