package ldslib

import (
	"strconv"
	"fmt"
	"net/url"
	"strings"
)

type Reference struct {
	bookName          string
	glPath            string
	item              CatalogItem
	chapter           int
	verseSelected     int
	versesHighlighted []int
}

func (ref Reference) String() string {
	return fmt.Sprintf("%v %v:%v", ref.bookName, ref.chapter, ref.verseSelected)
}

type RefParser struct {
	lib *Library
	cat *Catalog
	s    []string
	ref  Reference
}

func NewRefParser(lib *Library, cat *Catalog) RefParser {
	p := RefParser{}
	p.lib = lib
	p.cat = cat
	return p
}

func (p *RefParser) clean() {
	for i, s := range p.s {
		p.s[i] = strings.ToLower(strings.TrimSpace(s))
		if p.s[i] == "" {
			p.s = append(p.s[:i], p.s[i+1:]...)
		}
	}
}

func (p *RefParser) Load(s string) {
	p.s = strings.Split(s, " ")
	p.clean()
	p.parse()
}

func (p *RefParser) LoadSlice(s []string) {
	p.s = s
	p.clean()
	p.parse()
}

func (p *RefParser) LoadGlURL(glUrl string) {
	//TODO There is more parsing needed for verses
	// eg: /scriptures/dc-testament/dc/2.1-3
	p.ref.glPath = glUrl
}

func (p *RefParser) parse() {
	ref := Reference{}
	s := p.s
	
	if _, err := strconv.Atoi(s[0]); err == nil {
		// Nodes that start with a number (we don't check if the number is good)
		switch s[1] {
		case "nephi", "ne":
			ref.glPath = fmt.Sprintf("/scriptures/bofm/%v-ne", s[0])
		case "samuel", "sam":
			ref.glPath = fmt.Sprintf("/scriptures/ot/%v-sam", s[0])
		case "kings", "king", "kgs":
			ref.glPath = fmt.Sprintf("/scriptures/ot/%v-kgs", s[0])
		case "chronicles", "chron", "chr":
			ref.glPath = fmt.Sprintf("/scriptures/ot/%v-chr", s[0])
		case "corinthians", "corinth", "cor":
			ref.glPath = fmt.Sprintf("/scriptures/nt/%v-cor", s[0])
		case "thessalonians", "thes":
			ref.glPath = fmt.Sprintf("/scriptures/nt/%v-thes", s[0])
		case "timothy", "tim":
			ref.glPath = fmt.Sprintf("/scriptures/nt/%v-tim", s[0])
		case "peter", "pet":
			ref.glPath = fmt.Sprintf("/scriptures/nt/%v-pet", s[0])
		case "john", "jn":
			ref.glPath = fmt.Sprintf("/scriptures/nt/%v-jn", s[0])
		default:
			goto Done
		}
		s = s[2:]
	} else {
		w := s[0]
		s = s[1:]
		switch w {
		case "ot":
			ref.glPath =  "/scriptures/ot"
			goto Done
		case "nt":
			ref.glPath =  "/scriptures/nt"
			goto Done
		case "bom", "bofm":
			ref.glPath =  "/scriptures/bofm"
			goto Done
		case "dc", "d&c":
			if len(s) == 1 {
				ref.glPath = "/scriptures/dc-testament"
				goto Done
			} else {
				ref.glPath = "/scriptures/dc-testament/dc"
			}
		case "pgp":
			ref.glPath =  "/scriptures/pgp"
			goto Done
		case "1nephi", "1ne":
			ref.glPath = "/scriptures/bofm/1-ne"
		default:
			goto Done
		}
	}
	
	if len(s) > 0 {
		if chapter, err := strconv.Atoi(s[0]); err == nil {
			ref.chapter = chapter
			ref.glPath = fmt.Sprintf("%v/%v", ref.glPath, chapter)
			if len(s) > 1 {
				if verse, err := strconv.Atoi(s[1]); err == nil {
					ref.verseSelected = verse
					ref.versesHighlighted = []int{verse}
				}
			}
		} else {
			//TODO May have : seperator
		}
	}
	
Done:
	p.ref = ref
	p.s = s
	fmt.Println(ref.glPath, " ", ref)
}

func (p *RefParser) Book() (*Book, error) {
	return nil, nil
}

func (p *RefParser) Item() (CatalogItem, error) {
	return p.lib.lookupGlURI(p.ref.glPath, p.cat)
}

func (p *RefParser) Reference() (Reference, error) {
	item, err := p.lib.lookupGlURI(p.ref.glPath, p.cat)
	if err != nil {
		return Reference{}, err
	}
	ref := p.ref
	ref.item = item
	if ref.bookName == "" {
		ref.bookName = item.DisplayName()
	}
	return ref, nil
}

func (p *RefParser) glURI(s string) (ref Reference, err error) {
	uri, err := url.Parse(s)
	if err != nil {
		return ref, err
	}
	item, err := p.lib.lookupGlURI(s, p.cat)
	if err != nil {
		return ref, err
	}
	ref.item = item
	ref.verseSelected, err = strconv.Atoi(uri.Fragment)
	if err != nil {
		ref.verseSelected = 0
	}
	
	//TODO Parse Verses and other stuff to work with LDS.org URLs
	
	return
}