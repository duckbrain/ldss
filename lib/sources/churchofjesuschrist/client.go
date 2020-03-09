package churchofjesuschrist

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/duckbrain/ldss/lib"
)

var Default = &Client{
	BaseURL: "https://www.churchofjesuschrist.org/study/",
	Client:  http.DefaultClient,
}

type Client struct {
	BaseURL string
	Client  *http.Client
}

func (c Client) get(ctx context.Context, p string, params url.Values) (io.ReadCloser, error) {
	logger := ctx.Value(lib.CtxLogger).(lib.Logger)

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, p)
	u.RawQuery = params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	logger.Infof("downloading: %v", req.URL.String())
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, err
}

func (c Client) json(ctx context.Context, p string, params url.Values, out interface{}) error {
	body, err := c.get(ctx, p, params)
	if err != nil {
		return err
	}
	defer body.Close()
	decoder := json.NewDecoder(body)
	return decoder.Decode(out)
}

func (c Client) Dynamic(ctx context.Context, index lib.Index) (dynamic Dynamic, err error) {
	params := url.Values{}
	params.Set("lang", "eng") // TODO: Use real languages
	params.Set("uri", index.Path)
	err = c.json(ctx, "/api/v3/language-pages/type/dynamic", params, &dynamic)
	return
}

func (c Client) Load(ctx context.Context, store lib.Storer, index lib.Index) error {
	logger := ctx.Value(lib.CtxLogger).(lib.Logger)

	if index.Path == "/" || index.Path == "" {
		logger.Debug("skipping path \"/\"")
		return lib.ErrNotFound
	}

	dynamic, err := c.Dynamic(ctx, index)
	if err != nil {
		// TODO: Check for 404
		return err
	}

	item := dynamic.AsItem(index)
	logger.Debug("collection ", dynamic)

	err = store.Store(ctx, item)
	if err != nil {
		return err
	}

	return nil
}
