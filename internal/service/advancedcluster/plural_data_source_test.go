package advancedcluster_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
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

	t.Run("returns the zero value when cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		result, err := advancedcluster.RetryListOnClusterNotFound(ctx, func() (string, error) {
			cancel()
			return "partial", clusterNotFoundError()
		})
		require.Error(t, err)
		require.Empty(t, result)
	})
}

// clusterNotFoundError builds the SDK error Atlas returns for CLUSTER_NOT_FOUND.
func clusterNotFoundError() error {
	apiErr := &admin.GenericOpenAPIError{}
	apiErr.SetModel(admin.ApiError{ErrorCode: "CLUSTER_NOT_FOUND"})
	apiErr.SetError("cluster not found")
	return apiErr
}
