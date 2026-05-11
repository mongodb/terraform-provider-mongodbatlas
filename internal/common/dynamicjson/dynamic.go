// Package dynamicjson bridges Terraform's types.Dynamic and raw JSON.
//
// ToJSON walks an attr.Value tree and emits canonical JSON; FromJSON parses
// JSON back into a types.Dynamic, optionally honoring a prior attr.Type so
// the same payload keeps a stable Terraform type shape across applies.
package dynamicjson

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ToJSON serializes a types.Dynamic into canonical JSON.
// Null and Unknown (including underlying-unknown) become JSON null.
func ToJSON(d types.Dynamic) ([]byte, error) {
	if d.IsNull() || d.IsUnknown() || d.IsUnderlyingValueNull() || d.IsUnderlyingValueUnknown() {
		return []byte("null"), nil
	}
	v, err := attrValueToGo(d.UnderlyingValue())
	if err != nil {
		return nil, err
	}
	return marshalCanonical(v)
}

// FromJSON parses raw JSON into a types.Dynamic.
// If priorType is non-nil, the returned Dynamic mirrors that shape (Object
// vs Map, List vs Tuple, Int64 vs Float64 vs Number) where possible.
// Otherwise a best-effort inferred shape is used: objects→Object, arrays→Tuple,
// numbers→Number, bool→Bool, string→String, null→Dynamic null.
func FromJSON(data []byte, priorType attr.Type) (types.Dynamic, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return types.DynamicNull(), nil
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var raw any
	if err := dec.Decode(&raw); err != nil {
		return types.DynamicNull(), fmt.Errorf("dynamicjson: parse JSON: %w", err)
	}
	val, err := goToAttrValue(raw, priorType)
	if err != nil {
		return types.DynamicNull(), err
	}
	return types.DynamicValue(val), nil
}

// attrValueToGo converts an attr.Value to a plain Go value that
// encoding/json can serialize. Nulls/unknowns become nil.
func attrValueToGo(v attr.Value) (any, error) {
	if v == nil || v.IsNull() || v.IsUnknown() {
		return nil, nil
	}
	switch t := v.(type) {
	case basetypes.BoolValue:
		return t.ValueBool(), nil
	case basetypes.StringValue:
		return t.ValueString(), nil
	case basetypes.Int64Value:
		return json.Number(fmt.Sprintf("%d", t.ValueInt64())), nil
	case basetypes.Int32Value:
		return json.Number(fmt.Sprintf("%d", t.ValueInt32())), nil
	case basetypes.Float64Value:
		return json.Number(fmt.Sprintf("%v", t.ValueFloat64())), nil
	case basetypes.Float32Value:
		return json.Number(fmt.Sprintf("%v", t.ValueFloat32())), nil
	case basetypes.NumberValue:
		return jsonNumberFromBig(t.ValueBigFloat()), nil
	case basetypes.ListValue:
		return elementsToGo(t.Elements())
	case basetypes.SetValue:
		return elementsToGo(t.Elements())
	case basetypes.TupleValue:
		return elementsToGo(t.Elements())
	case basetypes.MapValue:
		return mapToGo(t.Elements())
	case basetypes.ObjectValue:
		return mapToGo(t.Attributes())
	case basetypes.DynamicValue:
		return attrValueToGo(t.UnderlyingValue())
	}
	return nil, fmt.Errorf("dynamicjson: unsupported attr.Value type %T", v)
}

func elementsToGo(elems []attr.Value) ([]any, error) {
	out := make([]any, len(elems))
	for i, e := range elems {
		gv, err := attrValueToGo(e)
		if err != nil {
			return nil, err
		}
		out[i] = gv
	}
	return out, nil
}

func mapToGo(elems map[string]attr.Value) (map[string]any, error) {
	out := make(map[string]any, len(elems))
	for k, v := range elems {
		gv, err := attrValueToGo(v)
		if err != nil {
			return nil, err
		}
		out[k] = gv
	}
	return out, nil
}

func jsonNumberFromBig(f *big.Float) json.Number {
	if f == nil {
		return "0"
	}
	if f.IsInt() {
		i, _ := f.Int64()
		return json.Number(fmt.Sprintf("%d", i))
	}
	return json.Number(f.Text('g', -1))
}

// marshalCanonical serializes a Go tree to JSON with sorted map keys.
func marshalCanonical(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := writeCanonical(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeCanonical(buf *bytes.Buffer, v any) error {
	switch t := v.(type) {
	case nil:
		buf.WriteString("null")
	case bool:
		if t {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case string:
		b, err := json.Marshal(t)
		if err != nil {
			return err
		}
		buf.Write(b)
	case json.Number:
		buf.WriteString(string(t))
	case []any:
		buf.WriteByte('[')
		for i, e := range t {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := writeCanonical(buf, e); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		buf.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			kb, err := json.Marshal(k)
			if err != nil {
				return err
			}
			buf.Write(kb)
			buf.WriteByte(':')
			if err := writeCanonical(buf, t[k]); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return err
		}
		buf.Write(b)
	}
	return nil
}

// goToAttrValue builds an attr.Value tree from JSON output, optionally honoring priorType.
func goToAttrValue(v any, priorType attr.Type) (attr.Value, error) {
	if v == nil {
		if priorType != nil {
			return nullForType(priorType)
		}
		return types.DynamicNull(), nil
	}
	switch t := v.(type) {
	case bool:
		return types.BoolValue(t), nil
	case string:
		return types.StringValue(t), nil
	case json.Number:
		return numberFromJSON(t, priorType)
	case []any:
		return arrayFromJSON(t, priorType)
	case map[string]any:
		return objectFromJSON(t, priorType)
	}
	return nil, fmt.Errorf("dynamicjson: unexpected JSON value type %T", v)
}

func numberFromJSON(n json.Number, priorType attr.Type) (attr.Value, error) {
	switch priorType.(type) {
	case basetypes.Int64Type:
		i, err := n.Int64()
		if err != nil {
			return nil, fmt.Errorf("dynamicjson: cannot parse %q as Int64: %w", n.String(), err)
		}
		return types.Int64Value(i), nil
	case basetypes.Float64Type:
		f, err := n.Float64()
		if err != nil {
			return nil, fmt.Errorf("dynamicjson: cannot parse %q as Float64: %w", n.String(), err)
		}
		return types.Float64Value(f), nil
	case basetypes.StringType:
		return types.StringValue(n.String()), nil
	}
	bf, _, err := big.ParseFloat(n.String(), 10, 512, big.ToNearestEven)
	if err != nil {
		return nil, fmt.Errorf("dynamicjson: cannot parse number %q: %w", n.String(), err)
	}
	return types.NumberValue(bf), nil
}

func arrayFromJSON(arr []any, priorType attr.Type) (attr.Value, error) {
	switch pt := priorType.(type) {
	case basetypes.ListType:
		return listFromJSON(arr, pt.ElemType, false)
	case basetypes.SetType:
		return listFromJSON(arr, pt.ElemType, true)
	case basetypes.TupleType:
		return tupleFromJSON(arr, pt.ElemTypes)
	}
	vals := make([]attr.Value, len(arr))
	elemTypes := make([]attr.Type, len(arr))
	for i, e := range arr {
		v, err := goToAttrValue(e, nil)
		if err != nil {
			return nil, err
		}
		vals[i] = v
		elemTypes[i] = v.Type(context.Background())
	}
	tv, diags := types.TupleValue(elemTypes, vals)
	if diags.HasError() {
		return nil, diagsToErr(diags)
	}
	return tv, nil
}

func listFromJSON(arr []any, elemType attr.Type, isSet bool) (attr.Value, error) {
	vals := make([]attr.Value, len(arr))
	for i, e := range arr {
		v, err := goToAttrValue(e, elemType)
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	if isSet {
		sv, diags := types.SetValue(elemType, vals)
		if diags.HasError() {
			return nil, diagsToErr(diags)
		}
		return sv, nil
	}
	lv, diags := types.ListValue(elemType, vals)
	if diags.HasError() {
		return nil, diagsToErr(diags)
	}
	return lv, nil
}

func tupleFromJSON(arr []any, elemTypes []attr.Type) (attr.Value, error) {
	if len(arr) != len(elemTypes) {
		// Length mismatch — fall back to inferred tuple.
		return arrayFromJSON(arr, nil)
	}
	vals := make([]attr.Value, len(arr))
	for i, e := range arr {
		v, err := goToAttrValue(e, elemTypes[i])
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	tv, diags := types.TupleValue(elemTypes, vals)
	if diags.HasError() {
		return nil, diagsToErr(diags)
	}
	return tv, nil
}

func objectFromJSON(obj map[string]any, priorType attr.Type) (attr.Value, error) {
	switch pt := priorType.(type) {
	case basetypes.ObjectType:
		return objectFromJSONTyped(obj, pt.AttrTypes)
	case basetypes.MapType:
		return mapFromJSON(obj, pt.ElemType)
	}
	attrs := make(map[string]attr.Value, len(obj))
	attrTypes := make(map[string]attr.Type, len(obj))
	for k, e := range obj {
		v, err := goToAttrValue(e, nil)
		if err != nil {
			return nil, err
		}
		attrs[k] = v
		attrTypes[k] = v.Type(context.Background())
	}
	ov, diags := types.ObjectValue(attrTypes, attrs)
	if diags.HasError() {
		return nil, diagsToErr(diags)
	}
	return ov, nil
}

func objectFromJSONTyped(obj map[string]any, attrTypes map[string]attr.Type) (attr.Value, error) {
	attrs := make(map[string]attr.Value, len(attrTypes))
	for name, at := range attrTypes {
		raw, ok := obj[name]
		if !ok {
			// Missing keys become null of the prior type.
			nullVal, err := nullForType(at)
			if err != nil {
				return nil, err
			}
			attrs[name] = nullVal
			continue
		}
		v, err := goToAttrValue(raw, at)
		if err != nil {
			return nil, err
		}
		attrs[name] = v
	}
	ov, diags := types.ObjectValue(attrTypes, attrs)
	if diags.HasError() {
		return nil, diagsToErr(diags)
	}
	return ov, nil
}

func mapFromJSON(obj map[string]any, elemType attr.Type) (attr.Value, error) {
	elems := make(map[string]attr.Value, len(obj))
	for k, raw := range obj {
		v, err := goToAttrValue(raw, elemType)
		if err != nil {
			return nil, err
		}
		elems[k] = v
	}
	mv, diags := types.MapValue(elemType, elems)
	if diags.HasError() {
		return nil, diagsToErr(diags)
	}
	return mv, nil
}

func nullForType(t attr.Type) (attr.Value, error) {
	switch tt := t.(type) {
	case basetypes.BoolType:
		return types.BoolNull(), nil
	case basetypes.StringType:
		return types.StringNull(), nil
	case basetypes.Int64Type:
		return types.Int64Null(), nil
	case basetypes.Float64Type:
		return types.Float64Null(), nil
	case basetypes.NumberType:
		return types.NumberNull(), nil
	case basetypes.ListType:
		return types.ListNull(tt.ElemType), nil
	case basetypes.SetType:
		return types.SetNull(tt.ElemType), nil
	case basetypes.MapType:
		return types.MapNull(tt.ElemType), nil
	case basetypes.ObjectType:
		return types.ObjectNull(tt.AttrTypes), nil
	case basetypes.TupleType:
		return types.TupleNull(tt.ElemTypes), nil
	case basetypes.DynamicType:
		return types.DynamicNull(), nil
	}
	return nil, fmt.Errorf("dynamicjson: cannot build null for type %T", t)
}

func diagsToErr(diags diag.Diagnostics) error {
	errs := diags.Errors()
	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		msgs = append(msgs, e.Summary()+": "+e.Detail())
	}
	return fmt.Errorf("dynamicjson: %v", msgs)
}
