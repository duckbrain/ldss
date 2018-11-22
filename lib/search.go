package lib

import (
	"sort"
	"sync"

	"github.com/duckbrain/ldss/lib/dl"
)

type SearchResult struct {
	Reference
	Weight int
}

type SearchResults []SearchResult

func (r SearchResults) Len() int {
	return len(r)
}

func (r SearchResults) Less(i, j int) bool {
	if r[i].Weight == r[j].Weight {
		return r[i].Path < r[j].Path
	} else {
		return r[i].Weight > r[j].Weight
	}
}

func (r SearchResults) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func SearchSort(item Item, keywords []string) SearchResults {
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
	if downloader, ok := item.(dl.Downloader); ok {
		if !downloader.Downloaded() {
			return
		}
	}
	if opener, ok := item.(Opener); ok {
		if err := opener.Open(); err != nil {
			panic(err)
		}
	}

	if node, ok := item.(Contenter); ok {
		if resultSet[item.Path()] {
			return
		}
		resultSet[item.Path()] = true
		waitGroup.Add(1)
		go func(node Contenter, item Item) {
			content := node.Content()
			result := content.Search(keywords)
			if result.Weight > 0 {
				result.Lang = item.Lang()
				result.Path = item.Path()
				result.Clean()
				c <- result
			}
			waitGroup.Done()
		}(node, item)
	}

	for _, child := range item.Children() {
		searchItem(child, keywords, c, waitGroup, resultSet)
	}
	return
}
