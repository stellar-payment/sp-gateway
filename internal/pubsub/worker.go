package pubsub

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-gateway/internal/inconst"
	"github.com/stellar-payment/sp-gateway/internal/indto"
	"github.com/stellar-payment/sp-gateway/internal/service"
	"github.com/stellar-payment/sp-gateway/internal/util/ctxutil"
)

type EventPubSub struct {
	logger       zerolog.Logger
	redis        *redis.Client
	service      service.Service
	secureRoutes []string
}

type NewEventPubSubParams struct {
	Logger       zerolog.Logger
	Redis        *redis.Client
	Service      service.Service
	SecureRoutes []string
}

func NewEventPubSub(params *NewEventPubSubParams) *EventPubSub {
	return &EventPubSub{
		logger:       params.Logger,
		redis:        params.Redis,
		service:      params.Service,
		secureRoutes: params.SecureRoutes,
	}
}

func (pb *EventPubSub) Listen() {
	ctx := context.Background()
	ctx = ctxutil.WrapCtx(ctx, inconst.SCOPE_CTX_KEY, indto.UserScopeMap{})

	subscriber := pb.redis.Subscribe(ctx, inconst.TOPIC_BROADCAST_SECURE_ROUTE)

	if err := pb.redis.Publish(context.Background(), inconst.TOPIC_REQUEST_SECURE_ROUTE, nil).Err(); err != nil {
		pb.logger.Error().Err(err).Send()
	}

	defer subscriber.Close()
	for msg := range subscriber.Channel() {
		switch msg.Channel {
		case inconst.TOPIC_BROADCAST_SECURE_ROUTE:
			pb.logger.Info().Str("event", msg.Channel).Str("msg", msg.Payload).Msg("incoming secure route msg")
			splitted := strings.Split(msg.Payload, ",")
			svcname := splitted[0]
			routes := splitted[1:]

			pb.service.UpsertSecureEndpoint(svcname, routes)
			continue
		}
	}
}
