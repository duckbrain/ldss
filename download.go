
package main

import (
	"io"
	"os"
	"net/http"
)

type Downloader struct {
	online *LDSContent
	offline *LocalContent
}

func (d Downloader) DownloadStatus() {

}

func (d Downloader) IsLanguagesDownloaded() bool {
	_, err := os.Stat("~/.ldss/languages.json");
	return os.IsNotExist(err);
}

func (d Downloader) DownloadMissing() {
}

func (d Downloader) downloadFile(get string, save string) {
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

func (d Downloader) DownloadLanguages() {
	d.downloadFile(d.online.GetLanguagesPath(), d.offline.GetLanguagesPath())
}

func (d Downloader) DownloadCatalog(languageId int) {
	d.downloadFile(d.online.GetCatalogPath(languageId), d.offline.GetCatalogPath(languageId))
}

func (d Downloader) DownloadBook(languageId int, bookId int) {
	
}

func (d Downloader) DownloadAllBooks(languageId int) {
}
