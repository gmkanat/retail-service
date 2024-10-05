package roundtripper

import (
	"golang.org/x/time/rate"
	"net/http"
)

type RateLimitedTransport struct {
	Transport http.RoundTripper
	Limiter   *rate.Limiter
}

func (t *RateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	err := t.Limiter.Wait(req.Context())
	if err != nil {
		return nil, err
	}

	return t.Transport.RoundTrip(req)
}
