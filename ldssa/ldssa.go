package ldssa

import (
	"net/http"

	"github.com/duckbrain/ldss/internal/web"
	"github.com/duckbrain/ldss/lib"
	"github.com/duckbrain/ldss/lib/sources/churchofjesuschrist"
	"github.com/duckbrain/ldss/lib/storages/filestore"
)

func init() {
	store, err := filestore.New("/storage/emulated/0/.ldss")
	if err != nil {
		panic(err)
	}
	library := lib.Default
	library.Store = store
	library.Index = store
	library.Register(churchofjesuschrist.Default)

	server := web.Server{
		Lang: lib.DefaultLang,
		Lib:  library,
	}

	http.Handle("/", server.Handler())
}
