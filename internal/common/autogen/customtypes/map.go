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
	Custom Map type used in auto-generated code to enable the generic marshal/unmarshal operations to access the elements' type during conversion.
	Custom types docs: https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/custom

	Usage:
		- Schema definition:
			"sample_string_map": schema.MapAttribute{
				...
				CustomType: customtypes.NewMapType[basetypes.StringValue](ctx),
				ElementType: types.StringType,
			}

		- TF Models:
			type TFModel struct {
				SampleStringMap customtypes.MapValue[basetypes.StringValue] `tfsdk:"sample_string_map"`
				...
			}
*/

var (
	_ basetypes.MapTypable  = MapType[basetypes.StringValue]{}
	_ basetypes.MapValuable = MapValue[basetypes.StringValue]{}
	_ MapValueInterface     = MapValue[basetypes.StringValue]{}
)

type MapType[T attr.Value] struct {
	basetypes.MapType
}

func NewMapType[T attr.Value](ctx context.Context) MapType[T] {
	elemType := getValueType[T](ctx)
	return MapType[T]{
		MapType: basetypes.MapType{ElemType: elemType},
	}
}

func (t MapType[T]) Equal(o attr.Type) bool {
	other, ok := o.(MapType[T])
	if !ok {
		return false
	}

	return t.MapType.Equal(other.MapType)
}

func (MapType[T]) String() string {
	var t T
	return fmt.Sprintf("MapType[%T]", t)
}

func (t MapType[T]) ValueFromMap(ctx context.Context, in basetypes.MapValue) (basetypes.MapValuable, diag.Diagnostics) {
	if in.IsNull() {
		return NewMapValueNull[T](ctx), nil
	}

	if in.IsUnknown() {
		return NewMapValueUnknown[T](ctx), nil
	}

	elemType := getValueType[T](ctx)
	baseMapValue, diags := basetypes.NewMapValue(elemType, in.Elements())
	if diags.HasError() {
		return nil, diags
	}

	return MapValue[T]{MapValue: baseMapValue}, nil
}

func (t MapType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.MapType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	mapValue, ok := attrValue.(basetypes.MapValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	mapValuable, diags := t.ValueFromMap(ctx, mapValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting MapValue to MapValuable: %v", diags)
	}

	return mapValuable, nil
}

func (t MapType[T]) ValueType(_ context.Context) attr.Value {
	return MapValue[T]{}
}

type MapValue[T attr.Value] struct {
	basetypes.MapValue
}

type MapValueInterface interface {
	basetypes.MapValuable
	NewMapValue(ctx context.Context, value map[string]attr.Value) MapValueInterface
	NewMapValueNull(ctx context.Context) MapValueInterface
	ElementType(ctx context.Context) attr.Type
	Elements() map[string]attr.Value
}

func (v MapValue[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.MapValue.ElementType(ctx) == nil {
		// MapValue created as a zero value (not explicitly initialized), initialize now so conversion does not panic.
		v.MapValue = NewMapValueNull[T](ctx).MapValue
	}
	return v.MapValue.ToTerraformValue(ctx)
}

func (v MapValue[T]) NewMapValue(ctx context.Context, value map[string]attr.Value) MapValueInterface {
	return NewMapValue[T](ctx, value)
}

func NewMapValue[T attr.Value](ctx context.Context, value map[string]attr.Value) MapValue[T] {
	elemType := getValueType[T](ctx)

	mapValue, diags := basetypes.NewMapValue(elemType, value)
	if diags.HasError() {
		return NewMapValueUnknown[T](ctx)
	}

	return MapValue[T]{MapValue: mapValue}
}

func (v MapValue[T]) NewMapValueNull(ctx context.Context) MapValueInterface {
	return NewMapValueNull[T](ctx)
}

func NewMapValueNull[T attr.Value](ctx context.Context) MapValue[T] {
	elemType := getValueType[T](ctx)
	return MapValue[T]{MapValue: basetypes.NewMapNull(elemType)}
}

func NewMapValueUnknown[T attr.Value](ctx context.Context) MapValue[T] {
	elemType := getValueType[T](ctx)
	return MapValue[T]{MapValue: basetypes.NewMapUnknown(elemType)}
}

func (v MapValue[T]) Equal(o attr.Value) bool {
	other, ok := o.(MapValue[T])
	if !ok {
		return false
	}
	return v.MapValue.Equal(other.MapValue)
}

func (v MapValue[T]) Type(ctx context.Context) attr.Type {
	return NewMapType[T](ctx)
}

func (v MapValue[T]) ElementType(ctx context.Context) attr.Type {
	return getValueType[T](ctx)
}

func (v MapValue[T]) Elements() map[string]attr.Value {
	return v.MapValue.Elements()
}
