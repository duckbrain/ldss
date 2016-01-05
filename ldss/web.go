package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"ldss/lib"
	"net/http"
	"strings"
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
	http.HandleFunc("/json/", app.handleJSON)

	app.initTemplates()

	app.efmt.Printf("Listening on port: %v\n", app.config.op.WebPort)
	http.ListenAndServe(fmt.Sprintf(":%v", app.config.op.WebPort), nil)
}

func (app *web) initTemplates() {
	app.templates = &webtemplates{}
	app.templates.layout = app.loadTemplate("layout.tpl")
	app.templates.nodeContent = app.loadTemplate("node-content.tpl")
	app.templates.nodeChildren = app.loadTemplate("node-children.tpl")
	app.templates.err = app.loadTemplate("403.tpl")
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
	// Need to create Executor interface to make this work
}

func (app *web) handleJSON(w http.ResponseWriter, r *http.Request) {
	lang, err := app.config.Library.Language(r.URL.Query().Get("lang"))
	if err != nil {
		lang = app.config.SelectedLanguage()
	}
	catalog := app.config.Catalog(lang)
	path := r.URL.Path[len("/json"):]
	fmt.Println("Looking for :" + path)
	item, err := app.config.Library.Lookup(path, catalog)
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
			switch rec.(type) {
			case lib.NotDownloadedBookErr:
				http.Redirect(w, r, "/download", http.StatusFound)
			case lib.NotDownloadedCatalogErr:
				http.Redirect(w, r, "/download", http.StatusFound)
			case lib.NotDownloadedLanguageErr:
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
			panic(err)
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
	data.Children, err = app.config.Library.Children(item)
	if err != nil {
		panic(err)
	}

	switch item.(type) {
	case lib.Node:
		node := item.(lib.Node)
		if node.HasContent {
			content, err := node.Content()
			if err != nil {
				panic(err)
			}
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
