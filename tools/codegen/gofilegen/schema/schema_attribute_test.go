package schema_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSchemaAttributes_CreateOnly(t *testing.T) {
	tests := map[string]struct {
		attribute       codespec.Attribute
		hasPlanModifier bool
	}{
		"No create_only - no plan modifiers": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_string",
				TFModelName:              "TestString",
				String:                   &codespec.StringAttribute{},
				ComputedOptionalRequired: codespec.Optional,
				CreateOnly:               false,
			},
			hasPlanModifier: false,
		},
		"String attribute with create_only - uses CreateOnly()": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_string",
				TFModelName:              "TestString",
				String:                   &codespec.StringAttribute{},
				ComputedOptionalRequired: codespec.Optional,
				CreateOnly:               true,
			},
			hasPlanModifier: true,
		},
		"Bool attribute with create_only but no default - uses CreateOnly()": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_bool",
				TFModelName:              "TestBool",
				Bool:                     &codespec.BoolAttribute{},
				ComputedOptionalRequired: codespec.Optional,
				CreateOnly:               true,
			},
			hasPlanModifier: true,
		},
		"Bool attribute with create_only and default true - uses CreateOnlyBoolWithDefault(true)": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_bool",
				TFModelName:              "TestBool",
				Bool:                     &codespec.BoolAttribute{Default: new(true)},
				ComputedOptionalRequired: codespec.ComputedOptional,
				CreateOnly:               true,
			},
			hasPlanModifier: true,
		},
		"Bool attribute with create_only and default false - uses CreateOnlyBoolWithDefault(false)": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_bool",
				TFModelName:              "TestBool",
				Bool:                     &codespec.BoolAttribute{Default: new(false)},
				ComputedOptionalRequired: codespec.ComputedOptional,
				CreateOnly:               true,
			},
			hasPlanModifier: true,
		},
		"Int64 attribute with create_only - uses CreateOnly()": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_int",
				TFModelName:              "TestInt",
				Int64:                    &codespec.Int64Attribute{},
				ComputedOptionalRequired: codespec.Optional,
				CreateOnly:               true,
			},
			hasPlanModifier: true,
		},
		"Computed attribute with create_only - uses CreateOnly() (model is enforced)": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_computed",
				TFModelName:              "TestComputed",
				String:                   &codespec.StringAttribute{},
				ComputedOptionalRequired: codespec.Computed,
				CreateOnly:               true,
			},
			hasPlanModifier: true,
		},
		"ComputedOptional attribute with create_only - uses CreateOnly()": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_computed_optional",
				TFModelName:              "TestComputedOptional",
				String:                   &codespec.StringAttribute{},
				ComputedOptionalRequired: codespec.ComputedOptional,
				CreateOnly:               true,
			},
			hasPlanModifier: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := schema.GenerateSchemaAttributes([]codespec.Attribute{tc.attribute}, nil)
			require.NoError(t, err)
			code := result.Code
			if !tc.hasPlanModifier {
				assert.NotContains(t, code, "PlanModifiers:")
				return
			}
			assert.Contains(t, code, "PlanModifiers:")
			if tc.attribute.Bool != nil && tc.attribute.Bool.Default != nil {
				expected := fmt.Sprintf("customplanmodifier.CreateOnlyBoolWithDefault(%t)", *tc.attribute.Bool.Default)
				assert.Contains(t, code, expected)
				return
			}
			if tc.attribute.CreateOnly {
				assert.Contains(t, code, "customplanmodifier.CreateOnly()")
			}
		})
	}
}

func TestGenerateSchemaAttributes_ImmutableComputed(t *testing.T) {
	tests := map[string]struct {
		attribute       codespec.Attribute
		hasPlanModifier bool
	}{
		"No ImmutableComputed - no plan modifiers": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_string",
				TFModelName:              "TestString",
				String:                   &codespec.StringAttribute{},
				ComputedOptionalRequired: codespec.Computed,
				ImmutableComputed:        false,
			},
			hasPlanModifier: false,
		},
		"String attribute with ImmutableComputed - uses UseStateForUnknown()": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_string",
				TFModelName:              "TestString",
				String:                   &codespec.StringAttribute{},
				ComputedOptionalRequired: codespec.Computed,
				ImmutableComputed:        true,
			},
			hasPlanModifier: true,
		},
		"Sensitive string attribute with ImmutableComputed - uses UseStateForUnknown()": {
			attribute: codespec.Attribute{
				TFSchemaName:             "secret",
				TFModelName:              "Secret",
				String:                   &codespec.StringAttribute{},
				ComputedOptionalRequired: codespec.Computed,
				Sensitive:                true,
				ImmutableComputed:        true,
			},
			hasPlanModifier: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := schema.GenerateSchemaAttributes([]codespec.Attribute{tc.attribute}, nil)
			require.NoError(t, err)
			code := result.Code
			if !tc.hasPlanModifier {
				assert.NotContains(t, code, "PlanModifiers:")
				return
			}
			assert.Contains(t, code, "PlanModifiers:")
			assert.Contains(t, code, "stringplanmodifier.UseStateForUnknown()")
		})
	}
}

func TestGenerateSchemaAttributes_ImmutableComputedNonStringReturnsError(t *testing.T) {
	tests := map[string]codespec.Attribute{
		"Bool attribute with ImmutableComputed": {
			TFSchemaName:             "test_bool",
			TFModelName:              "TestBool",
			Bool:                     &codespec.BoolAttribute{},
			ComputedOptionalRequired: codespec.Computed,
			ImmutableComputed:        true,
		},
		"Int64 attribute with ImmutableComputed": {
			TFSchemaName:             "test_int",
			TFModelName:              "TestInt",
			Int64:                    &codespec.Int64Attribute{},
			ComputedOptionalRequired: codespec.Computed,
			ImmutableComputed:        true,
		},
	}

	for name, attr := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := schema.GenerateSchemaAttributes([]codespec.Attribute{attr}, nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "immutableComputed is only supported for string attributes")
			assert.Contains(t, err.Error(), attr.TFSchemaName)
		})
	}
}

func TestGenerateSchemaAttributes_RequestOnlyRequiredOnCreate(t *testing.T) {
	tests := map[string]struct {
		attribute       codespec.Attribute
		hasPlanModifier bool
	}{
		"No RequestOnlyRequiredOnCreate - no plan modifiers": {
			attribute: codespec.Attribute{
				TFSchemaName:                "test_string",
				TFModelName:                 "TestString",
				String:                      &codespec.StringAttribute{},
				ComputedOptionalRequired:    codespec.Optional,
				RequestOnlyRequiredOnCreate: false,
			},
			hasPlanModifier: false,
		},
		"RequestOnlyRequiredOnCreate - uses RequestOnlyRequiredOnCreate()": {
			attribute: codespec.Attribute{
				TFSchemaName:                "test_string",
				TFModelName:                 "TestString",
				String:                      &codespec.StringAttribute{},
				ComputedOptionalRequired:    codespec.Optional,
				RequestOnlyRequiredOnCreate: true,
			},
			hasPlanModifier: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := schema.GenerateSchemaAttributes([]codespec.Attribute{tc.attribute}, nil)
			require.NoError(t, err)
			code := result.Code
			if !tc.hasPlanModifier {
				assert.NotContains(t, code, "PlanModifiers:")
				return
			}
			assert.Contains(t, code, "PlanModifiers:")
			if tc.attribute.RequestOnlyRequiredOnCreate {
				assert.Contains(t, code, "customplanmodifier.RequestOnlyRequiredOnCreate()")
			}
		})
	}
}

func TestGenerateSchemaAttributes_DiscriminatorValidator(t *testing.T) {
	disc := &codespec.Discriminator{
		PropertyName: codespec.DiscriminatorAttrName{APIName: "type", TFSchemaName: "type"},
		Mapping: map[string]codespec.DiscriminatorType{
			"Cluster": {
				Allowed:  []codespec.DiscriminatorAttrName{{TFSchemaName: "cluster_name"}, {TFSchemaName: "db_role"}},
				Required: []codespec.DiscriminatorAttrName{{TFSchemaName: "cluster_name"}},
			},
			"Https": {
				Allowed: []codespec.DiscriminatorAttrName{{TFSchemaName: "url"}, {TFSchemaName: "headers"}},
			},
			"Sample": {
				Allowed: []codespec.DiscriminatorAttrName{},
			},
		},
	}

	attr := codespec.Attribute{
		TFSchemaName:             "type",
		TFModelName:              "Type",
		String:                   &codespec.StringAttribute{},
		ComputedOptionalRequired: codespec.Required,
	}

	result, err := schema.GenerateSchemaAttributes([]codespec.Attribute{attr}, disc)
	require.NoError(t, err)
	code := result.Code

	assert.Contains(t, code, "Validators:")
	assert.Contains(t, code, "customvalidator.ValidateDiscriminator")
	assert.Contains(t, code, "customvalidator.DiscriminatorDefinition")
	assert.Contains(t, code, "customvalidator.VariantDefinition")

	assert.Contains(t, code, `"Cluster"`)
	assert.Contains(t, code, `"Https"`)
	assert.Contains(t, code, `"Sample"`)
	assert.Contains(t, code, `"cluster_name"`)
	assert.Contains(t, code, `"db_role"`)
	assert.Contains(t, code, `"url"`)
	assert.Contains(t, code, `"headers"`)

	// Verify sorted order: Cluster before Https before Sample
	clusterIdx := strings.Index(code, `"Cluster"`)
	httpsIdx := strings.Index(code, `"Https"`)
	sampleIdx := strings.Index(code, `"Sample"`)
	assert.Greater(t, httpsIdx, clusterIdx, "Cluster should come before Https")
	assert.Greater(t, sampleIdx, httpsIdx, "Https should come before Sample")

	// Verify Required is emitted for Cluster
	assert.Contains(t, code, "Required: []string{")

	// Verify imports
	assert.Contains(t, result.Imports, "github.com/hashicorp/terraform-plugin-framework/schema/validator")
	assert.Contains(t, result.Imports, "github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customvalidator")
}

func TestGenerateSchemaAttributes_DiscriminatorSkipValidation(t *testing.T) {
	disc := &codespec.Discriminator{
		PropertyName:   codespec.DiscriminatorAttrName{APIName: "type", TFSchemaName: "type"},
		SkipValidation: true,
		Mapping: map[string]codespec.DiscriminatorType{
			"AWS": {Allowed: []codespec.DiscriminatorAttrName{{TFSchemaName: "aws_field"}}},
		},
	}

	attr := codespec.Attribute{
		TFSchemaName:             "type",
		TFModelName:              "Type",
		String:                   &codespec.StringAttribute{},
		ComputedOptionalRequired: codespec.Required,
	}

	result, err := schema.GenerateSchemaAttributes([]codespec.Attribute{attr}, disc)
	require.NoError(t, err)
	assert.NotContains(t, result.Code, "Validators:")
	assert.NotContains(t, result.Code, "customvalidator")
}
