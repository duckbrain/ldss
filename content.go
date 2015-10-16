package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
)

type Content interface {
	GetLanguagesPath() string
	GetCatalogPath(language *Language) string
	GetBookPath(language *Language, glUri string) string
	OpenRead(path string) io.Reader
}

type LocalContent struct {
	BasePath string
}

func NewLocalContent() LocalContent {
	//TODO: Load path from config
	u, err := user.Current()

	if err != nil {
		panic(err)
	}

	return LocalContent{path.Join(u.HomeDir, ".ldss")}
}
func (c *LocalContent) GetLanguagesPath() string {
	os.MkdirAll(c.BasePath, os.ModeDir|os.ModePerm)
	return path.Join(c.BasePath, "languages.json")
}
func (c *LocalContent) GetCatalogPath(language *Language) string {
	os.MkdirAll(path.Join(c.BasePath, language.GlCode), os.ModeDir|os.ModePerm)
	return path.Join(c.BasePath, language.GlCode, "catalog.json")
}
func (c *LocalContent) GetBookPath(language *Language, glUri string) string {
	os.MkdirAll(path.Join(c.BasePath, language.GlCode, glUri), os.ModeDir|os.ModePerm)
	return path.Join(c.BasePath, language.GlCode, glUri, "contents.zbook")
}
func (c *LocalContent) OpenRead(path string) io.Reader {
	reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return reader
}

type LDSContent struct {
	BasePath string
}

func NewLDSContent() LDSContent {
	return LDSContent{"https://tech.lds.org/glweb"}
}
func (c *LDSContent) getAction(action string) string {
	return c.BasePath + "?action=" + action
}
func (c *LDSContent) GetLanguagesPath() string {
	return c.getAction("languages.query")
}
func (c *LDSContent) GetCatalogPath(language *Language) string {
	return c.getAction(fmt.Sprintf("catalog.query&platformid=17&languageid=%v", language.ID))
}
func (c *LDSContent) GetBookPath(language *Language, glUri string) string {
	return ""
}
func (c *LDSContent) OpenRead(path string) io.Reader {
	resp, err := http.Get(path)
	if err != nil {
		panic(err)
	}
	return resp.Body
}
