package web

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/duckbrain/ldss/lib"
)

const topicalGuidePath = "/scriptures/tg/"

func writeTitle(w io.Writer, tag, s, href string) {
	w.Write([]byte(fmt.Sprintf("<%v><a href=\"%v\">%v</a></%v>", tag, html.EscapeString(href), html.EscapeString(s), tag)))
}

// handleChristStudy renders the scriptures President Nelson encouraged us to study with the topical guide.
func handleChristStudy(w http.ResponseWriter, r *http.Request) {
	//defer handleError(w, r)
	defer r.Body.Close()

	lang := language(r)
	ref := lib.ParsePath(lang, "/scriptures/tg/jesus-christ")
	buff := new(bytes.Buffer)

	item, err := lib.AutoDownload(ref.Lookup)
	if err != nil {
		panic(err)
	}
	contenter, ok := item.(lib.Contenter)
	if !ok {
		panic(fmt.Errorf("Item %v has no content", item.Name()))
	}
	content, err := contenter.Content()
	if err != nil {
		panic(err)
	}
	traversedPaths := make(map[string]bool)

	writeTitle(buff, "h1", "Topical Study of Jesus Christ", ref.URL())

	for _, ref := range content.Links(lang) {
		if !strings.HasPrefix(ref.Path, topicalGuidePath) {
			continue
		}
		if traversedPaths[ref.URL()] {
			continue
		}
		traversedPaths[ref.URL()] = true
		item, err := lib.AutoDownload(ref.Lookup)
		if err != nil {
			panic(err)
		}
		if item == nil {
			continue
		}
		contenter, ok := item.(lib.Contenter)
		if !ok {
			continue
		}
		content, err := contenter.Content()
		if err != nil {
			panic(err)
		}
		content.Links(item.Lang())

		writeTitle(buff, "h2", item.Name(), ref.URL())

		traversedPaths := make(map[string]bool)

		for _, ref := range content.Links(lang) {
			if strings.HasPrefix(ref.Path, topicalGuidePath) {
				continue
			}
			if traversedPaths[ref.URL()] {
				continue
			}
			traversedPaths[ref.URL()] = true
			item, err := lib.AutoDownload(ref.Lookup)
			if err != nil {
				panic(err)
			}
			if item == nil {
				continue
			}
			contenter, ok := item.(lib.Contenter)
			if !ok {
				continue
			}
			content, err := contenter.Content()
			if err != nil {
				panic(err)
			}
			filteredContent := content.Filter(ref.VersesHighlighted)
			writeTitle(buff, "h3", ref.Name, ref.URL())
			buff.Write([]byte(filteredContent))

			//print(buff, r, ref, item, true)
		}
	}

	layout := webLayout{
		Title:   "Topical Study of Jesus Christ",
		Lang:    lang,
		Content: template.HTML(buff.String()),
		Breadcrumbs: []lib.Reference{
			{
				Lang: lang,
				Path: "/",
			},
		},
	}

	templates.layout.Execute(w, layout)
}
