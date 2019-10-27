package ref

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/duckbrain/ldss/lib"
)

type QueryFileFunc func(lang lib.Lang) ([]byte, error)
type queryTokenType int
type queryParseMode int

func (tt queryTokenType) String() string {
	switch tt {
	case tokenRef:
		return "ref"
	case tokenWord:
		return "word"
	case tokenChar:
		return "char"
	}
	return "Unknown token type"
}

func (m queryParseMode) String() string {
	switch m {
	case parseModeRef:
		return "ref"
	case parseModeChapter:
		return "chapter"
	case parseModeVerse:
		return "verse"
	case parseModeVerseRange:
		return "verseRange"
	}
	return "Unknown parse type"
}

const (
	tokenRef queryTokenType = iota
	tokenWord
	tokenChar
)

const (
	parseModeRef queryParseMode = iota
	parseModeChapter
	parseModeVerse
	parseModeVerseRange
)

var queryFileLoader QueryFileFunc
var queryParsers map[lib.Lang]*queryParser

func languageQueryParser(l lib.Lang) (*queryParser, error) {
	if queryParsers == nil {
		queryParsers = make(map[lib.Lang]*queryParser)
	}
	if parser, ok := queryParsers[l]; ok {
		return parser, nil
	}
	if queryFileLoader == nil {
		panic("You must call SetReferenceParseReader prior to loading the query parser")
	}
	file, err := queryFileLoader(l)
	if err != nil {
		return nil, err
	}
	return newQueryParser(l, file), nil
}

// Sets a function that will be called to get the ldss reference language file
// for a passed language. This will likely be from a file, but could be from
// another source, such as an embedded resource. This should not be set after
// the first call to Lookup()
func SetReferenceParseReader(open QueryFileFunc) {
	queryFileLoader = open
}

type queryParser struct {
	matchString map[string]string
	matchRegexp map[*regexp.Regexp]string
	matchFolder map[int]string
	parseClean  *regexp.Regexp
	lang        lib.Lang
}

func newQueryParser(lang lib.Lang, file []byte) *queryParser {
	p := &queryParser{
		matchFolder: make(map[int]string),
		matchString: make(map[string]string),
		matchRegexp: make(map[*regexp.Regexp]string),
		parseClean:  regexp.MustCompile("( |:)+"),
		lang:        lang,
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
			exp := "(?i)^" + tokens[0][1:len(tokens[0])-1]
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

func (p *queryParser) lookup(q string) []Reference {
	q = strings.ToLower(q)
	refs := make([]Reference, 0)
	ref := Reference{Lang: p.lang}
	var tt queryTokenType

	scanner := bufio.NewScanner(strings.NewReader(q))
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Skip leading spaces.
		start := 0
		for width := 0; start < len(data); start += width {
			var r rune
			r, width = utf8.DecodeRune(data[start:])
			if !unicode.IsSpace(r) {
				break
			}
		}
		// Find a reference
		adv, path := p.lookupBase(data[start:])
		if adv > 0 {
			advance = start + adv
			tt = tokenRef
			ref.Path = path
			return advance, data[start:advance], nil
		}
		// Scan until space or other token, marking end of word.
		for width, i := 0, start; i < len(data); i += width {
			var r rune
			r, width = utf8.DecodeRune(data[i:])
			if unicode.IsSpace(r) {
				tt = tokenWord
				return i + width, data[start:i], nil
			}
			switch r {
			case ':', ',', ';', '-', '–', '(', ')':
				if i == start {
					// This is the first character, return it as a token
					tt = tokenChar
					return width, data[start : start+width], nil
				} else {
					// Return the word before the character
					tt = tokenWord
					return i, data[start:i], nil
				}
			}
		}
		// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
		if atEOF && len(data) > start {
			tt = tokenWord
			return len(data), data[start:], nil
		}
		// Request more data.
		return start, nil, nil
	})

	var parseMode queryParseMode
	var inParenths bool
	var verseRangeStart int
	var chapter int

	finishReference := func() {
		path := ref.Path
		if chapter > 0 {
			ref.Path = fmt.Sprintf("%v/%v", ref.Path, chapter)
		}
		ref.Clean()
		refs = append(refs, ref)

		ref = Reference{
			Lang: p.lang,
			Path: path,
		}
		verseRangeStart = 0
		chapter = 0
		inParenths = false
		parseMode = parseModeChapter
	}

	for scanner.Scan() {
		text := scanner.Text()
		switch tt {
		case tokenChar:
			switch text {
			case ":":
				parseMode = parseModeVerse
			case ",":
				parseMode = parseModeVerse
			case ";":
				finishReference()
			case "-", "–":
				parseMode = parseModeVerseRange
			case "(":
				inParenths = true
			case ")":
				inParenths = false
			}
		case tokenRef:
			if i := strings.LastIndex(ref.Path, "#"); i != -1 {
				pathType := ref.Path[i:]
				switch pathType {
				case "#":
					parseMode = parseModeChapter
				case "#1":
					parseMode = parseModeChapter
					chapter = 1
				}
				ref.Path = ref.Path[:i]
				parseMode = parseModeChapter
			} else {
				parseMode = parseModeVerse
			}
		case tokenWord:
			if num, err := strconv.Atoi(text); err == nil {
				switch parseMode {
				case parseModeChapter:
					chapter = num
					parseMode = parseModeVerse
				case parseModeVerse:
					var v *[]int
					if !inParenths {
						v = &ref.VersesHighlighted
					} else {
						v = &ref.VersesExtra
					}
					verseRangeStart = num
					if *v == nil {
						*v = []int{num}
					} else {
						*v = append(*v, num)
					}
				case parseModeVerseRange:
					var v *[]int
					if !inParenths {
						v = &ref.VersesHighlighted
					} else {
						v = &ref.VersesExtra
					}
					for i := verseRangeStart; i <= num; i++ {
						*v = append(*v, i)
					}
				}
			} else {
				ref.Keywords = append(ref.Keywords, text)
			}
		}
	}

	finishReference()

	return refs
}

func (p *queryParser) lookupBase(q []byte) (advance int, path string) {
	for s, r := range p.matchString {
		slen := len(s)
		if advance < slen && strings.Index(string(q), s) == 0 {
			path = r
			advance = slen
		}
	}
	for s, r := range p.matchRegexp {
		if i := s.FindSubmatchIndex(q); i != nil {
			remTemp := s.ReplaceAll(q, []byte{})
			adv := len(q) - len(remTemp)
			if advance < adv {
				b := []byte{}
				path = string(s.Expand(b, []byte(r), q, i))
				advance = adv
			}
		}
	}
	return
}
