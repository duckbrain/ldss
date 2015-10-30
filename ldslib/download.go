package ldslib

import (
	"compress/zlib"
	"fmt"
	"io"
)

type Downloader struct {
	online  Source
	offline Source
}

func NewDownloader(online Source, offline Source) *Downloader {
	d := new(Downloader)
	d.online = online
	d.offline = offline
	return d
}

func (d *Downloader) Status() {

}

func (d *Downloader) Missing() {
}

func (d *Downloader) downloadFile(get string, save string, zlibDecompress bool) {
	var input io.Reader

	body, err := d.online.Open(get)
	if err != nil {
		panic(err)
	}
	defer body.Close()

	if zlibDecompress {
		input, err = zlib.NewReader(body)
		if err != nil {
			panic(err)
		}
	} else {
		input = body
	}

	file, err := d.offline.Create(save)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(file, input)
}

func (d *Downloader) Languages() {
	fmt.Println("Downloading language list")
	d.downloadFile(d.online.LanguagesPath(), d.offline.LanguagesPath(), false)
}

func (d *Downloader) Catalog(language *Language) {
	fmt.Println("Downloading \"" + language.Name + "\" catalog")
	d.downloadFile(d.online.CatalogPath(language), d.offline.CatalogPath(language), false)
}

func (d *Downloader) Book(book *Book) {
	fmt.Println("Downloading \"" + book.Name + "\"")
	d.downloadFile(d.online.BookPath(book), d.offline.BookPath(book), true)
}

func (d *Downloader) Books(languageId int) {
	fmt.Println("Not implemented")
}
