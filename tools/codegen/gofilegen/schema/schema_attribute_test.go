package schema_test

import (
	"fmt"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/schema"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSchemaAttributes_NonUpdatable(t *testing.T) {
	tests := map[string]struct {
		attribute       codespec.Attribute
		hasPlanModifier bool
	}{
		"No non-updatable - no plan modifiers": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_string",
				TFModelName:              "TestString",
				String:                   &codespec.StringAttribute{},
				ComputedOptionalRequired: codespec.Optional,
				NonUpdatable:             false,
			},
			hasPlanModifier: false,
		},
		"String attribute with non-updatable - uses CreateOnly()": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_string",
				TFModelName:              "TestString",
				String:                   &codespec.StringAttribute{},
				ComputedOptionalRequired: codespec.Optional,
				NonUpdatable:             true,
			},
			hasPlanModifier: true,
		},
		"Bool attribute with non-updatable but no default - uses CreateOnly()": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_bool",
				TFModelName:              "TestBool",
				Bool:                     &codespec.BoolAttribute{},
				ComputedOptionalRequired: codespec.Optional,
				NonUpdatable:             true,
			},
			hasPlanModifier: true,
		},
		"Bool attribute with non-updatable and default true - uses CreateOnlyBoolWithDefault(true)": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_bool",
				TFModelName:              "TestBool",
				Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(true)},
				ComputedOptionalRequired: codespec.ComputedOptional,
				NonUpdatable:             true,
			},
			hasPlanModifier: true,
		},
		"Bool attribute with non-updatable and default false - uses CreateOnlyBoolWithDefault(false)": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_bool",
				TFModelName:              "TestBool",
				Bool:                     &codespec.BoolAttribute{Default: conversion.Pointer(false)},
				ComputedOptionalRequired: codespec.ComputedOptional,
				NonUpdatable:             true,
			},
			hasPlanModifier: true,
		},
		"Int64 attribute with non-updatable - uses CreateOnly()": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_int",
				TFModelName:              "TestInt",
				Int64:                    &codespec.Int64Attribute{},
				ComputedOptionalRequired: codespec.Optional,
				NonUpdatable:             true,
			},
			hasPlanModifier: true,
		},
		"Computed attribute with non-updatable - uses CreateOnly() (model is enforced)": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_computed",
				TFModelName:              "TestComputed",
				String:                   &codespec.StringAttribute{},
				ComputedOptionalRequired: codespec.Computed,
				NonUpdatable:             true,
			},
			hasPlanModifier: true,
		},
		"ComputedOptional attribute with non-updatable - uses CreateOnly()": {
			attribute: codespec.Attribute{
				TFSchemaName:             "test_computed_optional",
				TFModelName:              "TestComputedOptional",
				String:                   &codespec.StringAttribute{},
				ComputedOptionalRequired: codespec.ComputedOptional,
				NonUpdatable:             true,
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
			if tc.attribute.NonUpdatable {
				assert.Contains(t, code, "customplanmodifier.CreateOnly()")
			}
		})
	}
}
