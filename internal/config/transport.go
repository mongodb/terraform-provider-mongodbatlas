package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// UserAgentExtra holds additional metadata to be appended to the User-Agent header and context.
type UserAgentExtra struct {
	Type           string // Type of the operation (e.g., "Resource", "Datasource", etc.)
	Name           string // Full name, for example mongodbatlas_database_user
	Operation      string // GrpcCall for example, ReadResource, see wrapped_provider_server.go for details
	ScriptLocation string // TODO: Support setting this field as opt-in on resources and datasources
}

// Combine returns a new UserAgentExtra by merging the receiver with another.
// Non-empty fields in 'other' take precedence over the receiver's fields.
func (e UserAgentExtra) Combine(other UserAgentExtra) UserAgentExtra {
	typeName := e.Type
	if other.Type != "" {
		typeName = other.Type
	}
	name := e.Name
	if other.Name != "" {
		name = other.Name
	}
	operation := e.Operation
	if other.Operation != "" {
		operation = other.Operation
	}
	scriptLocation := e.ScriptLocation
	if other.ScriptLocation != "" {
		scriptLocation = other.ScriptLocation
	}
	return UserAgentExtra{
		Type:           typeName,
		Name:           name,
		Operation:      operation,
		ScriptLocation: scriptLocation,
	}
}

// ToHeaderValue returns a string representation suitable for use as a User-Agent header value.
// If oldHeader is non-empty, it is prepended to the new value.
func (e UserAgentExtra) ToHeaderValue(oldHeader string) string {
	parts := []string{}
	addPart := func(key, part string) {
		if part == "" {
			return
		}
		parts = append(parts, fmt.Sprintf("%s/%s", key, part))
	}
	addPart("Type", e.Type)
	addPart("Name", e.Name)
	addPart("Operation", e.Operation)
	addPart("ScriptLocation", e.ScriptLocation)
	newPart := strings.Join(parts, " ")
	if oldHeader == "" {
		return newPart
	}
	return fmt.Sprintf("%s %s", oldHeader, newPart)
}

type UserAgentKey string

const (
	UserAgentExtraKey = UserAgentKey("user-agent-extra")
	UserAgentHeader   = "User-Agent"
)

// ReadUserAgentExtra retrieves the UserAgentExtra from the context if present.
// Returns a pointer to the UserAgentExtra, or nil if not set or of the wrong type.
// Logs a warning if the value is not of the expected type.
func ReadUserAgentExtra(ctx context.Context) *UserAgentExtra {
	extra := ctx.Value(UserAgentExtraKey)
	if extra == nil {
		return nil
	}
	if userAgentExtra, ok := extra.(UserAgentExtra); ok {
		return &userAgentExtra
	}
	log.Printf("[WARN] UserAgentExtra in context is not of type UserAgentExtra, got %v", extra)
	return nil
}

// AddUserAgentExtra returns a new context with UserAgentExtra merged into any existing value.
// If a UserAgentExtra is already present in the context, the fields of 'extra' will override non-empty fields.
func AddUserAgentExtra(ctx context.Context, extra UserAgentExtra) context.Context {
	oldExtra := ReadUserAgentExtra(ctx)
	if oldExtra == nil {
		return context.WithValue(ctx, UserAgentExtraKey, extra)
	}
	newExtra := oldExtra.Combine(extra)
	return context.WithValue(ctx, UserAgentExtraKey, newExtra)
}

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
	extra := ReadUserAgentExtra(req.Context())
	if extra != nil {
		userAgent := req.Header.Get(UserAgentHeader)
		newVar := extra.ToHeaderValue(userAgent)
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
