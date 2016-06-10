package lib

import (
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"os"
)

// The number of simultanious downloads when using DownloadChildren or DownloadAll
var DownloadLimit int = 6

// Downloads the content at the get path to the save path, optionally zlib decompressing it.
func downloadFile(get string, save string, zlibDecompress bool) (err error) {
	var input io.Reader

	response, err := http.Get(get)
	if err != nil {
		return
	}
	body := response.Body
	defer response.Body.Close()

	if zlibDecompress {
		input, err = zlib.NewReader(body)
		if err != nil {
			return
		}
	} else {
		input = body
	}

	file, err := os.Create(save)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(file, input)
	return
}

// Downloads the list of languages
func DownloadLanguages() error {
	return downloadFile(getServerAction("languages.query"), languagesPath(), false)
}

// Downloads the catalog for the passed language
func DownloadCatalog(language *Language) error {
	path := getServerAction(fmt.Sprintf("catalog.query&languageid=%v&platformid=%v", language.ID, platformID))
	return downloadFile(path, catalogPath(language), false)
}

// Downloads the passed book
func DownloadBook(book *Book) error {
	return downloadFile(book.URL(), bookPath(book), true)
}

// Recursively downloads all children of the passed Catalog or Folder.
func DownloadChildren(item Item, force bool) <-chan Message {
	// Find the catalog
	var catalog *Catalog
	parent := item
	for {
		if c, ok := parent.(*Catalog); ok {
			catalog = c
			break
		}
		parent = parent.Parent()
	}

	// Open a channel and start searching
	c := make(chan Message)
	go func() {
		lock := make(chan interface{})
		limit := make(chan interface{}, DownloadLimit)
		for i := 0; i < cap(limit); i++ {
			limit <- nil
		}
		for _, book := range catalog.booksById {
			go func(book *Book) {
				<-limit
				defer func() {
					lock <- nil
					limit <- nil
				}()
				if force || !fileExist(bookPath(book)) {
					// Skip if the book is not a child of the item
					for parent := book.Parent(); parent != item; parent = parent.Parent() {
						if parent == nil {
							return
						}
					}

					c <- MessageDownload{NotDownloadedBookErr{notDownloadedErr{}, book}}
					if err := DownloadBook(book); err != nil {
						//TODO Send warning message of error
					}
				}
			}(book)
		}
		for _ = range catalog.booksById {
			<-lock
		}
		//TODO Keep track of stats and send a message before closing
		c <- MessageDone{catalog}
		close(c)
	}()
	return c
}

// Downloads the catalog and all books for a language. If force is true, it will
// download these items even if they are already downloaded, replacing them.
func DownloadAll(lang *Language, force bool) <-chan Message {
	c := make(chan Message)
	go func() {
		if force || !fileExist(catalogPath(lang)) {
			c <- MessageDownload{NotDownloadedCatalogErr{notDownloadedErr{}, lang}}
			DownloadCatalog(lang)
		}

		catalog, err := lang.Catalog()
		if err != nil {
			c <- MessageError{err}
			close(c)
			return
		}

		for m := range DownloadChildren(catalog, force) {
			c <- m
		}
		close(c)
	}()
	return c
}
