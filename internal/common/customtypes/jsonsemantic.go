package customtypes

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// JSONSemanticNormalizedType is a custom string type that compares JSON values
// semantically using json.Unmarshal + reflect.DeepEqual, so "10" and "10.0" are
// treated as equal. This prevents spurious plan diffs when the Atlas API returns
// numeric values with unnecessary decimal points (e.g. 10.0 instead of 10).
type JSONSemanticNormalizedType struct {
	jsontypes.NormalizedType
}

func (t JSONSemanticNormalizedType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := t.NormalizedType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	norm, ok := val.(jsontypes.Normalized)
	if !ok {
		return nil, fmt.Errorf("unexpected value type %T", val)
	}
	return JSONSemanticNormalized{Normalized: norm}, nil
}

func (t JSONSemanticNormalizedType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	val, diags := t.NormalizedType.ValueFromString(ctx, in)
	if diags.HasError() {
		return nil, diags
	}
	norm, ok := val.(jsontypes.Normalized)
	if !ok {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("unexpected type", fmt.Sprintf("%T", val))}
	}
	return JSONSemanticNormalized{Normalized: norm}, nil
}

func (t JSONSemanticNormalizedType) ValueType(_ context.Context) attr.Value {
	return JSONSemanticNormalized{}
}

func (t JSONSemanticNormalizedType) Equal(o attr.Type) bool {
	_, ok := o.(JSONSemanticNormalizedType)
	return ok
}

// JSONSemanticNormalized is the value type for JSONSemanticNormalizedType.
type JSONSemanticNormalized struct {
	jsontypes.Normalized
}

func NewJSONSemanticNormalizedValue(s string) JSONSemanticNormalized {
	return JSONSemanticNormalized{Normalized: jsontypes.NewNormalizedValue(s)}
}

func (v JSONSemanticNormalized) Type(_ context.Context) attr.Type {
	return JSONSemanticNormalizedType{}
}

// StringSemanticEquals returns true if both JSON strings represent the same value
// when decoded. Unlike jsontypes.Normalized, this treats "10" and "10.0" as equal.
func (v JSONSemanticNormalized) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	newVal, ok := newValuable.(JSONSemanticNormalized)
	if !ok {
		return false, nil
	}
	var a, b any
	if err := json.Unmarshal([]byte(v.ValueString()), &a); err != nil {
		return false, nil
	}
	if err := json.Unmarshal([]byte(newVal.ValueString()), &b); err != nil {
		return false, nil
	}
	return reflect.DeepEqual(a, b), nil
}
