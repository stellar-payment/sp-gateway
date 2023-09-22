package ctxutil

import (
	"context"

	"github.com/nmluci/go-backend/internal/commonkey"
)

func WrapCtx(ctx context.Context, key commonkey.CtxKey, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}

func GetCtx[T any](ctx context.Context, key commonkey.CtxKey) (res T) {
	val := ctx.Value(key)
	res, _ = val.(T)

	return
}
