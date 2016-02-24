// +build !noweb

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"ldss/lib"
	"net/http"
	"path"
	"strconv"
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
	http.HandleFunc("/api/", app.handleJSON)
	http.HandleFunc("/lookup", app.handleLookup)
	http.HandleFunc("/favicon.ico", app.handleStatic)
	http.HandleFunc("/css", app.handleStatic)

	port := lib.Config().Get("WebPort").(int)
	app.initTemplates()
	app.efmt.Printf("Listening on port: %v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func (app *web) lang(r *http.Request) *lib.Language {
	lang, err := lib.LookupLanguage(r.URL.Query().Get("lang"))
	if err != nil {
		lang, err = lib.DefaultLanguage()
		if err != nil {
			panic(err)
		}
	}
	return lang
}

func (app *web) handleError(w http.ResponseWriter, r *http.Request) {
	if rec := recover(); rec != nil {
		var err error
		switch rec.(type) {
		case error:
			err = rec.(error)
		default:
			err = fmt.Errorf("%v", rec)
		}
		app.templates.err.Execute(w, err)
	}
}

func (app *web) handleLookup(w http.ResponseWriter, r *http.Request) {
	path, err := app.lang(r).Reference(r.URL.Query().Get("q"))
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, path, http.StatusFound)
}

func (app *web) handleStatic(w http.ResponseWriter, r *http.Request) {
	defer app.handleError(w, r)
	if err := app.static(w, r); err != nil {
		panic(err)
	}
}

func (app *web) static(w http.ResponseWriter, r *http.Request) error {
	data, err := Asset("data/web/static" + r.URL.Path)
	if err != nil {
		return err
	}

	switch path.Ext(r.URL.Path) {
	case ".ico":
		w.Header().Set("Content-type", "image/x-icon")
	case ".css":
		w.Header().Set("Content-type", "text/css")
	case ".js":
		w.Header().Set("Content-type", "application/x-javascript")
	case ".svg":
		w.Header().Set("Content-type", "image/svg+xml")
	default:
		panic(fmt.Errorf("Unknown extension"))
	}

	w.Write(data)
	return nil
}

func (app *web) handleJSON(w http.ResponseWriter, r *http.Request) {
	lang := app.lang(r)
	catalog, err := lang.Catalog()
	if err != nil {
		panic(err)
	}
	path := r.URL.Path[len("/json"):]
	fmt.Println("Looking for :" + path)
	item, err := catalog.LookupPath(path)
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
	defer app.handleError(w, r)

	if app.static(w, r) == nil {
		return
	}

	lang := app.lang(r)
	buff := new(bytes.Buffer)
	//TODO Remove for production
	app.initTemplates()

	catalog, err := lang.Catalog()
	if err != nil {
		panic(err)
	}

	item, err := catalog.LookupPath(r.URL.Path)
	if err != nil {
		panic(err)
	}
	app.print(buff, r, item)

	layout := struct {
		Title   string
		Content template.HTML
		Item    lib.Item
	}{
		Title:   "LDS Scriptures",
		Content: template.HTML(buff.String()),
		Item:    item,
	}
	app.templates.layout.Execute(w, layout)
}

func (app *web) print(w io.Writer, r *http.Request, item lib.Item) {
	var err error
	data := struct {
		Item     lib.Item
		Content  template.HTML
		Children []lib.Item
		LangCode string
	}{
		Item:     item,
		LangCode: item.Language().GlCode,
	}
	data.Children, err = item.Children()
	if err != nil {
		panic(err)
	}

	switch item.(type) {
	case *lib.Node:
		if content, err := item.(*lib.Node).Content(); err == nil {
			data.Content = template.HTML(content)
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
