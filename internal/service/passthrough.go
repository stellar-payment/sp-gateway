package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/stellar-payment/sp-gateway/internal/config"
	"github.com/stellar-payment/sp-gateway/internal/inconst"
	"github.com/stellar-payment/sp-gateway/internal/util/apiutil"
	"github.com/stellar-payment/sp-gateway/internal/util/svcutil"
	"github.com/stellar-payment/sp-gateway/pkg/dto"
	"github.com/stellar-payment/sp-gateway/pkg/errs"
)

func encryptRequest(ctx context.Context, partnerID string, data string) (res *dto.SecurityEncryptResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	apireq := &dto.SecurityEncryptPayload{
		Data:      base64.StdEncoding.EncodeToString([]byte(data)),
		PartnerID: partnerID,
	}

	bReq, err := json.Marshal(apireq)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal decrypt payload")
		return
	}

	caller := apiutil.NewRequester[dto.SecurityEncryptResponse]()
	apires, err := caller.SendRequest(ctx,
		conf.PassthroughConfig.SecurityPath+"/security/api/v1"+inconst.SEC_ENCRYPT,
		http.MethodPost,
		nil,
		bytes.NewBuffer(bReq),
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to process request")
		return
	}

	return apires, nil
}

func decryptRequest(ctx context.Context, payload *dto.PassthroughPayload) (res string, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	splitted := strings.Split(payload.Payload, ".")
	if len(splitted) != 3 {
		logger.Error().Err(errs.ErrBadRequest).Msg("payload must consist of 3 dot-separated data")
		return "", errs.ErrBadRequest
	}

	apireq := &dto.SecurityDecryptPayload{
		Data:        fmt.Sprintf("%s.%s", splitted[0], splitted[1]),
		PartnerID:   payload.Headers[inconst.HeaderXPartnerID][0],
		Tag:         splitted[2],
		KeypairHash: payload.Headers[inconst.HeaderXSecKeypair][0],
	}

	bReq, err := json.Marshal(apireq)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal decrypt payload")
		return
	}

	caller := apiutil.NewRequester[dto.SecurityDecryptResponse]()
	apires, err := caller.SendRequest(ctx,
		conf.PassthroughConfig.SecurityPath+"/security/api/v1"+inconst.SEC_DECRYPT,
		http.MethodPost,
		nil,
		bytes.NewBuffer(bReq),
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to process request")
		return
	}

	data, _ := base64.StdEncoding.DecodeString(apires.Data)
	return string(data), nil
}

func (s *service) PassthroughV1Request(ctx context.Context, payload *dto.PassthroughPayload) (res *dto.PassthroughResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	isSecuredRoute := false

	var basePath string
	switch payload.ServiceName {
	case inconst.SVC_ACCOUNT:
		basePath = conf.PassthroughConfig.AccountPath
	case inconst.SVC_SECURITY:
		basePath = conf.PassthroughConfig.SecurityPath
	case inconst.SVC_PAYMENT:
		basePath = conf.PassthroughConfig.PaymentPath
	default:
		logger.Error().Err(errs.ErrNotFound).Send()
		return nil, errs.ErrNotFound
	}

	if routes, ok := s.secureRouteStore[payload.ServiceName]; ok && !payload.OverrideSecurity {
		if _, ok := routes[strings.Split(payload.EndpointPath, "/")[2]]; ok {
			isSecuredRoute = true
		}
	}

	if isSecuredRoute && (payload.RequestMethod != http.MethodGet && payload.RequestMethod != http.MethodOptions) {
		var inres string
		inres, err = decryptRequest(ctx, payload)
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}

		payload.Payload = inres
	}

	fmt.Printf("%#+v\n", payload)

	params := &svcutil.SendRequestParams{
		Endpoint: fmt.Sprintf("%s/%s/%s", basePath, payload.ServiceName, payload.EndpointPath),
		Method:   payload.RequestMethod,
		Body:     payload.Payload,
		Queries:  payload.Queries,
	}

	params.Headers = make(map[string]string)
	for k, v := range payload.Headers {
		params.Headers[k] = strings.Join(v, ",")
	}

	caller := svcutil.NewRequester()
	svcres, err := caller.SendRequest(ctx, params)
	if err != nil {
		logger.Error().Err(err).Msg("failed to send request")
		return
	}

	res = &dto.PassthroughResponse{
		Status:  svcres.Status,
		Payload: svcres.Payload,
		Headers: make(map[string]string),
	}

	for k, v := range svcres.Headers {
		res.Headers[k] = strings.Join(v, ",")
	}

	if isSecuredRoute {
		outres, err := encryptRequest(ctx, payload.Headers[inconst.HeaderXPartnerID][0], res.Payload)
		if err != nil {
			logger.Error().Err(err).Send()
			return nil, err
		}

		out, err := json.Marshal(outres)
		if err != nil {
			logger.Error().Err(err).Send()
			return nil, err
		}

		res.Payload = string(out)
	}

	return
}
