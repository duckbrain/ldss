package ldsorg

var catalogsByLanguageId map[int]*Catalog
var booksByLangBookId map[langBookID]*Book

type langBookID struct {
	langID int
	bookID int
}

func init() {
	catalogsByLanguageId = make(map[int]*Catalog)
	booksByLangBookId = make(map[langBookID]*Book)
}

func catalog(langId int) (*Catalog, error) {

}
