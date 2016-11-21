package lib

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"golang.org/x/net/html"
)

type Footnote struct {
	item     Item
	Name     string        `json:"name"`
	LinkName string        `json:"linkName"`
	Content  template.HTML `json:"content"`
}

func (f *Footnote) References() (refs []Reference) {
	z := html.NewTokenizerFragment(strings.NewReader(string(f.Content)), "div")
	refs = make([]Reference, 0)
	lang := f.item.Language()

	for {
		ref := Reference{
			Language: lang,
			Content:  "",
		}

		switch z.Next() {
		case html.ErrorToken, html.EndTagToken:
			return
		case html.TextToken:
			ref.Content = string(z.Text())
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
						r := ParsePath(lang, string(val))
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
					return
				case html.TextToken:
					ref.Content = fmt.Sprintf("%v%v", ref.Content, string(z.Text()))
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
