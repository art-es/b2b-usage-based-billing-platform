package auth

import "context"

type ctxKey struct{}

func Set(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, ctxKey{}, val)
}

func Get(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(ctxKey{}).(string)
	return val, ok
}
