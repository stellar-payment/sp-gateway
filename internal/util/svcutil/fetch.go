package svcutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
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

type SendReqeuestParams struct {
	Endpoint string
	Method   string
	Headers  map[string]string
	Body     string
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

func (r *Requester) SendRequest(ctx context.Context, params SendReqeuestParams) (status int, res string, err error) {
	req, err := http.NewRequestWithContext(ctx, params.Method, params.Endpoint, bytes.NewBuffer([]byte(params.Body)))
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", "Stellar-Microservice by Misaki-chan")
	req.Header.Add("Content-Type", "application/json")

	for k, v := range params.Headers {
		req.Header.Add(k, v)
	}

	data, err := r.Client.Do(req)
	if err != nil {
		return
	}

	defer data.Body.Close()
	body, err := io.ReadAll(data.Body)
	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("failed to read response body")
	}

	return data.StatusCode, string(body), nil
}
