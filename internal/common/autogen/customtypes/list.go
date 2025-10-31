package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

/*
	Custom List type used in auto-generated code to enable the generic marshal/unmarshal operations to access the elements' type during conversion.
	Custom types docs: https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/custom

	Usage:
		- Schema definition:
			"sample_string_list": schema.ListAttribute{
				...
				CustomType: customtypes.NewListType[basetypes.StringValue](ctx),
				ElementType: types.StringType,
			}

		- TF Models:
			type TFModel struct {
				SampleStringList customtypes.ListValue[basetypes.StringValue] `tfsdk:"sample_string_list"`
				...
			}
*/

var (
	_ basetypes.ListTypable  = ListType[basetypes.StringValue]{}
	_ basetypes.ListValuable = ListValue[basetypes.StringValue]{}
	_ ListValueInterface     = ListValue[basetypes.StringValue]{}
)

type ListType[T attr.Value] struct {
	basetypes.ListType
}

func NewListType[T attr.Value](ctx context.Context) ListType[T] {
	elemType := getElemType[T](ctx)
	return ListType[T]{
		ListType: basetypes.ListType{ElemType: elemType},
	}
}

func (t ListType[T]) Equal(o attr.Type) bool {
	other, ok := o.(ListType[T])
	if !ok {
		return false
	}
	return t.ListType.Equal(other.ListType)
}

func (ListType[T]) String() string {
	var t T
	return fmt.Sprintf("ListType[%T]", t)
}

func (t ListType[T]) ValueFromList(ctx context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	if in.IsNull() {
		return NewListValueNull[T](ctx), nil
	}

	if in.IsUnknown() {
		return NewListValueUnknown[T](ctx), nil
	}

	elemType := getElemType[T](ctx)
	baseListValue, diags := basetypes.NewListValue(elemType, in.Elements())
	if diags.HasError() {
		return nil, diags
	}

	return ListValue[T]{ListValue: baseListValue}, nil
}

func (t ListType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.ListType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	listValue, ok := attrValue.(basetypes.ListValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	listValuable, diags := t.ValueFromList(ctx, listValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ListValue to ListValuable: %v", diags)
	}

	return listValuable, nil
}

func (t ListType[T]) ValueType(_ context.Context) attr.Value {
	return ListValue[T]{}
}

type ListValue[T attr.Value] struct {
	basetypes.ListValue
}

type ListValueInterface interface {
	basetypes.ListValuable
	NewListValue(ctx context.Context, value []attr.Value) ListValueInterface
	NewListValueNull(ctx context.Context) ListValueInterface
	ElementType(ctx context.Context) attr.Type
	Elements() []attr.Value
}

func (v ListValue[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.ListValue.ElementType(ctx) == nil {
		// ListValue created as a zero value (not explicitly initialized), initialize now so conversion does not panic.
		v.ListValue = NewListValueNull[T](ctx).ListValue
	}
	return v.ListValue.ToTerraformValue(ctx)
}

func (v ListValue[T]) NewListValue(ctx context.Context, value []attr.Value) ListValueInterface {
	return NewListValue[T](ctx, value)
}

func NewListValue[T attr.Value](ctx context.Context, value []attr.Value) ListValue[T] {
	elemType := getElemType[T](ctx)

	listValue, diags := basetypes.NewListValue(elemType, value)
	if diags.HasError() {
		return NewListValueUnknown[T](ctx)
	}

	return ListValue[T]{ListValue: listValue}
}

func (v ListValue[T]) NewListValueNull(ctx context.Context) ListValueInterface {
	return NewListValueNull[T](ctx)
}

func NewListValueNull[T attr.Value](ctx context.Context) ListValue[T] {
	elemType := getElemType[T](ctx)
	return ListValue[T]{ListValue: basetypes.NewListNull(elemType)}
}

func NewListValueUnknown[T attr.Value](ctx context.Context) ListValue[T] {
	elemType := getElemType[T](ctx)
	return ListValue[T]{ListValue: basetypes.NewListUnknown(elemType)}
}

func (v ListValue[T]) Equal(o attr.Value) bool {
	other, ok := o.(ListValue[T])
	if !ok {
		return false
	}
	return v.ListValue.Equal(other.ListValue)
}

func (v ListValue[T]) Type(ctx context.Context) attr.Type {
	return NewListType[T](ctx)
}

func (v ListValue[T]) ElementType(ctx context.Context) attr.Type {
	return getElemType[T](ctx)
}

func (v ListValue[T]) Elements() []attr.Value {
	return v.ListValue.Elements()
}

func getElemType[T attr.Value](ctx context.Context) attr.Type {
	var t T
	return t.Type(ctx)
}
