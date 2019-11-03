package ref

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/duckbrain/ldss/lib"
)

func Parse(lang lib.Lang, q string) []lib.Reference {
	ref := ParsePath(lang, q)
	if ref.Check() == nil {
		return []lib.Reference{ref}
	}
	if rp, err := languageQueryParser(lang); err == nil {
		return rp.lookup(q)
	}
	return []lib.Reference{}
}

func ParsePath(lang lib.Lang, p string) Reference {
	r := Reference{
		Lang: lang,
		Path: p,
	}

	if index := strings.IndexRune(r.Path, '.'); index != -1 {
		verseString := r.Path[index+1:]
		r.Path = r.Path[:index]

		if extraIndex := strings.IndexRune(verseString, '.'); extraIndex != -1 {
			r.VersesHighlighted = parseVerses(verseString[:extraIndex])
			r.VersesExtra = parseVerses(verseString[extraIndex+1:])
		} else {
			r.VersesHighlighted = parseVerses(verseString)
		}
	}

	r.Clean()

	return r
}

func (r Reference) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Path              string
		Lang              lib.Lang
		VerseSelected     int
		VersesHighlighted []int
		VersesExtra       []int
		Small, Content    string
		URL               string
	}{
		Path:              r.Path,
		Lang:              r.Lang,
		VerseSelected:     r.VerseSelected,
		VersesHighlighted: r.VersesHighlighted,
		VersesExtra:       r.VersesExtra,
		Small:             r.Small,
		Content:           r.Name,
		URL:               r.URL(),
	})
}

func (r *Reference) Clean() {
	r.Path = strings.TrimSpace(r.Path)
	r.Path = strings.TrimRight(r.Path, "/ ")
	if r.Path == "" {
		r.Path = "/"
	}

	r.VersesHighlighted = r.VersesHighlighted.Clean()
	r.VersesExtra = r.VersesExtra.Clean()

	if r.VerseSelected == 0 && len(r.VersesHighlighted) > 0 {
		r.VerseSelected = r.VersesHighlighted[0] - 1
	}
}

func parseVerses(s string) []int {
	verses := make([]int, 0)
	for _, span := range strings.Split(s, ",") {
		if verse, err := strconv.Atoi(span); err == nil {
			verses = append(verses, verse)
		} else {
			p := strings.Split(span, "-")
			if len(p) == 2 {
				vstart, estart := strconv.Atoi(p[0])
				vend, eend := strconv.Atoi(p[1])
				if estart == nil && eend == nil {
					for v := vstart; v <= vend; v++ {
						verses = append(verses, v)
					}
				}
			}
		}
	}
	return verses
}

func (r Reference) Check() error {
	if r.Lang == "" {
		return fmt.Errorf("Lang not set on reference")
	}
	if len(r.Path) == 0 || r.Path[0] != '/' {
		return fmt.Errorf("Path \"%v\" must start with '/'", r.Path)
	}

	return nil
}
