package config

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type RetryTransport struct {
	transport  http.RoundTripper
	maxRetries int
	retryDelay time.Duration
}

// retriableErrors contains error messages that indicate a request should be retried.
// These are typically transient network issues that may resolve on retry.
var retriableErrors = []string{
	"connection reset by peer",
	"connection refused",
	"no such host",
	"i/o timeout",
	"EOF",
	"context deadline exceeded",
}

func NewRetryTransport(transport http.RoundTripper, maxRetries int, retryDelay time.Duration) *RetryTransport {
	return &RetryTransport{
		transport:  transport,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// RoundTrip implements the http.RoundTripper interface with automatic retry logic.
// It attempts to execute the request up to maxRetries times with exponential backoff.
//
// The retry strategy works as follows:
// 1. First attempt is made immediately
// 2. On failure, if the error is retriable (see isRetriableError):
//   - Waits for retryDelay * attempt (exponential backoff)
//   - Retries up to maxRetries times
//
// 3. If the error is not retriable, returns immediately
// 4. If context is cancelled during retry, returns immediately
// 5. If all retries are exhausted, returns the last error
//
// The exponential backoff helps prevent overwhelming the server during issues
// while still maintaining responsiveness for transient failures.
func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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

	return nil, fmt.Errorf("max retries (%d) exceeded, last error: %w", t.maxRetries, lastErr)
}

// isRetriableError determines if an error should trigger a retry attempt.
// The function implements a hierarchical error checking strategy:
//
// 1. First checks for nil errors (never retriable)
// 2. Then checks for net.OpError specifically because:
//   - It's the primary error type for network operations in Go
//   - Often wraps the actual error in its Err field
//   - Common in HTTP client operations for transient issues
//   - Used in our test suite to simulate network failures
//
// 3. Then checks for general net.Error types (timeouts, temporary errors)
// 4. Finally checks error messages against known retriable patterns
//
// This approach prioritizes network-specific errors while maintaining
// compatibility with other error types through message pattern matching.
func isRetriableError(err error) bool {
	if err == nil {
		return false
	}

	// Extract the error message for checking against retriable patterns
	errStr := err.Error()

	// Special handling for *net.OpError type
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		// Check the inner error first as it's more specific and often contains
		// the actual error message from the network operation
		if opErr.Err != nil {
			innerErrStr := opErr.Err.Error()
			for _, msg := range retriableErrors {
				if strings.Contains(innerErrStr, msg) {
					return true
				}
			}
		}
		// If no match in inner error, check the main error
		for _, msg := range retriableErrors {
			if strings.Contains(errStr, msg) {
				return true
			}
		}
		return false
	}

	// Check if it's a network error
	var netErr net.Error
	if errors.As(err, &netErr) {
		// Retry on timeout or temporary errors
		return netErr.Timeout() || netErr.Temporary()
	}

	// Check against known retriable error messages
	for _, msg := range retriableErrors {
		if strings.Contains(errStr, msg) {
			return true
		}
	}

	return false
}
