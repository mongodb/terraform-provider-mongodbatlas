package config

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type retryTransport struct {
	transport  http.RoundTripper
	maxRetries int
	retryDelay time.Duration
}

var retriableErrors = map[string]bool{
	"connection reset by peer":  true,
	"connection refused":        true,
	"no such host":              true,
	"i/o timeout":               true,
	"EOF":                       true,
	"context deadline exceeded": true,
}

func newRetryTransport(transport http.RoundTripper, maxRetries int, retryDelay time.Duration) *retryTransport {
	return &retryTransport{
		transport:  transport,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= t.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-req.Context().Done():
				return nil, fmt.Errorf("request cancelled during retry: %w", req.Context().Err())
			case <-time.After(t.retryDelay * time.Duration(attempt)): // Exponential backoff
			}
		}

		resp, err := t.transport.RoundTrip(req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		if !isRetriableError(err) {
			return nil, fmt.Errorf("non-retriable error: %w", err)
		}

		if attempt < t.maxRetries {
			fmt.Printf("Retrying request due to error: %v (attempt %d/%d)\n", err, attempt+1, t.maxRetries)
		}
	}

	return nil, fmt.Errorf("max retries reached (%d), last error: %w", t.maxRetries, lastErr)
}

func isRetriableError(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's a network error
	if netErr, ok := err.(net.Error); ok {
		// Retry on timeout or temporary errors
		return netErr.Timeout() || netErr.Temporary()
	}

	// Check against known retriable error messages
	errStr := err.Error()
	for msg := range retriableErrors {
		if strings.Contains(errStr, msg) {
			return true
		}
	}

	return false
}
