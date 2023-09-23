package errs

import (
	"errors"
	"net/http"

	"github.com/stellar-payment/sp-gateway/pkg/constant"
	"github.com/stellar-payment/sp-gateway/pkg/dto"
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrBrokenUserReq       = errors.New("invalid request")
	ErrInvalidCred         = errors.New("invalid user credentials")
	ErrDuplicatedResources = errors.New("entity already existed")
	ErrNoAccess            = errors.New("user does not have required access privilege")
	ErrUnknown             = errors.New("internal server error")
	ErrNotFound            = errors.New("entity not found")
	ErrTokenExpired        = errors.New("user token already expired")
	ErrUserExisted         = errors.New("user already existed")
	ErrUserDeactivated     = errors.New("user is deactivated")
)

// Errcode: AAA-BB-C
// AAA => HTTP STATUS CODE
// BB = 01 Basic, 02 Business Logic
// C = ErrorID
// Ex: 403021 = 403 (Forbidden) - Business Logic - ID 1
const (
	ErrCodeUndefined           constant.ErrCode = 500011
	ErrCodeBadRequest          constant.ErrCode = 400012
	ErrCodeNoAccess            constant.ErrCode = 403013
	ErrCodeInvalidCred         constant.ErrCode = 401014
	ErrCodeDuplicatedResources constant.ErrCode = 400015
	ErrCodeBrokenUserReq       constant.ErrCode = 422016
	ErrCodeNotFound            constant.ErrCode = 404017
	ErrCodeTokenExpired        constant.ErrCode = 403021
	ErrCodeUserExisted         constant.ErrCode = 400022
	ErrCodeUserDeactivated     constant.ErrCode = 403023
)

const (
	ErrStatusUnknown     = http.StatusInternalServerError
	ErrStatusClient      = http.StatusBadRequest
	ErrStatusNotLoggedIn = http.StatusUnauthorized
	ErrStatusNoAccess    = http.StatusForbidden
	ErrStatusReqBody     = http.StatusUnprocessableEntity
	ErrStatusNotFound    = http.StatusNotFound
)

var errorMap = map[error]dto.ErrorResponse{
	ErrUnknown:             ErrorResponse(ErrStatusUnknown, ErrCodeUndefined, ErrUnknown),
	ErrBadRequest:          ErrorResponse(ErrStatusClient, ErrCodeBadRequest, ErrBadRequest),
	ErrInvalidCred:         ErrorResponse(ErrStatusNotLoggedIn, ErrCodeInvalidCred, ErrInvalidCred),
	ErrNoAccess:            ErrorResponse(ErrStatusNoAccess, ErrCodeNoAccess, ErrNoAccess),
	ErrDuplicatedResources: ErrorResponse(ErrStatusClient, ErrCodeDuplicatedResources, ErrDuplicatedResources),
	ErrBrokenUserReq:       ErrorResponse(ErrStatusReqBody, ErrCodeBrokenUserReq, ErrBrokenUserReq),
	ErrNotFound:            ErrorResponse(ErrStatusNotFound, ErrCodeNotFound, ErrNotFound),
	ErrTokenExpired:        ErrorResponse(ErrStatusNoAccess, ErrCodeTokenExpired, ErrNoAccess),
	ErrUserExisted:         ErrorResponse(ErrStatusClient, ErrCodeUserExisted, ErrDuplicatedResources),
	ErrUserDeactivated:     ErrorResponse(ErrStatusNoAccess, ErrCodeUserDeactivated, ErrUserDeactivated),
}

func ErrorResponse(status int, code constant.ErrCode, err error) dto.ErrorResponse {
	return dto.ErrorResponse{
		Status:  status,
		Code:    code,
		Message: err.Error(),
	}
}

func GetErrorResp(err error) (errResponse dto.ErrorResponse) {
	errResponse, ok := errorMap[err]
	if !ok {
		errResponse = errorMap[ErrUnknown]
	}

	return
}
