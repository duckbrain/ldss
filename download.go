package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Downloader struct {
	online  *LDSContent
	offline *LocalContent
}

func NewDownloader(online *LDSContent, offline *LocalContent) *Downloader {
	d := new(Downloader)
	d.online = online
	d.offline = offline
	return d
}

func (d *Downloader) DownloadStatus() {

}

func (d *Downloader) IsLanguagesDownloaded() bool {
	_, err := os.Stat("~/.ldss/languages.json")
	return os.IsNotExist(err)
}

func (d *Downloader) DownloadMissing() {
}

func (d *Downloader) downloadFile(get string, save string) {
	resp, err := http.Get(get)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	file, err := os.Create(save)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(file, resp.Body)
}

func (d *Downloader) DownloadLanguages() {
	fmt.Println("Downloading language list")
	d.downloadFile(d.online.GetLanguagesPath(), d.offline.GetLanguagesPath())
}

func (d *Downloader) DownloadCatalog(language *Language) {
	fmt.Println("Downloading \"" + language.Name + "\" catalog")
	d.downloadFile(d.online.GetCatalogPath(language), d.offline.GetCatalogPath(language))
}

func (d *Downloader) DownloadBook(book *Book) {
	fmt.Println("Downloading \"" + book.Name + "\"")
	d.downloadFile(d.online.GetBookPath(book), d.offline.GetBookPath(book))
}

func (d *Downloader) DownloadAllBooks(languageId int) {
}
