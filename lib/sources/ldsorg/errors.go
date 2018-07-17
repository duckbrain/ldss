package ldsorg

import (
	"fmt"
)

// An error that a resouce has not been downloaded
type notDownloadedErr interface {
	error
	fmt.Stringer
	InternalError() error
	Download() error
}

// Base commonality for the NotDownloadedErr interface
type notDownloadedErrBase struct {
	err error
}

// The error of the wrapped error
func (err notDownloadedErrBase) InternalError() error {
	return err.err
}

// An error that a book needs to be downloaded
type notDownloadedBookErr struct {
	notDownloadedErrBase
	book *Book
}

// An error that a language's catalog needs to be downloaded
type notDownloadedCatalogErr struct {
	notDownloadedErrBase
	lang Lang
}

// An Error that the language list has not been downloaded
type notDownloadedLanguageErr notDownloadedErrBase

// Error book is not downloaded
func (err notDownloadedBookErr) Error() string {
	return fmt.Sprintf("Book \"%v\" is not downloaded", err.book)
}

// Get the string representation of the Book
func (err notDownloadedBookErr) String() string {
	return err.book.String()
}

// Download the missing book
func (err notDownloadedBookErr) Download() error {
	return DownloadBook(err.book)
}

// Error catalog is not downloaded
func (err notDownloadedCatalogErr) Error() string {
	return fmt.Sprintf("Catalog for language \"%v\" is not downloaded", err.lang)
}

// Get the string representation of the language of the missing catalog
func (err notDownloadedCatalogErr) String() string {
	return err.lang.String()
}

// Download the missing catalog
func (err notDownloadedCatalogErr) Download() error {
	return DownloadCatalog(err.lang)
}

// Error languages are not downloaded
func (err notDownloadedLanguageErr) Error() string {
	return "The language list is not downloaded"
}

// Get the string representation of the language list, returns "Languages"
func (err notDownloadedLanguageErr) String() string {
	return "Languages"
}

// Download the language list
func (err notDownloadedLanguageErr) Download() error {
	return DownloadLanguages()
}
