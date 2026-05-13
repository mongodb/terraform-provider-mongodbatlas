//nolint:testpackage // accesses unexported helpers
package apiresource

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/responseproject"
)

// validateMutex inlines the body of urs.ValidateConfig that operates on TFModelUpdate.
// We exercise it directly rather than through the framework to keep the test fast and
// dependency-free. If ValidateConfig grows, factor out a pure helper and test it here.
func validateMutex(cfg TFModelUpdate) diag.Diagnostics { //nolint:gocritic // test helper; pass-by-value matches the framework's call site
	var diags diag.Diagnostics
	if !cfg.VersionHeader.IsNull() && !cfg.VersionHeader.IsUnknown() &&
		!cfg.Preview.IsNull() && !cfg.Preview.IsUnknown() && cfg.Preview.ValueBool() {
		diags.AddAttributeError(path.Root("preview"),
			"version_header and preview are mutually exclusive",
			"Set either `version_header` or `preview = true`, not both.")
	}
	if overlap := responseproject.PathsOverlap(
		exportPaths(cfg.ResponseExportValues), exportPaths(cfg.ResponseExportValuesSensitive),
	); len(overlap) > 0 {
		diags.AddAttributeError(path.Root("response_export_values_sensitive"),
			"path declared in both response_export_values and response_export_values_sensitive",
			fmt.Sprintf("each path must appear in only one list. Overlapping: %v", overlap))
	}
	return diags
}

func listOfStrings(vals ...string) types.List {
	elems := make([]attr.Value, 0, len(vals))
	for _, v := range vals {
		elems = append(elems, types.StringValue(v))
	}
	l, _ := types.ListValue(types.StringType, elems)
	return l
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

func TestUpdateResource_ValidateConfig_ExportValuesOverlap(t *testing.T) {
	tests := map[string]struct {
		cfg     TFModelUpdate
		wantErr bool
	}{
		"no overlap is fine": {
			cfg: TFModelUpdate{
				ResponseExportValues:          listOfStrings("apiKeyId", "createdAt"),
				ResponseExportValuesSensitive: listOfStrings("secret"),
			},
			wantErr: false,
		},
		"exact path overlap is rejected": {
			cfg: TFModelUpdate{
				ResponseExportValues:          listOfStrings("apiKeyId", "secret"),
				ResponseExportValuesSensitive: listOfStrings("secret"),
			},
			wantErr: true,
		},
		"both empty is fine": {
			cfg:     TFModelUpdate{},
			wantErr: false,
		},
		"only one list set is fine": {
			cfg: TFModelUpdate{
				ResponseExportValues: listOfStrings("apiKeyId"),
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
