package config_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// mockTransport implements http.RoundTripper for testing
type mockTransport struct {
	delays    []time.Duration // to test backoff timing
	responses []response
	current   int
}

type response struct {
	resp *http.Response
	err  error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.current >= len(m.responses) {
		return nil, errors.New("no more responses")
	}

	if m.delays != nil && m.current < len(m.delays) {
		time.Sleep(m.delays[m.current])
	}

	resp := m.responses[m.current]
	m.current++
	return resp.resp, resp.err
}

func TestRetryTransport_NoRetryNeeded(t *testing.T) {
	mock := &mockTransport{
		responses: []response{
			{
				resp: &http.Response{StatusCode: 200},
				err:  nil,
			},
		},
	}

	transport := config.NewRetryTransport(mock, 3, time.Millisecond)
	req, _ := http.NewRequestWithContext(t.Context(), "GET", "http://example.com", http.NoBody)
	resp, err := transport.RoundTrip(req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Error("expected response, got nil")
	}
	if mock.current != 1 {
		t.Errorf("expected 1 attempt, got %d", mock.current)
	}
}

func TestRetryTransport_SuccessAfterRetries(t *testing.T) {
	mock := &mockTransport{
		responses: []response{
			{resp: nil, err: &net.OpError{Err: errors.New("connection reset by peer")}},
			{resp: nil, err: &net.OpError{Err: errors.New("connection reset by peer")}},
			{resp: &http.Response{StatusCode: 200}, err: nil},
		},
	}

	transport := config.NewRetryTransport(mock, 3, time.Millisecond)
	req, _ := http.NewRequestWithContext(t.Context(), "GET", "http://example.com", http.NoBody)
	resp, err := transport.RoundTrip(req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Error("expected response, got nil")
	}
	if mock.current != 3 {
		t.Errorf("expected 3 attempts, got %d", mock.current)
	}
}

func TestRetryTransport_MaxRetriesExceeded(t *testing.T) {
	mock := &mockTransport{
		responses: []response{
			{resp: nil, err: &net.OpError{Err: errors.New("connection reset by peer")}},
			{resp: nil, err: &net.OpError{Err: errors.New("connection reset by peer")}},
			{resp: nil, err: &net.OpError{Err: errors.New("connection reset by peer")}},
			{resp: nil, err: &net.OpError{Err: errors.New("connection reset by peer")}},
		},
	}

	transport := config.NewRetryTransport(mock, 3, time.Millisecond)
	req, _ := http.NewRequestWithContext(t.Context(), "GET", "http://example.com", http.NoBody)
	resp, err := transport.RoundTrip(req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if resp != nil {
		t.Error("expected nil response, got response")
	}
	if mock.current != 4 { // initial attempt + 3 retries
		t.Errorf("expected 4 attempts, got %d", mock.current)
	}
}

func TestRetryTransport_NonRetriableError(t *testing.T) {
	mock := &mockTransport{
		responses: []response{
			{resp: nil, err: errors.New("non-retriable error")},
		},
	}

	transport := config.NewRetryTransport(mock, 3, time.Millisecond)
	req, _ := http.NewRequestWithContext(t.Context(), "GET", "http://example.com", http.NoBody)
	resp, err := transport.RoundTrip(req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if resp != nil {
		t.Error("expected nil response, got response")
	}
	if mock.current != 1 {
		t.Errorf("expected 1 attempt, got %d", mock.current)
	}
}

func TestRetryTransport_ContextCancellation(t *testing.T) {
	mock := &mockTransport{
		responses: []response{
			{resp: nil, err: &net.OpError{Err: errors.New("connection reset by peer")}},
			{resp: nil, err: &net.OpError{Err: errors.New("connection reset by peer")}},
		},
	}

	ctx, cancel := context.WithCancel(t.Context())
	transport := config.NewRetryTransport(mock, 3, 100*time.Millisecond)

	// Cancel after first attempt
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	req, _ := http.NewRequestWithContext(ctx, "GET", "http://example.com", http.NoBody)
	resp, err := transport.RoundTrip(req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if resp != nil {
		t.Error("expected nil response, got response")
	}
	if mock.current > 2 {
		t.Errorf("expected <= 2 attempts, got %d", mock.current)
	}
}

func TestRetryTransport_ExponentialBackoff(t *testing.T) {
	mock := &mockTransport{
		responses: []response{
			{resp: nil, err: &net.OpError{Err: errors.New("connection reset by peer")}},
			{resp: nil, err: &net.OpError{Err: errors.New("connection reset by peer")}},
			{resp: &http.Response{StatusCode: 200}, err: nil},
		},
		delays: []time.Duration{
			0 * time.Millisecond, // First attempt
			1 * time.Millisecond, // First retry
			2 * time.Millisecond, // Second retry
		},
	}

	start := time.Now()
	transport := config.NewRetryTransport(mock, 3, time.Millisecond)
	req, _ := http.NewRequestWithContext(t.Context(), "GET", "http://example.com", http.NoBody)
	resp, err := transport.RoundTrip(req)

	duration := time.Since(start)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Error("expected response, got nil")
	}

	// Check that the total time is at least the sum of delays plus backoff times
	expectedMinDuration := 3 * time.Millisecond // Sum of mock delays
	if duration < expectedMinDuration {
		t.Errorf("expected duration >= %v, got %v", expectedMinDuration, duration)
	}
}
