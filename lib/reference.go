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
	if strings.IndexRune(q, '/') == 0 {
		return q, nil
	}
	base, _, err := p.lookupBase(q)
	if err != nil {
		return "", err
	}
	//TODO: Parse remainder for chapter, verse, etc
	return base, nil
}

func (p *refParser) lookupBase(Q string) (string, string, error) {
	q := strings.ToLower(Q) + " "
	for s, r := range p.matchString {
		if strings.Index(q, s) == 0 {
			return r, q[len(s):], nil
		}
	}
	for s, r := range p.matchRegexp {
		if i := s.FindSubmatchIndex([]byte(q)); i != nil {
			b := make([]byte, 0)
			b = s.ExpandString(b, r, q, i)
			e := s.ReplaceAllString(q, "")
			return string(b), e, nil
		}
	}
	return "", "", errors.New("Query \"" + Q + "\" not found")
}
