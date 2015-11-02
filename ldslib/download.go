package ldslib

import (
	"compress/zlib"
	"errors"
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

func (d *Downloader) downloadFile(get string, save string, zlibDecompress bool) (err error) {
	var input io.Reader

	body, err := d.online.Open(get)
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

	file, err := d.offline.Create(save)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(file, input)
	return
}

func (d *Downloader) Languages() error {
	return d.downloadFile(d.online.LanguagesPath(), d.offline.LanguagesPath(), false)
}

func (d *Downloader) Catalog(language *Language) error {
	return d.downloadFile(d.online.CatalogPath(language), d.offline.CatalogPath(language), false)
}

func (d *Downloader) Book(book *Book) error {
	return d.downloadFile(d.online.BookPath(book), d.offline.BookPath(book), true)
}

func (d *Downloader) Books(languageId int) error {
	return errors.New("Not Implemented")
}
