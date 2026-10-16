package advancedcluster_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312026/admin"
)

func TestRemoveClusterNotFoundErrors(t *testing.T) {
	tests := map[string]struct {
		input    diag.Diagnostics
		expected diag.Diagnostics
	}{
		"empty diagnostics": {
			input:    diag.Diagnostics{},
			expected: diag.Diagnostics{},
		},
		"single cluster not found": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Cluster Not Found", "CLUSTER_NOT_FOUND"),
			},
			expected: diag.Diagnostics{},
		},
		"only cluster not found errors": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error", "CLUSTER_NOT_FOUND"),
				diag.NewErrorDiagnostic("Error", "CLUSTER_NOT_FOUND"),
			},
			expected: diag.Diagnostics{},
		},
		"mixed errors with cluster not found": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Other Error", "Some other error"),
				diag.NewErrorDiagnostic("Cluster Not Found", "CLUSTER_NOT_FOUND"),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic("Other Error", "Some other error"),
			},
		},
		"other errors only": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error 1", "Some error"),
				diag.NewErrorDiagnostic("Error 2", "Another error"),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error 1", "Some error"),
				diag.NewErrorDiagnostic("Error 2", "Another error"),
			},
		},
		"warnings with cluster not found error": {
			input: diag.Diagnostics{
				diag.NewWarningDiagnostic("Warning", "Some warning"),
				diag.NewErrorDiagnostic("Cluster Not Found", "CLUSTER_NOT_FOUND"),
			},
			expected: diag.Diagnostics{
				diag.NewWarningDiagnostic("Warning", "Some warning"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			advancedcluster.RemoveClusterNotFoundErrors(&tc.input)
			assert.Equal(t, tc.expected, tc.input)
		})
	}
}

func TestRetryListOnClusterNotFound(t *testing.T) {
	t.Run("retries transient CLUSTER_NOT_FOUND", func(t *testing.T) {
		calls := 0
		result, err := advancedcluster.RetryListOnClusterNotFound(t.Context(), func() (string, error) {
			calls++
			if calls < 3 {
				return "", clusterNotFoundError()
			}
			return "clusters", nil
		})
		require.NoError(t, err)
		require.Equal(t, "clusters", result)
		require.Equal(t, 3, calls)
	})

	t.Run("does not retry other errors", func(t *testing.T) {
		calls := 0
		_, err := advancedcluster.RetryListOnClusterNotFound(t.Context(), func() (string, error) {
			calls++
			return "", errors.New("unexpected error")
		})
		require.Error(t, err)
		require.Equal(t, 1, calls)
	})
}

// clusterNotFoundError builds the SDK error Atlas returns for CLUSTER_NOT_FOUND.
func clusterNotFoundError() error {
	apiErr := &admin.GenericOpenAPIError{}
	apiErr.SetModel(admin.ApiError{ErrorCode: "CLUSTER_NOT_FOUND"})
	apiErr.SetError("cluster not found")
	return apiErr
}

func TestPluralDataSourceEmptyProjectReturnsEmptyList(t *testing.T) {
	storageConfigCredentials(t)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { require.NoError(t, unit.MockConfigAdvancedCluster.RunBeforeEach()) },
		ProtoV6ProviderFactories: unit.TestAccProviderV6FactoriesWithMock(t, emptyProjectHTTPMock{}),
		Steps: []resource.TestStep{{
			Config: pluralDataSourceConfig,
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(dataSourcePluralName, tfjsonpath.New("results"), knownvalue.ListExact([]knownvalue.Check{})),
			},
		}},
	})
}

// emptyProjectHTTPMock returns an empty cluster list for every request.
type emptyProjectHTTPMock struct{}

func (emptyProjectHTTPMock) ModifyHTTPClient(c *http.Client) error {
	c.Transport = emptyProjectHTTPMock{}
	return nil
}

func (emptyProjectHTTPMock) ResetHTTPClient(*http.Client) {}

func (emptyProjectHTTPMock) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"results":[],"totalCount":0}`)),
		Request:    req,
	}, nil
}

const pluralDataSourceConfig = `
provider "mongodbatlas" {
  public_key  = "test-public-key"
  private_key = "test-private-key"
  base_url    = "https://atlas.invalid/"
}
data "mongodbatlas_advanced_clusters" "test" {
  project_id = "111111111111111111111111"
}
`
