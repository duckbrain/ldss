
package main

import (
	"io"
	"os"
	"net/http"
)

func DownloadStatus() {

}

func IsLanguagesDownloaded() bool {
	_, err := os.Stat("~/.ldss/languages.json");
	return os.IsNotExist(err);
}

func DownloadMissing() {
}

func downloadFile(get string, save string) {
	resp, err := http.Get(get)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	file, err := os.Create(save);
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(file, resp.Body)
}

func DownloadLanguages() {
	ldsContent := NewLDSContent()
	localContent := NewLocalContent()
	
	os.MkdirAll(localContent.BasePath, os.ModeDir);
	downloadFile(ldsContent.GetLanguagesPath(), localContent.GetLanguagesPath())
}

func DownloadCatalog(languageId int) {
	ldsContent := NewLDSContent()
	localContent := NewLocalContent()
	
	downloadFile(ldsContent.GetCatalogPath(languageId), localContent.GetCatalogPath(languageId))
}

func DownloadBook(languageId int, bookId int) {
}

func DownloadAllBooks(languageId int) {
}
