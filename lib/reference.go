package lib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Reference struct {
	Path              string
	Language          *Language
	VerseSelected     int
	VersesHighlighted []int
	VersesExtra       []int
	Small, Content    string
}

func Parse(lang *Language, q string) (r Reference, err error) {
	var rp *refParser
	r = ParsePath(lang, q)
	if r.Check() == nil {
		return r, nil
	}
	rp, err = lang.ref()
	if err == nil {
		r, err = rp.lookup(q)
		r.Language = lang
		return
	}
	return
}

func ParsePath(lang *Language, p string) Reference {
	r := Reference{
		Language: lang,
		Path:     p,
	}

	r.Path = strings.TrimSpace(r.Path)
	r.Path = strings.TrimRight(r.Path, "/ ")
	if r.Path == "" {
		r.Path = "/"
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

		if len(r.VersesHighlighted) > 0 {
			r.VerseSelected = r.VersesHighlighted[0]
		}
	}

	return r
}

func (r Reference) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Path              string
		Language          *Language
		VerseSelected     int
		VersesHighlighted []int
		VersesExtra       []int
		Small, Content    string
		URL               string
	}{
		Path:              r.Path,
		Language:          r.Language,
		VerseSelected:     r.VerseSelected,
		VersesHighlighted: r.VersesHighlighted,
		VersesExtra:       r.VersesExtra,
		Small:             r.Small,
		Content:           r.Content,
		URL:               r.URL(),
	})
}

func (r Reference) URL() string {
	p := r.Path
	p = fmt.Sprintf("%v%v%v", p, stringifyVerse(r.VersesHighlighted), stringifyVerse(r.VersesExtra))
	if r.Language != nil {
		p = fmt.Sprintf("%v?lang=%v", p, r.Language.GlCode)
	}
	if r.VerseSelected > 0 {
		p = fmt.Sprintf("%v#%v", p, r.VerseSelected)
	}
	return p
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

func stringifyVerse(verses []int) string {
	if verses == nil {
		return ""
	}
	p := ""
	var previousVerse, spanStart, verse int
	for _, verse = range verses {
		if previousVerse == 0 {
			p = fmt.Sprintf("%v.%v", p, verse)
			spanStart = verse
			previousVerse = verse
		} else if previousVerse == verse-1 {
			previousVerse = verse
		} else {
			p = fmt.Sprintf("%v-%v,%v", p, previousVerse, verse)
			spanStart = verse
			previousVerse = verse
		}
	}
	if verse != spanStart {
		p = fmt.Sprintf("%v-%v", p, verse)
	}
	return p
}

func (r Reference) Check() error {
	if r.Language == nil {
		return fmt.Errorf("Language not set on reference")
	}
	if len(r.Path) == 0 || r.Path[0] != '/' {
		return fmt.Errorf("Path \"%v\" must start with '/'", r.Path)
	}

	return nil
}

// Finds an Item by it's path. Expects a fully qualified path. "/" will
// return the catalog. Will return an error if there is an error
// loading the item or it is not downloaded.
func (r Reference) Lookup() (Item, error) {
	if err := r.Check(); err != nil {
		return nil, err
	}

	c, err := r.Language.Catalog()
	if err != nil {
		return nil, err
	}
	if r.Path == "/" {
		return c, nil
	}
	if folder, ok := c.foldersByPath[r.Path]; ok {
		return folder, nil
	}
	sections := strings.Split(r.Path, "/")
	if sections[0] != "" {
		return nil, fmt.Errorf("Invalid path \"%v\", must start with '/'", r.Path)
	}
	for i := 2; i <= len(sections); i++ {
		temppath := strings.Join(sections[0:i], "/")
		if book, ok := c.booksByPath[temppath]; ok {
			if r.Path == book.Path() {
				return book, nil
			}
			node := &Node{Book: book}
			db, err := book.db()
			if err != nil {
				return nil, err
			}
			err = db.stmtUri.QueryRow(r.Path).Scan(&node.id, &node.name, &node.path, &node.parentId, &node.hasContent, &node.childCount)
			if err != nil {
				return nil, fmt.Errorf("Path %v not found", r.Path)
			}
			return node, err

			return book.lookupPath(r.Path)
		}
	}
	return nil, fmt.Errorf("Path \"%v\" not found", r.Path)
}
