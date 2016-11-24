package lib

import (
	"bytes"
	"math"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// Content pulled from a node in the SQlite database. Is the content of the node
// formatted as HTML
type Content string

type contentParseMode int

const (
	parseTitleMode contentParseMode = iota
	parseSubtitleMode
	parseSummaryMode
	parseVerseMode
)

// A page parsed from a node's Content
type Page struct {
	Title, Subtitle, Summary string
	Verses                   []struct {
		Number int
		Text   string
	}
}

// Parse the content for a page. The page contains an structured representation
// of the content that can be displayed programattically in a variety of ways.
func (c Content) Page() *Page {
	page := new(Page)
	reader := strings.NewReader(string(c))

	doc, err := html.Parse(reader)
	if err != nil {
		return page
	}

	mode := parseTitleMode
	var verse struct {
		Number int
		Text   string
	}
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "type" && attr.Val == "chapter" {
					mode = parseTitleMode
				}
				if attr.Key == "class" && attr.Val == "studySummary" {
					mode = parseSummaryMode
				}
				if attr.Key == "class" && attr.Val == "bodyBlock" {
					mode = parseVerseMode
				}
				if mode == parseVerseMode && attr.Key == "id" {
					if verse.Number > 0 {
						page.Verses = append(page.Verses, verse)
					}
					verse.Number, err = strconv.Atoi(attr.Val)
					if err != nil {
						verse.Number = 0
					}
					verse.Text = ""
				}
			}
		}
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			switch mode {
			case parseTitleMode:
				page.Title += text
			case parseSubtitleMode:
				page.Subtitle += text
			case parseSummaryMode:
				page.Summary += text
			case parseVerseMode:
				text = strings.TrimLeft(text, " 1234567890")
				verse.Text += text + " "
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(doc)
	if verse.Number > 0 {
		page.Verses = append(page.Verses, verse)
	}
	return page
}

func (content Content) Filter(verses []int) Content {
	if len(verses) == 0 {
		return content
	}

	z := html.NewTokenizerFragment(strings.NewReader(string(content)), "div")
	verse := 0
	buffer := new(bytes.Buffer)
	verses = cleanVerses(verses)
	nextAllowedIndex := 0
	nextAllowed := verses[0]

	for {
		switch z.Next() {
		case html.ErrorToken:
			return Content(buffer.Bytes())
		case html.StartTagToken:
			raw := z.Raw()
			_, hasAttr := z.TagName()
			var key, val []byte
			for hasAttr {
				key, val, hasAttr = z.TagAttr()
				if string(key) == "id" {
					verse, _ = strconv.Atoi(string(val))
					if verse > nextAllowed {
						nextAllowedIndex++
						if nextAllowedIndex == len(verses) {
							nextAllowed = math.MaxInt32
						} else {
							nextAllowed = verses[nextAllowedIndex]
						}
					}
				}
			}
			if verse == nextAllowed {
				_, _ = buffer.Write(raw)
			}
		default:
			if verse == nextAllowed {
				_, _ = buffer.Write(z.Raw())
			}
		}
	}
	return Content(buffer.Bytes())
}

// Search the content for the given keywords and return a search result containing
// the verses in which the results were found and the a weighted score based on
// the number of occurances.
func (content Content) Search(keywords []string) SearchResult {
	z := html.NewTokenizerFragment(strings.NewReader(string(content)), "div")
	r := SearchResult{}
	verse := 0

	for {
		switch z.Next() {
		case html.ErrorToken:
			return r
		case html.TextToken:
			text := strings.ToLower(string(z.Text()))
			for _, k := range keywords {
				weight := strings.Count(text, k)
				if weight > 0 && verse > 0 {
					r.VersesHighlighted = append(r.VersesHighlighted, verse)
				}
				r.Weight += weight
			}
		case html.StartTagToken:
			_, hasAttr := z.TagName()
			var key, val []byte
			for hasAttr {
				key, val, hasAttr = z.TagAttr()
				if string(key) == "id" {
					verse, _ = strconv.Atoi(string(val))
				}
			}
		}
	}

	r.Clean()

	return r
}
