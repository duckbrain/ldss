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

func SearchSort(item Item, keywords []string) []SearchResult {
	c := make(chan SearchResult)
	Search(item, keywords, c)
	results := []SearchResult{}
	for result := range c {
		results = append(results, result)
	}
	sort.Sort(SearchResults(results))
	return results
}

func Search(item Item, keywords []string, c chan<- SearchResult) {
	waitGroup := new(sync.WaitGroup)
	resultSet := make(map[string]bool)

	searchItem(item, keywords, c, waitGroup, resultSet)

	go func() {
		waitGroup.Wait()
		close(c)
	}()
}

// TODO Change where the search result stores the content with highlights on words
func searchItem(item Item, keywords []string, c chan<- SearchResult, waitGroup *sync.WaitGroup, resultSet map[string]bool) {
	if node, ok := item.(Contenter); ok {
		if resultSet[item.Path()] {
			return
		}
		resultSet[item.Path()] = true
		waitGroup.Add(1)
		go func(node Contenter, item Item) {
			if content, err := node.Content(); err == nil {
				result := content.Search(keywords)
				if result.Weight > 0 {
					result.Language = item.Language()
					result.Path = item.Path()
					result.Clean()
					c <- result
				}
			}
			waitGroup.Done()
		}(node, item)
	}

	if children, err := item.Children(); err == nil {
		for _, child := range children {
			searchItem(child, keywords, c, waitGroup, resultSet)
		}
	}
	return
}
