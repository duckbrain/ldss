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
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			app.templates.err.Execute(w, err)
			//panic(err)
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

func (app *web) print(w http.ResponseWriter, r *http.Request, item ldslib.CatalogItem) {
	var err error
	switch item.(type) {
		case ldslib.Node:
			node := item.(ldslib.Node)
			if node.HasContent {
				data := struct{
					item ldslib.CatalogItem
					children []ldslib.CatalogItem
				}{}
				data.item = item
				data.children, err = item.Children()
				if err != nil {
					panic (err)
				}
				err = app.templates.nodeContent.Execute(w, item)
			} else {
				err = app.templates.nodeChildren.Execute(w, item)
			}
		default:
			err = app.templates.nodeChildren.Execute(w, item)
	}
	if err != nil {
		//TODO: Write 403 error
		panic(err)
	}
}