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
	item := lib.Item{}

	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketItems).Get(index.Hash())
		if data == nil {
			return lib.ErrNotFound
		}

		return s.unmarshaler(data, &item)
	})

	return item, err
}
func (s *Store) Store(ctx context.Context, item lib.Item) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		data, err := s.marshaller(item)
		if err != nil {
			return err
		}

		err = tx.Bucket(bucketItems).Put(item.Hash(), data)
		if err != nil {
			return nil
		}

		return s.index.Index(string(item.Hash()), item)
	})
}
func (s *Store) Header(ctx context.Context, index lib.Index) (lib.Header, error) {
	item := lib.Item{}

	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketItems).Get(index.Hash())
		if data == nil {
			return lib.ErrNotFound
		}

		return s.unmarshaler(data, &item)
	})

	return item.Header, err
}
func (s *Store) Metadata(ctx context.Context, index lib.Index, metadata interface{}) error {
	return s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketMetadata).Get(index.Hash())
		if data == nil {
			return lib.ErrNotFound
		}

		return s.unmarshaler(data, metadata)
	})
}
func (s *Store) SetMetadata(ctx context.Context, index lib.Index, metadata interface{}) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		data, err := s.marshaller(metadata)
		if err != nil {
			return err
		}

		return tx.Bucket(bucketMetadata).Put(index.Hash(), data)
	})
}
func (s *Store) Search(ctx context.Context, query string, results chan<- lib.Result) error {
	panic("not implemented")
}
