package ldslib

import (
	"fmt"
	"strings"
)

type RefParser struct {
	//r ReaderConnection
}

type Reference struct {
	bookName          string
	glPath            string
	node              CatalogItem
	chapter           int
	verseSelected     int
	versesHighlighted []int
}

func (ref Reference) String() string {
	return string(ref.bookName) + " " + string(ref.chapter) + ":" + string(ref.verseSelected)
}

func (p RefParser) ParsePath(path string) *Reference {
	panic("Not Implemented")
}

func (p RefParser) LookupPath(path string) {
	ref := p.ParsePath(path)
	fmt.Println(ref.String())
}

func ParseForBook(id string) string {
	if strings.HasPrefix(id, "/") {
		return id
	}
	id = strings.ToLower(id)
	switch id {
	case "ot":
		return "/scriptures/ot"
	case "nt":
		return "/scriptures/nt"
	case "bom", "bofm":
		return "/scriptures/bofm"
	case "dc", "d&c":
		return "/scriptures/dc-testament"
	case "pgp":
		return "/scriptures/pgp"
	default:
		return ""
	}
}

func ParseForNode(id string, c Library) string {
	if strings.HasPrefix(id, "/") {
		return id
	}
	id = strings.ToLower(id)

	//parts := strings.Split(id, " ")
	return id
}
