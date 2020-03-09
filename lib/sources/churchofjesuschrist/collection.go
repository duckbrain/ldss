package churchofjesuschrist

import (
	"encoding/json"
	"github.com/duckbrain/ldss/lib"
	"html/template"
)

type Dynamic struct {
	Collection *Collection
	Content    *Content
}

func (d Dynamic) String() string {
	j, _ := json.Marshal(d)
	return string(j)
}

func (d Dynamic) Title() string {
	switch {
	case d.Content != nil && d.Content.Meta.Title != "":
		return d.Content.Meta.Title
	}
	return ""
}

func (d Dynamic) AsItem(index lib.Index) (i lib.Item) {
	i.Index = index
	i.Merge(d.Content.AsItem(index))
	i.Merge(d.Collection.AsItem(index))
	return
}

type Collection struct {
	Item
	Section
	BreadCrumbs []struct {
		Title string
		URI   string
	}
}

func (c *Collection) AsItem(index lib.Index) (i lib.Item) {
	if c == nil {
		return
	}
	i.Header = c.Item.AsHeader(index.Lang, "")
	i.Merge(c.Section.AsItem(index))
	return
}

type Entry struct {
	Section *Section
	Item    *Item
}

type Section struct {
	Title   string
	Entries []Entry
}

func (s *Section) AsItem(index lib.Index) (l lib.Item) {
	if s == nil {
		return
	}
	for _, e := range s.Entries {
		if e.Item != nil {
			l.Children = append(l.Children, e.Item.AsHeader(index.Lang, s.Title))
		}
	}
	return
}

type Item struct {
	Position int64
	URI      string
	Title    string
	Src      string
	SrcSet   string
}

func (i *Item) AsHeader(lang lib.Lang, sectionName string) lib.Header {
	if i == nil {
		return lib.Header{}
	}
	return lib.Header{
		Index: lib.Index{
			Path: i.URI,
			Lang: lang,
		},
		Name: i.Title,
	}
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

func (c *Content) AsItem(index lib.Index) (i lib.Item) {
	if c == nil {
		return
	}
	i.Content = lib.Content(c.Content.Body)
	for _, f := range c.Content.Footnotes {
		i.Footnotes = append(i.Footnotes, lib.Footnote{
			Name:     f.Marker,
			LinkName: f.Context,
			Content:  template.HTML(f.Text),
		})
	}
	return
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
