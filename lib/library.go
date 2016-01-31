package lib

var catalogsByLanguageId map[int]*Catalog
var booksByLangBookId map[langBookID]*Book

type langBookID struct {
	langID int
	bookID int
}

func init() {
	//TODO Set source
	catalogsByLanguageId = make(map[int]*Catalog)
	booksByLangBookId = make(map[langBookID]*Book)
}

func autoDownload(open func() (Item, error), c chan Message) {
	item, err := open()
	var dlErr, preDlErr NotDownloadedErr
	dlErr, ok := err.(NotDownloadedErr)
	for ok {
		if dlErr == preDlErr {
			break
		}
		c <- MessageDownload{dlErr}
		err = dlErr.Download()
		if err != nil {
			c <- MessageError{err}
			return
		}
		item, err = open()
		preDlErr = dlErr
		dlErr, ok = err.(NotDownloadedErr)
	}

	if err == nil {
		c <- MessageDone{item}
	} else {
		c <- MessageError{err}
	}
}

func DefaultCatalog() <-chan Message {
	c := make(chan Message)
	go autoDownload(func() (Item, error) {
		lang, err := DefaultLanguage()
		if err != nil {
			return nil, err
		}
		catalog, err := lang.Catalog()
		if err != nil {
			return nil, err
		}
		return catalog, nil
	}, c)
	return c
}

func Lookup(q string, c *Catalog) (Item, error) {
	return nil, nil
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
