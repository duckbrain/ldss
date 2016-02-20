package lib

import (
	"fmt"
)

// An error that a resouce has not been downloaded
type NotDownloadedErr interface {
	error
	fmt.Stringer
	InternalError() error
	Download() error
}

// Base commonality for the NotDownloadedErr interface
type notDownloadedErr struct {
	err error
}

// The error of the wrapped error
func (err notDownloadedErr) InternalError() error {
	return err.err
}

// An error that a book needs to be downloaded
type NotDownloadedBookErr struct {
	notDownloadedErr
	book *Book
}

// An error that a language's catalog needs to be downloaded
type NotDownloadedCatalogErr struct {
	notDownloadedErr
	lang *Language
}

// An Error that the language list has not been downloaded
type NotDownloadedLanguageErr notDownloadedErr

// Error book is not downloaded
func (err NotDownloadedBookErr) Error() string {
	return fmt.Sprintf("Book \"%v\" is not downloaded", err.book)
}

// Get the string representation of the Book
func (err NotDownloadedBookErr) String() string {
	return err.book.String()
}

// Download the missing book
func (err NotDownloadedBookErr) Download() error {
	return DownloadBook(err.book)
}

// Error catalog is not downloaded
func (err NotDownloadedCatalogErr) Error() string {
	return fmt.Sprintf("Catalog for language \"%v\" is not downloaded", err.lang)
}

// Get the string representation of the language of the missing catalog
func (err NotDownloadedCatalogErr) String() string {
	return err.lang.String()
}

// Download the missing catalog
func (err NotDownloadedCatalogErr) Download() error {
	return DownloadCatalog(err.lang)
}

// Error languages are not downloaded
func (err NotDownloadedLanguageErr) Error() string {
	return "The language list is not downloaded"
}

// Get the string representation of the language list, returns "Languages"
func (err NotDownloadedLanguageErr) String() string {
	return "Languages"
}

// Download the language list
func (err NotDownloadedLanguageErr) Download() error {
	return DownloadLanguages()
}
