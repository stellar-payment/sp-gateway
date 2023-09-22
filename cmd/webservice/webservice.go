package webservice

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/go-backend/cmd/webservice/router"
	"github.com/nmluci/go-backend/internal/config"
	"github.com/nmluci/go-backend/internal/repository"
	"github.com/nmluci/go-backend/internal/service"
	"github.com/rs/zerolog"
)

const logTagStartWebservice = "[StartWebservice]"

func Start(conf *config.Config, logger zerolog.Logger) {
	// redis, err := component.InitRedis(&component.InitRedisParams{
	// 	Conf:   &conf.RedisConfig,
	// 	Logger: logger,
	// })

	// if err != nil {
	// 	logger.Fatalf("%s initalizing redis: %+v", logTagStartWebservice, err)
	// }

	ec := echo.New()
	ec.HideBanner = true
	ec.HidePort = true

	repo := repository.NewRepository(&repository.NewRepositoryParams{
		// MariaDB: db,
		// MongoDB:    mongo,
		// Redis: redis,
	})

	service := service.NewService(&service.NewServiceParams{
		Repository: repo,
	})

	router.Init(&router.InitRouterParams{
		Logger:  logger,
		Service: service,
		Ec:      ec,
		Conf:    conf,
	})

	logger.Info().Msgf("starting service, listening to: %s", conf.ServiceAddress)

	if err := ec.Start(conf.ServiceAddress); err != nil {
		logger.Error().Msgf("starting service, cause: %+v", err)
	}
}
