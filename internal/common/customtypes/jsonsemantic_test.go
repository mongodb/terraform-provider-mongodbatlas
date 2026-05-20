package customtypes_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONSemanticNormalized_IntegerEqualsFloat(t *testing.T) {
	v := customtypes.NewJSONSemanticNormalizedValue(`{"n":10}`)
	other := customtypes.NewJSONSemanticNormalizedValue(`{"n":10.0}`)
	equal, diags := v.StringSemanticEquals(context.Background(), other)
	require.False(t, diags.HasError())
	assert.True(t, equal, `"10" and "10.0" should be semantically equal`)
}

func TestJSONSemanticNormalized_DifferentValues(t *testing.T) {
	v := customtypes.NewJSONSemanticNormalizedValue(`{"n":10}`)
	other := customtypes.NewJSONSemanticNormalizedValue(`{"n":11}`)
	equal, diags := v.StringSemanticEquals(context.Background(), other)
	require.False(t, diags.HasError())
	assert.False(t, equal)
}

func TestJSONSemanticNormalized_WhitespaceIgnored(t *testing.T) {
	v := customtypes.NewJSONSemanticNormalizedValue(`{"a":1,"b":2}`)
	other := customtypes.NewJSONSemanticNormalizedValue(`{ "a": 1, "b": 2 }`)
	equal, diags := v.StringSemanticEquals(context.Background(), other)
	require.False(t, diags.HasError())
	assert.True(t, equal)
}

func TestJSONSemanticNormalized_NestedObjects(t *testing.T) {
	v := customtypes.NewJSONSemanticNormalizedValue(`[{"$match":{"size":10}}]`)
	other := customtypes.NewJSONSemanticNormalizedValue(`[{"$match":{"size":10.0}}]`)
	equal, diags := v.StringSemanticEquals(context.Background(), other)
	require.False(t, diags.HasError())
	assert.True(t, equal)
}

func TestJSONSemanticNormalized_NonJSONSemanticNormalizedType(t *testing.T) {
	v := customtypes.NewJSONSemanticNormalizedValue(`{"n":10}`)
	other := jsontypes.NewNormalizedValue(`{"n":10}`)
	equal, diags := v.StringSemanticEquals(context.Background(), other)
	require.False(t, diags.HasError())
	assert.False(t, equal, "different value type should not be equal")
}

func TestJSONSemanticNormalized_TypeReturnsCorrectType(t *testing.T) {
	v := customtypes.NewJSONSemanticNormalizedValue(`{}`)
	_, ok := v.Type(context.Background()).(customtypes.JSONSemanticNormalizedType)
	assert.True(t, ok)
}

func TestJSONSemanticNormalizedType_Equal(t *testing.T) {
	t1 := customtypes.JSONSemanticNormalizedType{}
	t2 := customtypes.JSONSemanticNormalizedType{}
	assert.True(t, t1.Equal(t2))
	assert.False(t, t1.Equal(jsontypes.NormalizedType{}))
}

func TestJSONSemanticNormalizedType_ValueFromTerraform(t *testing.T) {
	typ := customtypes.JSONSemanticNormalizedType{}
	val, err := typ.ValueFromTerraform(context.Background(), tftypes.NewValue(tftypes.String, `{"n":1}`))
	require.NoError(t, err)
	result, ok := val.(customtypes.JSONSemanticNormalized)
	require.True(t, ok, "expected JSONSemanticNormalized, got %T", val)
	assert.Equal(t, `{"n":1}`, result.ValueString())
}

func TestJSONSemanticNormalizedType_ValueType(t *testing.T) {
	typ := customtypes.JSONSemanticNormalizedType{}
	_, ok := typ.ValueType(context.Background()).(customtypes.JSONSemanticNormalized)
	assert.True(t, ok)
}
