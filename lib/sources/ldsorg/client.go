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

func (c Client) get(action string) (http.Response, error) {
	return c.Client.Get(c.BaseURL + "?action=" + action)
}

func (c Client) Languages(lang lib.Lang) ([]Lang, error) {
	res, err := c.get("languages.query")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	langs := []Lang{}
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(&catalog)
	return catalog, err
}

func (c Client) Catalog(lang lib.Lang) (Catalog, error) {
	catalog := Catalog{}
	res, err := c.get("catalog.query")
	if err != nil {
		return catalog, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(&catalog)
	return catalog, err
}

func (c Client) ZBook(path string) (*ZBook, error) {
	res, err := c.Client.Get(path)
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

func (c Client) Load(ctx context.Context, store lib.Store, index lib.Index) error {
	var item lib.Item
	var err error
	i := index
	for {
		item, err = store.Item(ctx, i)
		if err == nil {
			break
		}
		if err != lib.ErrNotFound {
			return err
		}
		if path == "/" || path == "." {
			break
		}
		i.Path = path.Dir(i.Path)
	}

}
