package dynamicjson_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dynamicjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToJSON_Primitives(t *testing.T) {
	cases := map[string]struct {
		in   types.Dynamic
		want string
	}{
		"null":    {types.DynamicNull(), "null"},
		"unknown": {types.DynamicUnknown(), "null"},
		"bool":    {types.DynamicValue(types.BoolValue(true)), "true"},
		"string":  {types.DynamicValue(types.StringValue("hello")), `"hello"`},
		"int64":   {types.DynamicValue(types.Int64Value(42)), "42"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := dynamicjson.ToJSON(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, string(got))
		})
	}
}

func TestToJSON_ObjectCanonical(t *testing.T) {
	obj, diags := types.ObjectValue(
		map[string]attr.Type{
			"b": types.Int64Type,
			"a": types.StringType,
		},
		map[string]attr.Value{
			"b": types.Int64Value(2),
			"a": types.StringValue("x"),
		},
	)
	require.False(t, diags.HasError())
	got, err := dynamicjson.ToJSON(types.DynamicValue(obj))
	require.NoError(t, err)
	assert.JSONEq(t, `{"a":"x","b":2}`, string(got))
}

func TestFromJSON_Inferred(t *testing.T) {
	got, err := dynamicjson.FromJSON([]byte(`{"a":"x","b":2,"c":[true,null]}`), nil)
	require.NoError(t, err)
	roundTrip, err := dynamicjson.ToJSON(got)
	require.NoError(t, err)
	assert.JSONEq(t, `{"a":"x","b":2,"c":[true,null]}`, string(roundTrip))
}

func TestFromJSON_PriorObjectType(t *testing.T) {
	priorType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":  types.StringType,
		"count": types.Int64Type,
		"extra": types.StringType,
	}}
	got, err := dynamicjson.FromJSON([]byte(`{"name":"x","count":3}`), priorType)
	require.NoError(t, err)
	// extra missing from JSON → null of prior type, count comes back as int64
	roundTrip, err := dynamicjson.ToJSON(got)
	require.NoError(t, err)
	assert.JSONEq(t, `{"count":3,"extra":null,"name":"x"}`, string(roundTrip))
}

func TestSemanticallyEqual(t *testing.T) {
	a, err := dynamicjson.FromJSON([]byte(`{"a":1,"b":[1,2]}`), nil)
	require.NoError(t, err)
	b, err := dynamicjson.FromJSON([]byte(`{"b":[1,2],"a":1}`), nil)
	require.NoError(t, err)
	assert.True(t, dynamicjson.SemanticallyEqual(a, b))

	c, err := dynamicjson.FromJSON([]byte(`{"a":1,"b":[2,1]}`), nil)
	require.NoError(t, err)
	assert.False(t, dynamicjson.SemanticallyEqual(a, c))
}
