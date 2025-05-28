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
	mockResp.Header.Set("Content-Type", "application/json")
	mockResp.Header.Set("X-Request-Id", "test-request-123")

	mockTransport := &mockTransport{
		response: mockResp,
		err:      nil,
	}
	transport := config.NewNetworkLoggingTransport("Test Service", mockTransport)
	req := httptest.NewRequest("GET", "https://api.example.com/test", http.NoBody)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)

	logStr := logOutput.String()
	assert.Contains(t, logStr, "Test Service Network Request Start")
	assert.Contains(t, logStr, "Test Service Network Request Complete")
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
	mockResp.Header.Set("Content-Type", "application/json")
	mockResp.Header.Set("X-Request-Id", "test-request-456")

	mockTransport := &mockTransport{
		response: mockResp,
		err:      nil,
	}
	transport := config.NewNetworkLoggingTransport("Test Service", mockTransport)
	req := httptest.NewRequest("POST", "https://api.example.com/test", http.NoBody)
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)
	require.Equal(t, 500, resp.StatusCode)

	logStr := logOutput.String()
	assert.Contains(t, logStr, "Test Service Network Request Start")
	assert.Contains(t, logStr, "Test Service Network Request Complete")
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

	transport := config.NewNetworkLoggingTransport("Test Service", mockTransport)
	req := httptest.NewRequest("GET", "https://api.example.com/test", http.NoBody)
	resp, err := transport.RoundTrip(req)
	require.Error(t, err)
	require.Equal(t, networkErr, err)
	require.Nil(t, resp)

	logStr := logOutput.String()
	assert.Contains(t, logStr, "Test Service Network Request Start")
	assert.Contains(t, logStr, "Test Service Network Request Failed")
	assert.Contains(t, logStr, "Network Timeout")
}

func TestAccNetworkLogging(t *testing.T) {
	acc.SkipInUnitTest(t)
	acc.PreCheckBasic(t)

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
	assert.Contains(t, logStr, "Atlas Network Request Start")
	assert.Contains(t, logStr, "Atlas Network Request Complete")
	assert.Contains(t, logStr, "Duration:")
}
