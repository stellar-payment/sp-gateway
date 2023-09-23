package ctxutil

import (
	"context"

	"github.com/stellar-payment/sp-gateway/internal/inconst"
)

func WrapCtx(ctx context.Context, key inconst.CtxKey, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}

func GetCtx[T any](ctx context.Context, key inconst.CtxKey) (res T) {
	val := ctx.Value(key)
	res, _ = val.(T)

	return
}
