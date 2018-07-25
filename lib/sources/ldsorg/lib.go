package ldsorg

var catalogsByLanguageId map[int]*catalog
var booksByLangBookId map[langBookID]*Book

type langBookID struct {
	langID int
	bookID int
}

func init() {
	catalogsByLanguageId = make(map[int]*catalog)
	booksByLangBookId = make(map[langBookID]*Book)
}
