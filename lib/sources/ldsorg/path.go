package ldsorg

import (
	"os"
	"path"

	"github.com/duckbrain/ldss/lib"
)

// The server to access the Gospel Library catalog and language lists from
var GospelLibraryServer = "https://tech.lds.org/glweb"

const mkdirMode = os.ModeDir | os.ModePerm

const platformID = 17

// Local Paths

func mkdirAndGetFile(paths ...string) string {
	os.MkdirAll(path.Join(paths[:len(paths)-1]...), mkdirMode)
	return path.Join(paths...)
}

func languagesPath() string {
	return mkdirAndGetFile(lib.DataDirectory, "languages.json")
}
func catalogPath(lang lib.Lang) string {
	return mkdirAndGetFile(lib.DataDirectory, lang.Code(), "catalog.json")
}
func bookPath(book *book) string {
	return mkdirAndGetFile(lib.DataDirectory,
		book.catalog.Lang().Code(),
		book.Path(), "contents.sqlite")
}

// Server Paths

func getServerAction(action string) string {
	return GospelLibraryServer + "?action=" + action
}
