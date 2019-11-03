package lib

import (
	"bytes"
	"html/template"

	"golang.org/x/net/html"
)

type Footnote struct {
	Name     string        `json:"name"`
	LinkName string        `json:"linkName"`
	Content  template.HTML `json:"content"`
}

// Parses the rest of the small tag, assuming the head has already been parsed
func (f *Footnote) parseSmall(z *html.Tokenizer, tag []byte) (small string) {
	depth := 1
	for depth > 0 {
		switch z.Next() {
		case html.ErrorToken:
			return ""
		case html.TextToken:
			small = string(z.Text())
		case html.StartTagToken:
			if startTag, _ := z.TagName(); bytes.Equal(startTag, tag) {
				depth++
			}
		case html.EndTagToken:
			if endTag, _ := z.TagName(); bytes.Equal(endTag, tag) {
				depth--
			}
		}
	}
	return
}

func continueToSmall(f *Footnote, z *html.Tokenizer, tag []byte) (small string) {
	depth := 1
	for depth > 0 {
		switch z.Next() {
		case html.ErrorToken:
			return ""
		case html.TextToken:
			small = string(z.Text())
		case html.StartTagToken:
			if startTag, _ := z.TagName(); bytes.Equal(startTag, tag) {
				depth++
			}
		case html.EndTagToken:
			if endTag, _ := z.TagName(); bytes.Equal(endTag, tag) {
				depth--
			}
		}
	}
	return
}
