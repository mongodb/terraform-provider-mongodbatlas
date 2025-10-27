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
	Custom Nested Set type used in auto-generated code to enable the generic marshal/unmarshal operations to access nested attribute struct tags during conversion.
	Custom types docs: https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/custom

	Usage:
		- Schema definition:
			"sample_nested_object_set": schema.SetNestedAttribute{
				...
				CustomType: customtypes.NewNestedSetType[TFSampleNestedObjectModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string_attribute": schema.StringAttribute{...},
					},
				},
			}

		- TF Models:
			type TFModel struct {
				SampleNestedObjectSet customtypes.NestedSetValue[TFSampleNestedObjectModel] `tfsdk:"sample_nested_object_set"`
				...
			}

			type TFSampleNestedObjectModel struct {
				StringAttribute types.String `tfsdk:"string_attribute"`
				...
			}
*/

var (
	_ basetypes.SetTypable    = NestedSetType[struct{}]{}
	_ basetypes.SetValuable   = NestedSetValue[struct{}]{}
	_ NestedSetValueInterface = NestedSetValue[struct{}]{}
)

type NestedSetType[T any] struct {
	basetypes.SetType
}

func NewNestedSetType[T any](ctx context.Context) NestedSetType[T] {
	elemType, diags := getElementType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating NestedSetType: %v", diags))
	}

	result := NestedSetType[T]{
		SetType: basetypes.SetType{ElemType: elemType},
	}
	return result
}

func (t NestedSetType[T]) Equal(o attr.Type) bool {
	other, ok := o.(NestedSetType[T])
	if !ok {
		return false
	}
	return t.SetType.Equal(other.SetType)
}

func (NestedSetType[T]) String() string {
	var t T
	return fmt.Sprintf("NestedSetType[%T]", t)
}

func (t NestedSetType[T]) ValueFromSet(ctx context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	if in.IsNull() {
		return NewNestedSetValueNull[T](ctx), nil
	}

	if in.IsUnknown() {
		return NewNestedSetValueUnknown[T](ctx), nil
	}

	elemType, diags := getElementType[T](ctx)
	if diags.HasError() {
		return nil, diags
	}

	baseSetValue, diags := basetypes.NewSetValue(elemType, in.Elements())
	if diags.HasError() {
		return nil, diags
	}

	return NestedSetValue[T]{SetValue: baseSetValue}, nil
}

func (t NestedSetType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t NestedSetType[T]) ValueType(_ context.Context) attr.Value {
	return NestedSetValue[T]{}
}

type NestedSetValue[T any] struct {
	basetypes.SetValue
}

type NestedSetValueInterface interface {
	basetypes.SetValuable
	NewNestedSetValue(ctx context.Context, value any) NestedSetValueInterface
	NewNestedSetValueNull(ctx context.Context) NestedSetValueInterface
	SlicePtrAsAny(ctx context.Context) (any, diag.Diagnostics)
	NewEmptySlicePtr() any
	Len() int
}

func (v NestedSetValue[T]) NewNestedSetValue(ctx context.Context, value any) NestedSetValueInterface {
	return NewNestedSetValue[T](ctx, value)
}

func NewNestedSetValue[T any](ctx context.Context, value any) NestedSetValue[T] {
	elemType, diags := getElementType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating NestedSetValue: %v", diags))
	}

	newValue, diags := basetypes.NewSetValueFrom(ctx, elemType, value)
	if diags.HasError() {
		return NewNestedSetValueUnknown[T](ctx)
	}

	return NestedSetValue[T]{SetValue: newValue}
}

func (v NestedSetValue[T]) NewNestedSetValueNull(ctx context.Context) NestedSetValueInterface {
	return NewNestedSetValueNull[T](ctx)
}

func NewNestedSetValueNull[T any](ctx context.Context) NestedSetValue[T] {
	elemType, diags := getElementType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating null NestedSetValue: %v", diags))
	}
	return NestedSetValue[T]{SetValue: basetypes.NewSetNull(elemType)}
}

func NewNestedSetValueUnknown[T any](ctx context.Context) NestedSetValue[T] {
	elemType, diags := getElementType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating unknown NestedSetValue: %v", diags))
	}
	return NestedSetValue[T]{SetValue: basetypes.NewSetUnknown(elemType)}
}

func (v NestedSetValue[T]) Equal(o attr.Value) bool {
	other, ok := o.(NestedSetValue[T])
	if !ok {
		return false
	}
	return v.SetValue.Equal(other.SetValue)
}

func (v NestedSetValue[T]) Type(ctx context.Context) attr.Type {
	return NewNestedSetType[T](ctx)
}

func (v NestedSetValue[T]) SlicePtrAsAny(ctx context.Context) (any, diag.Diagnostics) {
	valuePtr := new([]T)

	if v.IsNull() || v.IsUnknown() {
		return valuePtr, nil
	}

	diags := v.ElementsAs(ctx, valuePtr, false)
	if diags.HasError() {
		return nil, diags
	}

	return valuePtr, diags
}

func (v NestedSetValue[T]) NewEmptySlicePtr() any {
	return new([]T)
}

func (v NestedSetValue[T]) Len() int {
	return len(v.Elements())
}
