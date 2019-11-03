package lib

import (
	"context"
	"sort"
)

type Library struct {
	Store   Storer
	Index   Indexer
	Sources []Loader
	Logger  Logger
	Parser  *ReferenceParser
}

type indexStore struct {
	Storer
	Indexer
}

func (s indexStore) Store(ctx context.Context, item Item) error {
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
	store := indexStore{
		Storer:  l.Store,
		Indexer: l.Index,
	}
	ctx = context.WithValue(ctx, CtxLogger, logger)
	ctx = context.WithValue(ctx, CtxStore, store)
	ctx = context.WithValue(ctx, CtxIndex, l.Index)
	ctx = context.WithValue(ctx, CtxRefParser, l.Parser)

	return ctx
}

func (l Library) Lookup(ctx context.Context, index Index) (Item, error) {
	ctx = l.ctx(ctx)
	logger := ctx.Value(CtxLogger).(Logger)
	store := ctx.Value(CtxStore).(Storer)

	item, err := store.Item(ctx, index)
	if err != ErrNotFound {
		return item, err
	}

	logger.Debugf("index %v not found, trying to download", index)
	for _, source := range l.Sources {
		err := source.Load(ctx, store, index)
		if err == nil {
			return store.Item(ctx, index)
		}
		if err == ErrNotFound {
			continue
		}
		return item, err
	}

	return item, ErrNotFound
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
