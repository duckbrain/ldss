package churchofjesuschrist

import (
	"context"
	"path"

	"github.com/duckbrain/ldss/lib"
	"github.com/gobuffalo/packr/v2"
)

var refBox = packr.New("ldsorg_refs", "../../reference")

func (c Client) LoadParser(ctx context.Context, p *lib.ReferenceParser) {
	const ext = ".ldssref"

	logger := ctx.Value(lib.CtxLogger).(lib.Logger)

	files := refBox.List()

	if len(files) == 0 {
		panic("No files in box")
	}

	for _, filename := range files {
		if path.Ext(filename) != ext {
			continue
		}
		baseFilename := path.Base(filename)
		lang := lib.Lang(baseFilename[:len(baseFilename)-len(ext)])
		file, err := refBox.Open(filename)
		if err != nil {
			panic(err)
		}
		logger.Debugf("loading %v file %v\n", lang, filename)
		p.AppendFile(lang, file)
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}
}
