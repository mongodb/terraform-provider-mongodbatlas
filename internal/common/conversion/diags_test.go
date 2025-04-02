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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := conversion.FromTPFDiagsToSDKV2Diags(tc.inputDiags)
			assert.Equal(t, tc.expectedOutput, result)
		})
	}
}
func TestFormatDiags(t *testing.T) {
	testCases := map[string]struct {
		setupDiags   func() *diag.Diagnostics
		expectedText string
	}{
		"nil diagnostics": {
			setupDiags: func() *diag.Diagnostics {
				return nil
			},
			expectedText: "",
		},
		"empty diagnostics": {
			setupDiags: func() *diag.Diagnostics {
				var diags diag.Diagnostics
				return &diags
			},
			expectedText: "",
		},
		"single error": {
			setupDiags: func() *diag.Diagnostics {
				var diags diag.Diagnostics
				diags.AddError("Error summary", "Error detail")
				return &diags
			},
			expectedText: "Error summary\n\t detail: Error detail",
		},
		"multiple errors": {
			setupDiags: func() *diag.Diagnostics {
				var diags diag.Diagnostics
				diags.AddError("First error", "Error detail 1")
				diags.AddError("Second error", "Error detail 2")
				return &diags
			},
			expectedText: "First error\n\t detail: Error detail 1\nSecond error\n\t detail: Error detail 2",
		},
		"single warning": {
			setupDiags: func() *diag.Diagnostics {
				var diags diag.Diagnostics
				diags.AddWarning("Warning summary", "Warning detail")
				return &diags
			},
			expectedText: "Warnings:\nWarning summary\n\t detail: Warning detail",
		},
		"multiple warnings": {
			setupDiags: func() *diag.Diagnostics {
				var diags diag.Diagnostics
				diags.AddWarning("First warning", "Warning detail 1")
				diags.AddWarning("Second warning", "Warning detail 2")
				return &diags
			},
			expectedText: "Warnings:\nFirst warning\n\t detail: Warning detail 1\nSecond warning\n\t detail: Warning detail 2",
		},
		"errors and warnings": {
			setupDiags: func() *diag.Diagnostics {
				var diags diag.Diagnostics
				diags.AddError("Error summary", "Error detail")
				diags.AddWarning("Warning summary", "Warning detail")
				return &diags
			},
			expectedText: "Error summary\n\t detail: Error detail\n\nWarnings:\nWarning summary\n\t detail: Warning detail",
		},
		"multiple errors and warnings": {
			setupDiags: func() *diag.Diagnostics {
				var diags diag.Diagnostics
				diags.AddError("First error", "Error detail 1")
				diags.AddError("Second error", "Error detail 2")
				diags.AddWarning("First warning", "Warning detail 1")
				diags.AddWarning("Second warning", "Warning detail 2")
				return &diags
			},
			expectedText: "First error\n\t detail: Error detail 1\nSecond error\n\t detail: Error detail 2\n\nWarnings:\nFirst warning\n\t detail: Warning detail 1\nSecond warning\n\t detail: Warning detail 2",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			diags := tc.setupDiags()
			result := conversion.FormatDiags(diags)
			assert.Equal(t, tc.expectedText, result)
		})
	}
}
