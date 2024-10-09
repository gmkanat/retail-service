package roundtripper_test

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/kanat_9999/homework/cart/internal/pkg/product/roundtripper"
)

func TestCustomRateLimitedTransport(t *testing.T) {
	rateLimit := 10 // 10 RPS
	burstLimit := 10

	rateLimitedTransport, err := roundtripper.NewCustomRateLimitedTransport(http.DefaultTransport, rateLimit, burstLimit)
	require.NoError(t, err, "failed to set rate limiter")

	client := &http.Client{
		Transport: rateLimitedTransport,
	}

	numRequests := 30

	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			req, err := http.NewRequest("GET", "https://httpbin.org/anything", nil)
			require.NoError(t, err)

			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("request %d failed: %v", i, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("request %d returned status: %d", i, resp.StatusCode)
			} else {
				fmt.Printf("request %d: status code %d\n", i, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()

	totalTime := time.Since(startTime)
	fmt.Printf("total time of %d requests: %v\n", numRequests, totalTime)

	expectedDuration := time.Duration(numRequests/rateLimit) * time.Second
	if totalTime < expectedDuration {
		t.Errorf("finished quickly, expected minimum %v, but it was %v", expectedDuration, totalTime)
	}
}
