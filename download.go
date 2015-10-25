package main

import (
	"compress/zlib"
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

func (d *Downloader) downloadFile(get string, save string, zlibDecompress bool) {
	var input io.Reader

	resp, err := http.Get(get)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if zlibDecompress {
		input, err = zlib.NewReader(resp.Body)
		if err != nil {
			panic(err)
		}
	} else {
		input = resp.Body
	}

	file, err := os.Create(save)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(file, input)
}

func (d *Downloader) DownloadLanguages() {
	fmt.Println("Downloading language list")
	//d.downloadFile(d.online.GetLanguagesPath(), d.offline.GetLanguagesPath(), false)
	loader := NewJSONLanguageLoader(d.online)
	languages := loader.GetAll()
	cache := NewCacheConnection()
	cache.Open(d.offline.GetCachePath())
	cache.SaveLanguages(languages)
}

func (d *Downloader) DownloadCatalog(language *Language) {
	fmt.Println("Downloading \"" + language.Name + "\" catalog")
	d.downloadFile(d.online.GetCatalogPath(language), d.offline.GetCatalogPath(language), false)
}

func (d *Downloader) DownloadBook(book *Book) {
	fmt.Println("Downloading \"" + book.Name + "\"")
	d.downloadFile(d.online.GetBookPath(book), d.offline.GetBookPath(book), true)
}

func (d *Downloader) DownloadAllBooks(languageId int) {
}
