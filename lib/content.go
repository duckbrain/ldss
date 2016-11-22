package lib

import (
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
