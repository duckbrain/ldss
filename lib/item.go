package lib

import (
	"fmt"
)

type Item interface {
	Name() string
	Children() ([]Item, error)
	Path() string
	Language() *Lang
	Parent() Item
	Next() Item
	Previous() Item
	fmt.Stringer
}

type Contenter interface {
	Content() (Content, error)
	Subtitle() string
	SectionName() string
}

type Footnoter interface {
	Footnotes(verses []int) ([]Footnote, error)
}

type Reference struct {
	Path              string
	Language          *Lang
	VerseSelected     int
	VersesHighlighted []int
	VersesExtra       []int
	Small, Name       string
	Keywords          []string
}
