package schema_test

import (
	"fmt"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/schema"
	"github.com/stretchr/testify/assert"
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
				Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(true)},
				ComputedOptionalRequired: codespec.ComputedOptional,
				CreateOnly:               true,
			},
			hasPlanModifier: true,
		},
		"Bool attribute with create_only and default false - uses CreateOnlyBoolWithDefault(false)": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_bool",
				TFModelName:              "TestBool",
				Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(false)},
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
			result := schema.GenerateSchemaAttributes([]codespec.Attribute{tc.attribute})
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
