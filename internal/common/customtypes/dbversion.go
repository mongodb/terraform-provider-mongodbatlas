package customtypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/utility"
)

var (
	_ basetypes.StringTypable                    = DBVersionStringType{}
	_ basetypes.StringValuableWithSemanticEquals = DBVersionStringValue{}
)

// DBVersionStringType represents the type of an attribute that
// represents the database version for MongoDB. For example, "6" or "6.0".
type DBVersionStringType struct {
	basetypes.StringType
}

func (t DBVersionStringType) Equal(o attr.Type) bool {
	other, ok := o.(DBVersionStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t DBVersionStringType) String() string {
	return "DBVersionStringType"
}

func (t DBVersionStringType) ValueFromString(
	ctx context.Context,
	in basetypes.StringValue,
) (basetypes.StringValuable, diag.Diagnostics) {
	value := DBVersionStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t DBVersionStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t DBVersionStringType) ValueType(ctx context.Context) attr.Value {
	return DBVersionStringValue{}
}

// DBVersionStringValue represents a string that represents the database version for MongoDB. For example, "6" or "6.0"
// Additionally, this type implements Semantic equality for version value in different formats i.e. "6" == "6.0"
var _ basetypes.StringValuable = DBVersionStringValue{}

type DBVersionStringValue struct {
	basetypes.StringValue
}

func (v DBVersionStringValue) Equal(o attr.Value) bool {
	other, ok := o.(DBVersionStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v DBVersionStringValue) Type(ctx context.Context) attr.Type {
	return DBVersionStringType{}
}

func (v DBVersionStringValue) StringSemanticEquals(
	ctx context.Context,
	newValuable basetypes.StringValuable,
) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(DBVersionStringValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	priorVal := utility.FormatMongoDBMajorVersion(v.StringValue.ValueString())

	newVal := utility.FormatMongoDBMajorVersion(newValue.ValueString())

	return strings.EqualFold(newVal, priorVal), diags
}
