package ldsorg

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/duckbrain/ldss/lib"
)

type Legacy struct {
	Dir string
}

func (l Legacy) Languages(ctx context.Context) ([]Lang, error) {
	file, err := os.Open(path.Join(l.Dir, "languages.json"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	langs := []Lang{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&langs)
	return langs, err
}

func (l Legacy) Catalog(ctx context.Context, lang lib.Lang) (Catalog, error) {
	catalog := Catalog{}
	file, err := os.Open(path.Join(l.Dir, "languages.json"))
	if err != nil {
		return catalog, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&catalog)
	return catalog, err
}

func (l Legacy) ZBook(ctx context.Context, i lib.Index) (*ZBook, error) {
	file, err := os.Open(path.Join(l.Dir, string(i.Lang), i.Path, "contents.sqlite"))
	if err != nil {
		return nil, err
	}
	return NewZBook(file)
}

func (l Legacy) Load(ctx context.Context, store lib.Storer, index lib.Index) error {
	m := Metadata{}
	err := store.Metadata(ctx, index, &m)
	if err == lib.ErrNotFound {
		if index.Path == "/" {
			m.Type = TypeCatalog
		} else {
			i := index
			i.Path = path.Dir(i.Path)
			err := l.Load(ctx, store, i)
			if err != nil {
				return err
			}
		}
	} else if err != nil {
		return err
	}

	switch m.Type {
	case TypeBook:
		z, err := l.ZBook(ctx, index)
		if err != nil {
			return err
		}
		item, err := store.Item(ctx, index)
		if err != nil {
			return err
		}
		err = storeBook(ctx, store, z, &item, Node{})
		if err != nil {
			return err
		}
	case TypeCatalog:
		catalog, err := l.Catalog(ctx, index.Lang)
		if err != nil {
			return err
		}

		err = storeFolder(ctx, store, &lib.Item{}, catalog.Folder)
		if err != nil {
			return err
		}
	}

	return nil
}
