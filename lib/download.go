package lib

import (
	"compress/zlib"
	"errors"
	"io"
)

func downloadFile(get string, save string, zlibDecompress bool) (err error) {
	var input io.Reader

	body, err := server.Open(get)
	if err != nil {
		return
	}
	defer body.Close()

	if zlibDecompress {
		input, err = zlib.NewReader(body)
		if err != nil {
			return
		}
	} else {
		input = body
	}

	file, err := source.Create(save)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(file, input)
	return
}

func DownloadLanguages() error {
	return downloadFile(server.LanguagesPath(), source.LanguagesPath(), false)
}

func DownloadCatalog(language *Language) error {
	return downloadFile(server.CatalogPath(language), source.CatalogPath(language), false)
}

func DownloadBook(book *Book) error {
	return downloadFile(server.BookPath(book), source.BookPath(book), true)
}

func DownloadBooks(languageId int) error {
	return errors.New("Not Implemented")
}

func DownloadStatus() error {
	return nil
}

func DownloadMissing() error {
	return nil
}
