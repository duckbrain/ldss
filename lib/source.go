package lib

import (
	"fmt"

	"github.com/duckbrain/ldss/lib/dl"
)

const RootPath = "/"

var DataDirectory = ".ldss"

type Source interface {
	dl.Downloader
	Langs() []Lang

	Lookup(lang Lang, path string) (Item, error)
}

var srcs map[string]Source
var opened bool

func init() {
	srcs = make(map[string]Source)
}

func Register(name string, src Source) {
	if _, ok := srcs[name]; ok {
		panic(fmt.Errorf("Cannot have two sources with name %v", name))
	}
	srcs[name] = src
}

func Open() error {
	for name, src := range srcs {
		if x, ok := src.(dl.Downloader); ok {
			if ok := x.Downloaded(); !ok {
				err := dl.EnqueueAndWait(x)
				if err != nil {
					return err
				}
			}
		}
		if x, ok := src.(Opener); ok {
			if err := x.Open(); err != nil {
				return err
			}
		}
		registerLanguage(name, src.Langs())
	}
	opened = true
	return nil
}
