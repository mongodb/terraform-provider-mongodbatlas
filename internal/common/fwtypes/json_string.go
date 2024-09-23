package fwtypes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

var (
	_ basetypes.StringTypable                    = (*jsonStringType)(nil)
	_ basetypes.StringValuable                   = (*JSONString)(nil)
	_ basetypes.StringValuableWithSemanticEquals = (*JSONString)(nil)
)

type jsonStringType struct {
	basetypes.StringType
}

var (
	JSONStringType = jsonStringType{}
)

func (t jsonStringType) Equal(o attr.Type) bool {
	other, ok := o.(jsonStringType)
	if !ok {
		return false
	}
	return t.StringType.Equal(other.StringType)
}

func (t jsonStringType) String() string {
	return "jsonStringType"
}

func (t jsonStringType) ValueFromString(_ context.Context, in types.String) (basetypes.StringValuable, diag.Diagnostics) {
	var diags diag.Diagnostics
	if in.IsNull() {
		return JSONStringNull(), diags
	}
	if in.IsUnknown() {
		return JSONStringUnknown(), diags
	}
	return JSONString{StringValue: in}, diags
}

func (t jsonStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}
	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}
	return stringValuable, nil
}

func (t jsonStringType) ValueType(context.Context) attr.Value {
	return JSONString{}
}

func (t jsonStringType) Validate(ctx context.Context, in tftypes.Value, attrPath path.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	if !in.IsKnown() || in.IsNull() {
		return diags
	}
	var value string
	err := in.As(&value)
	if err != nil {
		diags.AddAttributeError(
			attrPath,
			"Invalid Terraform Value",
			"An unexpected error occurred while attempting to convert a Terraform value to a string. "+
				"This generally is an issue with the provider schema implementation. "+
				"Please contact the provider developers.\n\n"+
				"Path: "+attrPath.String()+"\n"+
				"Error: "+err.Error(),
		)
		return diags
	}
	if !json.Valid([]byte(value)) {
		diags.AddAttributeError(
			attrPath,
			"Invalid JSON String Value",
			"A string value was provided that is not valid JSON string format (RFC 7159).\n\n"+
				"Path: "+attrPath.String()+"\n"+
				"Given Value: "+value+"\n",
		)
		return diags
	}
	return diags
}

func JSONStringNull() JSONString {
	return JSONString{StringValue: basetypes.NewStringNull()}
}

func JSONStringUnknown() JSONString {
	return JSONString{StringValue: basetypes.NewStringUnknown()}
}

func JSONStringValue(value string) JSONString {
	return JSONString{StringValue: basetypes.NewStringValue(value)}
}

type JSONString struct {
	basetypes.StringValue
}

func (v JSONString) Equal(o attr.Value) bool {
	other, ok := o.(JSONString)
	if !ok {
		return false
	}
	return v.StringValue.Equal(other.StringValue)
}

func (v JSONString) Type(context.Context) attr.Type {
	return JSONStringType
}

func (v JSONString) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	newValue, ok := newValuable.(JSONString)
	if !ok {
		return false, diags
	}
	return schemafunc.EqualJSON(v.ValueString(), newValue.ValueString(), "JsonString"), diags
}
