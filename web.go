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
	"strings"
)

type web struct {
	appinfo
	templates *webtemplates
}

type webLayout struct {
	Title       string
	Content     template.HTML
	Footnotes   []lib.Footnote
	Lang        *lib.Language
	Item        lib.Item
	Breadcrumbs []lib.Item
	Query       string
}

func init() {
	apps["web"] = &web{}
}

func (app web) register(config *Configuration) {
	config.RegisterOption(ConfigOption{
		Name:     "WebPort",
		Default:  1830,
		ShortArg: 'p',
		Parse: func(arg string) (interface{}, error) {
			return strconv.Atoi(arg)
		},
	})
	config.RegisterOption(ConfigOption{
		Name:    "WebTemplatePath",
		Default: "",
	})
	app.config = config
}

func (app web) run() {
	http.HandleFunc("/", app.handler)
	http.HandleFunc("/api/", app.handleJSON)
	http.HandleFunc("/search", app.handleSearch)
	http.HandleFunc("/favicon.ico", app.handleStatic)
	http.HandleFunc("/css", app.handleStatic)

	port := app.config.Get("WebPort").(int)
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

func (app *web) handleError(w io.Writer, r *http.Request) {
	if rec := recover(); rec != nil {
		var err error
		switch rec.(type) {
		case error:
			err = rec.(error)
		default:
			err = fmt.Errorf("%", rec)
		}
		app.templates.err.Execute(w, err)
	}
}

func (app *web) handleSearch(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	lang := app.language(r)
	query := r.URL.Query().Get("q")
	refs := lib.Parse(lang, query)
	if len(refs) == 1 && len(refs[0].Keywords) == 0 {
		http.Redirect(w, r, refs[0].URL(), http.StatusFound)
		return
	}

	layout := webLayout{
		Title: "LDS Scriptures",
		Lang:  lang,
		Query: query,
	}
	buff := new(bytes.Buffer)
	for _, ref := range refs {

		if len(ref.Keywords) == 0 {
			func() {
				defer app.handleError(buff, r)
				item, err := ref.Lookup()
				if err != nil {
					panic(err)
				}
				app.print(buff, r, ref, item, true)

				if node, ok := item.(*lib.Node); ok {
					footnotes, err := node.Footnotes(ref.VersesHighlighted)
					if err == nil {
						layout.Footnotes = append(layout.Footnotes, footnotes...)
					}
				}

			}()
		} else {
			c := make(chan lib.SearchResult)
			err := ref.Search(c)
			if err != nil {
				func() {
					defer app.handleError(buff, r)
					panic(err)
				}()
			} else {
				//for r := range c {
				//TODO Print results
				//}
			}
		}

		//results = append(results, template.HTML(buff.String()))
	}
	layout.Content = template.HTML(buff.String())

	app.templates.layout.Execute(w, layout)
}

func (app *web) handleStatic(w http.ResponseWriter, r *http.Request) {
	defer app.handleError(w, r)
	defer r.Body.Close()
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
	//defer app.handleError(w, r)
	defer r.Body.Close()

	lang := app.language(r)
	path := r.URL.Path[len("/api"):]
	ref := lib.ParsePath(lang, path)
	item, err := ref.Lookup()
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
		data["footnotes"], _ = item.Footnotes(ref.VersesHighlighted)
	}

	if childItems, err := item.Children(); err == nil {
		children := make([]interface{}, len(childItems))
		for i, child := range childItems {
			children[i] = app.itemsRelativesPath(child)
		}
		data["children"] = children
	}

	breadcrumbs := make([]interface{}, 0)
	for p := item; p != nil; {
		parent := p.Parent()
		breadcrumbs = append([]interface{}{app.itemsRelativesPath(p)}, breadcrumbs...)
		p = parent
	}
	data["breadcrumbs"] = breadcrumbs

	j, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(j)
	if err != nil {
		panic(err)
	}
}

func (app *web) handler(w http.ResponseWriter, r *http.Request) {
	defer app.handleError(w, r)
	defer r.Body.Close()

	if app.static(w, r) == nil {
		return
	}

	lang := app.language(r)
	buff := new(bytes.Buffer)

	//TODO Remove for production
	app.initTemplates()

	ref := lib.ParsePath(lang, r.URL.Path)
	item, err := ref.Lookup()
	if err != nil {
		panic(err)
	}
	if children, err := item.Children(); err == nil {
		if len(children) == 1 {
			http.Redirect(w, r, children[0].Path(), 301)
			return
		}
	}
	app.print(buff, r, ref, item, false)

	layout := webLayout{
		Title:       "LDS Scriptures",
		Content:     template.HTML(buff.String()),
		Lang:        lang,
		Item:        item,
		Breadcrumbs: make([]lib.Item, 0),
	}

	// Get the footnote content
	if n, ok := item.(*lib.Node); ok {
		layout.Footnotes, err = n.Footnotes(ref.VersesHighlighted)
		if err != nil {
			panic(err)
		}
	}

	// Generate breadcrumbs
	for p := item; p != nil; p = p.Parent() {
		layout.Breadcrumbs = append([]lib.Item{p}, layout.Breadcrumbs...)
	}

	app.templates.layout.Execute(w, layout)
}

func (app *web) print(w io.Writer, r *http.Request, ref lib.Reference, item lib.Item, filter bool) {
	var err error
	data := struct {
		Item     lib.Item
		Content  template.HTML
		Children []lib.Item
		LangCode string
		HasTitle bool
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
			if filter {
				data.Content = template.HTML(content.Filter(ref.VersesHighlighted))
				data.HasTitle = false
			} else {
				data.Content = template.HTML(content)
				data.HasTitle = strings.Contains(string(content), "</h1>")
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
