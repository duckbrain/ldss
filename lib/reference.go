package lib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Sets a function that will be called to get the ldss reference language file
// for a passed language. This will likely be from a file, but could be from
// another source, such as an embedded resource.
func SetReferenceParseReader(open func(lang *Language) ([]byte, error)) {
	langs, err := Languages()
	if err != nil {
		panic(err)
	}
	for _, lang := range langs {
		func(l *Language) {
			l.reference.construct = func() (interface{}, error) {
				file, err := open(l)
				if err != nil {
					return nil, err
				}
				return newRefParser(file), nil
			}
		}(lang)
	}

}

type reference struct {
	bookName          string
	glPath            string
	item              Item
	chapter           int
	verseSelected     int
	versesHighlighted []int
}

func (ref reference) String() string {
	return fmt.Sprintf("%v %v:%v", ref.bookName, ref.chapter, ref.verseSelected)
}

type refParser struct {
	matchString map[string]string
	matchRegexp map[*regexp.Regexp]string
	matchFolder map[int]string
	parseClean  *regexp.Regexp
}

func newRefParser(file []byte) *refParser {
	p := &refParser{
		matchFolder: make(map[int]string),
		matchString: make(map[string]string),
		matchRegexp: make(map[*regexp.Regexp]string),
		parseClean:  regexp.MustCompile("( |:)+"),
	}
	s := bufio.NewScanner(bytes.NewReader(file))
	isRegex := regexp.MustCompile("^\\/.*\\/$")
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 || strings.IndexRune(line, '#') == 0 {
			continue
		}
		tokens := strings.Split(line, ":")
		path := tokens[len(tokens)-1]
		tokens = tokens[:len(tokens)-1]
		if id, err := strconv.Atoi(tokens[0]); err == nil {
			if v, ok := p.matchFolder[id]; ok {
				panic(fmt.Errorf("Token %v already used for %v", id, v))
			}
			p.matchFolder[id] = path
			tokens = tokens[1:]
		} else if len(tokens) == 1 && isRegex.MatchString(tokens[0]) {
			exp := "^" + tokens[0][1:len(tokens[0])-1] + " "
			r, err := regexp.Compile(exp)
			if err == nil {
				r.Longest()
				p.matchRegexp[r] = path
				continue
			}
		}
		for _, t := range tokens {
			t = strings.ToLower(t) + " "
			if v, ok := p.matchString[t]; ok {
				panic(fmt.Errorf("Token %v already used for %v", t, v))
			}
			p.matchString[t] = path
		}
	}
	return p
}

func (p *refParser) lookup(q string) (string, error) {
	// Clean q
	q = strings.TrimSpace(q)
	q = strings.ToLower(q) + " "

	if strings.IndexRune(q, '/') == 0 {
		return q, nil
	}

	q = p.parseClean.ReplaceAllString(q, " ")

	// Parse from the match maps
	base, remainder, err := p.lookupBase(q)
	if err != nil {
		return "", err
	}

	if i := strings.LastIndex(base, "#"); i != -1 {
		directive := string(base[i:])
		base = string(base[:i])
		if len(remainder) == 0 {
			return base, nil
		}
		// Parse remainder for chapter, verse, etc
		tokens := strings.Split(strings.TrimRight(remainder, " "), " ")
		switch directive {
		case "#":
			i, err := strconv.Atoi(tokens[0])
			if err != nil {
				return "", err
			}
			base = fmt.Sprintf("%v/%v", base, i)
		default:
			return "", fmt.Errorf("Unknown directive %v", directive)
		}
		return base, nil
	} else {
		if len(remainder) == 0 {
			return base, nil
		}
		//TODO Handle remainder
	}

	return base, nil
}

func (p *refParser) lookupBase(q string) (path, rem string, err error) {
	for s, r := range p.matchString {
		if strings.Index(q, s) == 0 && (len(rem) == 0 || len(rem) > len(q)-len(s)) {
			path = r
			rem = q[len(s):]
		}
	}
	for s, r := range p.matchRegexp {
		if i := s.FindSubmatchIndex([]byte(q)); i != nil {
			remTemp := s.ReplaceAllString(q, "")
			if len(rem) == 0 || len(rem) > len(remTemp) {
				rem = remTemp
				b := []byte{}
				path = string(s.ExpandString(b, r, q, i))
			}
		}
	}
	if path == "" {
		err = errors.New("Query \"" + q + "\" not found")
	}
	return
}
