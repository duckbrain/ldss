package gowebview

import (
	"log"
	"net"
	"net/http"
	"time"

	//importing the users package that will attach the handlers to the DefaultServeMux
	"github.com/duckbrain/ldss/lib/http"
	"github.com/duckbrain/ldss/lib"
	"github.com/duckbrain/ldss/lib/sources/churchofjesuschrist"
	"github.com/duckbrain/ldss/lib/storages/filestore"
)

var server = &http.Server{
	Addr:           "127.0.0.1:0",
	Handler:        http.DefaultServeMux,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}

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

//Start is called by the native portion of the webapp to start the web server.
//It returns the server root URL (without the trailing slash) and any errors.
func Start() string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		err = server.Serve(listener)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	return listener.Addr().String()
}

//Stop is called by the native portion of the webapp to stop the web server.
func Stop() {
	server.Close()
}
