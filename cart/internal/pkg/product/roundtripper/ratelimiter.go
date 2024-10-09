package roundtripper

import (
	"context"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/customerrors"
	"net/http"
	"time"
)

type CustomRateLimitedTransport struct {
	Transport http.RoundTripper
	tokens    chan struct{}
	cancel    context.CancelFunc
}

func NewCustomRateLimitedTransport(transport http.RoundTripper, rateLimit, burstLimit int) (*CustomRateLimitedTransport, error) {
	if rateLimit <= 0 || burstLimit <= 0 {
		return nil, customerrors.RateLimitSetupFail
	}

	tokens := make(chan struct{}, burstLimit)
	ticker := time.NewTicker(time.Second / time.Duration(rateLimit))

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case tokens <- struct{}{}:
				default:
					// chan is full, burstLimit reached
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return &CustomRateLimitedTransport{
		Transport: transport,
		tokens:    tokens,
		cancel:    cancel,
	}, nil
}

func (t *CustomRateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	<-t.tokens
	return t.Transport.RoundTrip(req)
}

func (t *CustomRateLimitedTransport) Shutdown() {
	t.cancel()
}
