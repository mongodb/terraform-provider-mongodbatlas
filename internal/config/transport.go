package config

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// NetworkLoggingTransport wraps an http.RoundTripper to provide enhanced logging
// for network operations, including timing, status codes, and error details.
type NetworkLoggingTransport struct {
	// Transport is the underlying transport to wrap
	Transport http.RoundTripper
	// Name is used to identify the transport in logs
	Name string
}

// NewNetworkLoggingTransport creates a new NetworkLoggingTransport that wraps
// the provided transport with enhanced network logging capabilities.
func NewNetworkLoggingTransport(name string, transport http.RoundTripper) *NetworkLoggingTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &NetworkLoggingTransport{
		Transport: transport,
		Name:      name,
	}
}

// RoundTrip implements the http.RoundTripper interface and adds enhanced logging
// around the HTTP request/response cycle.
func (t *NetworkLoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	startTime := time.Now()
	log.Printf("[DEBUG] %s Network Request Start: %s %s (started at %s)",
		t.Name, req.Method, req.URL.String(), startTime.Format(time.RFC3339Nano))
	resp, err := t.Transport.RoundTrip(req)
	duration := time.Since(startTime)
	if err != nil {
		log.Printf("[ERROR] %s Network Request Failed: %s %s - Duration: %v - Error: %v",
			t.Name, req.Method, req.URL.String(), duration, err)

		t.logNetworkErrorContext(err, req, duration)
		return resp, err
	}

	statusCode := resp.StatusCode
	statusClass := GetStatusClass(statusCode)

	log.Printf("[DEBUG] %s Network Request Complete: %s %s - Status: %d (%s) - Duration: %v",
		t.Name, req.Method, req.URL.String(), statusCode, statusClass, duration)

	if statusCode >= 300 {
		log.Printf("[WARN] %s HTTP Error Response: %s %s - Status: %d %s - Duration: %v - Content-Type: %s",
			t.Name, req.Method, req.URL.String(), statusCode, http.StatusText(statusCode),
			duration, resp.Header.Get("Content-Type"))
	}
	return resp, nil
}

// logNetworkErrorContext provides additional context for common network errors
func (t *NetworkLoggingTransport) logNetworkErrorContext(err error, req *http.Request, duration time.Duration) {
	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "timeout"):
		log.Printf("[ERROR] %s Network Timeout: %s %s - Duration: %v - This may indicate API server overload or network connectivity issues",
			t.Name, req.Method, req.URL.String(), duration)
	case strings.Contains(errStr, "connection refused"):
		log.Printf("[ERROR] %s Connection Refused: %s %s - Duration: %v - API server may be down or unreachable",
			t.Name, req.Method, req.URL.String(), duration)
	case strings.Contains(errStr, "no such host"):
		log.Printf("[ERROR] %s DNS Resolution Failed: %s %s - Duration: %v - Check DNS configuration and network connectivity",
			t.Name, req.Method, req.URL.String(), duration)
	case strings.Contains(errStr, "certificate"):
		log.Printf("[ERROR] %s TLS Certificate Error: %s %s - Duration: %v - Check certificate validity and trust chain",
			t.Name, req.Method, req.URL.String(), duration)
	case strings.Contains(errStr, "context deadline exceeded"):
		log.Printf("[ERROR] %s Request Deadline Exceeded: %s %s - Duration: %v - Request took longer than configured timeout",
			t.Name, req.Method, req.URL.String(), duration)
	case strings.Contains(errStr, "connection reset"):
		log.Printf("[ERROR] %s Connection Reset: %s %s - Duration: %v - Server closed connection unexpectedly",
			t.Name, req.Method, req.URL.String(), duration)
	default:
		log.Printf("[ERROR] %s Network Error: %s %s - Duration: %v - Error details: %v",
			t.Name, req.Method, req.URL.String(), duration, err)
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
