package lib

var myDummyLang Lang

type dummyItem struct{}
type dummyLang struct{}

func init() {
	myDummyLang = new(dummyLang)
}

func (i dummyItem) Name() string     { return "Dummy" }
func (i dummyItem) Children() []Item { return []Item{} }
func (i dummyItem) Path() string     { return "/dummy" }
func (i dummyItem) Lang() Lang       { return myDummyLang }
func (i dummyItem) Parent() Item     { return nil }
func (i dummyItem) Next() Item       { return nil }
func (i dummyItem) Prev() Item       { return nil }
func (i dummyItem) String() string   { return "{dummy}" }

func (l dummyLang) Name() string          { return "Demmy" }
func (l dummyLang) EnglishName() string   { return "Dummy" }
func (l dummyLang) Code() string          { return "dmy" }
func (l dummyLang) Matches(s string) bool { return s == "dmy" }
