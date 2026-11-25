package conversion_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateProjectID(t *testing.T) {
	require.NoError(t, conversion.ValidateProjectID("123456789012345678901234"))
	err := conversion.ValidateProjectID("invalid_project_id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid_project_id")
	assert.Contains(t, err.Error(), "project_id must be a 24 character hex string")
}

func TestValidateClusterName(t *testing.T) {
	require.NoError(t, conversion.ValidateClusterName("clusterName"))
	err := conversion.ValidateClusterName("_invalidClusterName")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "_invalidClusterName")
	assert.Contains(t, err.Error(), "cluster_name must be a string with length between 1 and 64, starting and ending with an alphanumeric character, and containing only alphanumeric characters and hyphens")
}

func TestImportSplit(t *testing.T) {
	tests := map[string]struct {
		importRaw     string
		parts         []string
		expectedParts int
		expected      bool
	}{
		"valid input": {
			importRaw:     "part1/part2/part3",
			expectedParts: 3,
			expected:      true,
			parts:         []string{"part1", "part2", "part3"},
		},
		"invalid input with more parts than expected": {
			importRaw:     "part1/part2/part3/part4",
			expectedParts: 3,
			expected:      false,
			parts:         nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ok, parts := conversion.ImportSplit(tc.importRaw, tc.expectedParts)
			if tc.expected {
				assert.Len(t, parts, tc.expectedParts)
			}
			assert.Equal(t, tc.expected, ok)
			assert.Equal(t, tc.parts, parts)
		})
	}
}

func TestImportStateProjectIDClusterName(t *testing.T) {
	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id":   schema.StringAttribute{},
			"cluster_name": schema.StringAttribute{},
		},
	}
	testType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"project_id":   tftypes.String,
			"cluster_name": tftypes.String,
		},
	}
	emptyStateValue := tftypes.NewValue(testType, map[string]tftypes.Value{
		"project_id":   tftypes.NewValue(tftypes.String, nil),
		"cluster_name": tftypes.NewValue(tftypes.String, nil),
	})

	tests := map[string]struct {
		importID      string
		projectID     string
		clusterName   string
		expectError   bool
		expectWarning bool
	}{
		"new slash format": {
			importID:    "507f1f77bcf86cd799439011/clusterName",
			projectID:   "507f1f77bcf86cd799439011",
			clusterName: "clusterName",
		},
		"legacy dash format": {
			importID:      "507f1f77bcf86cd799439011-clusterName",
			expectWarning: true,
			projectID:     "507f1f77bcf86cd799439011",
			clusterName:   "clusterName",
		},
		"legacy dash format with dashes in cluster name": {
			importID:      "507f1f77bcf86cd799439011-my-cluster-name",
			expectWarning: true,
			projectID:     "507f1f77bcf86cd799439011",
			clusterName:   "my-cluster-name",
		},
		"invalid format without separator": {
			importID:    "invalidID",
			expectError: true,
		},
		"invalid format with extra slash parts": {
			importID:    "507f1f77bcf86cd799439011/clusterName/extra",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := resource.ImportStateRequest{ID: tc.importID}
			resp := &resource.ImportStateResponse{
				State: tfsdk.State{
					Raw:    emptyStateValue,
					Schema: testSchema,
				},
			}

			conversion.ImportStateProjectIDClusterName(context.Background(), req, resp, "project_id", "cluster_name")

			if tc.expectError {
				require.True(t, resp.Diagnostics.HasError(), "expected errors, got none")
				errStr := ""
				for _, err := range resp.Diagnostics.Errors() {
					errStr += err.Summary() + " " + err.Detail() + " "
				}
				assert.Contains(t, errStr, "{project_id}/{cluster_name}")
				assert.Contains(t, errStr, "{project_id}-{cluster_name}")
				return
			}
			require.False(t, resp.Diagnostics.HasError(), "unexpected errors: %#v", resp.Diagnostics.Errors())

			var gotWarnings bool
			for _, warning := range resp.Diagnostics.Warnings() {
				if warning.Summary() == conversion.LegacyImportFormatWarningSummary {
					gotWarnings = true
				}
			}
			assert.Equal(t, tc.expectWarning, gotWarnings)

			var state struct {
				ProjectID   types.String `tfsdk:"project_id"`
				ClusterName types.String `tfsdk:"cluster_name"`
			}
			resp.State.Get(context.Background(), &state)
			require.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics after state get: %#v", resp.Diagnostics.Errors())
			assert.Equal(t, tc.projectID, state.ProjectID.ValueString())
			assert.Equal(t, tc.clusterName, state.ClusterName.ValueString())
		})
	}
}

func TestImportSplitN(t *testing.T) {
	tests := map[string]struct {
		importRaw     string
		parts         []string
		expectedParts int
		expected      bool
	}{
		"keeps slashes in the last part": {
			importRaw:     "project/10.0.0.0/8",
			expectedParts: 2,
			expected:      true,
			parts:         []string{"project", "10.0.0.0/8"},
		},
		"invalid input with fewer parts than expected": {
			importRaw:     "project",
			expectedParts: 2,
			expected:      false,
			parts:         nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ok, parts := conversion.ImportSplitN(tc.importRaw, tc.expectedParts)
			assert.Equal(t, tc.expected, ok)
			assert.Equal(t, tc.parts, parts)
		})
	}
}

func TestLegacyImportFormatWarningDetail(t *testing.T) {
	detail := conversion.LegacyImportFormatWarningDetail("{project_id}-{name}", "{project_id}/{name}")
	assert.Contains(t, detail, "{project_id}-{name}")
	assert.Contains(t, detail, "{project_id}/{name}")
	assert.Contains(t, detail, "deprecated")
}
