package svcutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// Rate Limiter
type Transport struct {
	base    http.RoundTripper
	limiter *rate.Limiter
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	res := t.limiter.Reserve()

	select {
	case <-time.After(res.Delay()):
		return t.base.RoundTrip(r)
	case <-r.Context().Done():
		res.Cancel()
		return nil, r.Context().Err()
	}
}

type SendRequestParams struct {
	Endpoint string
	Method   string
	Headers  map[string]string
	Queries  url.Values
	Body     string
}

type APIResponse struct {
	Status  int
	Payload string
	Headers map[string][]string
}

// Actual API
type Requester struct {
	Client http.Client
}

func NewRequester() Requester {
	t := Requester{
		Client: http.Client{
			Transport: &Transport{base: http.DefaultTransport, limiter: rate.NewLimiter(rate.Limit(50), 1)},
		},
	}
	return t
}

func (r *Requester) SendRequest(ctx context.Context, params *SendRequestParams) (res *APIResponse, err error) {
	req, err := http.NewRequestWithContext(ctx, params.Method, params.Endpoint, bytes.NewBuffer([]byte(params.Body)))
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", "Stellar-Microservice by Misaki-chan")
	req.Header.Add("Content-Type", "application/json")

	for k, v := range params.Headers {
		req.Header.Add(k, v)
	}

	queries := req.URL.Query()
	for k, v := range params.Queries {
		req.URL.Query().Add(k, strings.Join(v, ","))
	}
	req.URL.RawQuery = queries.Encode()

	data, err := r.Client.Do(req)
	if err != nil {
		return
	}

	res = &APIResponse{}

	defer data.Body.Close()
	body, err := io.ReadAll(data.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body")
	}

	res.Headers = data.Header
	res.Status = data.StatusCode
	res.Payload = string(body)

	return
}
