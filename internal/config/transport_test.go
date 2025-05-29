package config_test

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTransport struct {
	response *http.Response
	err      error
	delay    time.Duration
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	return m.response, m.err
}

func TestNetworkLoggingTransport_Success(t *testing.T) {
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)
	mockResp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
	}
	mockTransport := &mockTransport{
		response: mockResp,
		err:      nil,
	}
	transport := config.NewTransportWithNetworkLogging(mockTransport, true)
	req := httptest.NewRequest("GET", "https://api.example.com/test", http.NoBody)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	logStr := logOutput.String()
	assert.Contains(t, logStr, "Network Request Start")
	assert.Contains(t, logStr, "Network Request Complete")
	assert.Contains(t, logStr, "Status: 200 (Success)")
}

func TestNetworkLoggingTransport_HTTPError(t *testing.T) {
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)
	mockResp := &http.Response{
		StatusCode: 500,
		Header:     make(http.Header),
	}
	mockTransport := &mockTransport{
		response: mockResp,
		err:      nil,
	}
	transport := config.NewTransportWithNetworkLogging(mockTransport, true)
	req := httptest.NewRequest("POST", "https://api.example.com/test", http.NoBody)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, 500, resp.StatusCode)

	logStr := logOutput.String()
	assert.Contains(t, logStr, "Network Request Start")
	assert.Contains(t, logStr, "Network Request Complete")
	assert.Contains(t, logStr, "Status: 500 (Server Error)")
	assert.Contains(t, logStr, "HTTP Error Response")
}

func TestNetworkLoggingTransport_NetworkError(t *testing.T) {
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)
	networkErr := errors.New("connection timeout")
	mockTransport := &mockTransport{
		response: nil,
		err:      networkErr,
	}
	transport := config.NewTransportWithNetworkLogging(mockTransport, true)
	req := httptest.NewRequest("GET", "https://api.example.com/test", http.NoBody)
	resp, err := transport.RoundTrip(req)
	require.Error(t, err)
	require.Equal(t, networkErr, err)
	require.Nil(t, resp)

	logStr := logOutput.String()
	assert.Contains(t, logStr, "Network Request Start")
	assert.Contains(t, logStr, "Network Request Failed")
	assert.Contains(t, logStr, "Network Timeout")
}

func TestNetworkLoggingTransport_DigestAuthChallenge(t *testing.T) {
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)
	mockResp := &http.Response{
		StatusCode: 401,
		Header:     make(http.Header),
	}
	mockResp.Header.Set("WWW-Authenticate", "Digest realm=\"MongoDB Atlas\", nonce=\"abc123\"")
	mockTransport := &mockTransport{
		response: mockResp,
		err:      nil,
	}
	transport := config.NewTransportWithNetworkLogging(mockTransport, true)
	req := httptest.NewRequest("GET", "https://cloud.mongodb.com/api/atlas/v2/groups", http.NoBody)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)

	logStr := logOutput.String()
	assert.Contains(t, logStr, "Network Request Start")
	assert.Contains(t, logStr, "Network Request Complete")
	assert.Contains(t, logStr, "Status: 401 (Client Error)")
	assert.Contains(t, logStr, "Digest Authentication Challenge")
	assert.Contains(t, logStr, "Expected first request in digest authentication flow")
	// Should NOT contain the generic HTTP Error Response for 401
	assert.NotContains(t, logStr, "HTTP Error Response")
}

func TestNetworkLoggingTransport_Disabled(t *testing.T) {
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)
	mockResp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
	}
	mockTransport := &mockTransport{
		response: mockResp,
		err:      nil,
	}
	transport := config.NewTransportWithNetworkLogging(mockTransport, false)
	req := httptest.NewRequest("GET", "https://api.example.com/test", http.NoBody)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	logStr := logOutput.String()
	assert.Empty(t, logStr, "Expected no logs when network logging is disabled")
}

func TestAccNetworkLogging(t *testing.T) {
	acc.SkipInUnitTest(t)
	acc.PreCheckBasic(t)

	t.Setenv("TF_LOG", "DEBUG") // Enable debug logging for the test.
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)
	cfg := &config.Config{
		PublicKey:  os.Getenv("MONGODB_ATLAS_PUBLIC_KEY"),
		PrivateKey: os.Getenv("MONGODB_ATLAS_PRIVATE_KEY"),
		BaseURL:    os.Getenv("MONGODB_ATLAS_BASE_URL"),
	}
	clientInterface, err := cfg.NewClient(t.Context())
	require.NoError(t, err)
	client, ok := clientInterface.(*config.MongoDBClient)
	require.True(t, ok)

	// Make a simple API call that should trigger our enhanced logging
	_, _, err = client.AtlasV2.OrganizationsApi.ListOrganizations(t.Context()).Execute()
	require.NoError(t, err)
	logStr := logOutput.String()
	assert.Contains(t, logStr, "Network Request Start")
	assert.Contains(t, logStr, "Network Request Complete")
	assert.Contains(t, logStr, "Duration:")
}
