package ldslib

import ()

type LookupLoader struct {
	language *Language
	catalog  *catalogParser
	books    map[string]*Book
}

func NewLookupLoader(lang *Language, content Source) *LookupLoader {
	l := new(LookupLoader)
	l.language = lang
	l.catalog = newCatalogLoader(lang, content)
	l.books = make(map[string]*Book)
	return l
}
