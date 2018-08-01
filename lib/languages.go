package lib

import (
	"fmt"
)

var langs map[string]langMap

func init() {
	langs = make(map[string]langMap)
}

type Lang interface {
	Name() string
	EnglishName() string
	Code() string
	Matches(string) bool
}

type langMap map[string]Lang

func (m langMap) Name() string {
	for _, lang := range m {
		if name := lang.Name(); len(name) > 0 {
			return name
		}
	}
	return ""
}

func (m langMap) EnglishName() string {
	for _, lang := range m {
		if name := lang.EnglishName(); len(name) > 0 {
			return name
		}
	}
	return ""
}

func (m langMap) Code() string {
	for _, lang := range m {
		if code := lang.Code(); len(code) > 0 {
			return code
		}
	}
	return ""
}

func (m langMap) Matches(s string) bool {
	for _, lang := range m {
		if lang.Matches(s) {
			return true
		}
	}
	return false
}

func Languages() []Lang {
	if !opened {
		panic(fmt.Errorf("You must call lib.Open() before getting languages"))
	}
	res := make([]Lang, len(langs))
	i := 0
	for _, l := range langs {
		res[i] = l
		i++
	}
	return res
}

func languageFromSource(lang Lang, srcName string) Lang {
	return langs[lang.Code()][srcName]
}

func registerLanguage(srcName string, srcLangs []Lang) {
	for _, srcLang := range srcLangs {
		if lang, ok := langs[srcLang.Code()]; ok {
			lang[srcName] = srcLang
		} else {
			langs[srcLang.Code()] = langMap{srcName: srcLang}
		}
	}
}

// LookupLanguage finds a language by any of the accepted methods, compares ID, Code, and InternalCode
func LookupLanguage(id string) Lang {
	if !opened {
		panic(fmt.Errorf("You must call lib.Open() before looking up a language"))
	}
	for _, lang := range langs {
		if lang.Matches(id) {
			return lang
		}
	}
	return nil
}
