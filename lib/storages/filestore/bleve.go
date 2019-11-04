package filestore

import (
	"context"

	"github.com/blevesearch/bleve"
	"github.com/duckbrain/ldss/lib"
)

var _ lib.Indexer = FileStore{}

func (s FileStore) Index(ctx context.Context, item lib.Item) error {
	return s.index.Index(string(item.Hash()), item)
}

func (s FileStore) Search(ctx context.Context, ref lib.Reference, results chan<- lib.Result) error {
	req := bleve.NewSearchRequest(bleve.NewConjunctionQuery())
	_, err := s.index.SearchInContext(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
