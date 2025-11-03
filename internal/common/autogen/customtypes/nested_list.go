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
	Custom Nested List type used in auto-generated code to enable the generic marshal/unmarshal operations to access nested attribute struct tags during conversion.
	Custom types docs: https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/custom

	Usage:
		- Schema definition:
			"sample_nested_object_list": schema.ListNestedAttribute{
				...
				CustomType: customtypes.NewNestedListType[TFSampleNestedObjectModel](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string_attribute": schema.StringAttribute{...},
					},
				},
			}

		- TF Models:
			type TFModel struct {
				SampleNestedObjectList customtypes.NestedListValue[TFSampleNestedObjectModel] `tfsdk:"sample_nested_object_list"`
				...
			}

			type TFSampleNestedObjectModel struct {
				StringAttribute types.String `tfsdk:"string_attribute"`
				...
			}
*/

var (
	_ basetypes.ListTypable    = NestedListType[struct{}]{}
	_ basetypes.ListValuable   = NestedListValue[struct{}]{}
	_ NestedListValueInterface = NestedListValue[struct{}]{}
)

type NestedListType[T any] struct {
	basetypes.ListType
}

func NewNestedListType[T any](ctx context.Context) NestedListType[T] {
	elemType, diags := getNestedType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating NestedListType: %v", diags))
	}

	result := NestedListType[T]{
		ListType: basetypes.ListType{ElemType: elemType},
	}
	return result
}

func (t NestedListType[T]) Equal(o attr.Type) bool {
	other, ok := o.(NestedListType[T])
	if !ok {
		return false
	}
	return t.ListType.Equal(other.ListType)
}

func (NestedListType[T]) String() string {
	var t T
	return fmt.Sprintf("NestedListType[%T]", t)
}

func (t NestedListType[T]) ValueFromList(ctx context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	if in.IsNull() {
		return NewNestedListValueNull[T](ctx), nil
	}

	if in.IsUnknown() {
		return NewNestedListValueUnknown[T](ctx), nil
	}

	elemType, diags := getNestedType[T](ctx)
	if diags.HasError() {
		return nil, diags
	}

	baseListValue, diags := basetypes.NewListValue(elemType, in.Elements())
	if diags.HasError() {
		return nil, diags
	}

	return NestedListValue[T]{ListValue: baseListValue}, nil
}

func (t NestedListType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t NestedListType[T]) ValueType(_ context.Context) attr.Value {
	return NestedListValue[T]{}
}

type NestedListValue[T any] struct {
	basetypes.ListValue
}

type NestedListValueInterface interface {
	basetypes.ListValuable
	NewNestedListValue(ctx context.Context, value any) NestedListValueInterface
	NewNestedListValueNull(ctx context.Context) NestedListValueInterface
	SlicePtrAsAny(ctx context.Context) (any, diag.Diagnostics)
	NewEmptySlicePtr() any
	Len() int
}

func (v NestedListValue[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.ElementType(ctx) == nil {
		// NestedListValue created as a zero value (not explicitly initialized), initialize now so conversion does not panic.
		v.ListValue = NewNestedListValueNull[T](ctx).ListValue
	}
	return v.ListValue.ToTerraformValue(ctx)
}

func (v NestedListValue[T]) NewNestedListValue(ctx context.Context, value any) NestedListValueInterface {
	return NewNestedListValue[T](ctx, value)
}

func NewNestedListValue[T any](ctx context.Context, value any) NestedListValue[T] {
	elemType, diags := getNestedType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating NestedListValue: %v", diags))
	}

	newValue, diags := basetypes.NewListValueFrom(ctx, elemType, value)
	if diags.HasError() {
		return NewNestedListValueUnknown[T](ctx)
	}

	return NestedListValue[T]{ListValue: newValue}
}

func (v NestedListValue[T]) NewNestedListValueNull(ctx context.Context) NestedListValueInterface {
	return NewNestedListValueNull[T](ctx)
}

func NewNestedListValueNull[T any](ctx context.Context) NestedListValue[T] {
	elemType, diags := getNestedType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating null NestedListValue: %v", diags))
	}
	return NestedListValue[T]{ListValue: basetypes.NewListNull(elemType)}
}

func NewNestedListValueUnknown[T any](ctx context.Context) NestedListValue[T] {
	elemType, diags := getNestedType[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating unknown NestedListValue: %v", diags))
	}
	return NestedListValue[T]{ListValue: basetypes.NewListUnknown(elemType)}
}

func (v NestedListValue[T]) Equal(o attr.Value) bool {
	other, ok := o.(NestedListValue[T])
	if !ok {
		return false
	}
	return v.ListValue.Equal(other.ListValue)
}

func (v NestedListValue[T]) Type(ctx context.Context) attr.Type {
	return NewNestedListType[T](ctx)
}

func (v NestedListValue[T]) SlicePtrAsAny(ctx context.Context) (any, diag.Diagnostics) {
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

func (v NestedListValue[T]) NewEmptySlicePtr() any {
	return new([]T)
}

func (v NestedListValue[T]) Len() int {
	return len(v.Elements())
}
