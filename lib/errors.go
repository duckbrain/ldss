package lib

import (
	"fmt"
)

type NotDownloadedErr interface {
	error
	fmt.Stringer
	InternalError() error
	Download() error
}

type notDownloadedErr struct {
	err error
}

func (err *notDownloadedErr) InternalError() error {
	return err.err
}

type NotDownloadedBookErr struct {
	notDownloadedErr
	book *Book
}

type NotDownloadedCatalogErr struct {
	notDownloadedErr
	lang *Language
}

type NotDownloadedLanguageErr notDownloadedErr

func (err NotDownloadedBookErr) Error() string {
	return fmt.Sprintf("Book \"%v\" is not downloaded", err.book)
}

func (err NotDownloadedBookErr) String() string {
	return err.book.String()
}

func (err NotDownloadedBookErr) Download() error {
	return DownloadBook(err.book)
}

func (err NotDownloadedCatalogErr) Error() string {
	return fmt.Sprintf("Catalog for language \"%v\" is not downloaded", err.lang)
}

func (err NotDownloadedCatalogErr) String() string {
	return err.lang.String()
}

func (err NotDownloadedCatalogErr) Download() error {
	return DownloadCatalog(err.lang)
}

func (err NotDownloadedLanguageErr) Error() string {
	return "The language list is not downloaded"
}

func (err NotDownloadedLanguageErr) String() string {
	return "Languages"
}

func (err NotDownloadedLanguageErr) Download() error {
	return DownloadLanguages()
}
