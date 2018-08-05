package lib

type Opener interface {
	Open() error
}
type Closer interface {
	Close() error
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
	Content() Content
	Subtitle() string
	SectionName() string
	Footnotes(verses []int) []Footnote
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
