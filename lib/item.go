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
	// Changing this to content not have the error return may be desierable.
	// I'd want to make sure that caching content in memory forever is not
	// too big. Chances are, we don't want to do that, but expirement first.
	Content() (Content, error)
	Subtitle() string
	SectionName() string
}

type Footnoter interface {
	// Same as Contenter
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
