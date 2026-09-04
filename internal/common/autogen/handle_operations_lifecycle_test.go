package autogen_test

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

const (
	lifecycleResourceName = "mongodbatlas_ai_model_api_key.test"
	lifecycleProjectID    = "111111111111111111111111"
	// #nosec G101 -- test identifier, not a credential
	lifecycleAPIKeyID = "mocked-api-key-id"
)

// TestAPIErrorsDuringLifecycle drives a full Terraform lifecycle against canned Atlas responses to assert
// API errors that real Atlas does not produce on demand: every autogen delete endpoint answers 404 for a
// missing resource, never a 4xx/5xx. mongodbatlas_ai_model_api_key is used as it has no long-running wait.
func TestAPIErrorsDuringLifecycle(t *testing.T) {
	mock := newMockAtlas()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, mock),
		Steps: []resource.TestStep{
			{ // Create against a healthy API.
				Config: lifecycleConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(lifecycleResourceName, "api_key_id", lifecycleAPIKeyID),
					resource.TestCheckResourceAttr(lifecycleResourceName, "status", "ACTIVE"),
				),
			},
			{ // A read error must surface instead of being treated as a deleted resource.
				PreConfig:   func() { mock.setReadStatus(http.StatusInternalServerError) },
				Config:      lifecycleConfig(),
				ExpectError: regexp.MustCompile(`Error calling API in Read`),
			},
			{ // The failed read must have left the resource in state, so the plan is empty once the API recovers.
				PreConfig: func() { mock.setReadStatus(http.StatusOK) },
				Config:    lifecycleConfig(),
				PlanOnly:  true,
			},
			{ // A delete error must fail the destroy instead of reporting success.
				PreConfig:   func() { mock.failNextDelete(http.StatusConflict) },
				Config:      lifecycleConfig(),
				Destroy:     true,
				ExpectError: regexp.MustCompile(`Error calling API in Delete`),
			},
		},
	})
}

func lifecycleConfig() string {
	return fmt.Sprintf(`
		provider "mongodbatlas" {
			public_key  = "dummy-public-key"
			private_key = "dummy-private-key"
		}

		resource "mongodbatlas_ai_model_api_key" "test" {
			project_id = %[1]q
			cloud      = "ANY"
			geography  = "ANY"
			name       = "mocked-api-key"
		}
	`, lifecycleProjectID)
}

// mockAtlas answers the ai model API key endpoints, and lets a test step arm a failure for the next
// read or delete. It is both the HTTP client modifier expected by unit.TestAccProviderV6FactoriesWithMock
// and the round tripper it installs.
type mockAtlas struct {
	mu             sync.Mutex
	readStatus     int
	deleteStatus   int
	deleteFailures int
}

func newMockAtlas() *mockAtlas {
	return &mockAtlas{readStatus: http.StatusOK}
}

func (m *mockAtlas) setReadStatus(status int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.readStatus = status
}

// failNextDelete makes only the next delete fail, so the test cleanup destroy can still succeed.
func (m *mockAtlas) failNextDelete(status int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deleteStatus = status
	m.deleteFailures = 1
}

func (m *mockAtlas) ModifyHTTPClient(httpClient *http.Client) error {
	httpClient.Transport = m
	return nil
}

func (m *mockAtlas) ResetHTTPClient(_ *http.Client) {}

func (m *mockAtlas) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	switch req.Method {
	case http.MethodPost:
		return mockResponse(http.StatusCreated, apiKeyBody(true)), nil
	case http.MethodGet:
		if m.readStatus != http.StatusOK {
			return mockResponse(m.readStatus, atlasErrorBody(m.readStatus, "read failed")), nil
		}
		return mockResponse(http.StatusOK, apiKeyBody(false)), nil
	case http.MethodDelete:
		if m.deleteFailures > 0 {
			m.deleteFailures--
			return mockResponse(m.deleteStatus, atlasErrorBody(m.deleteStatus, "delete failed")), nil
		}
		return mockResponse(http.StatusNoContent, ""), nil
	default:
		return mockResponse(http.StatusMethodNotAllowed, atlasErrorBody(http.StatusMethodNotAllowed, "unexpected method")), nil
	}
}

func mockResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// apiKeyBody mirrors the Atlas response; the secret is only returned when the API key is created.
func apiKeyBody(withSecret bool) string {
	secret := ""
	if withSecret {
		// #nosec G101 -- fake value in a canned API response
		secret = `"secret":"mdb_ai_secret_value",`
	}
	return fmt.Sprintf(`{
		"apiKeyId": %[1]q,
		"cloud": "ANY",
		"geography": "ANY",
		"name": "mocked-api-key",
		"endpoint": "https://ai.example.com",
		"createdAt": "2026-01-01T00:00:00Z",
		"createdBy": "tester",
		"lastUsedAt": "2026-01-01T00:00:00Z",
		"maskedSecret": "mdb_ai_...value",
		%[2]s
		"status": "ACTIVE"
	}`, lifecycleAPIKeyID, secret)
}

func atlasErrorBody(status int, detail string) string {
	return fmt.Sprintf(`{"detail":%q,"error":%d,"errorCode":"TEST_ERROR","reason":%q}`, detail, status, http.StatusText(status))
}
