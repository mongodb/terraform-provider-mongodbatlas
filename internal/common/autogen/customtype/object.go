package customtype

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

/*
	Custom Object type used in auto-generated code to enable the generic marshal/unmarshal operations to access nested attribute struct tags during conversion.
	Custom types docs: https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/custom

	Usage:
		- Schema definition:
			"sample_nested_object": schema.SingleNestedAttribute{
				...
				CustomType: customtype.NewObjectType[TFSampleNestedObjectModel](ctx),
				Attributes: map[string]schema.Attribute{
					"string_attribute": schema.StringAttribute{...},
				},
			}

		- TF Models:
			type TFModel struct {
				SampleNestedObject customtype.ObjectValue[TFSampleNestedObjectModel] `tfsdk:"sample_nested_object"`
				...
			}

			type TFSampleNestedObjectModel struct {
				StringAttribute types.String `tfsdk:"string_attribute"`
				...
			}
*/

var (
	_ basetypes.ObjectTypable  = ObjectType[struct{}]{}
	_ basetypes.ObjectValuable = ObjectValue[struct{}]{}
	_ ObjectValueInterface     = ObjectValue[struct{}]{}
)

type ObjectType[T any] struct {
	basetypes.ObjectType
}

func NewObjectType[T any](ctx context.Context) ObjectType[T] {
	result := ObjectType[T]{}

	attrTypes, diags := getAttributeTypes[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating ObjectType: %v", diags))
	}

	result.ObjectType = basetypes.ObjectType{AttrTypes: attrTypes}
	return result
}

func (t ObjectType[T]) Equal(o attr.Type) bool {
	other, ok := o.(ObjectType[T])
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (ObjectType[T]) String() string {
	var t T
	return fmt.Sprintf("ObjectType[%T]", t)
}

func (t ObjectType[T]) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	if in.IsNull() {
		return NewObjectValueNull[T](ctx), nil
	}

	if in.IsUnknown() {
		return NewObjectValueUnknown[T](ctx), nil
	}

	attrTypes, diags := getAttributeTypes[T](ctx)
	if diags.HasError() {
		return nil, diags
	}

	baseObjectValue, diags := basetypes.NewObjectValue(attrTypes, in.Attributes())
	if diags.HasError() {
		return nil, diags
	}

	return ObjectValue[T]{ObjectValue: baseObjectValue}, nil
}

func (t ObjectType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.ObjectType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	objectValue, ok := attrValue.(basetypes.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	objectValuable, diags := t.ValueFromObject(ctx, objectValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ObjectValue to ObjectValuable: %v", diags)
	}

	return objectValuable, nil
}

func (t ObjectType[T]) ValueType(_ context.Context) attr.Value {
	return ObjectValue[T]{}
}

type ObjectValue[T any] struct {
	basetypes.ObjectValue
}

type ObjectValueInterface interface {
	basetypes.ObjectValuable
	ValuePtrAsAny(ctx context.Context) (any, diag.Diagnostics)
	NewObjectValue(ctx context.Context, value any) ObjectValueInterface
	NewObjectValueNull(ctx context.Context) ObjectValueInterface
}

func (v ObjectValue[T]) NewObjectValue(ctx context.Context, value any) ObjectValueInterface {
	return NewObjectValue[T](ctx, value)
}

func NewObjectValue[T any](ctx context.Context, value any) ObjectValue[T] {
	attrTypes, diags := getAttributeTypes[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating ObjectValue: %v", diags))
	}

	newValue, diags := basetypes.NewObjectValueFrom(ctx, attrTypes, value)
	if diags.HasError() {
		return NewObjectValueUnknown[T](ctx)
	}

	return ObjectValue[T]{ObjectValue: newValue}
}

func (v ObjectValue[T]) NewObjectValueNull(ctx context.Context) ObjectValueInterface {
	return NewObjectValueNull[T](ctx)
}

func NewObjectValueNull[T any](ctx context.Context) ObjectValue[T] {
	attrTypes, diags := getAttributeTypes[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating null ObjectValue: %v", diags))
	}
	return ObjectValue[T]{ObjectValue: basetypes.NewObjectNull(attrTypes)}
}

func NewObjectValueUnknown[T any](ctx context.Context) ObjectValue[T] {
	attrTypes, diags := getAttributeTypes[T](ctx)
	if diags.HasError() {
		panic(fmt.Errorf("error creating unknown ObjectValue: %v", diags))
	}
	return ObjectValue[T]{ObjectValue: basetypes.NewObjectUnknown(attrTypes)}
}

func (v ObjectValue[T]) Equal(o attr.Value) bool {
	other, ok := o.(ObjectValue[T])
	if !ok {
		return false
	}
	return v.ObjectValue.Equal(other.ObjectValue)
}

func (v ObjectValue[T]) Type(ctx context.Context) attr.Type {
	return NewObjectType[T](ctx)
}

func (v ObjectValue[T]) ValuePtrAsAny(ctx context.Context) (any, diag.Diagnostics) {
	valuePtr := new(T)

	if v.IsNull() || v.IsUnknown() {
		return valuePtr, nil
	}

	diags := v.As(ctx, valuePtr, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	return valuePtr, diags
}

func getAttributeTypes[T any](ctx context.Context) (map[string]attr.Type, diag.Diagnostics) {
	var t T
	return valueToAttributeTypes(ctx, reflect.ValueOf(t))
}

func valueToAttributeTypes(ctx context.Context, value reflect.Value) (map[string]attr.Type, diag.Diagnostics) {
	valueType := value.Type()

	if valueType.Kind() != reflect.Struct {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
			"Error getting value attribute types",
			fmt.Sprintf(`%T has usupported type: %s`, value.Interface(), valueType),
		)}
	}

	attributeTypes := make(map[string]attr.Type)
	for i := range valueType.NumField() {
		typeField := valueType.Field(i)
		valueField := value.Field(i)

		tfName := typeField.Tag.Get(`tfsdk`)
		if tfName == "" {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
				"Error getting value attribute types",
				fmt.Sprintf(`%T has no tfsdk tag on field %s`, value.Interface(), typeField.Name),
			)}
		}

		attrValue, ok := valueField.Interface().(attr.Value)
		if !ok {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
				"Error getting value attribute types",
				fmt.Sprintf(`%T has unsupported type in field %s: %T`, value.Interface(), typeField.Name, valueField.Interface()),
			)}
		}

		attributeTypes[tfName] = attrValue.Type(ctx)
	}

	return attributeTypes, nil
}
