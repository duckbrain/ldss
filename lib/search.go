package lib

import (
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type SearchResult struct {
	Reference
	Weight int
}

func (r Reference) Search(c chan<- SearchResult) error {
	item, err := r.Lookup()
	if err != nil {
		return nil
	}

	return r.searchItem(item, c)
}

func (r Reference) searchItem(item Item, c chan<- SearchResult) error {
	if node, ok := item.(*Node); ok {
		go func() {
			if content, err := node.Content(); err == nil {
				result := content.Search(r.Keywords)
				if result.Weight > 0 {
					result.Language = item.Language()
					result.Path = item.Path()
					c <- result
				}
			}
		}()
	}
	children, err := item.Children()
	if err != nil {
		return err
	}
	for _, child := range children {
		err = r.searchItem(child, c)
		if err != nil {
			return err
		}
	}
	return nil
}

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
