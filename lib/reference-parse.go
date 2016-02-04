package lib

import (
	"bufio"
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

func newRefParser(file []byte) *refParser {
	p := &refParser{
		matchFolder: make(map[int]string),
		matchString: make(map[string]string),
		matchRegexp: make(map[*regexp.Regexp]string),
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
			p.matchFolder[id] = path
			tokens = tokens[1:]
		} else if len(tokens) == 1 && isRegex.MatchString(tokens[0]) {
			r, err := regexp.Compile(tokens[0][1 : len(tokens[0])-2])
			if err == nil {
				p.matchRegexp[r] = path
				continue
			}
		}
		for _, t := range tokens {
			p.matchString[strings.ToLower(t)] = path
		}
	}
	return p
}

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
