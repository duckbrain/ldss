package lib

import (
	"bytes"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// Content pulled from a node in the SQlite database. Is the content of the node
// formatted as HTML
type Content string
type ContentParser struct {
	content Content
	z       *html.Tokenizer

	// Paragraph info
	paragraphStyle     ParagraphStyle
	verse              int
	justFoundParagraph bool

	// For determining if you are done with a paragraph
	paragraphTag string
	depth        int

	// Text info
	textContent string
	textStyle   TextStyle
	href        string
}
type ParagraphStyle int
type TextStyle int

const (
	ParagraphStyleNormal ParagraphStyle = iota
	ParagraphStyleTitle
	ParagraphStyleChapter
	ParagraphStyleSummary
)

func (p ParagraphStyle) String() string {
	switch p {
	case ParagraphStyleNormal:
		return "Normal"
	case ParagraphStyleTitle:
		return "Title"
	case ParagraphStyleChapter:
		return "Chapter"
	case ParagraphStyleSummary:
		return "Summary"
	}
	return ""
}

const (
	TextStyleNormal TextStyle = iota
	TextStyleLink
	TextStyleFootnote
)

func (t TextStyle) String() string {
	switch t {
	case TextStyleNormal:
		return "Normal"
	case TextStyleLink:
		return "Link"
	case TextStyleFootnote:
		return "Footnote"
	}
	return ""
}

func (c Content) Parse() *ContentParser {
	return &ContentParser{content: c}
}

func (c *ContentParser) NextParagraph() bool {
	if c.z == nil {
		c.z = html.NewTokenizerFragment(strings.NewReader(string(c.content)), "div")
	}
	log.Println("Reading next paragraph")
	for {
		switch c.z.Next() {
		case html.ErrorToken:
			log.Println(c.z.Err())
			return false
		case html.TextToken:
			log.Printf("Paragraph found %v %v\n", c.paragraphStyle, c.verse)
			c.depth = 1
			c.textStyle = TextStyleNormal
			c.justFoundParagraph = true
			return true
		case html.StartTagToken:
			var key, val []byte
			tag, hasAttr := c.z.TagName()
			c.verse = 0
			log.Printf("Paragraph tag found %v\n", string(tag))
			switch string(tag) {
			case "h1":
				c.paragraphStyle = ParagraphStyleTitle
			case "p":
				c.paragraphStyle = ParagraphStyleNormal
				for hasAttr {
					key, val, hasAttr = c.z.TagAttr()
					switch string(key) {
					case "class":
						switch string(val) {
						case "titleNumber":
							c.paragraphStyle = ParagraphStyleChapter
						case "studySummary":
							c.paragraphStyle = ParagraphStyleSummary
						}
					case "id":
						if verse, err := strconv.Atoi(string(val)); err == nil {
							log.Println("Setting verse ", verse)
							c.verse = verse
						}
					}
				}
			case "video":
				//TODO
			}
			c.paragraphTag = string(tag)
		case html.EndTagToken:
		case html.SelfClosingTagToken:
		case html.CommentToken:
		case html.DoctypeToken:
		}
	}
	return false
}
func (c *ContentParser) ParagraphStyle() ParagraphStyle {
	return c.paragraphStyle
}
func (c *ContentParser) ParagraphVerse() int {
	return c.verse
}

func (c *ContentParser) NextText() bool {
	if c.justFoundParagraph {
		c.justFoundParagraph = false
		return true
	}
	for {
		switch c.z.Next() {
		case html.ErrorToken:
			return false
		case html.TextToken:
			return true
		case html.StartTagToken:
			var key, val []byte
			tag, hasAttr := c.z.TagName()
			switch string(tag) {
			case c.paragraphTag:
				c.depth++
			case "sup":
				c.textStyle = TextStyleFootnote
			case "a":
				c.textStyle = TextStyleLink
				for hasAttr {
					key, val, hasAttr = c.z.TagAttr()
					switch string(key) {
					case "href":
						c.href = string(val)
					}
				}
			}
		case html.EndTagToken:
			tag, _ := c.z.TagName()
			switch string(tag) {
			case c.paragraphTag:
				c.depth--
				if c.depth == 0 {
					return false
				}
			case "sup":
				c.textStyle = TextStyleNormal
			case "a":
				c.textStyle = TextStyleLink
			}
		case html.SelfClosingTagToken:
		case html.CommentToken:
		case html.DoctypeToken:
		}
	}
}
func (c *ContentParser) TextStyle() TextStyle {
	return c.textStyle
}
func (c *ContentParser) Text() string {
	text := string(c.z.Text())
	return text
}

/*
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
*/

func (content Content) Filter(verses []int) Content {
	if len(verses) == 0 {
		return content
	}

	hasZero := sort.SearchInts(verses, 0) != -1

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
			continue
		default:
			if verse == nextAllowed || hasZero && verse == 0 {
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
