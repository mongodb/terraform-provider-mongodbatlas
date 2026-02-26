package customvalidator_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	schemavalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customvalidator"
)

func TestValidateDiscriminator(t *testing.T) {
	rootObjType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"type":           tftypes.String,
			"aws_specific":   tftypes.String,
			"azure_specific": tftypes.String,
		},
	}

	innerObjType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"type":           tftypes.String,
			"aws_specific":   tftypes.String,
			"azure_specific": tftypes.String,
		},
	}

	nestedObjType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"nested": innerObjType,
		},
	}

	listNestedType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"items": tftypes.List{ElementType: innerObjType},
		},
	}

	mapNestedType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"items": tftypes.Map{ElementType: innerObjType},
		},
	}

	setNestedType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"items": tftypes.Set{ElementType: innerObjType},
		},
	}

	def := customvalidator.DiscriminatorDefinition{
		Mapping: map[string]customvalidator.VariantDefinition{
			"AWS": {
				Allowed:  []string{"aws_specific"},
				Required: []string{"aws_specific"},
			},
			"AZURE": {
				Allowed:  []string{"azure_specific"},
				Required: []string{"azure_specific"},
			},
		},
	}

	tests := []struct {
		raw            tftypes.Value
		def            customvalidator.DiscriminatorDefinition
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
				"type":           tftypes.NewValue(tftypes.String, "AWS"),
				"aws_specific":   tftypes.NewValue(tftypes.String, nil),
				"azure_specific": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors:   1,
			expectInDetail: []string{`"aws_specific" must be set when type is "AWS"`},
		},
		{
			name:        "disallowed attribute present emits not-allowed diagnostic",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":           tftypes.NewValue(tftypes.String, "AWS"),
				"aws_specific":   tftypes.NewValue(tftypes.String, "some-arn"),
				"azure_specific": tftypes.NewValue(tftypes.String, "some-url"),
			}),
			expectErrors:   1,
			expectInDetail: []string{`"azure_specific" is not allowed when type is "AWS"`},
		},
		{
			name:        "unknown discriminator value skips all checks",
			def:         def,
			configValue: types.StringValue("GCP"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":           tftypes.NewValue(tftypes.String, "GCP"),
				"aws_specific":   tftypes.NewValue(tftypes.String, "some-arn"),
				"azure_specific": tftypes.NewValue(tftypes.String, "some-url"),
			}),
			expectErrors: 0,
		},
		{
			name:        "null discriminator config value skips all checks",
			def:         def,
			configValue: types.StringNull(),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":           tftypes.NewValue(tftypes.String, nil),
				"aws_specific":   tftypes.NewValue(tftypes.String, nil),
				"azure_specific": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors: 0,
		},
		{
			name:        "unknown discriminator config value skips all checks",
			def:         def,
			configValue: types.StringUnknown(),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":           tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"aws_specific":   tftypes.NewValue(tftypes.String, nil),
				"azure_specific": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors: 0,
		},
		{
			name:        "valid configuration for active variant emits no diagnostics",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":           tftypes.NewValue(tftypes.String, "AWS"),
				"aws_specific":   tftypes.NewValue(tftypes.String, "arn:aws:iam::role"),
				"azure_specific": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors: 0,
		},
		// Unknown values in config mean the practitioner set the attribute to an expression whose value is not yet resolved.
		// Unset Optional+Computed attributes are null (not unknown) during validation, they only become unknown later during PlanResourceChange.
		{
			name:        "unrelated sibling with unknown value emits no diagnostics",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"type":           tftypes.String,
					"aws_specific":   tftypes.String,
					"azure_specific": tftypes.String,
					"name":           tftypes.String,
				},
			}, map[string]tftypes.Value{
				"type":           tftypes.NewValue(tftypes.String, "AWS"),
				"aws_specific":   tftypes.NewValue(tftypes.String, "arn:aws:iam::role"),
				"azure_specific": tftypes.NewValue(tftypes.String, nil),
				"name":           tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}),
			expectErrors: 0,
		},
		{
			name:        "unknown disallowed sibling emits not-allowed diagnostic",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":           tftypes.NewValue(tftypes.String, "AWS"),
				"aws_specific":   tftypes.NewValue(tftypes.String, "arn:aws:iam::role"),
				"azure_specific": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}),
			expectErrors:   1,
			expectInDetail: []string{`"azure_specific" is not allowed when type is "AWS"`},
		},
		{
			name:        "unknown required sibling is accepted as set by user",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":           tftypes.NewValue(tftypes.String, "AWS"),
				"aws_specific":   tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"azure_specific": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors: 0,
		},
		{
			name:        "both required missing and disallowed present emit multiple diagnostics",
			def:         def,
			configValue: types.StringValue("AZURE"),
			configPath:  path.Root("type"),
			raw: tftypes.NewValue(rootObjType, map[string]tftypes.Value{
				"type":           tftypes.NewValue(tftypes.String, "AZURE"),
				"aws_specific":   tftypes.NewValue(tftypes.String, "some-arn"),
				"azure_specific": tftypes.NewValue(tftypes.String, nil),
			}),
			expectErrors: 2,
			expectInDetail: []string{
				`"azure_specific" must be set when type is "AZURE"`,
				`"aws_specific" is not allowed when type is "AZURE"`,
			},
		},
		{
			name:        "discriminator nested inside object element",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("nested").AtName("type"),
			raw: tftypes.NewValue(nestedObjType, map[string]tftypes.Value{
				"nested": tftypes.NewValue(innerObjType, map[string]tftypes.Value{
					"type":           tftypes.NewValue(tftypes.String, "AWS"),
					"aws_specific":   tftypes.NewValue(tftypes.String, nil),
					"azure_specific": tftypes.NewValue(tftypes.String, nil),
				}),
			}),
			expectErrors:   1,
			expectInDetail: []string{`"aws_specific" must be set when type is "AWS"`},
		},
		{
			name:        "discriminator nested inside list element emits error",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("items").AtListIndex(0).AtName("type"),
			raw: tftypes.NewValue(listNestedType, map[string]tftypes.Value{
				"items": tftypes.NewValue(tftypes.List{ElementType: innerObjType}, []tftypes.Value{
					tftypes.NewValue(innerObjType, map[string]tftypes.Value{
						"type":           tftypes.NewValue(tftypes.String, "AWS"),
						"aws_specific":   tftypes.NewValue(tftypes.String, nil),
						"azure_specific": tftypes.NewValue(tftypes.String, nil),
					}),
				}),
			}),
			expectErrors:   1,
			expectInDetail: []string{`"aws_specific" must be set when type is "AWS"`},
		},
		{
			name:        "discriminator nested inside map element emits error",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath:  path.Root("items").AtMapKey("conn1").AtName("type"),
			raw: tftypes.NewValue(mapNestedType, map[string]tftypes.Value{
				"items": tftypes.NewValue(tftypes.Map{ElementType: innerObjType}, map[string]tftypes.Value{
					"conn1": tftypes.NewValue(innerObjType, map[string]tftypes.Value{
						"type":           tftypes.NewValue(tftypes.String, "AWS"),
						"aws_specific":   tftypes.NewValue(tftypes.String, nil),
						"azure_specific": tftypes.NewValue(tftypes.String, nil),
					}),
				}),
			}),
			expectErrors:   1,
			expectInDetail: []string{`"aws_specific" must be set when type is "AWS"`},
		},
		{
			name:        "discriminator nested inside set element emits error",
			def:         def,
			configValue: types.StringValue("AWS"),
			configPath: path.Root("items").AtSetValue(
				types.ObjectValueMust(map[string]attr.Type{
					"type":           types.StringType,
					"aws_specific":   types.StringType,
					"azure_specific": types.StringType,
				}, map[string]attr.Value{
					"type":           types.StringValue("AWS"),
					"aws_specific":   types.StringNull(),
					"azure_specific": types.StringNull(),
				}),
			).AtName("type"),
			raw: tftypes.NewValue(setNestedType, map[string]tftypes.Value{
				"items": tftypes.NewValue(tftypes.Set{ElementType: innerObjType}, []tftypes.Value{
					tftypes.NewValue(innerObjType, map[string]tftypes.Value{
						"type":           tftypes.NewValue(tftypes.String, "AWS"),
						"aws_specific":   tftypes.NewValue(tftypes.String, nil),
						"azure_specific": tftypes.NewValue(tftypes.String, nil),
					}),
				}),
			}),
			expectErrors:   1,
			expectInDetail: []string{`"aws_specific" must be set when type is "AWS"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := customvalidator.ValidateDiscriminator(tt.def)

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
