package gowebview

import (
	"log"
	"net"
	"net/http"
	"time"

	//importing the users package that will attach the handlers to the DefaultServeMux
	_ "github.com/duckbrain/ldss/ldssa"
)

var server = &http.Server{
	Addr:           "127.0.0.1:0",
	Handler:        http.DefaultServeMux,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
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
