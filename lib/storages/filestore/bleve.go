package filestore

import (
	"context"

	"github.com/blevesearch/bleve"
	"github.com/duckbrain/ldss/lib"
)

var _ lib.Indexer = BleveIndex{}

type BleveIndex struct {
	BleveIndex bleve.Index
}

func (i BleveIndex) Index(ctx context.Context, item lib.Item) error {
	return i.BleveIndex.Index(string(item.Hash()), item)
}

func (i BleveIndex) Search(ctx context.Context, ref lib.Reference, results chan<- lib.Result) error {
	index := i.BleveIndex
	req := bleve.NewSearchRequest(bleve.NewConjunctionQuery())
	_, err := index.SearchInContext(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
