package ldslib

import "encoding/json"

type CatalogParser struct {
	content      Content
	foldersById  map[int]*Folder
	booksById    map[int]*Book
	booksByGlURI map[string]*Book
	catalog      *Catalog
	language     *Language
}

func NewCatalogLoader(lang *Language, content Content) *CatalogLoader {
	c := new(CatalogLoader)
	c.language = lang
	c.content = content
	return c
}

type glCatalogDescrpition struct {
	Catalog         *Catalog `json:"catalog"`
	Success         bool     `json:"success"`
	CoverArtBaseUrl string   `json:"cover_art_base_url"`
}

func (l *CatalogLoader) populateIfNeeded() {
	if l.catalog != nil {
		return
	}

	var description glCatalogDescrpition
	file := l.content.OpenRead(l.content.GetCatalogPath(l.language))
	dec := json.NewDecoder(file)
	err := dec.Decode(&description)
	if err != nil {
		panic(err)
	}

	l.foldersById = make(map[int]*Folder)
	l.booksById = make(map[int]*Book)
	l.booksByGlURI = make(map[string]*Book)

	l.catalog = description.Catalog
	l.addFolders(description.Catalog.Folders)
	l.addBooks(description.Catalog.Books)
}

func (l *CatalogLoader) addFolders(folders []*Folder) {
	for _, f := range folders {
		l.foldersById[f.ID] = f
		l.addFolders(f.Folders)
		l.addBooks(f.Books)
	}
}

func (l *CatalogLoader) addBooks(books []*Book) {
	for _, b := range books {
		l.booksById[b.ID] = b
		l.booksByGlURI[b.GlURI] = b
		b.Language = l.language
	}
}

func (l *CatalogLoader) GetCatalog() *Catalog {
	l.populateIfNeeded()
	return l.catalog
}

func (l *CatalogLoader) GetFolderById(id int) *Folder {
	l.populateIfNeeded()
	return l.foldersById[id]
}

func (l *CatalogLoader) GetBookById(id int) *Book {
	l.populateIfNeeded()
	return l.booksById[id]
}

func (l *CatalogLoader) GetBookByGlURI(glUri string) *Book {
	l.populateIfNeeded()
	return l.booksByGlURI[glUri]
}

func (l *CatalogLoader) GetBookByUnknown(id string) *Book {
	gl := ParseForBook(id)

	if gl != "" {
		return l.GetBookByGlURI(gl)
	}

	return l.GetBookByGlURI(id)
}
