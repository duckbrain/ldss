package ldsorg

import (
	"compress/zlib"
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/duckbrain/ldss/lib"
)

var Default = Client{
	BaseURL:    "https://tech.lds.org/glweb",
	PlatformID: 17,
	Client:     http.DefaultClient,
}

type Client struct {
	BaseURL    string
	PlatformID int
	Client     *http.Client
}

func (c Client) get(ctx context.Context, action string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"?action="+action, nil)
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
	langs := []Lang{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&langs)
	return langs, err
}

func (c Client) Catalog(ctx context.Context, lang Lang) (Catalog, error) {
	catalog := Catalog{}
	res, err := c.get(ctx, "catalog.query")
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
	m := Metadata{}
	err := store.Metadata(ctx, index, &m)
	if err == lib.ErrNotFound {
		if index.Path == "/" {
			m.Type = TypeCatalog
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

	switch m.Type {
	case TypeBook:
		z, err := c.ZBook(ctx, m.DownloadURL)
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
		lang, err := c.Lang(ctx, store, index.Lang)
		if err != nil {
			return err
		}
		catalog, err := c.Catalog(ctx, lang)
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

func (c Client) Lang(ctx context.Context, store lib.Storer, libLang lib.Lang) (Lang, error) {
	m := Metadata{}
	index := lib.Index{Lang: libLang, Path: "/"}
	err := store.Metadata(ctx, index, &m)
	if err != nil {
		return Lang{}, err
	}
	if len(m.Languages) == 0 {
		langs, err := c.Languages(ctx)
		if err != nil {
			return Lang{}, err
		}
		m.Languages = make(map[string]Lang)
		for _, lang := range langs {
			m.Languages[lang.Code] = lang
		}
		if err := store.SetMetadata(ctx, index, m); err != nil {
			return Lang{}, err
		}
	}
	lang, ok := m.Languages[string(libLang)]
	if !ok {
		return Lang{}, lib.ErrNotFound
	}
	return lang, nil
}

func storeFolder(ctx context.Context, store lib.Storer, item *lib.Item, folder Folder) error {
	lang := item.Lang

	item.Header = folder.Header(ctx, lang)

	item.Children = make([]lib.Index, len(folder.Folders)+len(folder.Books))
	for i, childFolder := range folder.Folders {
		childItem := lib.Item{}
		childItem.Parent = item.Index
		if i > 0 {
			childItem.Prev = item.Children[i-1]
		}
		if i < len(item.Children)-1 {
			childItem.Next = folder.Folders[i+1].Header(ctx, lang).Index
		}

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
	item.Children = make([]lib.Index, len(children))
	for i, childNode := range children {
		footnotes, err := z.Footnotes(ctx, childNode.ID)
		if err != nil {
			return err
		}

		childItem := lib.Item{
			Header:    childNode.Header(lang),
			Content:   node.Content,
			Parent:    item.Index,
			Footnotes: footnotes,
		}
		if i > 0 {
			childItem.Prev = item.Children[i-1]
		}
		if i < len(children)-1 {
			childItem.Next = children[i+1].Header(lang).Index
		}

		// This populates the children of the child item and stores it
		err = storeBook(ctx, store, z, &childItem, childNode)
		if err != nil {
			return err
		}
	}

	return store.Store(ctx, *item)
}
