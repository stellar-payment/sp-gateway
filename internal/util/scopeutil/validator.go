package scopeutil

import (
	"context"

	"github.com/stellar-payment/sp-gateway/internal/inconst"
	"github.com/stellar-payment/sp-gateway/internal/indto"
	"github.com/stellar-payment/sp-gateway/internal/util/ctxutil"
)

func ValidateScope(ctx context.Context, scope string) (ok bool) {
	usrScope := ctxutil.GetCtx[indto.UserScopeMap](ctx, inconst.SCOPE_CTX_KEY)
	_, ok = usrScope[scope]
	return
}
