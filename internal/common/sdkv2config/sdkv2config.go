// Package sdkv2config reads Terraform HCL from SDKv2 GetRawConfig as Terraform
// tri-state values: null (unset in HCL), unknown (expression not yet resolved),
// or set. Use it when Get/GetOk/GetOkExists cannot tell unset from a zero value,
// or when Optional+Computed Get still carries prior state.
package sdkv2config

import (
	"math/big"

	"github.com/hashicorp/go-cty/cty"
)

// resourceView is satisfied by *schema.ResourceData and *schema.ResourceDiff.
type resourceView interface {
	GetRawConfig() cty.Value
	HasChange(key string) bool
}

// state classifies one attribute read from raw config.
type state struct {
	null, unknown, changed bool
}

// IsNull reports whether the attribute is unset in HCL.
func (s state) IsNull() bool { return s.null }

// IsUnknown reports whether the config expression is not yet resolved.
// Apply-time config is usually resolved, but the ApplyResourceChange contract
// permits unknown values, so request code must not fabricate a value.
func (s state) IsUnknown() bool { return s.unknown }

// HasChange reports ResourceData/ResourceDiff.HasChange for this attribute.
func (s state) HasChange() bool { return s.changed }

// Set reports a known value set in HCL. Use for Optional+Computed attributes:
// an omitted or unknown value is left out of the request.
func (s state) Set() bool { return !s.unknown && !s.null }

// Removed reports an attribute dropped from HCL while state still holds a
// value. Use to send an explicit API nil for the cleared field.
func (s state) Removed() bool { return s.null && s.changed }

// SetOrRemoved reports whether a request should carry the attribute's value:
// set in HCL, or removed from it (send the zero value). Use for Optional-only
// attributes whose value is always sent; skipping a removal leaves the server
// value in place and drifts. Unknown is skipped so no value is fabricated.
func (s state) SetOrRemoved() bool { return !s.unknown && (!s.null || s.changed) }

// StringValue is a raw-config string: null, unknown, or set.
type StringValue struct {
	value string
	state
}

// ValueString returns the HCL value; the Go zero value when null or unknown.
func (v StringValue) ValueString() string { return v.value }

// BoolValue is a raw-config bool: null, unknown, or set.
type BoolValue struct {
	state
	value bool
}

// ValueBool returns the HCL value; false when null or unknown.
func (v BoolValue) ValueBool() bool { return v.value }

// Int64Value is a raw-config int: null, unknown, or set.
type Int64Value struct {
	state
	value int64
}

// ValueInt64 returns the HCL value; 0 when null or unknown.
func (v Int64Value) ValueInt64() int64 { return v.value }

// String reads name from raw config. Missing, JSON null, and wrong type are null.
// Empty string is a set value. Unknown stays unknown.
func String(d resourceView, name string) StringValue {
	v := StringValue{state: state{changed: d.HasChange(name)}}
	raw := attrAt(d.GetRawConfig(), name)
	switch {
	case !raw.IsKnown():
		v.unknown = true
	case raw.IsNull() || raw.Type() != cty.String:
		v.null = true
	default:
		v.value = raw.AsString()
	}
	return v
}

// Bool reads name from raw config. Missing, JSON null, and wrong type are null.
// false is a set value. ValueBool on null is false, so a Removed value PATCHes false.
func Bool(d resourceView, name string) BoolValue {
	v := BoolValue{state: state{changed: d.HasChange(name)}}
	raw := attrAt(d.GetRawConfig(), name)
	switch {
	case !raw.IsKnown():
		v.unknown = true
	case raw.IsNull() || raw.Type() != cty.Bool:
		v.null = true
	default:
		v.value = raw.True()
	}
	return v
}

// Int64 reads name from raw config. Schema TypeInt is a cty Number.
// Missing, JSON null, and wrong type are null. 0 is a set value.
// A number that is not an exact int64 (fractional or out of range) is unknown,
// so no predicate can send a truncated or clamped value.
func Int64(d resourceView, name string) Int64Value {
	v := Int64Value{state: state{changed: d.HasChange(name)}}
	raw := attrAt(d.GetRawConfig(), name)
	switch {
	case !raw.IsKnown():
		v.unknown = true
	case raw.IsNull() || raw.Type() != cty.Number:
		v.null = true
	default:
		if n, acc := raw.AsBigFloat().Int64(); acc == big.Exact {
			v.value = n
		} else {
			v.unknown = true
		}
	}
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
