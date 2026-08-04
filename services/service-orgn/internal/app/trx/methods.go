package trx

import (
	"context"
	"errors"
	"sync"
)

type ctxKey struct{}

type trx struct {
	rollbackFuncs []func() error
	commitFuncs   []func() error
	values        map[any]any
	mu            sync.Mutex
}

func newTrx() *trx {
	return &trx{values: make(map[any]any)}
}

func Exists(ctx context.Context) bool {
	_, ok := ctx.Value(ctxKey{}).(*trx)
	return ok
}

func Begin(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, newTrx())
}

func AddRollback(ctx context.Context, fn func() error) {
	t, ok := ctx.Value(ctxKey{}).(*trx)
	if !ok {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.rollbackFuncs = append(t.rollbackFuncs, fn)
}

func Rollback(ctx context.Context) error {
	t, ok := ctx.Value(ctxKey{}).(*trx)
	if !ok {
		return errors.New("transaction not started")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	var err error
	for _, fn := range t.rollbackFuncs {
		err = errors.Join(err, fn())
	}

	return err
}

func AddCommit(ctx context.Context, fn func() error) {
	t, ok := ctx.Value(ctxKey{}).(*trx)
	if !ok {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.commitFuncs = append(t.commitFuncs, fn)
}

func Commit(ctx context.Context) error {
	t, ok := ctx.Value(ctxKey{}).(*trx)
	if !ok {
		return errors.New("transaction not started")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	var err error
	for _, fn := range t.commitFuncs {
		err = errors.Join(err, fn())
	}

	return err
}

func GetValue(ctx context.Context, key any) (any, bool) {
	t, ok := ctx.Value(ctxKey{}).(*trx)
	if !ok {
		return nil, false
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	val, ok := t.values[key]
	return val, ok
}

func SetValue(ctx context.Context, key, val any) {
	t, ok := ctx.Value(ctxKey{}).(*trx)
	if !ok {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.values[key] = val
}
