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
	Custom Set type used in auto-generated code to enable the generic marshal/unmarshal operations to access the elements' type during conversion.
	Custom types docs: https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/custom

	Usage:
		- Schema definition:
			"sample_string_set": schema.SetAttribute{
				...
				CustomType: customtypes.NewSetType[basetypes.StringValue](ctx),
				ElementType: types.StringType,
			}

		- TF Models:
			type TFModel struct {
				SampleStringSet customtypes.SetValue[basetypes.StringValue] `tfsdk:"sample_string_set"`
				...
			}
*/

var (
	_ basetypes.SetTypable  = SetType[basetypes.StringValue]{}
	_ basetypes.SetValuable = SetValue[basetypes.StringValue]{}
	_ SetValueInterface     = SetValue[basetypes.StringValue]{}
)

type SetType[T attr.Value] struct {
	basetypes.SetType
}

func NewSetType[T attr.Value](ctx context.Context) SetType[T] {
	elemType := getValueElementType[T](ctx)
	return SetType[T]{
		SetType: basetypes.SetType{ElemType: elemType},
	}
}

func (t SetType[T]) Equal(o attr.Type) bool {
	other, ok := o.(SetType[T])
	if !ok {
		return false
	}
	return t.SetType.Equal(other.SetType)
}

func (SetType[T]) String() string {
	var t T
	return fmt.Sprintf("SetType[%T]", t)
}

func (t SetType[T]) ValueFromSet(ctx context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	if in.IsNull() {
		return NewSetValueNull[T](ctx), nil
	}

	if in.IsUnknown() {
		return NewSetValueUnknown[T](ctx), nil
	}

	elemType := getValueElementType[T](ctx)
	baseSetValue, diags := basetypes.NewSetValue(elemType, in.Elements())
	if diags.HasError() {
		return nil, diags
	}

	return SetValue[T]{SetValue: baseSetValue}, nil
}

func (t SetType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.SetType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	setValue, ok := attrValue.(basetypes.SetValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	setValuable, diags := t.ValueFromSet(ctx, setValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting SetValue to SetValuable: %v", diags)
	}

	return setValuable, nil
}

func (t SetType[T]) ValueType(_ context.Context) attr.Value {
	return SetValue[T]{}
}

type SetValue[T attr.Value] struct {
	basetypes.SetValue
}

type SetValueInterface interface {
	basetypes.SetValuable
	NewSetValue(ctx context.Context, value []attr.Value) SetValueInterface
	NewSetValueNull(ctx context.Context) SetValueInterface
	ElementType(ctx context.Context) attr.Type
	Elements() []attr.Value
}

func (v SetValue[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.SetValue.ElementType(ctx) == nil {
		// SetValue created as a zero value (not explicitly initialized), initialize now so conversion does not panic.
		v.SetValue = NewSetValueNull[T](ctx).SetValue
	}
	return v.SetValue.ToTerraformValue(ctx)
}

func (v SetValue[T]) NewSetValue(ctx context.Context, value []attr.Value) SetValueInterface {
	return NewSetValue[T](ctx, value)
}

func NewSetValue[T attr.Value](ctx context.Context, value []attr.Value) SetValue[T] {
	elemType := getValueElementType[T](ctx)

	setValue, diags := basetypes.NewSetValue(elemType, value)
	if diags.HasError() {
		return NewSetValueUnknown[T](ctx)
	}

	return SetValue[T]{SetValue: setValue}
}

func (v SetValue[T]) NewSetValueNull(ctx context.Context) SetValueInterface {
	return NewSetValueNull[T](ctx)
}

func NewSetValueNull[T attr.Value](ctx context.Context) SetValue[T] {
	elemType := getValueElementType[T](ctx)
	return SetValue[T]{SetValue: basetypes.NewSetNull(elemType)}
}

func NewSetValueUnknown[T attr.Value](ctx context.Context) SetValue[T] {
	elemType := getValueElementType[T](ctx)
	return SetValue[T]{SetValue: basetypes.NewSetUnknown(elemType)}
}

func (v SetValue[T]) Equal(o attr.Value) bool {
	other, ok := o.(SetValue[T])
	if !ok {
		return false
	}
	return v.SetValue.Equal(other.SetValue)
}

func (v SetValue[T]) Type(ctx context.Context) attr.Type {
	return NewSetType[T](ctx)
}

func (v SetValue[T]) ElementType(ctx context.Context) attr.Type {
	return getValueElementType[T](ctx)
}

func (v SetValue[T]) Elements() []attr.Value {
	return v.SetValue.Elements()
}
