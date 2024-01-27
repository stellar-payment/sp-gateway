package dto

import "net/url"

type PassthroughPayload struct {
	ServiceName      string `param:"svc"`
	EndpointPath     string `param:"path"`
	Queries          url.Values
	Headers          map[string][]string
	OverrideSecurity bool
	RequestMethod    string
	Payload          string
}

type PassthroughResponse struct {
	Status  int
	Payload string
	Headers map[string]string
}
