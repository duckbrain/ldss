package ldslib

import (
	"strconv"
	"fmt"
	"net/url"
	"strings"
)

type RefParser struct {
	l    Library
	cat *Catalog
}

type Reference struct {
	bookName          string
	glPath            string
	item              CatalogItem
	chapter           int
	verseSelected     int
	versesHighlighted []int
}

func (ref Reference) String() string {
	return string(ref.bookName) + " " + string(ref.chapter) + ":" + string(ref.verseSelected)
}

func (p RefParser) GlURI(path string) (ref Reference, err error) {
	uri, err := url.Parse(path)
	if err != nil {
		return ref, err
	}
	item, err := p.l.lookupGlURI(path, p.cat)
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

func (p *RefParser) clean(s string) ([]string) {
	s = strings.ToLower(s)
	s = strings.Replace(s, "  ", " ", 0)
	return strings.Split(s, " ")
}

func (p *RefParser) findBook(id string, allowExtra bool) (*Book, error) {
	return nil, nil
}

func (p *RefParser) Book(id string) (*Book, error) {
	return p.findBook(id, false)
}

func (p *RefParser) Item(id string) (CatalogItem, error) {
	return nil, nil
}

func (p RefParser) ParsePath(path string) *Reference {
	panic("Not Implemented")
}

func (p RefParser) LookupPath(path string) {
	ref := p.ParsePath(path)
	fmt.Println(ref.String())
}

func (p RefParser) BookPath(id string) (string, error) {
	if strings.HasPrefix(id, "/") {
		return id, nil
	}
	id = strings.ToLower(id)
	switch id {
	case "ot":
		return "/scriptures/ot", nil
	case "nt":
		return "/scriptures/nt", nil
	case "bom", "bofm":
		return "/scriptures/bofm", nil
	case "dc", "d&c":
		return "/scriptures/dc-testament", nil
	case "pgp":
		return "/scriptures/pgp", nil
	default:
		return "", fmt.Errorf("Path %v not recognized")
	}
}

func (p RefParser) NodePath(id string, c Library) string {
	if strings.HasPrefix(id, "/") {
		return id
	}
	id = strings.ToLower(id)
	switch id {
	case "1ne", "1nephi":
		return "/scriptures/bofm/1ne"
	}

	//parts := strings.Split(id, " ")
	return id
}
