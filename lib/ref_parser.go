package lib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/net/html"
)

type ReferenceParser struct {
	descMap map[Lang]*langRefDesc
}
type langRefDesc struct {
	matchString map[string]string
	matchRegexp map[*regexp.Regexp]string
	matchFolder map[int]string
}

var parseClean = regexp.MustCompile("( |:)+")

func NewReferenceParser() *ReferenceParser {
	return &ReferenceParser{
		descMap: make(map[Lang]*langRefDesc),
	}
}

func (p ReferenceParser) AppendFile(lang Lang, file io.Reader) {
	d := &langRefDesc{
		matchFolder: make(map[int]string),
		matchString: make(map[string]string),
		matchRegexp: make(map[*regexp.Regexp]string),
	}
	s := bufio.NewScanner(file)
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
			if v, ok := d.matchFolder[id]; ok {
				panic(fmt.Errorf("Token %v already used for %v", id, v))
			}
			d.matchFolder[id] = path
			tokens = tokens[1:]
		} else if len(tokens) == 1 && isRegex.MatchString(tokens[0]) {
			exp := "(?i)^" + tokens[0][1:len(tokens[0])-1]
			r, err := regexp.Compile(exp)
			if err == nil {
				r.Longest()
				d.matchRegexp[r] = path
				continue
			}
		}
		for _, t := range tokens {
			t = strings.ToLower(t) + " "
			if v, ok := d.matchString[t]; ok {
				panic(fmt.Errorf("Token %v already used for %v", t, v))
			}
			d.matchString[t] = path
		}
	}
	if _, ok := p.descMap[lang]; ok {
		panic("merge not implemented")
	}
	p.descMap[lang] = d
}

func (p ReferenceParser) Parse(lang Lang, query string) ([]Reference, error) {
	type queryTokenType int
	const (
		tokenRef queryTokenType = iota
		tokenWord
		tokenChar
	)
	type queryParseMode int
	const (
		parseModeRef queryParseMode = iota
		parseModeChapter
		parseModeVerse
		parseModeVerseRange
	)

	desc := p.descMap[lang]
	if desc == nil {
		return nil, errors.New("language not found")
	}

	q := strings.ToLower(query)
	refs := make([]Reference, 0)
	ref := Reference{}
	ref.Lang = lang
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
		adv, path := desc.lookup(string(data[start:]))
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
	var queryWords []string

	finishReference := func() {
		path := ref.Path
		if chapter > 0 {
			ref.Path = fmt.Sprintf("%v/%v", ref.Path, chapter)
		}
		ref.Clean()
		ref.Query = strings.Join(queryWords, " ")
		refs = append(refs, ref)

		ref.Path = path
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
					var v *Verses
					if !inParenths {
						v = &ref.VersesHighlighted
					} else {
						v = &ref.VersesExtra
					}
					verseRangeStart = num
					if *v == nil {
						*v = Verses{num}
					} else {
						*v = append(*v, num)
					}
				case parseModeVerseRange:
					var v *Verses
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
				queryWords = append(queryWords, text)
			}
		}
	}

	finishReference()

	return refs, nil
}

func (p *ReferenceParser) PathFromID(id int) string {
	desc, ok := p.descMap[DefaultLang]
	if !ok {
		return ""
	}
	return desc.matchFolder[id]
}

func (p *ReferenceParser) ParsePath(lang Lang, path string) Reference {
	r := Reference{}
	r.Lang = lang
	r.Path = path

	if index := strings.IndexRune(r.Path, '.'); index != -1 {
		verseString := r.Path[index+1:]
		r.Path = r.Path[:index]

		if extraIndex := strings.IndexRune(verseString, '.'); extraIndex != -1 {
			r.VersesHighlighted = parseVerses(verseString[:extraIndex])
			r.VersesExtra = parseVerses(verseString[extraIndex+1:])
		} else {
			r.VersesHighlighted = parseVerses(verseString)
		}
	}

	r.Clean()

	return r
}

func parseVerses(s string) []int {
	verses := make([]int, 0)
	for _, span := range strings.Split(s, ",") {
		if verse, err := strconv.Atoi(span); err == nil {
			verses = append(verses, verse)
		} else {
			p := strings.Split(span, "-")
			if len(p) == 2 {
				vstart, estart := strconv.Atoi(p[0])
				vend, eend := strconv.Atoi(p[1])
				if estart == nil && eend == nil {
					for v := vstart; v <= vend; v++ {
						verses = append(verses, v)
					}
				}
			}
		}
	}
	return verses
}

func (p *ReferenceParser) ParseFootnote(lang Lang, f Footnote) []Reference {
	z := html.NewTokenizerFragment(strings.NewReader(string(f.Content)), "div")
	refs := make([]Reference, 0)

loop:
	for {
		ref := Reference{}
		ref.Lang = lang

		switch z.Next() {
		case html.ErrorToken, html.EndTagToken:
			break loop
		case html.TextToken:
			ref.Name = string(z.Text())
		case html.SelfClosingTagToken:

		case html.StartTagToken:
			tag, hasAttr := z.TagName()
			depth := 1

			switch string(tag) {
			case "a":
				for hasAttr {
					var key, val []byte
					key, val, hasAttr = z.TagAttr()
					switch string(key) {
					case "href":
						r := p.ParsePath(lang, string(val))
						ref.Path = r.Path
						ref.VerseSelected = r.VerseSelected
						ref.VersesHighlighted = r.VersesHighlighted
						ref.VersesExtra = r.VersesExtra
					}
				}
			case "span":
				for hasAttr {
					var key, val []byte
					key, val, hasAttr = z.TagAttr()
					switch string(key) {
					case "class":
						if string(val) == "small" {
							ref.Small = f.parseSmall(z, tag)
							depth--
						}
					}
				}
			}

			for depth > 0 {
				switch z.Next() {
				case html.ErrorToken:
					break loop
				case html.TextToken:
					ref.Name = fmt.Sprintf("%v%v", ref.Name, string(z.Text()))
				case html.StartTagToken:
					if startTag, _ := z.TagName(); bytes.Equal(startTag, tag) {
						depth++
					} else if "small" == string(startTag) {
						ref.Small = f.parseSmall(z, startTag)
					}
				case html.EndTagToken:
					endTag, _ := z.TagName()
					if bytes.Equal(endTag, tag) {
						depth--
					}
				}
			}
		}

		refs = append(refs, ref)
	}

	cleanRefs := []Reference{}
	oldRef := refs[0]
	oldRef.Name = ""
	for _, ref := range refs {
		if oldRef.Path == ref.Path && oldRef.VerseSelected == ref.VerseSelected && oldRef.Small == ref.Small {
			oldRef.Name += ref.Name
		} else {
			cleanRefs = append(cleanRefs, oldRef)
			oldRef = ref
		}
	}
	cleanRefs = append(cleanRefs, oldRef)
	return cleanRefs
}

func (d *langRefDesc) lookup(q string) (advance int, path string) {
	for s, r := range d.matchString {
		slen := len(s)
		if advance < slen && strings.Index(string(q), s) == 0 {
			path = r
			advance = slen
		}
	}
	for s, r := range d.matchRegexp {
		if i := s.FindStringSubmatchIndex(q); i != nil {
			remTemp := s.ReplaceAllString(q, "")
			adv := len(q) - len(remTemp)
			if advance < adv {
				b := []byte{}
				path = string(s.ExpandString(b, r, q, i))
				advance = adv
			}
		}
	}
	return
}
