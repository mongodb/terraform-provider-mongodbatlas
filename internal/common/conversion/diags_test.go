package conversion_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	sdkv2diag "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
)

func TestFromTPFDiagsToSDKV2Diags(t *testing.T) {
	tests := []struct {
		opts           *conversion.DiagsOptions
		name           string
		inputDiags     []diag.Diagnostic
		expectedOutput sdkv2diag.Diagnostics
	}{
		{
			name:           "Nil slice",
			inputDiags:     nil,
			expectedOutput: nil,
		},
		{
			name:           "Empty slice",
			inputDiags:     []diag.Diagnostic{},
			expectedOutput: nil,
		},
		{
			name: "Single error diagnostic",
			inputDiags: []diag.Diagnostic{
				diag.NewErrorDiagnostic("Error summary", "Error detail"),
			},
			expectedOutput: []sdkv2diag.Diagnostic{
				{
					Severity: sdkv2diag.Error,
					Summary:  "Error summary",
					Detail:   "Error detail",
				},
			},
		},
		{
			name: "Mixed error and warning diagnostics",
			inputDiags: []diag.Diagnostic{
				diag.NewErrorDiagnostic("Error summary", "Error detail"),
				diag.NewWarningDiagnostic("Warning summary", "Warning detail"),
			},
			expectedOutput: []sdkv2diag.Diagnostic{
				{
					Severity: sdkv2diag.Error,
					Summary:  "Error summary",
					Detail:   "Error detail",
				},
				{
					Severity: sdkv2diag.Warning,
					Summary:  "Warning summary",
					Detail:   "Warning detail",
				},
			},
		},
		{
			name: "With options",
			inputDiags: []diag.Diagnostic{
				diag.NewErrorDiagnostic("Error summary", "Error detail"),
			},
			expectedOutput: []sdkv2diag.Diagnostic{
				{
					Severity: sdkv2diag.Error,
					Summary:  "Custom summary",
					Detail:   "Error detail",
				},
			},
			opts: &conversion.DiagsOptions{
				Summary: "Custom summary",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var result sdkv2diag.Diagnostics
			if tc.opts != nil {
				result = conversion.FromTPFDiagsToSDKV2Diags(tc.inputDiags, *tc.opts)
			} else {
				result = conversion.FromTPFDiagsToSDKV2Diags(tc.inputDiags)
			}
			assert.Equal(t, tc.expectedOutput, result)
		})
	}
}
