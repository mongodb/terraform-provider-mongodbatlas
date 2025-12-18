package autogen_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

func TestBuildQueryParamMap(t *testing.T) {
	ctx := context.Background()
	listValue, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("WORKFORCE"),
		types.StringValue("WORKLOAD"),
	})
	setValue, _ := types.SetValue(types.StringType, []attr.Value{
		types.StringValue("SAML"),
		types.StringValue("OIDC"),
	})
	listValue2, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("TYPE1"),
	})
	listIntValue, _ := types.ListValue(types.Int64Type, []attr.Value{
		types.Int64Value(1),
		types.Int64Value(2),
		types.Int64Value(3),
	})
	listBoolValue, _ := types.ListValue(types.BoolType, []attr.Value{
		types.BoolValue(true),
		types.BoolValue(false),
	})
	setIntValue, _ := types.SetValue(types.Int64Type, []attr.Value{
		types.Int64Value(10),
		types.Int64Value(20),
	})
	setBoolValue, _ := types.SetValue(types.BoolType, []attr.Value{
		types.BoolValue(true),
	})

	tests := []struct {
		expectedOutput map[string]string
		customAssert   func(*testing.T, map[string]string)
		name           string
		args           []autogen.QueryParamArg
	}{
		{
			name:           "should return empty map when no args provided",
			args:           []autogen.QueryParamArg{},
			expectedOutput: map[string]string{},
		},
		{
			name: "should handle string type",
			args: []autogen.QueryParamArg{
				{APIName: "status", Value: types.StringValue("active")},
			},
			expectedOutput: map[string]string{"status": "active"},
		},
		{
			name: "should handle int64 type",
			args: []autogen.QueryParamArg{
				{APIName: "pageSize", Value: types.Int64Value(100)},
			},
			expectedOutput: map[string]string{"pageSize": "100"},
		},
		{
			name: "should handle bool type",
			args: []autogen.QueryParamArg{
				{APIName: "includeDeleted", Value: types.BoolValue(true)},
			},
			expectedOutput: map[string]string{"includeDeleted": "true"},
		},
		{
			name: "should handle list type",
			args: []autogen.QueryParamArg{
				{APIName: "idpType", Value: listValue},
			},
			expectedOutput: map[string]string{"idpType": "WORKFORCE,WORKLOAD"},
		},
		{
			name: "should handle set type",
			args: []autogen.QueryParamArg{
				{APIName: "protocol", Value: setValue},
			},
			customAssert: func(t *testing.T, result map[string]string) {
				t.Helper()
				assert.Contains(t, result["protocol"], "SAML")
				assert.Contains(t, result["protocol"], "OIDC")
			},
		},
		{
			name: "should handle list of integers",
			args: []autogen.QueryParamArg{
				{APIName: "ids", Value: listIntValue},
			},
			expectedOutput: map[string]string{"ids": "1,2,3"},
		},
		{
			name: "should handle list of bools",
			args: []autogen.QueryParamArg{
				{APIName: "flags", Value: listBoolValue},
			},
			expectedOutput: map[string]string{"flags": "true,false"},
		},
		{
			name: "should handle set of integers",
			args: []autogen.QueryParamArg{
				{APIName: "ports", Value: setIntValue},
			},
			customAssert: func(t *testing.T, result map[string]string) {
				t.Helper()
				assert.Contains(t, result["ports"], "10")
				assert.Contains(t, result["ports"], "20")
			},
		},
		{
			name: "should handle set of bools",
			args: []autogen.QueryParamArg{
				{APIName: "enabled", Value: setBoolValue},
			},
			expectedOutput: map[string]string{"enabled": "true"},
		},
		{
			name: "should skip null values",
			args: []autogen.QueryParamArg{
				{APIName: "status", Value: types.StringNull()},
				{APIName: "pageSize", Value: types.Int64Null()},
				{APIName: "enabled", Value: types.BoolNull()},
			},
			expectedOutput: map[string]string{},
		},
		{
			name: "should skip unknown values",
			args: []autogen.QueryParamArg{
				{APIName: "status", Value: types.StringUnknown()},
			},
			expectedOutput: map[string]string{},
		},
		{
			name: "should handle mixed types in single call",
			args: []autogen.QueryParamArg{
				{APIName: "name", Value: types.StringValue("test")},
				{APIName: "limit", Value: types.Int64Value(50)},
				{APIName: "active", Value: types.BoolValue(true)},
				{APIName: "types", Value: listValue2},
				{APIName: "optional", Value: types.StringNull()},
			},
			expectedOutput: map[string]string{
				"name":   "test",
				"limit":  "50",
				"active": "true",
				"types":  "TYPE1",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := autogen.BuildQueryParamMap(ctx, tc.args)

			if tc.customAssert != nil {
				tc.customAssert(t, result)
			} else {
				assert.Equal(t, tc.expectedOutput, result)
			}
		})
	}
}
