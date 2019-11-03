package filestore

import (
	"context"

	"github.com/duckbrain/ldss/lib"
	bolt "github.com/etcd-io/bbolt"
)

func (s *FileStore) BulkRead(fn func(lib.Storer) error) error {
	return s.db.View(func(tx *bolt.Tx) error {
		return fn(bulkStore{s, tx, true})
	})
}
func (s *FileStore) BulkEdit(fn func(lib.Storer) error) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return fn(bulkStore{s, tx, false})
	})
}

// bulkStore implements the Store interface, but does all of its operations in a single transaction
type bulkStore struct {
	*FileStore
	tx       *bolt.Tx
	readonly bool
}

func (s bulkStore) Item(ctx context.Context, index lib.Index) (lib.Item, error) {
	data := s.tx.Bucket(bucketItems).Get(index.Hash())
	if data == nil {
		return lib.Item{}, lib.ErrNotFound
	}

	item := lib.Item{}
	err := s.unmarshaler(data, &item)
	return item, err
}
func (s bulkStore) Store(ctx context.Context, item lib.Item) error {
	if s.readonly {
		panic("cannot Store in read-only")
	}
	data, err := s.marshaller(item)
	if err != nil {
		return err
	}

	err = s.tx.Bucket(bucketItems).Put(item.Hash(), data)
	if err != nil {
		return nil
	}

	return s.index.Index(string(item.Hash()), item)
}
func (s bulkStore) Header(ctx context.Context, index lib.Index) (lib.Header, error) {
	data := s.tx.Bucket(bucketItems).Get(index.Hash())
	if data == nil {
		return lib.Header{}, lib.ErrNotFound
	}

	item := lib.Item{}
	err := s.unmarshaler(data, &item)
	return item.Header, err
}
func (s bulkStore) Metadata(ctx context.Context, index lib.Index, metadata interface{}) error {
	data := s.tx.Bucket(bucketMetadata).Get(index.Hash())
	if data == nil {
		return lib.ErrNotFound
	}

	return s.unmarshaler(data, metadata)
}
func (s bulkStore) SetMetadata(ctx context.Context, index lib.Index, metadata interface{}) error {
	if s.readonly {
		panic("cannot SetMetadata in read-only")
	}
	data, err := s.marshaller(metadata)
	if err != nil {
		return err
	}

	return s.tx.Bucket(bucketMetadata).Put(index.Hash(), data)
}
