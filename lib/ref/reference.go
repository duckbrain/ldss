package ref

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Reference struct {
	Path              string
	Lang              Lang
	VerseSelected     int
	VersesHighlighted []int
	VersesExtra       []int
	Small, Name       string
	Keywords          []string
}

func Parse(lang Lang, q string) []Reference {
	ref := ParsePath(lang, q)
	if ref.Check() == nil {
		return []Reference{ref}
	}
	if rp, err := languageQueryParser(lang); err == nil {
		return rp.lookup(q)
	}
	return []Reference{}
}

func ParsePath(lang Lang, p string) Reference {
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
		Lang              Lang
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

	r.VersesHighlighted = cleanVerses(r.VersesHighlighted)
	r.VersesExtra = cleanVerses(r.VersesExtra)

	if r.VerseSelected == 0 && len(r.VersesHighlighted) > 0 {
		r.VerseSelected = r.VersesHighlighted[0] - 1
	}
}

func (r Reference) URL() string {
	p := r.Path
	if r.VersesHighlighted != nil {
		p = fmt.Sprintf("%v.%v", p, stringifyVerses(r.VersesHighlighted))
	}
	if r.VersesExtra != nil {
		p = fmt.Sprintf("%v.%v", p, stringifyVerses(r.VersesExtra))
	}
	if r.Lang != nil {
		p = fmt.Sprintf("%v?lang=%v", p, r.Lang.Code())
	}
	if r.VerseSelected > 0 {
		p = fmt.Sprintf("%v#%v", p, r.VerseSelected)
	}
	return p
}

func (r Reference) String() string {
	s := r.Path
	if item, err := r.Lookup(); err == nil {
		s = item.Name()
	}
	if r.VersesHighlighted != nil {
		s = fmt.Sprintf("%v:%v", s, stringifyVerses(r.VersesHighlighted))
	}
	if r.VersesExtra != nil {
		s = fmt.Sprintf("%v (%v)", s, stringifyVerses(r.VersesHighlighted))
	}
	if len(s) == 0 {
		s = fmt.Sprintf("\"%v\"", r.Name)
	}
	return s
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

func stringifyVerses(verses []int) string {
	if verses == nil {
		return ""
	}
	p := ""
	var previousVerse, spanStart, verse int
	for _, verse = range verses {
		if previousVerse == 0 {
			p = fmt.Sprintf("%v", verse)
			spanStart = verse
		} else if previousVerse == verse-1 {
		} else if previousVerse != spanStart {
			p = fmt.Sprintf("%v-%v,%v", p, previousVerse, verse)
			spanStart = verse
		} else {
			p = fmt.Sprintf("%v,%v", p, verse)
			spanStart = verse
		}
		previousVerse = verse
	}
	if verse != spanStart {
		p = fmt.Sprintf("%v-%v", p, verse)
	}
	return p
}

func cleanVerses(a []int) []int {
	sort.Sort(sort.IntSlice(a))

	l := 0
	for i := 0; i < len(a); i++ {
		v := a[i]
		if v <= l {
			a = append(a[:i], a[i+1:]...)
			i--
		}
		l = v
	}
	return a
}

func (r Reference) Check() error {
	if r.Lang == nil {
		return fmt.Errorf("Lang not set on reference")
	}
	if len(r.Path) == 0 || r.Path[0] != '/' {
		return fmt.Errorf("Path \"%v\" must start with '/'", r.Path)
	}

	return nil
}
