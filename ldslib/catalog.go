package ldslib

import "encoding/json"

type catalogParser struct {
	source       Source
	foldersById  map[int]*Folder
	booksById    map[int]*Book
	booksByGlURI map[string]*Book
	catalog      *Catalog
	language     *Language
}

func newCatalogLoader(lang *Language, content Source) *catalogParser {
	c := new(catalogParser)
	c.language = lang
	c.source = content
	return c
}

type glCatalogDescrpition struct {
	Catalog         *Catalog `json:"catalog"`
	Success         bool     `json:"success"`
	CoverArtBaseUrl string   `json:"cover_art_base_url"`
}

func (l *catalogParser) populateIfNeeded() {
	if l.catalog != nil {
		return
	}

	var description glCatalogDescrpition
	file, err := l.source.Open(l.source.CatalogPath(l.language))
	if err != nil {
		panic(err)
	}
	dec := json.NewDecoder(file)
	err = dec.Decode(&description)
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

func (l *catalogParser) addFolders(folders []*Folder) {
	for _, f := range folders {
		l.foldersById[f.ID] = f
		l.addFolders(f.Folders)
		l.addBooks(f.Books)
	}
}

func (l *catalogParser) addBooks(books []*Book) {
	for _, b := range books {
		l.booksById[b.ID] = b
		l.booksByGlURI[b.GlURI] = b
		b.Language = l.language
	}
}

func (l *catalogParser) GetCatalog() *Catalog {
	l.populateIfNeeded()
	return l.catalog
}

func (l *catalogParser) GetFolderById(id int) *Folder {
	l.populateIfNeeded()
	return l.foldersById[id]
}

func (l *catalogParser) GetBookById(id int) *Book {
	l.populateIfNeeded()
	return l.booksById[id]
}

func (l *catalogParser) GetBookByGlURI(glUri string) *Book {
	l.populateIfNeeded()
	return l.booksByGlURI[glUri]
}

func (l *catalogParser) GetBookByUnknown(id string) *Book {
	gl := ParseForBook(id)

	if gl != "" {
		return l.GetBookByGlURI(gl)
	}

	return l.GetBookByGlURI(id)
}
