package filestore

import (
	"context"
	"os"
	"path"

	"github.com/duckbrain/ldss/lib"
	bolt "github.com/etcd-io/bbolt"
	"github.com/pkg/errors"
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
func (s bulkStore) Clear(ctx context.Context) error {
	if s.readonly {
		panic("cannot Clear in read-only")
	}

	// Delete and re-create the buckets in the db
	for _, name := range buckets {
		if err := s.tx.DeleteBucket(name); err != nil {
			return errors.Wrapf(err, "bbolt delete %v", name)
		}
	}
	if err := createBuckets(s.tx); err != nil {
		return err
	}

	err := s.index.Close()
	if err != nil {
		return err
	}
	err = os.RemoveAll(path.Join(s.dir, "search.bleve"))
	if err != nil {
		return err
	}
	s.index, err = createBleveIndex(s.dir)
	if err != nil {
		return err
	}
	return nil
}
