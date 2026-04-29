package advancedcluster_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

func TestAdvancedCluster_OverrideAttributesWithPrevStateValue_mongoDBMajorVersionWarning(t *testing.T) {
	testCases := []struct {
		name          string
		beforeVersion string
		afterVersion  string
		expectWarning bool
	}{
		{
			name:          "warns when major version changed outside of Terraform",
			beforeVersion: "7.0",
			afterVersion:  "8.0",
			expectWarning: true,
		},
		{
			name:          "no warning for formatting difference",
			beforeVersion: "8",
			afterVersion:  "8.0",
			expectWarning: false,
		},
		{
			name:          "no warning when Atlas version differs within same major",
			beforeVersion: "8.0",
			afterVersion:  "8.2",
			expectWarning: false,
		},
		{
			name:          "no warning when versions are equal",
			beforeVersion: "8.0",
			afterVersion:  "8.0",
			expectWarning: false,
		},
		{
			name:          "warns when major version downgraded outside of Terraform",
			beforeVersion: "8.0",
			afterVersion:  "7.0",
			expectWarning: true,
		},
	}

	t.Run("no-op when prior state has no version set", func(t *testing.T) {
		modelIn := &advancedcluster.TFModel{
			MongoDBMajorVersion: types.StringNull(),
			Labels:              types.MapNull(types.StringType),
			Tags:                types.MapNull(types.StringType),
		}
		modelOut := &advancedcluster.TFModel{
			MongoDBMajorVersion: types.StringValue("8.0"),
			Labels:              types.MapNull(types.StringType),
			Tags:                types.MapNull(types.StringType),
		}
		var diags diag.Diagnostics

		advancedcluster.OverrideAttributesWithPrevStateValue(modelIn, modelOut, &diags)

		assert.Equal(t, 0, diags.WarningsCount())
		assert.Equal(t, "8.0", modelOut.MongoDBMajorVersion.ValueString())
	})

	// afterVersion unknown: override must still apply (no regression from old behavior)
	t.Run("preserves prior state when Atlas returns Unknown version", func(t *testing.T) {
		modelIn := &advancedcluster.TFModel{
			MongoDBMajorVersion: types.StringValue("8.0"),
			Labels:              types.MapNull(types.StringType),
			Tags:                types.MapNull(types.StringType),
		}
		modelOut := &advancedcluster.TFModel{
			MongoDBMajorVersion: types.StringUnknown(),
			Labels:              types.MapNull(types.StringType),
			Tags:                types.MapNull(types.StringType),
		}
		var diags diag.Diagnostics

		advancedcluster.OverrideAttributesWithPrevStateValue(modelIn, modelOut, &diags)

		assert.Equal(t, 0, diags.WarningsCount())
		assert.Equal(t, "8.0", modelOut.MongoDBMajorVersion.ValueString())
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			modelIn := &advancedcluster.TFModel{
				MongoDBMajorVersion: types.StringValue(tc.beforeVersion),
				Labels:              types.MapNull(types.StringType),
				Tags:                types.MapNull(types.StringType),
			}
			modelOut := &advancedcluster.TFModel{
				MongoDBMajorVersion: types.StringValue(tc.afterVersion),
				Labels:              types.MapNull(types.StringType),
				Tags:                types.MapNull(types.StringType),
			}
			var diags diag.Diagnostics

			advancedcluster.OverrideAttributesWithPrevStateValue(modelIn, modelOut, &diags)

			assert.False(t, diags.HasError())
			if tc.expectWarning {
				assert.Equal(t, 1, diags.WarningsCount())
			} else {
				assert.Equal(t, 0, diags.WarningsCount())
			}
			assert.Equal(t, tc.beforeVersion, modelOut.MongoDBMajorVersion.ValueString())
		})
	}
}
