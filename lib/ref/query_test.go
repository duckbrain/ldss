package ref

import "testing"

const file = `
# This file describes the format for reference files

# Blank lines are ignored, lines with # in front are ignored

# Lookup mapping
1 Nephi:1 ne:1ne:1nephi:/scriptures/bofm/1-ne#

# Regex matching
/([1-4])( |-)?nephi/:/scriptures/bofm/${1}-ne

# Folder mapping
42762:/music

# Joseph Smith History was having issues
/(joseph smith|js)( |\-|\-\-|—)?h(istory)?/:/scriptures/pgp/js-h/1

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
	p := newQueryParser(myDummyLang, []byte(file))
	if p.matchFolder[42762] != "/music" {
		t.Error("Unable to reverse lookup /music folder")
	}
	if p.matchString["1 ne "] != "/scriptures/bofm/1-ne#" {
		t.Error("Unable to match \"1 ne\" with /scriptures/bofm/1-ne#")
	}
	if len(p.matchRegexp) != 2 {
		// There are two regular expressions in file
		t.Errorf("Wrong number of regular expressions, found %v", len(p.matchRegexp))
	}
}

func TestReferenceParseDuplicate(t *testing.T) {
	test := func(code string) {
		defer func() {
			if recover() == nil {
				t.Fail()
			}
		}()
		newQueryParser(myDummyLang, []byte(code))
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
	p := newQueryParser(myDummyLang, []byte(file))

	testQuery := func(in string, r ...Reference) {
		t.Logf("Testing string \"%v\" for match %v", in, r)
		p := p.lookup(in)

		testReferences(t, p, r...)
	}

	test := func(in, out string, verses ...int) {
		testQuery(in, Reference{
			Path:              out,
			VersesHighlighted: verses,
		})
	}

	test("1ne 3", "/scriptures/bofm/1-ne/3")
	test("1ne 3:4", "/scriptures/bofm/1-ne/3", 4)
	test("1ne 3:4-5", "/scriptures/bofm/1-ne/3", 4, 5)
	test("1ne 3:4-6", "/scriptures/bofm/1-ne/3", 4, 5, 6)
	test("1ne 3:4,6", "/scriptures/bofm/1-ne/3", 4, 6)
	test("1ne 3:4-6,6", "/scriptures/bofm/1-ne/3", 4, 5, 6)
	test("1ne 3:4-6,6-8, 2", "/scriptures/bofm/1-ne/3", 2, 4, 5, 6, 7, 8)
	test("1 nephi 3:4-6,6-8, 2", "/scriptures/bofm/1-ne/3", 2, 4, 5, 6, 7, 8)
	testQuery("1ne 3:4 (2-6)", Reference{
		Path:              "/scriptures/bofm/1-ne/3",
		VersesHighlighted: []int{4},
		VersesExtra:       []int{2, 3, 4, 5, 6},
	})
	testQuery("1ne 3:4 (2-6); 4:5", Reference{
		Path:              "/scriptures/bofm/1-ne/3",
		VersesHighlighted: []int{4},
		VersesExtra:       []int{2, 3, 4, 5, 6},
	}, Reference{
		Path:              "/scriptures/bofm/1-ne/4",
		VersesHighlighted: []int{5},
	})

	test("Joseph Smith—History", "/scriptures/pgp/js-h/1")
}
