package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ReferenceFileFunc func(lang *Language) ([]byte, error)

var refParserFileLoader ReferenceFileFunc
var refParsers map[*Language]*refParser

func referenceParser(l *Language) (*refParser, error) {
	if refParsers == nil {
		refParsers = make(map[*Language]*refParser)
	}
	if parser, ok := refParsers[l]; ok {
		return parser, nil
	}
	file, err := refParserFileLoader(l)
	if err != nil {
		return nil, err
	}
	return newRefParser(file), nil
}

// Sets a function that will be called to get the ldss reference language file
// for a passed language. This will likely be from a file, but could be from
// another source, such as an embedded resource. This should not be set after
// the first call to Lookup()
func SetReferenceParseReader(open ReferenceFileFunc) {
	refParserFileLoader = open
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

func (p *refParser) lookup(q string) []Reference {
	//refs := make([]Reference, 0)
	ref := Reference{}
	var err error

	// Clean q
	q = strings.TrimSpace(q)
	q = strings.ToLower(q) + " "
	q = p.parseClean.ReplaceAllString(q, " ")

	// Parse from the match maps
	var remainder string
	ref.Path, remainder = p.lookupBase(q)

	if i := strings.LastIndex(ref.Path, "#"); i != -1 {
		directive := string(ref.Path[i:])
		ref.Path = string(ref.Path[:i])
		if len(remainder) == 0 {
			return []Reference{ref}
		}
		// Parse remainder for chapter, verse, etc
		tokens := strings.Split(strings.TrimRight(remainder, " "), " ")
		switch directive {
		case "#":
			var i int
			i, err = strconv.Atoi(tokens[0])
			if err != nil {
				return []Reference{ref}
			} else {
				ref.Path = fmt.Sprintf("%v/%v", ref.Path, i)
			}
		default:
			err = fmt.Errorf("Unknown directive %v", directive)
		}
	} else {
		if len(remainder) == 0 {
			return []Reference{ref}
		}
		err = fmt.Errorf("Unknown extra characters %v", remainder)
	}
	return []Reference{ref}
}

func (p *refParser) lookupBase(q string) (path, rem string) {
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
	return
}
