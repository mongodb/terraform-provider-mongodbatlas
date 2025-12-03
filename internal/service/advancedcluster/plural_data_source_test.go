package advancedcluster_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/assert"
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
