package filestore

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/blevesearch/bleve"
	"github.com/duckbrain/ldss/lib"
	bolt "github.com/etcd-io/bbolt"
)

var _ lib.Store = &Store{}

var (
	bucketItems    = []byte("Items")
	bucketMetadata = []byte("Metadata")
)

type Store struct {
	index       bleve.Index
	db          *bolt.DB
	marshaller  func(in interface{}) (out []byte, err error)
	unmarshaler func(in []byte, out interface{}) (err error)
}

func New(dir string) (*Store, error) {
	mapping := bleve.NewIndexMapping()

	mapping.AddDocumentMapping("item", bleve.NewDocumentMapping())

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	index, err := bleve.New(path.Join(dir, "search.bleve"), mapping)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(path.Join(dir, "store.bbolt"), 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		for _, name := range [][]byte{bucketItems, bucketMetadata} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})

	store := &Store{
		index:       index,
		db:          db,
		marshaller:  json.Marshal,
		unmarshaler: json.Unmarshal,
	}
	return store, err
}

func (s *Store) Item(ctx context.Context, index lib.Index) (lib.Item, error) {
	return s.BulkRead(func(s lib.Store) {
		return s.Item(ctx, index)
	})
}
func (s *Store) Store(ctx context.Context, item lib.Item) error {
	return s.BulkEdit(func(s lib.Store) {
		return s.Store(ctx, item)
	})
}
func (s *Store) Header(ctx context.Context, index lib.Index) (header lib.Header, err error) {
	s.BulkRead(func(s lib.Store) {
		header, err = s.Header(ctx, index)
		return err
	})
	return
}
func (s *Store) Metadata(ctx context.Context, index lib.Index, metadata interface{}) error {
	return s.BulkRead(func(s lib.Store) {
		return s.Metadata(ctx, index, metadata)
	})
}
func (s *Store) SetMetadata(ctx context.Context, index lib.Index, metadata interface{}) error {
	return s.BulkEdit(func(s lib.Store) {
		return s.SetMetadata(ctx, index, metadata)
	})
}
func (s *Store) Search(ctx context.Context, query string, results chan<- lib.Result) error {
	panic("not implemented")
}
