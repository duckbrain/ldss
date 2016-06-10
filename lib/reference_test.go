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
	test := func(in, out string) {
		if p, err := p.lookup(in); err != nil || p != out {
			t.Errorf("Expected \"%v\"->\"%v\" received \"%v\"", in, out, p)
		}
	}

	test("1ne 3", "/scriptures/bofm/1-ne/3")
	test("1ne 3:4", "/scriptures/bofm/1-ne/3.4")
}
