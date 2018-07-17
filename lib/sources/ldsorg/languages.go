package ldsorg

type jsonLangRoot struct {
	Languages []*jsonLang `json:"languages"`
	Success   bool        `json:"success"`
}

// Lang defines a language as from the server. The fields should not be modified.
type jsonLang struct {
	// The Gospel Library ID for the language. Used for downloads.
	ID int `json:"id"`

	// Native representation of the language in the language observed
	JsonName string `json:"name"`

	// English representation of the language
	JsonEnglishName string `json:"eng_name"`

	// The internationalization (i18n) code used in most programs
	JsonCode string `json:"code"`

	// Gospel Library language code, seen in the urls of https://lds.org
	GlCode string `json:"code_three"`
}

func (l jsonLang) Name() string {
	return l.Name
}

func (l jsonLang) EnglishName() string {
	return l.EnglishName
}

func (l jsonLang) Code() string {
	return l.Code
}

func (l jsonLang) Matches(s string) bool {
	s = strings.ToLower(s)
	return s == strings.ToLower(l.Name)
}
