package main

import (
	"html/template"
	"strings"
	"fmt"
	"net/http"
	"ldslib"
)

type web struct {
	appinfo
	templates *webtemplates
}

type webtemplates struct {
	nodeChildren, nodeContent, layout, err *template.Template
}

type webtemplateData struct {
	data interface{}
	
}

func (app web) run() {
	http.HandleFunc("/", app.handler)
	
	app.initTemplates()
	
	app.efmt.Printf("Listening on port: %v\n", app.config.op.WebPort)
	http.ListenAndServe(fmt.Sprintf(":%v", app.config.op.WebPort), nil)
}

func (app *web) initTemplates() {
	app.templates = &webtemplates{}
	app.templates.layout = app.loadTemplate("layout.html")
	app.templates.nodeContent = app.loadTemplate("node-content.html")
	app.templates.nodeChildren = app.loadTemplate("node-children.html")
	app.templates.err = app.loadTemplate("403.html")
}

func (app *web) loadTemplate(path string) *template.Template {
	data, err := Asset("data/web/templates/" + path)
	if err != nil {
		panic(err)
	}
	temp := template.New(path)
	temp, err = temp.Parse(string(data))
	if err != nil {
		panic(err)
	}
	return temp
}

func (app *web) loadLayoutTemplate(path string, layout *template.Template) {
	//temp := app.loadTemplate(path)
}

func (app *web) handler(w http.ResponseWriter, r *http.Request) {
	//TODO Remove for production
	app.initTemplates()
	
	defer func() {
		if rec := recover(); rec != nil {
			switch rec.(type) {
				case ldslib.NotDownloadedBookErr:
					http.Redirect(w, r, "/download", http.StatusFound)
				case ldslib.NotDownloadedCatalogErr:
					http.Redirect(w, r, "/download", http.StatusFound)
				case ldslib.NotDownloadedLanguageErr:
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
	
	lang, err := app.config.Library.Language(r.URL.Query().Get("lang"))
	if err != nil {
		lang = app.config.SelectedLanguage()
	}
	catalog := app.config.Catalog(lang)
	
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
				panic (err)
			}
			w.Write(data)
		default:
			item, err := app.config.Library.Lookup(r.URL.Path, catalog)
			if err != nil {
				panic(err)
			}
			app.print(w, r, item)
	}
	
}

func (app *web) print(w http.ResponseWriter, r *http.Request, item ldslib.Item) {
	var err error
	data := struct{
		Item ldslib.Item
		Content template.HTML
		Children []ldslib.Item
	}{}
	data.Item = item
	data.Children, err = app.config.Library.Children(item)
	if err != nil {
		panic (err)
	}
	
	switch item.(type) {
		case ldslib.Node:
			node := item.(ldslib.Node)
			if node.HasContent {
				if data.Content, err = node.Content(); err != nil {
					panic(err)
				}
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