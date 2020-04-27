package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/duckbrain/ldss/lib"
	packr "github.com/gobuffalo/packr/v2"
)

var staticBox = packr.New("ldss_web_static", "./static")

type webLayout struct {
	Title      string
	Content    template.HTML
	Item       lib.Item
	NextHeader lib.Header
	PrevHeader lib.Header
	Query      string
}

type Server struct {
	Lang lib.Lang
	Lib  *lib.Library
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

func (s Server) language(r *http.Request) lib.Lang {
	lang, err := lib.ParseLang(r.URL.Query().Get("lang"))
	if err == lib.ErrNotFound {
		return s.Lang
	}
	return lang
}

func handleError(w io.Writer, r *http.Request, err error) {
	if err == nil {
		if rec := recover(); rec != nil {
			switch r := rec.(type) {
			case error:
				err = r
			default:
				err = fmt.Errorf("%v", rec)
			}
		}
	}
	if err == nil {
		return
	}

	err = templates.err.Execute(w, err)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	initTemplates()

	defer r.Body.Close()
	lang := s.language(r)
	query := r.URL.Query().Get("q")
	refs, err := s.Lib.Parser.Parse(lang, query)
	if err != nil {
		handleError(w, r, err)
		return
	}

	s.Lib.Logger.Infof(`search "%v" -> %v`, query, refs)

	if len(refs) == 1 && refs[0].Query == "" {
		http.Redirect(w, r, refs[0].URL(), http.StatusFound)
		return
	}

	layout := webLayout{
		Title: "LDS Scriptures",
		Query: query,
	}
	buff := new(bytes.Buffer)
	for _, ref := range refs {
		if ref.Query == "" {
			item, err := s.Lib.LookupReference(r.Context(), &ref)
			if err != nil {
				panic(err)
			}
			print(buff, r, ref, item, true)
		} else {
			item, err := s.Lib.LookupReference(r.Context(), &ref)
			if err != nil {
				handleError(buff, r, err)
				return
			} else {
				results, err := s.Lib.SearchSlice(r.Context(), ref)
				if err != nil {
					handleError(buff, r, err)
					return
				}
				err = templates.searchResults.Execute(buff, struct {
					Item          lib.Item
					SearchString  string
					SearchResults lib.Results
				}{
					Item:          item,
					SearchString:  ref.Query,
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

	err = templates.layout.Execute(w, layout)
	if err != nil {
		panic(err)
	}
}

func (s *Server) handleJSON(w http.ResponseWriter, r *http.Request) {
	defer handleError(w, r, nil)
	defer r.Body.Close()

	lang := s.language(r)
	path := r.URL.Path[len("/api"):]
	ref := s.Lib.Parser.ParsePath(lang, path)
	item, err := s.Lib.LookupReference(r.Context(), &ref)
	if err != nil {
		panic(err)
	}

	j, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(j)
	if err != nil {
		panic(err)
	}
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer handleError(w, r, nil)
	defer r.Body.Close()

	lang := s.language(r)
	buff := new(bytes.Buffer)

	s.Lib.Logger.Infof("http %v %v", r.Method, r.URL.Path)

	//TODO Remove for production
	initTemplates()

	ref := s.Lib.Parser.ParsePath(lang, r.URL.Path)

	item, err := s.Lib.LookupReference(ctx, &ref)
	if err != nil {
		panic(err)
	}

	if len(item.Children) == 1 {
		http.Redirect(w, r, item.Children[0].Path, 301)
		return
	}
	print(buff, r, ref, item, false)

	prev, err := s.Lib.Sibling(ctx, item, -1)
	if err != nil {
		panic(err)
	}
	next, err := s.Lib.Sibling(ctx, item, 1)
	if err != nil {
		panic(err)
	}

	layout := webLayout{
		Title:      item.Name,
		Content:    template.HTML(buff.String()),
		Item:       item,
		PrevHeader: prev,
		NextHeader: next,
		Query:      ref.String(),
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
		HasTitle  bool
		Filtered  bool
	}{
		Item:      item,
		Reference: ref,
		Content:   template.HTML(item.Content),
	}

	if err := templates.item.Execute(w, data); err != nil {
		panic(err)
	}
}