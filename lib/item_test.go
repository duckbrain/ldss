package lib

var dummyLang Lang

type dummyItem struct{}

func init() {
	dummyLang = new(Lang)
}

func (i dummyItem) Name() string              { return "Dummy" }
func (i dummyItem) Children() ([]Item, error) { return nil, nil }
func (i dummyItem) Path() string              { return "/dummy" }
func (i dummyItem) Language() Lang            { return dummyLang }
func (i dummyItem) Parent() Item              { return nil }
func (i dummyItem) Next() Item                { return nil }
func (i dummyItem) Previous() Item            { return nil }
func (i dummyItem) String() string            { return "{dummy}" }
