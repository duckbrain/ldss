package ldssa

import (
	"github.com/duckbrain/ldss/internal/web"
	"github.com/duckbrain/ldss/lib"
)

func init() {
	lib.DataDirectory = "/storage/emulated/0/.ldss"
	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			defer web.HandleError(w, r)
			lib.DataDirectory = "/storage/emulated/0/.ldss"
			dir, err := filepath.Abs(lib.DataDirectory)
			w.Write([]byte(dir))
			if err != nil {
				panic(err)
			}

			w.Write([]byte("\n"))

			lang, err := lib.LookupLanguage("eng")
			if err != nil {
				panic(err)
			}
			w.Write([]byte(lang.String()))
		})
	*/

	lang, err := lib.LookupLanguage("eng")
	if err != nil {
		panic(err)
	}
	web.Handle(lang)
}
