package churchofjesuschrist

import (
	"github.com/duckbrain/ldss/lib"
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

func (d Dynamic) Item() lib.Item {
	i := lib.Item{}
	i.Name = d.Title()
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
