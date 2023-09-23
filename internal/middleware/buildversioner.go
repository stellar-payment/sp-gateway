package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-gateway/internal/config"
	"github.com/stellar-payment/sp-gateway/internal/inconst"
)

func ServiceVersioner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		conf := config.Get()

		crid := uuid.NewString()
		c.Request().Header.Set(inconst.HeaderXCorrelationID, crid)
		c.Response().Header().Set(inconst.HeaderXCorrelationID, crid)
		c.Response().Header().Set("BUILD-TIME", conf.BuildTime)
		c.Response().Header().Set("BUILD-VER", conf.BuildVer)

		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}
