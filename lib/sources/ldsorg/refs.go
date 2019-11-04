package ldsorg

import (
	"path"

	"github.com/duckbrain/ldss/lib"
	"github.com/gobuffalo/packr/v2"
)

var refBox = packr.New("ldsorg_refs", "../../reference")

var _ interface{ LoadParser(*lib.ReferenceParser) } = &Client{}

func (c Client) LoadParser(p *lib.ReferenceParser) {
	const ext = ".ldssref"

	for _, filename := range refBox.List() {
		if path.Ext(filename) != ext {
			continue
		}
		baseFilename := path.Base(filename)
		lang := lib.Lang(baseFilename[:len(baseFilename)-len(ext)])
		file, err := refBox.Open(filename)
		if err != nil {
			panic(err)
		}
		p.AppendFile(lang, file)
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	lib.RegisterLoader(Default)
}
