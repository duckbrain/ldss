package lib

import (
	"fmt"
)

type notDownloadedErr struct {
	err  error
	item fmt.Stringer
}

func (err *notDownloadedErr) InternalError() error {
	return err.err
}

type NotDownloadedBookErr notDownloadedErr

type NotDownloadedCatalogErr notDownloadedErr

type NotDownloadedLanguageErr notDownloadedErr

func (err *NotDownloadedBookErr) Error() string {
	return fmt.Sprintf("Book \"%v\" is not downloaded", err.item)
}

func (err *NotDownloadedCatalogErr) Error() string {
	return fmt.Sprintf("Catalog for language \"%v\" is not downloaded", err.item)
}

func (err *NotDownloadedLanguageErr) Error() string {
	return "The language list is not downloaded"
}
