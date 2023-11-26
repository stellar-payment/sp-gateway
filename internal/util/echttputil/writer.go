package echttputil

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-gateway/pkg/dto"
	"github.com/stellar-payment/sp-gateway/pkg/errs"
)

func WriteSuccessResponse(ec echo.Context, data interface{}) error {
	return ec.JSON(http.StatusOK, dto.BaseResponse{
		Data:   data,
		Errors: nil,
	})
}

func WritePassthroughResponse(ec echo.Context, res *dto.PassthroughResponse) error {
	for k, v := range res.Headers {
		ec.Response().Header().Set(k, v)
	}

	ec.Response().Header().Set(echo.HeaderContentType, "application/json; charset=UTF-8")
	ec.Response().Header().Set(echo.HeaderContentLength, strconv.Itoa(len(res.Payload)))
	return ec.String(res.Status, res.Payload)
}

func WriteErrorResponse(ec echo.Context, err error) error {
	errResp := errs.GetErrorResp(err)
	return ec.JSON(errResp.Status, dto.BaseResponse{
		Data:   nil,
		Errors: errResp,
	})
}

func WriteFileAttachment(ec echo.Context, path string, filename string) error {
	return ec.Attachment(path, filename)
}
