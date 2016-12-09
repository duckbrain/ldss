package lib

import (
	"sort"
	"sync"
)

type SearchResult struct {
	Reference
	Weight int
}

type SearchResults []SearchResult

func (r SearchResults) Len() int {
	return len([]SearchResult(r))
}

func (r SearchResults) Less(i, j int) bool {
	rs := []SearchResult(r)
	if rs[i].Weight == rs[j].Weight {
		return rs[i].Path < rs[j].Path
	} else {
		return rs[i].Weight > rs[j].Weight
	}

}

func (r SearchResults) Swap(i, j int) {
	rs := []SearchResult(r)
	t := rs[i]
	rs[i] = rs[j]
	rs[j] = t
}

func (r Reference) SearchSort() ([]SearchResult, error) {
	c := make(chan SearchResult)
	err := r.Search(c)
	if err != nil {
		return nil, err
	}
	results := []SearchResult{}
	for result := range c {
		results = append(results, result)
	}
	sort.Sort(SearchResults(results))
	return results, nil
}

func (r Reference) Search(c chan<- SearchResult) error {
	item, err := r.Lookup()
	if err != nil {
		return nil
	}

	waitGroup := new(sync.WaitGroup)
	resultSet := make(map[string]bool)

	r.searchItem(item, c, waitGroup, resultSet)

	go func() {
		waitGroup.Wait()
		close(c)
	}()

	return nil
}

func (r Reference) searchItem(item Item, c chan<- SearchResult, waitGroup *sync.WaitGroup, resultSet map[string]bool) {
	if node, ok := item.(*Node); ok {
		if resultSet[node.Path()] {
			return
		}
		resultSet[node.Path()] = true
		waitGroup.Add(1)
		go func(node *Node) {
			if content, err := node.Content(); err == nil {
				result := content.Search(r.Keywords)
				if result.Weight > 0 {
					result.Language = node.Language()
					result.Path = node.Path()
					result.Clean()
					c <- result
				}
			}
			waitGroup.Done()
		}(node)
	}

	if children, err := item.Children(); err == nil {
		for _, child := range children {
			r.searchItem(child, c, waitGroup, resultSet)
		}
	}
	return
}
