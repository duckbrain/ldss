package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/duckbrain/ldss/internal/assets"
	"github.com/duckbrain/ldss/lib"
)

var defaultLanguage lib.Lang

type webLayout struct {
	Title       string
	Content     template.HTML
	Footnotes   []lib.Footnote
	Lang        lib.Lang
	Item        lib.Item
	Breadcrumbs []lib.Reference
	Query       string
}

// Handle attaches events to the net/http package, but does not start the web server
func Handle(lang lib.Lang) {
	defaultLanguage = lang

	http.HandleFunc("/", handler)
	http.HandleFunc("/api/", handleJSON)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/favicon.ico", handleStatic)
	http.HandleFunc("/manifest.webmanifest", handleStatic)
	http.HandleFunc("/special/jesus-christ", handleChristStudy)
	http.HandleFunc("/css", handleStatic)

	initTemplates()

}

// Run starts listening on the given port
func Run(port int, lang lib.Lang) {
	Handle(lang)
	log.Printf("Listening on port: %v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func language(r *http.Request) lib.Lang {
	lang, err := lib.LookupLanguage(r.URL.Query().Get("lang"))
	if err != nil {
		return defaultLanguage
	}
	return lang
}

func handleError(w io.Writer, r *http.Request) {
	if rec := recover(); rec != nil {
		var err error
		switch rec.(type) {
		case error:
			err = rec.(error)
		default:
			err = fmt.Errorf("%v", rec)
		}
		templates.err.Execute(w, err)
	}
}

// HandleError writes error information from a panic to the web stream
func HandleError(w io.Writer, r *http.Request) {
	if rec := recover(); rec != nil {
		var err error
		switch rec.(type) {
		case error:
			err = rec.(error)
		default:
			err = fmt.Errorf("%v", rec)
		}
		w.Write([]byte(err.Error()))
	}
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	initTemplates()

	defer r.Body.Close()
	lang := language(r)
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
		Breadcrumbs: []lib.Reference{
			{
				Language: lang,
				Path:     "/",
			},
		},
	}
	buff := new(bytes.Buffer)
	for _, ref := range refs {

		if len(ref.Keywords) == 0 {
			func() {
				defer handleError(buff, r)
				item, err := ref.Lookup()
				if err != nil {
					panic(err)
				}
				print(buff, r, ref, item, true)
				ref.Name = item.Name()
				layout.Breadcrumbs = append(layout.Breadcrumbs, ref)

				if node, ok := item.(*lib.Node); ok {
					footnotes, err := node.Footnotes(ref.VersesHighlighted)
					if err == nil {
						layout.Footnotes = append(layout.Footnotes, footnotes...)
					}
				}

			}()
		} else {
			item, err := ref.Lookup()
			if err != nil {
				func() {
					defer handleError(buff, r)
					panic(err)
				}()
			} else {
				ref.Name = item.Name()
				results := lib.SearchSort(item, ref.Keywords)
				layout.Breadcrumbs = append(layout.Breadcrumbs, ref)
				templates.searchResults.Execute(buff, struct {
					Item          lib.Item
					Keywords      []string
					SearchString  string
					SearchResults []lib.SearchResult
				}{
					Item:          item,
					Keywords:      ref.Keywords,
					SearchString:  strings.Join(ref.Keywords, " "),
					SearchResults: results,
				})
			}
		}

		//results = append(results, template.HTML(buff.String()))
	}
	layout.Content = template.HTML(buff.String())

	templates.layout.Execute(w, layout)
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	defer handleError(w, r)
	defer r.Body.Close()
	if err := static(w, r); err != nil {
		panic(err)
	}
}

func static(w http.ResponseWriter, r *http.Request) error {
	data, err := assets.Asset("data/web/static" + r.URL.Path)
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

func itemsRelativesPath(item lib.Item) interface{} {
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
	}
	return nil
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	defer handleError(w, r)
	defer r.Body.Close()

	lang := language(r)
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
	data["parent"] = itemsRelativesPath(item.Parent())
	data["next"] = itemsRelativesPath(item.Next())
	data["previous"] = itemsRelativesPath(item.Previous())

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
			children[i] = itemsRelativesPath(child)
		}
		data["children"] = children
	}

	breadcrumbs := make([]interface{}, 0)
	for p := item; p != nil; {
		parent := p.Parent()
		breadcrumbs = append([]interface{}{itemsRelativesPath(p)}, breadcrumbs...)
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

func handler(w http.ResponseWriter, r *http.Request) {
	defer handleError(w, r)
	defer r.Body.Close()

	if static(w, r) == nil {
		return
	}

	lang := language(r)
	buff := new(bytes.Buffer)

	//TODO Remove for production
	initTemplates()

	ref := lib.ParsePath(lang, r.URL.Path)
	var children []lib.Item

	item, err := lib.AutoDownload(func() (item lib.Item, err error) {
		item, err = ref.Lookup()
		if err != nil {
			return
		}
		children, err = item.Children()
		if err != nil {
			return
		}

		return
	})
	if err != nil {
		panic(err)
	}

	if len(children) == 1 {
		http.Redirect(w, r, children[0].Path(), 301)
		return
	}
	print(buff, r, ref, item, false)

	layout := webLayout{
		Title:       item.Name(),
		Content:     template.HTML(buff.String()),
		Lang:        lang,
		Item:        item,
		Breadcrumbs: make([]lib.Reference, 0),
		Query:       ref.String(),
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
		layout.Breadcrumbs = append([]lib.Reference{{
			Path:     p.Path(),
			Name:     p.Name(),
			Language: p.Language(),
		}}, layout.Breadcrumbs...)
	}

	templates.layout.Execute(w, layout)
}

func print(w io.Writer, r *http.Request, ref lib.Reference, item lib.Item, filter bool) {
	var err error
	data := struct {
		Item      lib.Item
		Reference lib.Reference
		Content   template.HTML
		Children  []lib.Item
		LangCode  string
		HasTitle  bool
		Filtered  bool
	}{
		Item:      item,
		Reference: ref,
		LangCode:  item.Language().GlCode,
	}
	data.Children, err = item.Children()
	if err != nil {
		panic(err)
	}

	switch item.(type) {
	case *lib.Node:
		if content, err := item.(*lib.Node).Content(); err == nil {
			if filter {
				content = content.Filter(ref.VersesHighlighted)
				data.Content = template.HTML(content)
				data.HasTitle = false
				data.Filtered = true
			} else {
				content = content.Highlight(ref.VersesHighlighted, "highlight")
				content = content.Highlight(ref.VersesExtra, "highlight")
				data.Content = template.HTML(content)
				data.HasTitle = strings.Contains(string(content), "</h1>")
			}
			err = templates.nodeContent.Execute(w, data)
		} else {
			err = templates.nodeChildren.Execute(w, data)
		}
	default:
		err = templates.nodeChildren.Execute(w, data)
	}

	if err != nil {
		panic(err)
	}
}
