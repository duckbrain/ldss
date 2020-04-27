package ldsorg

import (
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/pkg/errors"

	"github.com/duckbrain/ldss/lib"
)

var Default = &Client{
	BaseURL:    "https://tech.lds.org/glweb",
	PlatformID: 17,
	Client:     http.DefaultClient,
}

type Client struct {
	BaseURL    string
	PlatformID int
	Client     *http.Client
}

func (c Client) get(ctx context.Context, action string, v ...interface{}) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"?action="+fmt.Sprintf(action, v...), nil)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c Client) Languages(ctx context.Context) ([]Lang, error) {
	res, err := c.get(ctx, "languages.query")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var result struct {
		Languages []Lang
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	return result.Languages, err
}

func (c Client) Catalog(ctx context.Context, lang lib.Lang) (Catalog, error) {
	catalog := Catalog{}
	res, err := c.get(ctx, "catalog.query&languageid=%v&platformid=%v", lang.GLCode(), c.PlatformID)
	if err != nil {
		return catalog, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&catalog)
	return catalog, err
}

func (c Client) ZBook(ctx context.Context, url string) (*ZBook, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	r, err := zlib.NewReader(res.Body)
	if err != nil {
		return nil, err
	}
	return NewZBook(r)
}

func (c Client) Load(ctx context.Context, store lib.Storer, index lib.Index) error {
	logger := ctx.Value(lib.CtxLogger).(lib.Logger)
	item, err := store.Item(ctx, index)
	logger.Debugf("load %v item %v: %v", index, item, err)
	if err == lib.ErrNotFound {
		if index.Path == "/" {
			item.Metadata["ldsorg_type"] = string(TypeCatalog)
		} else {
			i := index
			i.Path = path.Dir(i.Path)
			err := c.Load(ctx, store, i)
			if err != nil {
				return err
			}
		}
	} else if err != nil {
		return err
	}

	switch ItemType(item.Metadata["ldsorg_type"]) {
	case TypeBook:
		downloadURL := item.Metadata["ldsorg_download_url"]
		logger.Debugf("download book %v", downloadURL)
		z, err := c.ZBook(ctx, downloadURL)
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
		logger.Debugf("download catalog %v", index.Lang)
		catalog, err := c.Catalog(ctx, index.Lang)
		if err != nil {
			return errors.Wrap(err, "catalog: download")
		}

		err = storeFolder(ctx, store, &lib.Item{}, catalog.Folder)
		if err != nil {
			return errors.Wrap(err, "catalog: store")
		}
	}

	return nil
}

func storeFolder(ctx context.Context, store lib.Storer, item *lib.Item, folder Folder) error {
	lang := item.Lang

	item.Header = folder.Header(ctx, lang)

	item.Children = make([]lib.Header, len(folder.Folders)+len(folder.Books))
	for _, childFolder := range folder.Folders {
		childItem := lib.Item{}
		childItem.Breadcrumbs = append(item.Breadcrumbs, item.Header)

		err := storeFolder(ctx, store, &childItem, childFolder)
		if err != nil {
			return err
		}
	}

	return store.Store(ctx, *item)
}

func storeBook(ctx context.Context, store lib.Storer, z *ZBook, item *lib.Item, node Node) error {
	lang := item.Lang

	children, err := z.Children(ctx, node.ID)
	if err != nil {
		return err
	}
	item.Children = make([]lib.Header, len(children))
	for _, childNode := range children {
		footnotes, err := z.Footnotes(ctx, childNode.ID)
		if err != nil {
			return err
		}

		childItem := lib.Item{
			Header:      childNode.Header(lang),
			Content:     node.Content,
			Breadcrumbs: append(item.Breadcrumbs, item.Header),
			Footnotes:   footnotes,
		}

		// This populates the children of the child item and stores it
		err = storeBook(ctx, store, z, &childItem, childNode)
		if err != nil {
			return err
		}
	}

	return store.Store(ctx, *item)
}
