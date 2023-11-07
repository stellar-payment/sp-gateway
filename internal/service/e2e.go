package service

import "github.com/stellar-payment/sp-gateway/internal/component"

func (s *service) UpsertSecureEndpoint(service string, routes []string) {
	logger := component.GetLogger()

	_, ok := s.secureRouteStore[service]
	if ok {
		delete(s.secureRouteStore, service)
	}

	s.secureRouteStore[service] = make(map[string]struct{})

	for _, route := range routes {
		s.secureRouteStore[service][route] = struct{}{}
	}

	logger.Info().Str("secure-service", service).Any("secure-routes", routes).Msg("updated secure endpoint mapping")
}
