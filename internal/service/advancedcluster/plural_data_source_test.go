package advancedcluster_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
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
			result := advancedcluster.DiagsHasOnlyClusterNotFoundErrors(&tc.diags)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestResetClusterNotFoundErrors(t *testing.T) {
	tests := map[string]struct {
		input    diag.Diagnostics
		expected diag.Diagnostics
	}{
		"empty diagnostics": {
			input:    diag.Diagnostics{},
			expected: diag.Diagnostics{},
		},
		"only cluster not found errors": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error", "CLUSTER_NOT_FOUND"),
				diag.NewErrorDiagnostic("Error", "CLUSTER_NOT_FOUND"),
			},
			expected: diag.Diagnostics{},
		},
		"mixed errors": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error", "CLUSTER_NOT_FOUND"),
				diag.NewErrorDiagnostic("Other Error", "Some other error"),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic("Other Error", "Some other error"),
			},
		},
		"no cluster not found errors": {
			input: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error 1", "Some error"),
				diag.NewErrorDiagnostic("Error 2", "Another error"),
			},
			expected: diag.Diagnostics{
				diag.NewErrorDiagnostic("Error 1", "Some error"),
				diag.NewErrorDiagnostic("Error 2", "Another error"),
			},
		},
		"warnings with cluster not found": {
			input: diag.Diagnostics{
				diag.NewWarningDiagnostic("Warning", "Some warning"),
				diag.NewErrorDiagnostic("Error", "CLUSTER_NOT_FOUND"),
			},
			expected: diag.Diagnostics{
				diag.NewWarningDiagnostic("Warning", "Some warning"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := advancedcluster.ResetClusterNotFoundErrors(&tc.input)
			assert.Equal(t, tc.expected, *result)
		})
	}
}
