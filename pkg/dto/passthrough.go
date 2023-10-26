package dto

type PassthroughPayload struct {
	ServiceName   string `param:"svc"`
	EndpointPath  string `param:"path"`
	Headers       map[string][]string
	RequestMethod string
	Payload       string
}

type PassthroughResponse struct {
	Status  int
	Payload string
	Headers map[string]string
}
