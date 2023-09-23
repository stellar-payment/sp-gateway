package service

import (
	"bytes"
	"context"
	"encoding/json"
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
		PartnerID:   structutil.StringToUint64(payload.PartnerID),
		Tag:         splitted[1],
		KeypairHash: payload.KeypairHash,
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
	case inconst.SVC_AUTH:
		basePath = conf.PassthroughConfig.AuthPath
	case inconst.SVC_ACCOUNT:
		basePath = conf.PassthroughConfig.AccountPath
	case inconst.SVC_SECURITY:
		basePath = conf.PassthroughConfig.SecurityPath
	case inconst.SVC_PAYMENT:
		basePath = conf.PassthroughConfig.PaymentPath
	}

	if payload.ServiceName == inconst.SVC_ACCOUNT || payload.ServiceName == inconst.SVC_PAYMENT {
		var inres string
		inres, err = decryptRequest(ctx, payload)
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}

		payload.Payload = inres
	}

	params := svcutil.SendReqeuestParams{
		Endpoint: basePath + payload.EndpointPath,
		Method:   payload.RequestMethod,
		Body:     payload.Payload,
	}

	headers := make(map[string]string)
	headers["x-correlation-id"] = payload.CorrelationID

	caller := svcutil.NewRequester()
	apires, err := caller.SendRequest(ctx, params)
	if err != nil {
		logger.Error().Err(err).Msg("failed to send request")
		return
	}

	res = &dto.PassthroughResponse{
		Payload: apires,
		Headers: headers,
	}

	if payload.ServiceName == inconst.SVC_ACCOUNT || payload.ServiceName == inconst.SVC_PAYMENT {
		var outres, sk string
		outres, sk, err = encryptRequest(ctx, structutil.StringToUint64(payload.PartnerID), apires)
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}

		res.Payload = outres
		res.Headers["x-sec-keypair"] = sk
	}

	return
}