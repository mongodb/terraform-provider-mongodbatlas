package advancedclustertpf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/stretchr/testify/assert"
)

func TestDiagsHasOnlyClusterNotFound(t *testing.T) {
	tests := map[string]struct {
		diags    diag.Diagnostics
		expected bool
	}{
		"empty diagnostics": {
			diags:    diag.Diagnostics{},
			expected: true,
		},
		"single cluster not found": {
			diags: diag.Diagnostics{
				diag.NewErrorDiagnostic("Cluster Not Found", "CLUSTER_NOT_FOUND"),
			},
			expected: true,
		},
		"multiple errors with cluster not found": {
			diags: diag.Diagnostics{
				diag.NewErrorDiagnostic("Other Error", "Some other error"),
				diag.NewErrorDiagnostic("Cluster Not Found", "CLUSTER_NOT_FOUND"),
			},
			expected: false,
		},
		"other errors only": {
			diags: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error 1", "Some error"),
				diag.NewErrorDiagnostic("Error 2", "Another error"),
			},
			expected: false,
		},
		"warnings with cluster not found error": {
			diags: diag.Diagnostics{
				diag.NewWarningDiagnostic("Warning", "Some warning"),
				diag.NewErrorDiagnostic("Cluster Not Found", "CLUSTER_NOT_FOUND"),
			},
			expected: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := advancedclustertpf.DiagsHasOnlyClusterNotFound(&tc.diags)
			assert.Equal(t, tc.expected, result)
		})
	}
}
