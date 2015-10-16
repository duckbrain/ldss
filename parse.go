package main

import (
	"fmt"
	"strings"
)

type Reference struct {
	bookName          string
	glPath            string
	node              int
	chapter           int
	verseSelected     int
	versesHighlighted []int
}

func (ref Reference) String() string {
	return string(ref.bookName) + " " + string(ref.chapter) + ":" + string(ref.verseSelected)
}

func ParsePath(path string) *Reference {
	panic("Not Implemented")
}

func LookupPath(path string) {
	ref := ParsePath(path)
	fmt.Println(ref.String())
}

func ParseForBook(id string) string {
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
