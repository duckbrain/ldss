package lib

import (
	"html/template"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type Content struct {
	rawHTML string
}

func (c *Content) HTML() template.HTML {
	return template.HTML(c.rawHTML)
}

type parseMode int

const (
	parseTitleMode parseMode = iota
	parseSubtitleMode
	parseSummaryMode
	parseVerseMode
)

type Page struct {
	Title, Subtitle, Summary string
	Verses                   []Verse
	originalHtml             template.HTML
}

func (c *Page) String() string {
	return string(c.originalHtml)
}

func (c *Page) HTML() template.HTML {
	return c.originalHtml
}

type Verse struct {
	Number int
	Text   string
}

type VerseReference struct {
	Verse  int
	Letter string
}

func (c *Content) Page() (*Page, error) {
	reader := strings.NewReader(c.rawHTML)
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	mode := parseTitleMode
	page := new(Page)
	var verse Verse
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
	return page, nil
}
