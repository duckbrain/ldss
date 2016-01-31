package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

var source, server Source

func init() {
	source = NewOfflineSource(".ldss")
	server = NewOnlineSource("https://tech.lds.org/glweb")
}

type Source interface {
	LanguagesPath() string
	CatalogPath(language *Language) string
	BookPath(book *Book) string
	Open(path string) (io.ReadCloser, error)
	Create(path string) (io.WriteCloser, error)
	Exist(path string) bool
}

type localSource struct {
	BasePath  string
	MkdirMode os.FileMode
}

func NewOfflineSource(path string) Source {
	return localSource{path, os.ModeDir | os.ModePerm}
}

func mkdirAndGetFile(paths ...string) string {
	os.MkdirAll(path.Join(paths[:len(paths)-1]...), os.ModeDir|os.ModePerm)
	return path.Join(paths...)
}

func (c localSource) CachePath() string {
	return mkdirAndGetFile(c.BasePath, "cache.sqlite")
}
func (c localSource) ConfigPath() string {
	return mkdirAndGetFile(c.BasePath, "config.json")
}
func (c localSource) LanguagesPath() string {
	return mkdirAndGetFile(c.BasePath, "languages.json")
}
func (c localSource) CatalogPath(language *Language) string {
	return mkdirAndGetFile(c.BasePath, language.GlCode, "catalog.json")
}
func (c localSource) BookPath(book *Book) string {
	return mkdirAndGetFile(c.BasePath, book.catalog.Language().GlCode, book.Path(), "contents.sqlite")
}
func (c localSource) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}
func (c localSource) Create(path string) (io.WriteCloser, error) {
	return os.Create(path)
}
func (c localSource) Exist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

type ldsSource struct {
	BasePath   string
	PlatformID int
}

func NewOnlineSource(path string) Source {
	return ldsSource{path, 17}
}
func (c ldsSource) getAction(action string) string {
	return c.BasePath + "?action=" + action
}
func (c ldsSource) LanguagesPath() string {
	return c.getAction("languages.query")
}
func (c ldsSource) CatalogPath(language *Language) string {
	return c.getAction(fmt.Sprintf("catalog.query&languageid=%v&platformid=%v", language.ID, c.PlatformID))
}
func (c ldsSource) BookPath(book *Book) string {
	return book.URL()
}
func (c ldsSource) Open(path string) (io.ReadCloser, error) {
	resp, err := http.Get(path)
	return resp.Body, err
}
func (c ldsSource) Create(path string) (io.WriteCloser, error) {
	return nil, errors.New("Can't create an online source")
}
func (c ldsSource) Exist(path string) bool {
	return true // Assume that it is all there
}
