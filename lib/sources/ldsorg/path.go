package ldsorg

import (
	"os"
	"path"
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
	return mkdirAndGetFile(DataDirectory, "languages.json")
}
func catalogPath(language Lang) string {
	return mkdirAndGetFile(DataDirectory, language.GlCode, "catalog.json")
}
func bookPath(book *Book) string {
	return mkdirAndGetFile(DataDirectory,
		book.catalog.Language().GlCode,
		book.Path(), "contents.sqlite")
}
func fileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Server Paths

func getServerAction(action string) string {
	return GospelLibraryServer + "?action=" + action
}
