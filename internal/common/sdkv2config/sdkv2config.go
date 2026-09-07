// Package sdkv2config reads Terraform HCL from SDKv2 GetRawConfig as TPF-like values.
// Use it when Get/GetOk/GetOkExists cannot tell unset from a zero value, or when
// Optional+Computed Get still carries prior state.
package sdkv2config

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceView interface {
	GetRawConfig() cty.Value
	HasChange(key string) bool
}

type hasChange bool

// HasChange reports ResourceData/ResourceDiff.HasChange for this attribute.
func (h hasChange) HasChange() bool { return bool(h) }

// StringValue embeds TPF types.String so IsNull, IsUnknown, and ValueString promote.
type StringValue struct {
	types.String
	hasChange
}

// BoolValue embeds TPF types.Bool so IsNull, IsUnknown, and ValueBool promote.
type BoolValue struct {
	types.Bool
	hasChange
}

// Int64Value embeds TPF types.Int64 so IsNull, IsUnknown, and ValueInt64 promote.
type Int64Value struct {
	types.Int64
	hasChange
}

// String reads name from raw config. Missing, JSON null, and wrong type are null.
// Empty string is a set value. Unknown stays unknown.
func String(d resourceView, name string) StringValue {
	v := StringValue{hasChange: hasChange(d.HasChange(name))}
	raw := attrAt(d.GetRawConfig(), name)
	if !raw.IsKnown() {
		v.String = types.StringUnknown()
		return v
	}
	if raw.IsNull() || raw.Type() != cty.String {
		v.String = types.StringNull()
		return v
	}
	v.String = types.StringValue(raw.AsString())
	return v
}

// Bool reads name from raw config. Missing, JSON null, and wrong type are null.
// false is a set value. ValueBool on null is false, so null+changed PATCHes false.
func Bool(d resourceView, name string) BoolValue {
	v := BoolValue{hasChange: hasChange(d.HasChange(name))}
	raw := attrAt(d.GetRawConfig(), name)
	if !raw.IsKnown() {
		v.Bool = types.BoolUnknown()
		return v
	}
	if raw.IsNull() || raw.Type() != cty.Bool {
		v.Bool = types.BoolNull()
		return v
	}
	v.Bool = types.BoolValue(raw.True())
	return v
}

// Int64 reads name from raw config. Schema TypeInt is a cty Number.
// Missing, JSON null, and wrong type are null. 0 is a set value.
func Int64(d resourceView, name string) Int64Value {
	v := Int64Value{hasChange: hasChange(d.HasChange(name))}
	raw := attrAt(d.GetRawConfig(), name)
	if !raw.IsKnown() {
		v.Int64 = types.Int64Unknown()
		return v
	}
	if raw.IsNull() || raw.Type() != cty.Number {
		v.Int64 = types.Int64Null()
		return v
	}
	n, _ := raw.AsBigFloat().Int64()
	v.Int64 = types.Int64Value(n)
	return v
}

// CollectionEmpty is true when the named list or set is not a known non-empty collection:
// attr null, attr unknown, or known length 0. If the resource raw-config object is null or
// unknown, it returns false: the attribute cannot be read, so this is not an HCL omit.
func CollectionEmpty(d resourceView, name string) bool {
	raw := d.GetRawConfig()
	if raw.IsNull() || !raw.IsKnown() {
		return false
	}
	return knownCollectionLen(attrAt(raw, name)) == 0
}

// NestedCollectionLen is the known length of list[index].attr in raw config.
// Parent list must be a list or tuple. Nested sets are measured with LengthInt only.
// Returns 0 when the list is omitted, unknown, out of range, or the nested attr is
// omitted, unknown, or empty.
func NestedCollectionLen(d resourceView, list string, index int, attr string) int {
	listVal := attrAt(d.GetRawConfig(), list)
	if listVal.IsNull() || !listVal.IsKnown() {
		return 0
	}
	ty := listVal.Type()
	if !ty.IsListType() && !ty.IsTupleType() {
		return 0
	}
	if index < 0 || index >= listVal.LengthInt() {
		return 0
	}
	return knownCollectionLen(attrAt(listVal.Index(cty.NumberIntVal(int64(index))), attr))
}

// attrAt walks one object attribute. It returns null on a missing name, null root, or
// non-object, and unknown when the parent is unknown. It does not panic.
func attrAt(raw cty.Value, name string) cty.Value {
	if !raw.IsKnown() {
		return cty.UnknownVal(cty.DynamicPseudoType)
	}
	if raw.IsNull() || !raw.Type().IsObjectType() {
		return cty.NullVal(cty.DynamicPseudoType)
	}
	if _, ok := raw.Type().AttributeTypes()[name]; !ok {
		return cty.NullVal(cty.DynamicPseudoType)
	}
	return raw.GetAttr(name)
}

func knownCollectionLen(v cty.Value) int {
	if v.IsNull() || !v.IsKnown() || !v.CanIterateElements() {
		return 0
	}
	return v.LengthInt()
}
