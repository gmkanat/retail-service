package transport

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type RetryRoundTripper struct {
	Next http.RoundTripper
}

func NewRetryRoundTripper(next http.RoundTripper) *RetryRoundTripper {
	return &RetryRoundTripper{Next: next}
}

func (r *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	const (
		maxRetries     = 3
		initialBackoff = 500 * time.Millisecond
	)

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	var resp *http.Response
	var err error

	for retries := 0; retries < maxRetries; retries++ {
		if retries > 0 {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		resp, err = r.Next.RoundTrip(req)
		if err == nil && !shouldRetry(resp) {
			return resp, nil
		}

		if err != nil || shouldRetry(resp) {
			if retries < maxRetries-1 {
				backoff := time.Duration(1<<retries) * initialBackoff
				log.Printf("Retry attempt %d/%d, status code: %d, retrying...",
					retries+1, maxRetries, resp.StatusCode)

				time.Sleep(backoff)
				continue
			}
			break
		}
	}

	return nil, fmt.Errorf("request failed after %d retries", maxRetries)
}

func shouldRetry(resp *http.Response) bool {
	return resp != nil && (resp.StatusCode == 420 || resp.StatusCode == 429)
}
