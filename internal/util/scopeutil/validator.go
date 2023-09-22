package scopeutil

import (
	"context"

	"github.com/nmluci/go-backend/internal/commonkey"
	"github.com/nmluci/go-backend/internal/indto"
	"github.com/nmluci/go-backend/internal/util/ctxutil"
)

func ValidateScope(ctx context.Context, scope string) (ok bool) {
	usrScope := ctxutil.GetCtx[indto.UserScopeMap](ctx, commonkey.SCOPE_CTX_KEY)
	_, ok = usrScope[scope]
	return
}
