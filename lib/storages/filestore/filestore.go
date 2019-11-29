package filestore

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/blevesearch/bleve"
	"github.com/duckbrain/ldss/lib"
	bolt "github.com/etcd-io/bbolt"
	"github.com/pkg/errors"
)

var _ lib.Storer = &FileStore{}

var (
	bucketItems    = []byte("Items")
	bucketMetadata = []byte("Metadata")
)
var buckets = [][]byte{bucketItems, bucketMetadata}

type FileStore struct {
	index       bleve.Index
	db          *bolt.DB
	marshaller  func(in interface{}) (out []byte, err error)
	unmarshaler func(in []byte, out interface{}) (err error)
}

func New(dir string) (*FileStore, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	index, err := bleve.Open(path.Join(dir, "search.bleve"))
	if err != nil {
		mapping := bleve.NewIndexMapping()
		mapping.AddDocumentMapping("item", bleve.NewDocumentMapping())
		index, err = bleve.New(path.Join(dir, "search.bleve"), mapping)
		if err != nil {
			return nil, errors.Wrap(err, "bleve")
		}
	}

	db, err := bolt.Open(path.Join(dir, "store.bbolt"), 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "bbolt open")
	}

	err = db.Update(createBuckets)

	store := &FileStore{
		index:       index,
		db:          db,
		marshaller:  json.Marshal,
		unmarshaler: json.Unmarshal,
	}
	return store, err
}

func createBuckets(tx *bolt.Tx) error {
	for _, name := range buckets {
		if _, err := tx.CreateBucketIfNotExists(name); err != nil {
			return errors.Wrapf(err, "bbolt create %v", name)
		}
	}
	return nil
}

func (s *FileStore) Item(ctx context.Context, index lib.Index) (item lib.Item, err error) {
	err = s.BulkRead(func(s lib.Storer) error {
		item, err = s.Item(ctx, index)
		return err
	})
	return item, err
}
func (s *FileStore) Store(ctx context.Context, item lib.Item) error {
	return s.BulkEdit(func(s lib.Storer) error {
		return s.Store(ctx, item)
	})
}
func (s *FileStore) Header(ctx context.Context, index lib.Index) (header lib.Header, err error) {
	s.BulkRead(func(s lib.Storer) error {
		header, err = s.Header(ctx, index)
		return err
	})
	return
}
func (s *FileStore) Clear(ctx context.Context) error {
	return s.BulkEdit(func(s lib.Storer) error {
		return s.Clear(ctx)
	})
}
