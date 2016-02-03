// +build !noweb

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"ldss/lib"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type web struct {
	appinfo
	templates *webtemplates
}

func init() {
	apps["web"] = &web{}
	lib.Config().RegisterOption(lib.ConfigOption{
		Name:     "WebPort",
		Default:  1830,
		ShortArg: 'p',
		Parse: func(arg string) (interface{}, error) {
			return strconv.Atoi(arg)
		},
	})
	lib.Config().RegisterOption(lib.ConfigOption{
		Name:    "WebTemplatePath",
		Default: "",
	})
}

func (app web) run() {
	http.HandleFunc("/", app.handler)
	http.HandleFunc("/json/", app.handleJSON)

	port := lib.Config().Get("WebPort").(int)

	app.initTemplates()

	app.efmt.Printf("Listening on port: %v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func (app *web) handleJSON(w http.ResponseWriter, r *http.Request) {
	lang, err := lib.LookupLanguage(r.URL.Query().Get("lang"))
	if err != nil {
		lang, err = lib.DefaultLanguage()
		if err != nil {
			panic(err)
		}
	}
	catalog, err := lang.Catalog()
	if err != nil {
		panic(err)
	}
	path := r.URL.Path[len("/json"):]
	fmt.Println("Looking for :" + path)
	item, err := catalog.Lookup(path)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(data)
	if err != nil {
		panic(err)
	}
}

func (app *web) handler(w http.ResponseWriter, r *http.Request) {
	//TODO Remove for production
	app.initTemplates()

	defer func() {
		if rec := recover(); rec != nil {
			app.debug.Println(reflect.TypeOf(rec))
			switch rec.(type) {
			case *lib.NotDownloadedBookErr:
				http.Redirect(w, r, "/download", http.StatusFound)
			case *lib.NotDownloadedCatalogErr:
				http.Redirect(w, r, "/download", http.StatusFound)
			case *lib.NotDownloadedLanguageErr:
				http.Redirect(w, r, "/download", http.StatusFound)
			case error:
				err := rec.(error)
				app.templates.err.Execute(w, err)
			default:
				err := fmt.Errorf("%v", rec)
				app.templates.err.Execute(w, err)
			}
		}
	}()

	lang, err := lib.LookupLanguage(r.URL.Query().Get("lang"))
	if err != nil {
		lang, err = lib.DefaultLanguage()
		if err != nil {
			panic(err)
		}
	}
	catalog, err := lang.Catalog()
	if err != nil {
		panic(err)
	}

	switch strings.ToLower(r.URL.Path) {
	case "/download":
		switch r.Method {
		case "GET":
			fmt.Fprintf(w, "Download Page Here")
		case "POST":
			fmt.Fprintf(w, "Download content here")
		}
	case "/":
		app.print(w, r, catalog)
	case "/favicon.ico":
		w.Header().Set("Content-type", "image/x-icon")
		data, err := Asset("data/web/favicon.ico")
		if err != nil {
			panic(err)
		}
		w.Write(data)
	default:
		item, err := catalog.Lookup(r.URL.Path)
		if err != nil {
			panic(err)
		}
		app.print(w, r, item)
	}

}

func (app *web) print(w http.ResponseWriter, r *http.Request, item lib.Item) {
	var err error
	data := struct {
		Item     lib.Item
		Content  template.HTML
		Children []lib.Item
		LangCode string
	}{}
	data.Item = item
	data.LangCode = item.Language().GlCode
	data.Children, err = item.Children()
	if err != nil {
		panic(err)
	}

	switch item.(type) {
	case *lib.Node:
		node := item.(*lib.Node)
		if content, err := node.Content(); err == nil {
			data.Content = content.HTML()
			err = app.templates.nodeContent.Execute(w, data)
		} else {
			err = app.templates.nodeChildren.Execute(w, data)
		}
	default:
		err = app.templates.nodeChildren.Execute(w, data)
	}
	if err != nil {
		panic(err)
	}
}
