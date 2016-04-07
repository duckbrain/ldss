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
	Config().RegisterOption(ConfigOption{
		Name:     "WebPort",
		Default:  1830,
		ShortArg: 'p',
		Parse: func(arg string) (interface{}, error) {
			return strconv.Atoi(arg)
		},
	})
	Config().RegisterOption(ConfigOption{
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

	port := Config().Get("WebPort").(int)
	app.initTemplates()
	app.efmt.Printf("Listening on port: %v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func (app *web) language(r *http.Request) *lib.Language {
	lang, err := lib.LookupLanguage(r.URL.Query().Get("lang"))
	if err != nil {
		return app.lang
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
	path, err := app.language(r).Reference(r.URL.Query().Get("q"))
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

func (app *web) itemsRelativesPath(item lib.Item) interface{} {
	if item != nil {
		data := struct {
			Name string `json:"name"`
			Type string `json:"type"`
			Path string `json:"path"`
		}{item.Name(), "", item.Path()}

		switch item.(type) {
		case *lib.Catalog:
			data.Type = "catalog"
		case *lib.Folder:
			data.Type = "folder"
		case *lib.Book:
			data.Type = "book"
		case *lib.Node:
			data.Type = "node"
		}

		return data
	} else {
		return nil
	}
}

func (app *web) handleJSON(w http.ResponseWriter, r *http.Request) {
	defer app.handleError(w, r)

	lang := app.language(r)
	catalog, err := lang.Catalog()
	if err != nil {
		panic(err)
	}
	path := r.URL.Path[len("/api"):]
	item, err := catalog.LookupPath(path)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{}

	data["name"] = item.Name()
	data["path"] = item.Path()
	data["language"] = item.Language().GlCode
	data["parent"] = app.itemsRelativesPath(item.Parent())
	data["next"] = app.itemsRelativesPath(item.Next())
	data["previous"] = app.itemsRelativesPath(item.Previous())

	switch item := item.(type) {
	case *lib.Catalog:
		data["type"] = "catalog"
	case *lib.Folder:
		data["type"] = "folder"
	case *lib.Book:
		data["type"] = "book"
	case *lib.Node:
		data["type"] = "node"
		data["content"], _ = item.Content()
	}

	if childItems, err := item.Children(); err == nil {
		children := make([]interface{}, len(childItems))
		for i, child := range childItems {
			children[i] = app.itemsRelativesPath(child)
		}
		data["children"] = children
	}

	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(json)
	if err != nil {
		panic(err)
	}
}

func (app *web) handler(w http.ResponseWriter, r *http.Request) {
	defer app.handleError(w, r)

	if app.static(w, r) == nil {
		return
	}

	lang := app.language(r)
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
		Lang    *lib.Language
		Item    lib.Item
	}{
		Title:   "LDS Scriptures",
		Content: template.HTML(buff.String()),
		Lang:    lang,
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

type webtemplates struct {
	nodeChildren, nodeContent, layout, err *template.Template
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
