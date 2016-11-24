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

	err = r.searchItem(item, c)
	close(c)
	return err
}

func (r Reference) searchItem(item Item, c chan<- SearchResult) error {
	if node, ok := item.(*Node); ok {
		if content, err := node.Content(); err == nil {
			result := content.Search(r.Keywords)
			if result.Weight > 0 {
				result.Language = item.Language()
				result.Path = item.Path()
				c <- result
			}
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
