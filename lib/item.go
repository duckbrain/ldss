package lib

type Opener interface {
	Open() error
}

type Item interface {
	Name() string
	Children() []Item
	Path() string
	Lang() Lang
	Parent() Item
	Next() Item
	Prev() Item
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
	Lang              Lang
	VerseSelected     int
	VersesHighlighted []int
	VersesExtra       []int
	Small, Name       string
	Keywords          []string
}
