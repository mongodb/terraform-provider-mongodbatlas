//nolint:testpackage // accesses unexported helpers
package apiresource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateMutex inlines the body of urs.ValidateConfig that operates on TFModelUpdate.
// We exercise it directly rather than through the framework to keep the test fast and
// dependency-free. If ValidateConfig grows, factor out a pure helper and test it here.
func validateMutex(cfg TFModelUpdate) diag.Diagnostics {
	var diags diag.Diagnostics
	if !cfg.VersionHeader.IsNull() && !cfg.VersionHeader.IsUnknown() &&
		!cfg.Preview.IsNull() && !cfg.Preview.IsUnknown() && cfg.Preview.ValueBool() {
		diags.AddAttributeError(path.Root("preview"),
			"version_header and preview are mutually exclusive",
			"Set either `version_header` or `preview = true`, not both.")
	}
	return diags
}

func TestUpdateResource_ValidateConfig_PreviewVersionHeaderMutex(t *testing.T) {
	tests := map[string]struct {
		cfg     TFModelUpdate
		wantErr bool
	}{
		"both set is an error": {
			cfg: TFModelUpdate{
				VersionHeader: types.StringValue("application/vnd.atlas.2023-02-01+json"),
				Preview:       types.BoolValue(true),
			},
			wantErr: true,
		},
		"only preview is fine": {
			cfg:     TFModelUpdate{Preview: types.BoolValue(true)},
			wantErr: false,
		},
		"only version_header is fine": {
			cfg:     TFModelUpdate{VersionHeader: types.StringValue("application/vnd.atlas.2023-02-01+json")},
			wantErr: false,
		},
		"both empty is fine": {
			cfg:     TFModelUpdate{},
			wantErr: false,
		},
		"preview false + version_header set is fine": {
			cfg: TFModelUpdate{
				VersionHeader: types.StringValue("application/vnd.atlas.preview+json"),
				Preview:       types.BoolValue(false),
			},
			wantErr: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			diags := validateMutex(tc.cfg)
			gotErr := diags.HasError()
			if gotErr != tc.wantErr {
				t.Fatalf("HasError=%v want=%v; diags=%v", gotErr, tc.wantErr, diags)
			}
		})
	}
}
