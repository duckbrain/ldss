package lib

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
	Item
	Content() (Content, error)
	Subtitle() string
	SectionName() string
}

type Footnoter interface {
	Footnotes(verses []int) ([]Footnote, error)
}
