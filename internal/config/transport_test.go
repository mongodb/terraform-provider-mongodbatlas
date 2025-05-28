package config_test

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
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
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got: %d", resp.StatusCode)
	}

	logStr := logOutput.String()
	if !strings.Contains(logStr, "Test Service Network Request Start") {
		t.Error("Expected start log message not found")
	}
	if !strings.Contains(logStr, "Test Service Network Request Complete") {
		t.Error("Expected completion log message not found")
	}
	if !strings.Contains(logStr, "Status: 200 (Success)") {
		t.Error("Expected success status log not found")
	}
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
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if resp.StatusCode != 500 {
		t.Fatalf("Expected status 500, got: %d", resp.StatusCode)
	}
	logStr := logOutput.String()
	if !strings.Contains(logStr, "Test Service Network Request Start") {
		t.Error("Expected start log message not found")
	}
	if !strings.Contains(logStr, "Test Service Network Request Complete") {
		t.Error("Expected completion log message not found")
	}
	if !strings.Contains(logStr, "Status: 500 (Server Error)") {
		t.Error("Expected server error status log not found")
	}
	if !strings.Contains(logStr, "HTTP Error Response") {
		t.Error("Expected HTTP error response log not found")
	}
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
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err != networkErr {
		t.Fatalf("Expected original error, got: %v", err)
	}
	if resp != nil {
		t.Fatal("Expected nil response on error")
	}
	logStr := logOutput.String()
	if !strings.Contains(logStr, "Test Service Network Request Start") {
		t.Error("Expected start log message not found")
	}
	if !strings.Contains(logStr, "Test Service Network Request Failed") {
		t.Error("Expected failure log message not found")
	}
	if !strings.Contains(logStr, "Network Timeout") {
		t.Error("Expected timeout context log not found")
	}
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
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client, ok := clientInterface.(*config.MongoDBClient)
	if !ok {
		t.Fatal("Failed to cast client to MongoDBClient")
	}

	// Make a simple API call that should trigger our enhanced logging
	// We'll try to list organizations, which is a basic read operation
	_, resp, err := client.AtlasV2.OrganizationsApi.ListOrganizations(t.Context()).Execute()

	// We don't care if the API call fails (could be due to permissions, etc.)
	// We just want to verify that our logging is working
	_ = resp // Ignore response
	_ = err  // Ignore error

	logStr := logOutput.String()
	if !strings.Contains(logStr, "MongoDB Atlas Network Request Start") {
		t.Error("Expected to find 'MongoDB Atlas Network Request Start' in logs")
	}

	hasCompletion := strings.Contains(logStr, "MongoDB Atlas Network Request Complete")
	if hasCompletion && !strings.Contains(logStr, "Duration:") {
		t.Error("Expected to find duration information in completion logs")
	}
}
