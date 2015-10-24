package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

type Content interface {
	GetLanguagesPath() string
	GetCatalogPath(language *Language) string
	GetBookPath(book *Book) string
	OpenRead(path string) io.Reader
}

type LocalContent struct {
	BasePath  string
	MkdirMode os.FileMode
}

func NewLocalContent(path string) *LocalContent {
	c := new(LocalContent)
	c.BasePath = path
	c.MkdirMode = os.ModeDir | os.ModePerm
	return c
}

func mkdirAndGetFile(paths ...string) string {
	os.MkdirAll(path.Join(paths[:len(paths)-1]...), os.ModeDir|os.ModePerm)
	return path.Join(paths...)
}

func (c LocalContent) GetCachePath() string {
	return mkdirAndGetFile(c.BasePath, "cache.sqlite")
}
func (c LocalContent) GetConfigPath() string {
	return mkdirAndGetFile(c.BasePath, "config.json")
}
func (c LocalContent) GetLanguagesPath() string {
	return mkdirAndGetFile(c.BasePath, "languages.json")
}
func (c LocalContent) GetCatalogPath(language *Language) string {
	return mkdirAndGetFile(c.BasePath, language.GlCode, "catalog.json")
}
func (c LocalContent) GetBookPath(book *Book) string {
	return mkdirAndGetFile(c.BasePath, book.Language.GlCode, book.GlURI, "contents.sqlite")
}
func (c LocalContent) OpenRead(path string) io.Reader {
	reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return reader
}

type LDSContent struct {
	BasePath   string
	PlatformID int
}

func NewLDSContent(path string, platformId int) *LDSContent {
	c := new(LDSContent)
	c.BasePath = path
	c.PlatformID = platformId
	return c
}
func (c *LDSContent) getAction(action string) string {
	return c.BasePath + "?action=" + action
}
func (c *LDSContent) GetLanguagesPath() string {
	return c.getAction("languages.query")
}
func (c *LDSContent) GetCatalogPath(language *Language) string {
	return c.getAction(fmt.Sprintf("catalog.query&languageid=%v&platformid=%v", language.ID, c.PlatformID))
}
func (c *LDSContent) GetBookPath(book *Book) string {
	return book.URL
}
func (c *LDSContent) OpenRead(path string) io.Reader {
	resp, err := http.Get(path)
	if err != nil {
		panic(err)
	}
	return resp.Body
}
