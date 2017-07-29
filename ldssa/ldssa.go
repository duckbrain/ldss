package ldssa

import (
	"github.com/duckbrain/ldss/lib"
	"github.com/duckbrain/ldss/web"
)

func init() {
	lib.DataDirectory = "/storage/emulated/0/.ldss"
	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			defer web.HandleError(w, r)

			user, _ := user.Current()

			w.Write([]byte(user.HomeDir))
		})
		//*/

	//*
	lang, err := lib.LookupLanguage("eng")
	if err != nil {
		panic(err)
	}
	web.Handle(lang)
	//*/
}
