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

func TestSDKV2DiagsToError(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr string
		diags       sdkv2diag.Diagnostics
	}{
		{
			name:        "empty diagnostics",
			diags:       sdkv2diag.Diagnostics{},
			expectedErr: "",
		},
		{
			name: "warnings only returns nil",
			diags: sdkv2diag.Diagnostics{
				{Severity: sdkv2diag.Warning, Summary: "just a warning"},
			},
			expectedErr: "",
		},
		{
			name: "single error without detail",
			diags: sdkv2diag.Diagnostics{
				{Severity: sdkv2diag.Error, Summary: "only error"},
			},
			expectedErr: "only error",
		},
		{
			name: "single error with detail",
			diags: sdkv2diag.Diagnostics{
				{Severity: sdkv2diag.Error, Summary: "summary", Detail: "detail info"},
			},
			expectedErr: "summary: detail info",
		},
		{
			name: "mixed warnings and errors joins error entries",
			diags: sdkv2diag.Diagnostics{
				{Severity: sdkv2diag.Warning, Summary: "just a warning"},
				{Severity: sdkv2diag.Error, Summary: "first error", Detail: "detail about first"},
				{Severity: sdkv2diag.Error, Summary: "second error"},
			},
			expectedErr: "first error: detail about first\nsecond error",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := conversion.SDKV2DiagsToError(tc.diags)
			if tc.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
