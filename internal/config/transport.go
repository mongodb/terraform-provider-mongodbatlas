package config

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// UserAgentTransport wraps an http.RoundTripper to add User-Agent header with additional metadata.
type UserAgentTransport struct {
	Transport http.RoundTripper
	Enabled   bool
}

func NewUserAgentTransport(transport http.RoundTripper, enabled bool) *UserAgentTransport {
	return &UserAgentTransport{
		Transport: transport,
		Enabled:   enabled,
	}
}

func (t *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !t.Enabled {
		return t.Transport.RoundTrip(req)
	}
	ctx := req.Context()
	extra := ReadUserAgentExtra(ctx)
	if extra != nil {
		userAgent := req.Header.Get(UserAgentHeader)
		newVar := extra.ToHeaderValue(ctx, userAgent)
		req.Header.Set(UserAgentHeader, newVar)
	}
	resp, err := t.Transport.RoundTrip(req)
	return resp, err
}

// NetworkLoggingTransport wraps an http.RoundTripper to provide enhanced logging
// for network operations, including timing, status codes, and error details.
type NetworkLoggingTransport struct {
	Transport http.RoundTripper
	Enabled   bool
}

// NewTransportWithNetworkLogging creates a new NetworkLoggingTransport that wraps
// the provided transport with enhanced network logging capabilities.
func NewTransportWithNetworkLogging(transport http.RoundTripper, enabled bool) *NetworkLoggingTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &NetworkLoggingTransport{
		Transport: transport,
		Enabled:   enabled,
	}
}

// RoundTrip implements the http.RoundTripper interface and adds enhanced logging
// around the HTTP request/response cycle.
func (t *NetworkLoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !t.Enabled {
		return t.Transport.RoundTrip(req)
	}

	startTime := time.Now()
	log.Printf("[DEBUG] Network Request Start: %s %s (started at %s)",
		req.Method, req.URL.String(), startTime.Format(time.RFC3339Nano))

	resp, err := t.Transport.RoundTrip(req)
	duration := time.Since(startTime)
	if err != nil {
		log.Printf("[ERROR] Network Request Failed: %s %s - Duration: %v - Error: %v",
			req.Method, req.URL.String(), duration, err)

		t.logNetworkErrorContext(err, req, duration)
		return resp, err
	}
	statusCode := resp.StatusCode
	statusClass := GetStatusClass(statusCode)

	log.Printf("[DEBUG] Network Request Complete: %s %s - Status: %d (%s) - Duration: %v",
		req.Method, req.URL.String(), statusCode, statusClass, duration)

	if statusCode == http.StatusUnauthorized {
		log.Printf("[DEBUG] Digest Authentication Challenge: %s %s - Status: 401 - Expected first request in digest authentication flow",
			req.Method, req.URL.String())
	} else if statusCode >= 300 {
		log.Printf("[WARN] HTTP Error Response: %s %s - Status: %d %s - Duration: %v - Content-Type: %s",
			req.Method, req.URL.String(), statusCode, http.StatusText(statusCode),
			duration, resp.Header.Get("Content-Type"))
	}
	return resp, nil
}

// logNetworkErrorContext provides additional context for common network errors
func (t *NetworkLoggingTransport) logNetworkErrorContext(err error, req *http.Request, duration time.Duration) {
	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "timeout"):
		log.Printf("[ERROR] Network Timeout: %s %s - Duration: %v - This may indicate API server overload or network connectivity issues",
			req.Method, req.URL.String(), duration)
	case strings.Contains(errStr, "connection refused"):
		log.Printf("[ERROR] Connection Refused: %s %s - Duration: %v - API server may be down or unreachable",
			req.Method, req.URL.String(), duration)
	case strings.Contains(errStr, "no such host"):
		log.Printf("[ERROR] DNS Resolution Failed: %s %s - Duration: %v - Check DNS configuration and network connectivity",
			req.Method, req.URL.String(), duration)
	case strings.Contains(errStr, "certificate"):
		log.Printf("[ERROR] TLS Certificate Error: %s %s - Duration: %v - Check certificate validity and trust chain",
			req.Method, req.URL.String(), duration)
	case strings.Contains(errStr, "context deadline exceeded"):
		log.Printf("[ERROR] Request Deadline Exceeded: %s %s - Duration: %v - Request took longer than configured timeout",
			req.Method, req.URL.String(), duration)
	case strings.Contains(errStr, "connection reset"):
		log.Printf("[ERROR] Connection Reset: %s %s - Duration: %v - Server closed connection unexpectedly",
			req.Method, req.URL.String(), duration)
	default:
		log.Printf("[ERROR] Network Error: %s %s - Duration: %v - Error details: %v",
			req.Method, req.URL.String(), duration, err)
	}
}

// GetStatusClass returns a human-readable status class for the HTTP status code
func GetStatusClass(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "Success"
	case statusCode >= 300 && statusCode < 400:
		return "Redirection"
	case statusCode >= 400 && statusCode < 500:
		return "Client Error"
	case statusCode >= 500 && statusCode < 600:
		return "Server Error"
	default:
		return "Unknown"
	}
}
