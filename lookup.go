package main

import ()

type LookupLoader struct {
	language *Language
	catalog  *CatalogLoader
	books    map[string]*Book
}

func NewLookupLoader(lang *Language, content *Content) {
	l := new(LookupLoader)
	l.language = lang
	l.catalog = NewCatalogLoader(l, content)
	l.books = make(map[string]*Book)
}
