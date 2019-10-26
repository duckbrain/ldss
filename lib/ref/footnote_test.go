package lib

import (
	"html/template"
	"testing"
)

func TestFootnoteParse(t *testing.T) {

	test := func(s string, b ...Reference) {
		f := Footnote{
			Content: template.HTML(s),
			Item:    dummyItem{},
		}
		testReferences(t, f.References(), b...)
	}

	test(`<a href="/scriptures/nt/rom/16.24.20-24" class="scriptureRef">Rom. 16:24 (20–24)</a>.`, Reference{
		Path:              "/scriptures/nt/rom/16",
		VersesHighlighted: []int{24},
		VersesExtra:       []int{20, 21, 22, 23, 24},
		Name:              "Rom. 16:24 (20–24)",
	}, Reference{
		Name: ".",
	})
	test(`<a href="/scriptures/tg/jesus-christ-lord" class="scriptureRef"><small>TG</small> Jesus Christ, Lord</a>`, Reference{
		Path:  "/scriptures/tg/jesus-christ-lord",
		Small: "TG",
		Name:  "Jesus Christ, Lord",
	})
	test(`<span class="small">JST</span> Rev. 20:6 Blessed and holy <em>are they who have</em> part in the first resurrection`, Reference{
		Small: "JST",
	}, Reference{
		Name: "Rev. 20:6 Blessed and holy are they who have part in the first resurrection",
	})
}
