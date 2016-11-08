package lib

import (
	"fmt"
	"strings"
)

type Reference struct {
	GlPath                  string
	Language                *Language
	VerseSelected           int
	VersesHighlighted       []int
	Name, LinkName, Content string
}

func Parse(lang *Language, q string) (r Reference, err error) {
	var ref *refParser
	r, err = ParsePath(lang, q)
	if err == nil {
		return
	}
	ref, err = lang.ref()
	if err == nil {
		r, err = ref.lookup(q)
		r.Language = lang
		return
	}
	return
}

func ParsePath(lang *Language, p string) (Reference, error) {
	return Reference{
		Language: lang,
		GlPath:   p,
	}, nil
}

func (r Reference) URL() string {
	p := r.GlPath
	if r.VersesHighlighted != nil {
		//TODO Add verses
		var previousVerse, spanStart, verse int
		for _, verse = range r.VersesHighlighted {
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
	}
	if r.Language != nil {
		p = fmt.Sprintf("%v?lang=%v", p, r.Language.GlCode)
	}
	if r.VerseSelected > 0 {
		p = fmt.Sprintf("%v#%v", p, r.VerseSelected)
	}
	return p
}

func (r Reference) Check() error {
	if r.Language == nil {
		panic(fmt.Errorf("Language not set on reference"))
	}
	return nil
}

// Finds an Item by it's path. Expects a fully qualified path. An empty string
// or "/" will return this catalog. Will return an error if there is an error
// loading the item or it is not downloaded.
func (r Reference) Lookup() (Item, error) {
	c, err := r.Language.Catalog()
	if err != nil {
		return nil, err
	}
	r.GlPath = strings.TrimSpace(r.GlPath)
	if r.GlPath == "" || r.GlPath == "/" {
		return c, nil
	}
	r.GlPath = strings.TrimRight(r.GlPath, "/ ")
	if folder, ok := c.foldersByPath[r.GlPath]; ok {
		return folder, nil
	}
	sections := strings.Split(r.GlPath, "/")
	if sections[0] != "" {
		return nil, fmt.Errorf("Invalid path \"%v\", must start with '/'", r.GlPath)
	}
	for i := 2; i <= len(sections); i++ {
		temppath := strings.Join(sections[0:i], "/")
		if book, ok := c.booksByPath[temppath]; ok {
			if r.GlPath == book.Path() {
				return book, nil
			}
			node := &Node{Book: book}
			db, err := book.db()
			if err != nil {
				return nil, err
			}
			err = db.stmtUri.QueryRow(r.GlPath).Scan(&node.id, &node.name, &node.path, &node.parentId, &node.hasContent, &node.childCount)
			if err != nil {
				return nil, fmt.Errorf("Path %v not found", r.GlPath)
			}
			return node, err

			return book.lookupPath(r.GlPath)
		}
	}
	return nil, fmt.Errorf("Path \"%v\" not found", r.GlPath)
}
