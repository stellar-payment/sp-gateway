package webservice

import (
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-gateway/cmd/webservice/router"
	"github.com/stellar-payment/sp-gateway/internal/component"
	"github.com/stellar-payment/sp-gateway/internal/config"
	"github.com/stellar-payment/sp-gateway/internal/pubsub"
	"github.com/stellar-payment/sp-gateway/internal/repository"
	"github.com/stellar-payment/sp-gateway/internal/service"
)

const logTagStartWebservice = "[StartWebservice]"

func Start(conf *config.Config, logger zerolog.Logger) {
	redis, err := component.InitRedis(&component.InitRedisParams{
		Conf:   &conf.RedisConfig,
		Logger: logger,
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initalize redis")
	}

	ec := echo.New()
	ec.HideBanner = true
	ec.HidePort = true

	repo := repository.NewRepository(&repository.NewRepositoryParams{
		// MariaDB: db,
		// MongoDB:    mongo,
		Redis: redis,
	})

	service := service.NewService(&service.NewServiceParams{
		Repository: repo,
	})

	psWorker := pubsub.NewEventPubSub(&pubsub.NewEventPubSubParams{
		Logger:  logger,
		Redis:   redis,
		Service: service,
	})

	router.Init(&router.InitRouterParams{
		Logger:  logger,
		Service: service,
		Ec:      ec,
		Conf:    conf,
	})

	wg := &sync.WaitGroup{}

	logger.Info().Msgf("starting service, listening to: %s", conf.ServiceAddress)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := ec.Start(conf.ServiceAddress); err != nil {
			logger.Error().Msgf("starting service, cause: %+v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		psWorker.Listen()
	}()

	wg.Wait()
}
