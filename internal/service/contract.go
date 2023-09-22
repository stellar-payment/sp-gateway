package service

import (
	"github.com/nmluci/go-backend/internal/repository"
	"github.com/nmluci/go-backend/pkg/dto"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)
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
