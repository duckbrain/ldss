package lib

import (
	"strings"
	"html/template"
	"golang.org/x/net/html"
)

type Footnote struct {
	Name     string        `json:"name"`
	LinkName string        `json:"linkName"`
	Content  template.HTML `json:"content"`
}

func (f *Footnote) References() ([]Reference, error) {
	doc, err := html.Parse(strings.NewReader(string(f.Content)))
	if err != nil {
		return nil, err
	}

	refs := make([]Reference, 0)

	//Adapt for parsing refs
	for n := doc.FirstChild; n != nil; n = n.NextSibling {
		//TODO: Parse HTML for ref
	}

	return refs, nil
}
