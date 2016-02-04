package lib

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Reference struct {
	bookName          string
	glPath            string
	item              Item
	chapter           int
	verseSelected     int
	versesHighlighted []int
}

func (ref Reference) String() string {
	return fmt.Sprintf("%v %v:%v", ref.bookName, ref.chapter, ref.verseSelected)
}

type refParser struct {
	matchString map[string]string
	matchRegexp map[*regexp.Regexp]string
	matchFolder map[int]string
}

func (p *refParser) lookup(q string) (string, error) {
	q = strings.ToLower(q)
	if strings.IndexRune(q, '/') == 0 {
		return q, nil
	}
	for s, r := range p.matchString {
		if strings.Index(q, s) == 0 {
			return r, nil
		}
	}
	for s, r := range p.matchRegexp {
		if index := s.FindSubmatchIndex([]byte(q)); index != nil {
			b := make([]byte, 0)
			b = s.ExpandString(b, r, q, index)
			return string(b), nil
		}
	}
	return "", errors.New("Query \"" + q + "\" not found")
}
