package ldslib

import (
	
)

type Reader struct {
	c Content
	languages LanguageLoader
	catalogsByLanguageId map[int] *CatalogLoader
}

func NewReader(content Content) Reader {
	return &JSONConnection{
		content, 
		NewJSONLanguageLoader(content), 
		make(map[int]*CatalogLoader)
	}
}

func (r *Reader) Language(id string) *Language {
	return r.languages.GetByUnknown(id)
}
func (r *Reader) Languages() []Language {
	return r.languages.GetAll();
}
	Index(lang *Language) []CatalogItem
	Children(item CatalogItem) []CatalogItem
	Lookup(path string) CatalogItem