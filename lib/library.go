package lib

import "context"

type Library struct {
	Store   Store
	Sources []Source
	Logger  Logger
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
	CtxLogger ContextKey = "logger"
	CtxStore  ContextKey = "store"
)

func (l Library) ctx(ctx context.Context) context.Context {
	logger := l.Logger
	if logger == nil {
		logger = dummyLogger{}
	}
	ctx = context.WithValue(ctx, CtxLogger, logger)
	ctx = context.WithValue(ctx, CtxStore, l.Store)

	return ctx
}

func (l Library) Lookup(ctx context.Context, index Index) (Item, error) {
	ctx = l.ctx(ctx)
	logger := ctx.Value(CtxLogger).(Logger)

	item, err := l.Store.Item(ctx, index)
	if err != ErrNotFound {
		return item, err
	}

	logger.Debugf("index %v not found, trying to download", index)
	for _, source := range l.Sources {
		err := source.Load(ctx, l.Store, index)
		if err == nil {
			return l.Store.Item(ctx, index)
		}
		if err == ErrNotFound {
			continue
		}
		return item, err
	}

	return item, ErrNotFound
}
