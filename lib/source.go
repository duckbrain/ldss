package lib

import (
	"fmt"

	"github.com/duckbrain/ldss/lib/download"
)

const RootPath = "/"

type Source interface {
	download.Downloader
	Langs() ([]Lang, error)

	Lookup(lang Lang, path string) (Item, error)
}

var srcs map[string]Source

func init() {
	srcs = make(map[string]Source)
}

func Register(name string, src Source) {
	if _, ok := srcs[name]; ok {
		panic(fmt.Errorf("Cannot have two sources with name %v", name))
	}
	srcs[name] = src
}
