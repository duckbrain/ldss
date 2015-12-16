package ldslib

import (
	"golang.org/x/net/html"
	"html/template"
	"strconv"
	"strings"
)

type ContentParser struct {
	node        Node
	contentHtml string
	content     *Content
}

type parseMode int

const (
	parseTitleMode parseMode = iota
	parseSubtitleMode
	parseSummaryMode
	parseVerseMode
)

type Content struct {
	Title, Subtitle, Summary string
	Verses                   []Verse
	originalHtml             template.HTML
}

func (c *Content) String() string {
	return string(c.originalHtml)
}

func (c *Content) HTML() template.HTML {
	return c.originalHtml
}

type Verse struct {
	Number int
	Text   string
}

type VerseReference struct {
	Verse  int
	Letter rune
}

func (p *ContentParser) parse() error {
	if p.content != nil {
		return nil
	}
	reader := strings.NewReader(p.contentHtml)
	doc, err := html.Parse(reader)
	if err != nil {
		return err
	}

	mode := parseTitleMode
	content := new(Content)
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
						content.Verses = append(content.Verses, verse)
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
				content.Title += text
			case parseSubtitleMode:
				content.Subtitle += text
			case parseSummaryMode:
				content.Summary += text
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
		content.Verses = append(content.Verses, verse)
	}
	content.originalHtml = template.HTML(p.contentHtml)

	p.content = content
	return nil
}

func (p *ContentParser) Content() (*Content, error) {
	err := p.parse()
	return p.content, err
}

func (p *ContentParser) OriginalHTML() string {
	return p.contentHtml
}
