package transport

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "gitlab.ozon.dev/kanat_9999/homework/cart/internal/config"
)

type RetryRoundTripper struct {
	Next           http.RoundTripper
	MaxRetries     int
	InitialBackoff time.Duration
}

func NewRetryRoundTripper(next http.RoundTripper, maxRetries int, initialBackoff time.Duration) *RetryRoundTripper {
	return &RetryRoundTripper{Next: next, MaxRetries: maxRetries, InitialBackoff: initialBackoff}
}

func (r *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	var resp *http.Response
	var err error

	for retries := 0; retries < r.MaxRetries; retries++ {
		if retries > 0 {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		resp, err = r.Next.RoundTrip(req)
		if err == nil && !shouldRetry(resp) {
			return resp, nil
		}

		if err != nil || shouldRetry(resp) {
			if retries < r.MaxRetries-1 {
				backoff := time.Duration(1<<retries) * r.InitialBackoff
				log.Printf("Retry attempt %d/%d, status code: %d, retrying...",
					retries+1, r.MaxRetries, resp.StatusCode)

				time.Sleep(backoff)
				continue
			}
			break
		}
	}

	return nil, fmt.Errorf("request failed after %d retries", r.MaxRetries)
}

func shouldRetry(resp *http.Response) bool {
	return resp != nil && (resp.StatusCode == 420 || resp.StatusCode == 429)
}
