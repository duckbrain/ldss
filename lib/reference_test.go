package lib

import "testing"

func TestParsePath(t *testing.T) {
	testParse := func(in string, r Reference) {
		t.Logf("Testing string \"%v\" for match %v", in, r.URL())
		testReference(t, ParsePath(nil, in), r)
	}

	test := func(in, out string, verses ...int) {
		testParse(in, Reference{
			Path:              out,
			VersesHighlighted: verses,
		})
	}

	test("/scriptures/bofm/1-ne/3", "/scriptures/bofm/1-ne/3")
	test("/scriptures/bofm/1-ne/3.4", "/scriptures/bofm/1-ne/3", 4)
	test("/scriptures/bofm/1-ne/3.4-5", "/scriptures/bofm/1-ne/3", 4, 5)
	test("/scriptures/bofm/1-ne/3.4-6", "/scriptures/bofm/1-ne/3", 4, 5, 6)
	test("/scriptures/bofm/1-ne/3.4,6", "/scriptures/bofm/1-ne/3", 4, 6)
	test("/scriptures/bofm/1-ne/3.4-6,6", "/scriptures/bofm/1-ne/3", 4, 5, 6)
	test("/scriptures/bofm/1-ne/3.4-6,6-8,2", "/scriptures/bofm/1-ne/3", 2, 4, 5, 6, 7, 8)
	testParse("/scriptures/bofm/1-ne/3.4.2-6", Reference{
		Path:              "/scriptures/bofm/1-ne/3",
		VersesHighlighted: []int{4},
		VersesExtra:       []int{2, 3, 4, 5, 6},
	})
}

func TestReferenceClean(t *testing.T) {
	testReference(t, Reference{
		VersesHighlighted: []int{1, 3, 4, 6, 7},
	}, Reference{
		VersesHighlighted: []int{1, 1, 3, 4, 1, 7, 6},
	})
}

func testVerses(t *testing.T, a, b []int) {
	if len(a) != len(b) {
		t.Errorf("    Verse range %v is not the same length as %v", a, b)
	} else {
		for i, x := range a {
			if b[i] != x {
				t.Errorf("    Verse at index %v->%v does not match  %v", i, x, b[i])
			}
		}
	}
}

func testReference(t *testing.T, p, r Reference) {
	p.Clean()
	r.Clean()

	if p.Path != r.Path {
		t.Errorf("    Paths don't match %v != %v", p.Path, r.Path)
	}

	testVerses(t, p.VersesHighlighted, r.VersesHighlighted)
	testVerses(t, p.VersesExtra, r.VersesExtra)

	if p.VerseSelected != r.VerseSelected {
		t.Errorf("    VerseSelected doesn't match %v != %v", p.VerseSelected, r.VerseSelected)
	}
}
