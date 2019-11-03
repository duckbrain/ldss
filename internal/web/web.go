package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/duckbrain/ldss/lib"
	packr "github.com/gobuffalo/packr/v2"
)

var staticBox = packr.New("ldss_web_static", "./static")

type webLayout struct {
	Title       string
	Content     template.HTML
	Footnotes   []lib.Footnote
	Lang        lib.Lang
	Item        lib.Item
	Breadcrumbs []lib.Header
	Query       string
}

type Server struct {
	Lang lib.Lang
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	handleStatic := http.FileServer(staticBox)

	mux.HandleFunc("/api/", s.handleJSON)
	mux.HandleFunc("/search", s.handleSearch)
	mux.Handle("/favicon.ico", handleStatic)
	mux.Handle("/manifest.webmanifest", handleStatic)
	mux.Handle("/css/", handleStatic)
	mux.Handle("/js/", handleStatic)
	mux.Handle("/svg/", handleStatic)
	mux.HandleFunc("/", s.handler)

	return mux
}

// Run starts listening on the given port
func Run(port int, lang lib.Lang) {
}

func (s Server) language(r *http.Request) lib.Lang {
	lang := lib.Lang(r.URL.Query().Get("lang"))
	if lang == "" {
		return s.Lang
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
		err = templates.err.Execute(w, err)
		if err != nil {
			fmt.Println(err)
		}
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

func (s Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	initTemplates()

	defer r.Body.Close()
	lang := s.language(r)
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
		Breadcrumbs: []lib.Header{
			{
				Lang: lang,
				Path: "/",
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

				if x, ok := item.(lib.Contenter); ok {
					f := x.Footnotes(ref.VersesHighlighted)
					layout.Footnotes = append(layout.Footnotes, f...)
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
				err := templates.searchResults.Execute(buff, struct {
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
				if err != nil {
					panic(err)
				}
			}
		}

		//results = append(results, template.HTML(buff.String()))
	}
	layout.Content = template.HTML(buff.String())

	err := templates.layout.Execute(w, layout)
	if err != nil {
		panic(err)
	}
}

func itemsRelativesPath(item lib.Item) interface{} {
	if item != nil {
		data := struct {
			Name string `json:"name"`
			Path string `json:"path"`
		}{item.Name(), item.Path()}

		return data
	}
	return nil
}

func (s Server) handleJSON(w http.ResponseWriter, r *http.Request) {
	defer handleError(w, r)
	defer r.Body.Close()

	lang := s.language(r)
	path := r.URL.Path[len("/api"):]
	ref := lib.ParsePath(lang, path)
	item, err := ref.Lookup()
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{}

	data["name"] = item.Name()
	data["path"] = item.Path()
	data["language"] = item.Lang().Code()
	data["parent"] = itemsRelativesPath(item.Parent())
	data["next"] = itemsRelativesPath(item.Next())
	data["previous"] = itemsRelativesPath(item.Prev())

	if x, ok := item.(lib.Contenter); ok {
		data["content"] = x.Content()
		data["footnotes"] = x.Footnotes(ref.VersesHighlighted)
	}

	childItems := item.Children()
	children := make([]interface{}, len(childItems))
	for i, child := range childItems {
		children[i] = itemsRelativesPath(child)
	}
	data["children"] = children

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

func (s Server) handler(w http.ResponseWriter, r *http.Request) {
	defer handleError(w, r)
	defer r.Body.Close()

	lang := s.language(r)
	buff := new(bytes.Buffer)

	//TODO Remove for production
	initTemplates()

	ref := lib.ParsePath(lang, r.URL.Path)
	var children []lib.Item

	item, err := ref.Lookup()
	if err != nil {
		panic(err)
	}

	children = item.Children()

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
	if x, ok := item.(lib.Contenter); ok {
		layout.Footnotes = x.Footnotes(ref.VersesHighlighted)
		if err != nil {
			panic(err)
		}
	}

	// Generate breadcrumbs
	for p := item; p != nil; p = p.Parent() {
		layout.Breadcrumbs = append([]lib.Reference{{
			Path: p.Path(),
			Name: p.Name(),
			Lang: p.Lang(),
		}}, layout.Breadcrumbs...)
	}

	err = templates.layout.Execute(w, layout)
	if err != nil {
		panic(err)
	}
}

func print(w io.Writer, r *http.Request, ref lib.Reference, item lib.Item, filter bool) {
	data := struct {
		Item      lib.Item
		Reference lib.Reference
		Content   template.HTML
		LangCode  string
		HasTitle  bool
		Filtered  bool
	}{
		Item:      item,
		Reference: ref,
		LangCode:  item.Lang().Code(),
	}

	// TODO: Support having content and children
	var err error
	var hasContent bool
	if x, ok := item.(lib.Contenter); ok {
		if content := x.Content(); len(content) > 0 {
			hasContent = true
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
		}
	}
	if !hasContent {
		err = templates.nodeChildren.Execute(w, data)
	}

	if err != nil {
		panic(err)
	}
}
