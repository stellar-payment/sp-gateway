package pubsub

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/nmluci/go-backend/internal/commonkey"
	"github.com/nmluci/go-backend/internal/indto"
	"github.com/nmluci/go-backend/internal/service"
	"github.com/nmluci/go-backend/internal/util/ctxutil"
	"github.com/rs/zerolog"
)

var (
	tagLoggerPBListen = "[PubSub-Listen]"
)

type FilePubSub struct {
	logger  zerolog.Logger
	redis   *redis.Client
	service service.Service
}

type NewFilePubSubParams struct {
	Logger  zerolog.Logger
	Redis   *redis.Client
	Service service.Service
}

func NewFileSub(params NewFilePubSubParams) *FilePubSub {
	return &FilePubSub{
		logger:  params.Logger,
		redis:   params.Redis,
		service: params.Service,
	}
}

func (pb *FilePubSub) Listen() {
	ctx := context.Background()
	ctx = ctxutil.WrapCtx(ctx, commonkey.SCOPE_CTX_KEY, indto.UserScopeMap{})

	subscriber := pb.redis.Subscribe(ctx, "")

	defer subscriber.Close()
	for msg := range subscriber.Channel() {
		fmt.Print(msg)
	}
}
