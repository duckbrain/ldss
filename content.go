
package main

import (
	"os"
	"io"
	"os/user"
	"path"
	"net/http"
)

type Content interface {
	GetLanguagesPath() string
	GetCatalogPath(languageId int) string
	GetBookPath(languageId int, glUri string) string
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
func (c LocalContent) GetLanguagesPath() string {
	os.MkdirAll(c.BasePath, os.ModeDir | os.ModePerm)
	return path.Join(c.BasePath, "languages.json")
}
func (c LocalContent) GetCatalogPath(languageId int) string {
	os.MkdirAll(path.Join(c.BasePath, string(languageId)), os.ModeDir | os.ModePerm)
	return path.Join(c.BasePath, string(languageId), "catalog.json")
}
func (c LocalContent) GetBookPath(languageId int, glUri string) string {
	os.MkdirAll(path.Join(c.BasePath, string(languageId), glUri), os.ModeDir | os.ModePerm)
	return path.Join(c.BasePath, string(languageId), glUri, "contents.zbook")
}
func (c LocalContent) OpenRead(path string) io.Reader {
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
func (c LDSContent) getAction(action string) string {
	return c.BasePath + "?action=" + action
}
func (c LDSContent) GetLanguagesPath() string {
	return c.getAction("languages.query")
}
func (c LDSContent) GetCatalogPath(languageId int) string {
	return c.getAction("catalog.query&platformid=17&languageid=" + string(languageId));
}
func (c LDSContent) GetBookPath(languageId int, glUri string) string {
	return ""
}
func (c LDSContent) OpenRead(path string) io.Reader {
	resp, err := http.Get(path)
	if err != nil {
		panic(err)
	}
	return resp.Body
}
