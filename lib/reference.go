package lib

import (
	"fmt"
	"sort"
	"strings"
)

type Reference struct {
	Index

	VerseSelected     int
	VersesHighlighted Verses
	VersesExtra       Verses

	// Query is the search query. If an empty string, indicates no search.
	Query string

	// Small and Name are used for references to display, meaningless for lookup.
	Small, Name string
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

func (r Reference) String() string {
	s := r.Name
	if len(s) == 0 {
		s = r.Path
	}
	if r.VersesHighlighted != nil {
		s = fmt.Sprintf("%v:%v", s, r.VersesHighlighted)
	}
	if r.VersesExtra != nil {
		s = fmt.Sprintf("%v (%v)", s, r.VersesHighlighted)
	}
	return s
}

func (r Reference) URL() string {
	p := r.Path
	if r.VersesHighlighted != nil {
		p = fmt.Sprintf("%v.%v", p, r.VersesHighlighted)
	}
	if r.VersesExtra != nil {
		p = fmt.Sprintf("%v.%v", p, r.VersesExtra)
	}
	if r.Lang != "" {
		p = fmt.Sprintf("%v?lang=%v", p, r.Lang)
	}
	if r.VerseSelected > 0 {
		p = fmt.Sprintf("%v#%v", p, r.VerseSelected)
	}
	return p
}

type Verses []int

func (verses Verses) String() string {
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

func (a Verses) Clean() Verses {
	sort.Sort(sort.IntSlice([]int(a)))

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
