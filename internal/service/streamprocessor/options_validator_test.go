package streamprocessor_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprocessor"
	"github.com/stretchr/testify/assert"
)

func TestOptionsValidator(t *testing.T) {
	testCases := map[string]struct {
		value   types.Object
		wantErr bool
	}{
		"null options": {
			value: types.ObjectNull(streamprocessor.OptionsObjectType.AttributeTypes()),
		},
		"empty options": {
			value: types.ObjectValueMust(streamprocessor.OptionsObjectType.AttributeTypes(), map[string]attr.Value{
				"dlq":                    types.ObjectNull(streamprocessor.DlqObjectType.AttributeTypes()),
				"autoscaling":            types.ObjectNull(streamprocessor.AutoscalingObjectType.AttributeTypes()),
				"resume_from_checkpoint": types.BoolNull(),
			}),
			wantErr: true,
		},
		"options with DLQ": {
			value: types.ObjectValueMust(streamprocessor.OptionsObjectType.AttributeTypes(), map[string]attr.Value{
				"dlq": types.ObjectValueMust(streamprocessor.DlqObjectType.AttributeTypes(), map[string]attr.Value{
					"coll":            types.StringValue("dlq"),
					"connection_name": types.StringValue("connection"),
					"db":              types.StringValue("db"),
				}),
				"autoscaling":            types.ObjectNull(streamprocessor.AutoscalingObjectType.AttributeTypes()),
				"resume_from_checkpoint": types.BoolNull(),
			}),
		},
		"options with autoscaling": {
			value: types.ObjectValueMust(streamprocessor.OptionsObjectType.AttributeTypes(), map[string]attr.Value{
				"dlq": types.ObjectNull(streamprocessor.DlqObjectType.AttributeTypes()),
				"autoscaling": types.ObjectValueMust(streamprocessor.AutoscalingObjectType.AttributeTypes(), map[string]attr.Value{
					"min_tier": types.StringValue("SP10"),
					"max_tier": types.StringValue("SP30"),
				}),
				"resume_from_checkpoint": types.BoolNull(),
			}),
		},
		"options with only resume_from_checkpoint": {
			value: types.ObjectValueMust(streamprocessor.OptionsObjectType.AttributeTypes(), map[string]attr.Value{
				"dlq":                    types.ObjectNull(streamprocessor.DlqObjectType.AttributeTypes()),
				"autoscaling":            types.ObjectNull(streamprocessor.AutoscalingObjectType.AttributeTypes()),
				"resume_from_checkpoint": types.BoolValue(false),
			}),
		},
		"options with a future attribute": {
			value: types.ObjectValueMust(map[string]attr.Type{
				"dlq":         streamprocessor.DlqObjectType,
				"autoscaling": streamprocessor.AutoscalingObjectType,
				"future":      types.StringType,
			}, map[string]attr.Value{
				"dlq":         types.ObjectNull(streamprocessor.DlqObjectType.AttributeTypes()),
				"autoscaling": types.ObjectNull(streamprocessor.AutoscalingObjectType.AttributeTypes()),
				"future":      types.StringValue("value"),
			}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			resp := &validator.ObjectResponse{Diagnostics: diag.Diagnostics{}}
			streamprocessor.OptionsValidator().ValidateObject(t.Context(), validator.ObjectRequest{ConfigValue: tc.value}, resp)
			assert.Equal(t, tc.wantErr, resp.Diagnostics.HasError())
		})
	}
}
