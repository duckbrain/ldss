package ldsorg

import (
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"os"
)

// The number of simultaneous downloads when using DownloadChildren or DownloadAll
var DownloadLimit int = 6

// Downloads the catalog for the passed language
func DownloadCatalog(language Lang) error {
	lang := language.(*jsonLang)
	path := getServerAction(fmt.Sprintf("catalog.query&languageid=%v&platformid=%v", lang.ID, platformID))
	return downloadFile(path, catalogPath(language), false)
}

// Downloads the passed book
func DownloadBook(book *Book) error {
	return downloadFile(book.URL(), bookPath(book), true)
}
