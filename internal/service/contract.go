package service

import (
	"context"

	"github.com/stellar-payment/sp-gateway/internal/repository"
	"github.com/stellar-payment/sp-gateway/pkg/dto"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)

	PassthroughV1Request(context.Context, *dto.PassthroughPayload) (*dto.PassthroughResponse, error)
}

type service struct {
	conf       *serviceConfig
	repository repository.Repository
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Repository repository.Repository
}

func NewService(params *NewServiceParams) Service {
	return &service{
		conf:       &serviceConfig{},
		repository: params.Repository,
	}
}
