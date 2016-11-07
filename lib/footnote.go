package lib

import (
	"html/template"
)

type Footnote struct {
	Name     string        `json:"name"`
	LinkName string        `json:"linkName"`
	Content  template.HTML `json:"content"`
}

func (f *Footnote) References([]Reference, error) {
	panic("References not implemented")
}
