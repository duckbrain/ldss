package lib

import (
	"testing"
)

const file = `
# This file describes the format for reference files

# Blank lines are ignored, lines with # in front are ignored

# Lookup mapping
1 Nephi:1 ne:1ne:1nephi:/scriptures/bofm/1-ne#

# Regex matching
/([1-4])( |-)?nephi/:/scriptures/bofm/${1}-ne

# Folder mapping
42762:/music
`

func testReferences(t *testing.T, a []Reference, b ...Reference) {
	if len(a) != len(b) {
		t.Errorf("Number of references differs between sets %v and %v", a, b)
	} else {
		for i, x := range a {
			testReference(t, x, b[i])
		}
	}
}

func TestReferenceParseBasic(t *testing.T) {
	p := newRefParser([]byte(file))
	if p.matchFolder[42762] != "/music" {
		t.Fail()
	}
	if p.matchString["1 ne "] != "/scriptures/bofm/1-ne#" {
		t.Fail()
	}
	if len(p.matchRegexp) != 1 {
		t.Fail()
	}
}

func TestReferenceParseDuplicate(t *testing.T) {
	test := func(code string) {
		defer func() {
			if recover() == nil {
				t.Fail()
			}
		}()
		newRefParser([]byte(code))
	}
	test(`
# Lookup mapping (has two "1ne")
1 Nephi:1ne:1ne:1nephi:/scriptures/bofm/1-ne
`)
	test(`
# Folder mapping (two of same number)
42762:/music
42762:/music2
`)
}

func TestReferenceLookup(t *testing.T) {
	p := newRefParser([]byte(file))

	testQuery := func(in string, r ...Reference) {
		t.Logf("Testing string \"%v\" for match %v", in, r)
		p := p.lookup(in)

		testReferences(t, p, r...)
	}

	test := func(in, out string, verses ...int) {
		s := 0
		if len(verses) > 0 {
			s = verses[0]
		}
		testQuery(in, Reference{
			Path:              out,
			VersesHighlighted: verses,
			VerseSelected:     s,
		})
	}

	test("1ne 3", "/scriptures/bofm/1-ne/3")
	test("1ne 3:4", "/scriptures/bofm/1-ne/3", 4)
	test("1ne 3:4-5", "/scriptures/bofm/1-ne/3", 4, 5)
	test("1ne 3:4-6", "/scriptures/bofm/1-ne/3", 4, 5, 6)
	test("1ne 3:4,6", "/scriptures/bofm/1-ne/3", 4, 6)
	test("1ne 3:4-6,6", "/scriptures/bofm/1-ne/3", 4, 5, 6)
	test("1ne 3:4-6,6-8, 2", "/scriptures/bofm/1-ne/3", 2, 4, 5, 6, 7, 8)
	testQuery("1ne 3:4 (2-6)", Reference{
		Path:              "/scriptures/bofm/1-ne/3",
		VerseSelected:     4,
		VersesHighlighted: []int{4},
		VersesExtra:       []int{2, 3, 4, 5, 6},
	})
}
