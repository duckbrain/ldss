package lib

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
		if content, err := node.Content(); err == nil {
			weight := content.Search(r.Keywords)
			if weight > 0 {
				c <- SearchResult{
					Reference: r,
					Weight:    weight,
				}
			}
			return nil
		}
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

func (content *Content) Search(keywords []string) int {
	return 0
}

//func Search(i Item, terms []string
