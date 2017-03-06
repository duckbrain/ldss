package ldssa

import (
	"github.com/duckbrain/ldss/lib"
	"github.com/duckbrain/ldss/web"
)

func init() {
	lang, err := lib.LookupLanguage("eng")
	if err != nil {
		panic(err)
	}
	web.Handle(lang)
}
