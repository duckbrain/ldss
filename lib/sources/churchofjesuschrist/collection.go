package churchofjesuschrist

import (
	"github.com/duckbrain/ldss/lib"
	"html/template"
)

type Dynamic struct {
	Collection *Collection
	Content    *Content
}

func (d Dynamic) Title() string {
	switch {
	case d.Content != nil && d.Content.Meta.Title != "":
		return d.Content.Meta.Title
	}
	return ""
}

func (d Dynamic) Item(index lib.Index) lib.Item {
	i := lib.Item{}
	i.Index = index
	i.Name = d.Title()
	if d.Content != nil {
		i.Content = lib.Content(d.Content.Content.Body)
		for _, f := range d.Content.Content.Footnotes {
			i.Footnotes = append(i.Footnotes, lib.Footnote{
				Name:     f.Marker,
				LinkName: f.Context,
				Content:  template.HTML(f.Text),
			})
		}
	}
	if d.Collection != nil {
		for _, e := range d.Collection.Entries {
			if e.Section != nil {

			}
			if e.Item != nil {
				i.Children = append(i.Children, lib.Index{Path: e.Item.URI, Lang: index.Lang})
			}
		}
	}
	return i
}

type Collection struct {
	Title       string
	URI         string
	Src         string
	SrcSet      string
	BreadCrumbs []struct {
		Title string
		URI   string
	}
	Entries []Entry
}

type Entry struct {
	Section *Section
	Item    *Item
}

type Section struct {
	Title   string
	Entries []Entry
}

type Item struct {
	Position int64
	URI      string
	Title    string
	Src      string
	SrcSet   string
}

type Content struct {
	Meta struct {
		Title       string
		ContentType string
		Audio       []struct {
			MediaURL string
			Voice    struct {
				Name string
			}
		}
		PDF            string
		PageAttributes map[string]string // Attributes to set on the body?
	}
	Content struct {
		Head struct {
			// Contains things for the head of the HTML
		}
		Body      string
		Footnotes map[string]Footnote
	}
}

type Footnote struct {
	ID            string
	Marker        string
	Context       string
	Text          string
	ReferenceURIs []struct {
		Type string
		Href string
		Text string
	}
}

type TOC struct {
}
