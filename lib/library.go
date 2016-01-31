package lib

var source Source
var catalogsByLanguageId map[int]*Catalog
var booksByLangBookId map[langBookID]*bookParser

type langBookID struct {
	langID int
	bookID int
}

func init() {
	//TODO Set source
	catalogsByLanguageId = make(map[int]*Catalog)
	booksByLangBookId = make(map[langBookID]*bookParser)
}

/*
func (l *Library) populateCatalog(lang *Language) (*Catalog, error) {
	if c, ok := l.catalogsByLanguageId[lang.ID]; ok {
		return c, nil
	}
	c, err := newCatalog(lang, l.source)
	if err != nil {
		return nil, err
	}
	l.catalogsByLanguageId[lang.ID] = c
	return c, nil
}

func (l *Library) populateBook(book *Book) *bookParser {
	id := langBookID{book.catalog.language.ID, book.ID}
	b, ok := l.booksByLangBookId[id]
	if !ok {
		b = newBookParser(book, l.source)
		l.booksByLangBookId[id] = b
		book.parser = b
	}
	return b
}

func (l *Library) FindLanguage(id string) (*Language, error) {
	if err := l.populateLanguages(); err != nil {
		return nil, err
	}
	for _, lang := range l.languages {
		if lang.Name == id || fmt.Sprintf("%v", lang.ID) == id || lang.EnglishName == id || lang.Code == id || lang.GlCode == id {
			return &lang, nil
		}
	}
	return nil, errors.New("Language not found")
}

func (l *Library) Languages() ([]Language, error) {
	return l.languages, l.populateLanguages()
}

func (l *Library) Catalog(lang *Language) (*Catalog, error) {
	return l.populateCatalog(lang)
}
func (l *Library) Book(path string, catalog *Catalog) (*Book, error) {
	return l.populateCatalog(catalog.Language()).BookByUnknown(path)
}





func (l *Library) Children(item Item) ([]Item, error) {
	switch item.(type) {
	case *Book:
		l.populateBook(item.(*Book))
		return item.Children()
	default:
		return item.Children()
	}
}

func (l *Library) Content(node Node) (*Page, error) {
	rawContent, err := l.populateBook(node.Book).Content(node)
	if err != nil {
		return nil, err
	}
	parser := ContentParser{contentHtml: rawContent}
	//return parser.Content()
	return nil, nil
}
*/
//	Index(lang *Language) []CatalogItem
//	Children(item CatalogItem) []CatalogItem
