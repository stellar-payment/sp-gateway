package dto

type PassthroughPayload struct {
	ServiceName   string `param:"svc"`
	EndpointPath  string `param:"path"`
	PartnerID     string `header:"X-Partner-Id"`
	KeypairHash   string `header:"X-Sec-Keypair"`
	CorrelationID string `header:"X-Correlation-Id"`
	ExternalID    string `header:"X-External-Id"`
	RequestMethod string
	Payload       string
}

type PassthroughResponse struct {
	Status  int
	Payload string
	Headers map[string]string
}
