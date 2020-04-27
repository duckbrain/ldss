package lib

import (
	"context"
	"sort"

	"github.com/pkg/errors"
)

type Library struct {
	Store   Storer
	Index   Indexer
	Sources []Loader
	Logger  Logger
	Parser  *ReferenceParser
}

type smartStore struct {
	Storer
	Indexer
}

func (s smartStore) Store(ctx context.Context, item Item) error {
	if item.Path == "" {
		return errors.New("item has no path")
	}
	if item.Lang == "" {
		return errors.New("item has no lang")
	}
	err := s.Storer.Store(ctx, item)
	if err != nil {
		return err
	}
	return s.Index(ctx, item)
}

type dummyLogger struct{}

func (dummyLogger) Debugf(string, ...interface{}) {}
func (dummyLogger) Infof(string, ...interface{})  {}
func (dummyLogger) Printf(string, ...interface{}) {}
func (dummyLogger) Warnf(string, ...interface{})  {}
func (dummyLogger) Errorf(string, ...interface{}) {}
func (dummyLogger) Fatalf(string, ...interface{}) {}
func (dummyLogger) Debug(...interface{})          {}
func (dummyLogger) Info(...interface{})           {}
func (dummyLogger) Warn(...interface{})           {}
func (dummyLogger) Error(...interface{})          {}
func (dummyLogger) Fatal(...interface{})          {}
func (dummyLogger) Panic(...interface{})          {}

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Printf(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})
}

type ContextKey string

const (
	CtxLogger    ContextKey = "logger"
	CtxStore     ContextKey = "store"
	CtxIndex     ContextKey = "index"
	CtxRefParser ContextKey = "reference-parser"
)

func (l Library) ctx(ctx context.Context) context.Context {
	logger := l.Logger
	if logger == nil {
		logger = dummyLogger{}
	}
	// indexStore ensures that any newly stored items will be indexed as well
	store := smartStore{
		Storer:  l.Store,
		Indexer: l.Index,
	}
	ctx = context.WithValue(ctx, CtxLogger, logger)
	ctx = context.WithValue(ctx, CtxStore, store)
	ctx = context.WithValue(ctx, CtxIndex, l.Index)
	ctx = context.WithValue(ctx, CtxRefParser, l.Parser)

	return ctx
}

func (l Library) LookupAndDownload(ctx context.Context, index Index) (Item, error) {
	ctx = l.ctx(ctx)
	logger := ctx.Value(CtxLogger).(Logger)
	store := ctx.Value(CtxStore).(Storer)

	item, err := store.Item(ctx, index)
	if err != ErrNotFound {
		return item, errors.Wrap(err, "initial lookup")
	}

	logger.Debugf("index %v not found, trying to download", index)
	for _, source := range l.Sources {
		err := source.Load(ctx, store, index)
		if err == nil {
			item, err := store.Item(ctx, index)
			return item, errors.Wrapf(err, "second lookup %v", source)
		}
		if err == ErrNotFound {
			logger.Debugf("skipping %v, not found", source)
			continue
		}
		return item, errors.Wrapf(err, "load %v", source)
	}

	return item, ErrNotFound
}

func (l Library) Parse(ctx context.Context, lang Lang, q string) ([]Reference, error) {
	if q != "" && q[0] == '/' {
		ref := l.Parser.ParsePath(lang, q)
		return []Reference{ref}, nil
	}
	return l.Parser.Parse(lang, q)
}

func (l Library) Lookup(ctx context.Context, index Index) (Item, error) {
	ctx = l.ctx(ctx)
	store := ctx.Value(CtxStore).(Storer)

	return store.Item(ctx, index)
}
func (l Library) LookupReference(ctx context.Context, ref *Reference) (Item, error) {
	item, err := l.LookupAndDownload(ctx, ref.Index)
	if err == nil {
		ref.Name = item.Name
	}
	return item, err
}
func (l Library) Sibling(ctx context.Context, item Item, offset int) (Header, error) {
	if offset == 0 {
		return item.Header, nil
	}
	parentHeader := item.Parent()
	if !parentHeader.Valid() {
		return Header{}, nil
	}
	parent, err := l.Lookup(ctx, item.Parent().Index)
	if err != nil {
		return Header{}, err
	}
	for i, child := range parent.Children {
		if child.Index == item.Index {
			j := i + offset
			if j >= 0 && j < len(parent.Children) {
				return parent.Children[j], nil
			} else {
				return Header{}, nil
			}
		}
	}
	return Header{}, errors.New("could not find self in parent's children")
}

func (l Library) Download(ctx context.Context, index Index) error {
	ctx = l.ctx(ctx)
	store := ctx.Value(CtxStore).(Storer)

	for _, source := range l.Sources {
		err := source.Load(ctx, store, index)
		if err == nil {
			return nil
		}
		if err == ErrNotFound {
			continue
		}
		return err
	}
	return ErrNotFound
}

func (lib *Library) Register(l Loader) {
	lib.Sources = append(lib.Sources, l)
	ctx := lib.ctx(context.Background())
	if x, ok := l.(interface {
		LoadParser(context.Context, *ReferenceParser)
	}); ok {
		x.LoadParser(ctx, lib.Parser)
	}
}

func (l Library) Search(ctx context.Context, ref Reference, results chan<- Result) error {
	ctx = l.ctx(ctx)
	return l.Index.Search(ctx, ref, results)
}

func (l Library) SearchSlice(ctx context.Context, ref Reference) (Results, error) {
	out := make(chan Result)
	results := make(Results, 0)
	go func() {
		for result := range out {
			// TODO not thread safe
			results = append(results, result)
		}
	}()
	err := l.Search(ctx, ref, out)
	close(out)
	sort.Sort(results)
	return results, err
}
