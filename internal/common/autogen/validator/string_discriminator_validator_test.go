package validator_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemavalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/validator"
)

func TestValidateDiscriminator(t *testing.T) {
	rootObjType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"type":        tftypes.String,
			"role_arn":    tftypes.String,
			"service_url": tftypes.String,
		},
	}

	innerObjType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"type":        tftypes.String,
			"role_arn":    tftypes.String,
			"service_url": tftypes.String,
		},
	}

	nestedObjType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"nested": innerObjType,
		},
	}

	def := validator.DiscriminatorDefinition{
		Mapping: map[string]validator.VariantDefinition{
			"AWS": {
				Allowed:  []string{"role_arn"},
				Required: []string{"role_arn"},
			},
			"AZURE": {
				Allowed:  []string{"service_url"},
				Required: []string{"service_url"},
			},
		},
	}

	tests := []struct {
		raw            tftypes.Value
		def            validator.DiscriminatorDefinition
		configValue    basetypes.StringValue
		name           string
		configPath     path.Path
		expectInDetail []string
		expectErrors   int
	}{
		{
			name:        "required attribute missing emits must-be-set diagnostic",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":        tftypes.NewValue(tftypes.String, "AWS"),
				"role_arn":    tftypes.NewValue(tftypes.String, nil),
				"service_url": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors:   1,
			expectInDetail: []string{`"role_arn" must be set when type is "AWS"`},
		},
		{
			name:        "disallowed attribute present emits not-allowed diagnostic",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":        tftypes.NewValue(tftypes.String, "AWS"),
				"role_arn":    tftypes.NewValue(tftypes.String, "some-arn"),
				"service_url": tftypes.NewValue(tftypes.String, "some-url"),
			}),
			expectErrors:   1,
			expectInDetail: []string{`"service_url" is not allowed when type is "AWS"`},
		},
		{
			name:        "unknown discriminator value skips all checks",
			def:         def,
			configValue: types.StringValue("GCP"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":        tftypes.NewValue(tftypes.String, "GCP"),
				"role_arn":    tftypes.NewValue(tftypes.String, "some-arn"),
				"service_url": tftypes.NewValue(tftypes.String, "some-url"),
			}),
			expectErrors: 0,
		},
		{
			name:        "null discriminator config value skips all checks",
			def:         def,
			configValue: types.StringNull(),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":        tftypes.NewValue(tftypes.String, nil),
				"role_arn":    tftypes.NewValue(tftypes.String, nil),
				"service_url": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors: 0,
		},
		{
			name:        "unknown discriminator config value skips all checks",
			def:         def,
			configValue: types.StringUnknown(),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":        tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"role_arn":    tftypes.NewValue(tftypes.String, nil),
				"service_url": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors: 0,
		},
		{
			name:        "valid configuration for active variant emits no diagnostics",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":        tftypes.NewValue(tftypes.String, "AWS"),
				"role_arn":    tftypes.NewValue(tftypes.String, "arn:aws:iam::role"),
				"service_url": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors: 0,
		},
		{
			name:        "nested path sibling resolution works",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("nested").AtName("type"),
			raw: tftypes.NewValue(nestedObjType, map[string]tftypes.Value{
				"nested": tftypes.NewValue(innerObjType, map[string]tftypes.Value{
					"type":        tftypes.NewValue(tftypes.String, "AWS"),
					"role_arn":    tftypes.NewValue(tftypes.String, nil),
					"service_url": tftypes.NewValue(tftypes.String, nil),
				}),
			}),
			expectErrors:   1,
			expectInDetail: []string{`"role_arn" must be set when type is "AWS"`},
		},
		{
			name:        "unrelated sibling with unknown value emits no diagnostics",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"type":        tftypes.String,
					"role_arn":    tftypes.String,
					"service_url": tftypes.String,
					"name":        tftypes.String,
				},
			}, map[string]tftypes.Value{
				"type":        tftypes.NewValue(tftypes.String, "AWS"),
				"role_arn":    tftypes.NewValue(tftypes.String, "arn:aws:iam::role"),
				"service_url": tftypes.NewValue(tftypes.String, nil),
				"name":        tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}),
			expectErrors: 0,
		},
		{
			name:        "unknown disallowed sibling skips not-allowed check",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":        tftypes.NewValue(tftypes.String, "AWS"),
				"role_arn":    tftypes.NewValue(tftypes.String, "arn:aws:iam::role"),
				"service_url": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}),
			expectErrors: 0,
		},
		{
			name:        "unknown required sibling emits must-be-set diagnostic",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":        tftypes.NewValue(tftypes.String, "AWS"),
				"role_arn":    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"service_url": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors:   1,
			expectInDetail: []string{`"role_arn" must be set when type is "AWS"`},
		},
		{
			name:        "both required missing and disallowed present emit multiple diagnostics",
			def:         def,
			configValue: types.StringValue("AZURE"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":        tftypes.NewValue(tftypes.String, "AZURE"),
				"role_arn":    tftypes.NewValue(tftypes.String, "some-arn"),
				"service_url": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors: 2,
			expectInDetail: []string{
				`"service_url" must be set when type is "AZURE"`,
				`"role_arn" is not allowed when type is "AZURE"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := validator.ValidateDiscriminator(tt.def)

			req := schemavalidator.StringRequest{
				ConfigValue: tt.configValue,
				Path:        tt.configPath,
				Config: tfsdk.Config{
					Raw: tt.raw,
				},
			}
			resp := &schemavalidator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}

			v.ValidateString(t.Context(), req, resp)

			errors := resp.Diagnostics.Errors()
			if len(errors) != tt.expectErrors {
				t.Fatalf("expected %d errors, got %d: %v", tt.expectErrors, len(errors), errors)
			}

			for _, msg := range tt.expectInDetail {
				if !diagnosticContains(errors, msg) {
					t.Errorf("expected diagnostic detail containing %q, got: %v", msg, errors)
				}
			}
		})
	}
}

func diagnosticContains(diagnostics diag.Diagnostics, substr string) bool {
	for _, d := range diagnostics {
		if strings.Contains(d.Detail(), substr) {
			return true
		}
	}
	return false
}
