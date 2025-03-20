package config_test

import (
	"errors"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			{resp: &http.Response{StatusCode: 200}, err: nil},
		},
	}

	transport := config.NewRetryTransport(mock, 3, time.Millisecond)
	req, _ := http.NewRequestWithContext(t.Context(), "GET", "http://example.com", http.NoBody)
	resp, err := transport.RoundTrip(req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, mock.current)
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

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 3, mock.current)
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

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, 4, mock.current)
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

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, 1, mock.current)
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

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.GreaterOrEqual(t, duration, 3*time.Millisecond) // Sum of mock delays
}
