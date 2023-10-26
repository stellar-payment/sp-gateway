package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/stellar-payment/sp-gateway/internal/config"
	"github.com/stellar-payment/sp-gateway/internal/inconst"
	"github.com/stellar-payment/sp-gateway/internal/util/apiutil"
	"github.com/stellar-payment/sp-gateway/internal/util/structutil"
	"github.com/stellar-payment/sp-gateway/internal/util/svcutil"
	"github.com/stellar-payment/sp-gateway/pkg/dto"
	"github.com/stellar-payment/sp-gateway/pkg/errs"
)

func encryptRequest(ctx context.Context, partnerID uint64, data string) (res string, sk string, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	apireq := &dto.SecurityEncryptPayload{
		Data:      data,
		PartnerID: partnerID,
	}

	bReq, err := json.Marshal(apireq)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal decrypt payload")
		return
	}

	caller := apiutil.NewRequester[dto.SecurityEncryptResponse]()
	apires, err := caller.SendRequest(ctx,
		conf.PassthroughConfig.SecurityPath+inconst.SEC_DECRYPT,
		http.MethodPost,
		nil,
		bytes.NewBuffer(bReq),
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to process request")
		return
	}

	return strings.Join([]string{apires.Data, apires.Tag}, "."), apires.SecretKey, nil
}

func decryptRequest(ctx context.Context, payload *dto.PassthroughPayload) (res string, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	splitted := strings.Split(payload.Payload, ".")
	if len(splitted) != 2 {
		logger.Error().Err(errs.ErrBadRequest).Msg("payload must consist of 2 dot-separated data")
		return "", errs.ErrBadRequest
	}

	apireq := &dto.SecurityDecryptPayload{
		Data:        splitted[0],
		PartnerID:   structutil.StringToUint64(payload.Headers[inconst.HeaderXPartnerID][0]),
		Tag:         splitted[1],
		KeypairHash: payload.Headers[inconst.HeaderXSecKeypair][0],
	}

	bReq, err := json.Marshal(apireq)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal decrypt payload")
		return
	}

	caller := apiutil.NewRequester[dto.SecurityDecryptResponse]()
	apires, err := caller.SendRequest(ctx,
		conf.PassthroughConfig.SecurityPath+inconst.SEC_DECRYPT,
		http.MethodPost,
		nil,
		bytes.NewBuffer(bReq),
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to process request")
		return
	}

	return apires.Data, nil
}

func (s *service) PassthroughV1Request(ctx context.Context, payload *dto.PassthroughPayload) (res *dto.PassthroughResponse, err error) {
	logger := log.Ctx(ctx)
	conf := config.Get()

	var basePath string
	switch payload.ServiceName {
	case inconst.SVC_ACCOUNT:
		basePath = conf.PassthroughConfig.AccountPath
	case inconst.SVC_SECURITY:
		basePath = conf.PassthroughConfig.SecurityPath
	case inconst.SVC_PAYMENT:
		basePath = conf.PassthroughConfig.PaymentPath
	}

	// Todo: add per-path whitelist
	// if payload.ServiceName == inconst.SVC_PAYMENT {
	// 	var inres string
	// 	inres, err = decryptRequest(ctx, payload)
	// 	if err != nil {
	// 		logger.Error().Err(err).Send()
	// 		return
	// 	}

	// 	payload.Payload = inres
	// }

	params := &svcutil.SendRequestParams{
		Endpoint: fmt.Sprintf("%s/%s/api/%s", basePath, payload.ServiceName, payload.EndpointPath),
		Method:   payload.RequestMethod,
		Body:     payload.Payload,
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

	// if payload.ServiceName == inconst.SVC_PAYMENT {
	// 	var outres, sk string
	// 	outres, sk, err = encryptRequest(ctx, structutil.StringToUint64(res.Headers[inconst.HeaderXPartnerID]), res.Payload)
	// 	if err != nil {
	// 		logger.Error().Err(err).Send()
	// 		return
	// 	}

	// 	res.Payload = outres
	// 	res.Headers["x-sec-keypair"] = sk
	// }

	return
}
