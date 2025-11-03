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
	Custom Nested Map type used in auto-generated code to enable the generic marshal/unmarshal operations to access nested attribute struct tags during conversion.
	Custom types docs: https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/custom

	Usage:
		- Schema definition:
			"sample_nested_object_map": schema.MapNestedAttribute{
				...
				CustomType: customtypes.NewNestedMapType[TFSampleNestedObjectModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string_attribute": schema.StringAttribute{...},
					},
				},
			}

		- TF Models:
			type TFModel struct {
				SampleNestedObjectMap customtypes.NestedMapValue[TFSampleNestedObjectModel] `tfsdk:"sample_nested_object_map"`
				...
			}

			type TFSampleNestedObjectModel struct {
				StringAttribute types.String `tfsdk:"string_attribute"`
				...
			}
*/

var (
	_ basetypes.MapTypable    = NestedMapType[struct{}]{}
	_ basetypes.MapValuable   = NestedMapValue[struct{}]{}
	_ NestedMapValueInterface = NestedMapValue[struct{}]{}
)

type NestedMapType[T any] struct {
	basetypes.MapType
}

func NewNestedMapType[T any](ctx context.Context) NestedMapType[T] {
	elemType, diags := getElementType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating NestedMapType: %v", diags))
	}

	result := NestedMapType[T]{
		MapType: basetypes.MapType{ElemType: elemType},
	}
	return result
}

func (t NestedMapType[T]) Equal(o attr.Type) bool {
	other, ok := o.(NestedMapType[T])
	if !ok {
		return false
	}
	return t.MapType.Equal(other.MapType)
}

func (NestedMapType[T]) String() string {
	var t T
	return fmt.Sprintf("NestedMapType[%T]", t)
}

func (t NestedMapType[T]) ValueFromMap(ctx context.Context, in basetypes.MapValue) (basetypes.MapValuable, diag.Diagnostics) {
	if in.IsNull() {
		return NewNestedMapValueNull[T](ctx), nil
	}

	if in.IsUnknown() {
		return NewNestedMapValueUnknown[T](ctx), nil
	}

	elemType, diags := getElementType[T](ctx)
	if diags.HasError() {
		return nil, diags
	}

	baseMapValue, diags := basetypes.NewMapValue(elemType, in.Elements())
	if diags.HasError() {
		return nil, diags
	}

	return NestedMapValue[T]{MapValue: baseMapValue}, nil
}

func (t NestedMapType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t NestedMapType[T]) ValueType(_ context.Context) attr.Value {
	return NestedMapValue[T]{}
}

type NestedMapValue[T any] struct {
	basetypes.MapValue
}

type NestedMapValueInterface interface {
	basetypes.MapValuable
	NewNestedMapValue(ctx context.Context, value any) NestedMapValueInterface
	NewNestedMapValueNull(ctx context.Context) NestedMapValueInterface
	MapPtrAsAny(ctx context.Context) (any, diag.Diagnostics)
	NewEmptyMapPtr() any
}

func (v NestedMapValue[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.ElementType(ctx) == nil {
		// NestedMapValue created as a zero value (not explicitly initialized), initialize now so conversion does not panic.
		v.MapValue = NewNestedMapValueNull[T](ctx).MapValue
	}
	return v.MapValue.ToTerraformValue(ctx)
}

func (v NestedMapValue[T]) NewNestedMapValue(ctx context.Context, value any) NestedMapValueInterface {
	return NewNestedMapValue[T](ctx, value)
}

func NewNestedMapValue[T any](ctx context.Context, value any) NestedMapValue[T] {
	elemType, diags := getElementType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating NestedMapValue: %v", diags))
	}

	newValue, diags := basetypes.NewMapValueFrom(ctx, elemType, value)
	if diags.HasError() {
		return NewNestedMapValueUnknown[T](ctx)
	}

	return NestedMapValue[T]{MapValue: newValue}
}

func (v NestedMapValue[T]) NewNestedMapValueNull(ctx context.Context) NestedMapValueInterface {
	return NewNestedMapValueNull[T](ctx)
}

func NewNestedMapValueNull[T any](ctx context.Context) NestedMapValue[T] {
	elemType, diags := getElementType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating null NestedMapValue: %v", diags))
	}
	return NestedMapValue[T]{MapValue: basetypes.NewMapNull(elemType)}
}

func NewNestedMapValueUnknown[T any](ctx context.Context) NestedMapValue[T] {
	elemType, diags := getElementType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating unknown NestedMapValue: %v", diags))
	}
	return NestedMapValue[T]{MapValue: basetypes.NewMapUnknown(elemType)}
}

func (v NestedMapValue[T]) Equal(o attr.Value) bool {
	other, ok := o.(NestedMapValue[T])
	if !ok {
		return false
	}
	return v.MapValue.Equal(other.MapValue)
}

func (v NestedMapValue[T]) Type(ctx context.Context) attr.Type {
	return NewNestedMapType[T](ctx)
}

func (v NestedMapValue[T]) MapPtrAsAny(ctx context.Context) (any, diag.Diagnostics) {
	valuePtr := new(map[string]T)

	if v.IsNull() || v.IsUnknown() {
		return valuePtr, nil
	}

	diags := v.ElementsAs(ctx, valuePtr, false)
	if diags.HasError() {
		return nil, diags
	}

	return valuePtr, diags
}

func (v NestedMapValue[T]) NewEmptyMapPtr() any {
	return new(map[string]T)
}
