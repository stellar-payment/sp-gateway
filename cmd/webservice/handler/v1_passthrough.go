package handler

import (
	"context"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-gateway/internal/util/echttputil"
	"github.com/stellar-payment/sp-gateway/pkg/dto"
)

type PassthroughV1Handler func(context.Context, *dto.PassthroughPayload) (*dto.PassthroughResponse, error)

func HandlePassthroughV1(handler PassthroughV1Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &dto.PassthroughPayload{
			ServiceName:   c.Param("svc"),
			EndpointPath:  c.Param("path"),
			Headers:       c.Request().Header,
			RequestMethod: c.Request().Method,
		}

		// extract request body
		reqBody := []byte{}
		if c.Request().Body != nil {
			reqBody, _ = io.ReadAll(c.Request().Body)
		}
		req.Payload = string(reqBody)

		res, err := handler(c.Request().Context(), req)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WritePassthroughResponse(c, res)
	}
}
