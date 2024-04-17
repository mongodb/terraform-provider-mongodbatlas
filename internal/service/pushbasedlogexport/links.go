package pushbasedlogexport

// import (
// 	"context"
// 	"fmt"
// 	"strings"
//

// 	"github.com/hashicorp/terraform-plugin-framework/attr"
// 	"github.com/hashicorp/terraform-plugin-framework/diag"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// 	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
// 	"github.com/hashicorp/terraform-plugin-go/tftypes"
// )

// // TODO: remove this file
// var _ basetypes.ObjectTypable = LinksType{}

// type LinksType struct {
// 	basetypes.ObjectType
// }

// func (t LinksType) Equal(o attr.Type) bool {
// 	other, ok := o.(LinksType)

// 	if !ok {
// 		return false
// 	}

// 	return t.ObjectType.Equal(other.ObjectType)
// }

// func (t LinksType) String() string {
// 	return "LinksType"
// }

// func (t LinksType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
// 	var diags diag.Diagnostics

// 	attributes := in.Attributes()

// 	hrefAttribute, ok := attributes["href"]

// 	if !ok {
// 		diags.AddError(
// 			"Attribute Missing",
// 			`href is missing from object`)

// 		return nil, diags
// 	}

// 	hrefVal, ok := hrefAttribute.(basetypes.StringValue)

// 	if !ok {
// 		diags.AddError(
// 			"Attribute Wrong Type",
// 			fmt.Sprintf(`href expected to be basetypes.StringValue, was: %T`, hrefAttribute))
// 	}

// 	relAttribute, ok := attributes["rel"]

// 	if !ok {
// 		diags.AddError(
// 			"Attribute Missing",
// 			`rel is missing from object`)

// 		return nil, diags
// 	}

// 	relVal, ok := relAttribute.(basetypes.StringValue)

// 	if !ok {
// 		diags.AddError(
// 			"Attribute Wrong Type",
// 			fmt.Sprintf(`rel expected to be basetypes.StringValue, was: %T`, relAttribute))
// 	}

// 	if diags.HasError() {
// 		return nil, diags
// 	}

// 	return LinksValue{
// 		Href:  hrefVal,
// 		Rel:   relVal,
// 		state: attr.ValueStateKnown,
// 	}, diags
// }

// func NewLinksValueNull() LinksValue {
// 	return LinksValue{
// 		state: attr.ValueStateNull,
// 	}
// }

// func NewLinksValueUnknown() LinksValue {
// 	return LinksValue{
// 		state: attr.ValueStateUnknown,
// 	}
// }

// func NewLinksValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (LinksValue, diag.Diagnostics) {
// 	var diags diag.Diagnostics

// 	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
// 	ctx := context.Background()

// 	for name, attributeType := range attributeTypes {
// 		attribute, ok := attributes[name]

// 		if !ok {
// 			diags.AddError(
// 				"Missing LinksValue Attribute Value",
// 				"While creating a LinksValue value, a missing attribute value was detected. "+
// 					"A LinksValue must contain values for all attributes, even if null or unknown. "+
// 					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
// 					fmt.Sprintf("LinksValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
// 			)

// 			continue
// 		}

// 		if !attributeType.Equal(attribute.Type(ctx)) {
// 			diags.AddError(
// 				"Invalid LinksValue Attribute Type",
// 				"While creating a LinksValue value, an invalid attribute value was detected. "+
// 					"A LinksValue must use a matching attribute type for the value. "+
// 					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
// 					fmt.Sprintf("LinksValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
// 					fmt.Sprintf("LinksValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
// 			)
// 		}
// 	}

// 	for name := range attributes {
// 		_, ok := attributeTypes[name]

// 		if !ok {
// 			diags.AddError(
// 				"Extra LinksValue Attribute Value",
// 				"While creating a LinksValue value, an extra attribute value was detected. "+
// 					"A LinksValue must not contain values beyond the expected attribute types. "+
// 					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
// 					fmt.Sprintf("Extra LinksValue Attribute Name: %s", name),
// 			)
// 		}
// 	}

// 	if diags.HasError() {
// 		return NewLinksValueUnknown(), diags
// 	}

// 	hrefAttribute, ok := attributes["href"]

// 	if !ok {
// 		diags.AddError(
// 			"Attribute Missing",
// 			`href is missing from object`)

// 		return NewLinksValueUnknown(), diags
// 	}

// 	hrefVal, ok := hrefAttribute.(basetypes.StringValue)

// 	if !ok {
// 		diags.AddError(
// 			"Attribute Wrong Type",
// 			fmt.Sprintf(`href expected to be basetypes.StringValue, was: %T`, hrefAttribute))
// 	}

// 	relAttribute, ok := attributes["rel"]

// 	if !ok {
// 		diags.AddError(
// 			"Attribute Missing",
// 			`rel is missing from object`)

// 		return NewLinksValueUnknown(), diags
// 	}

// 	relVal, ok := relAttribute.(basetypes.StringValue)

// 	if !ok {
// 		diags.AddError(
// 			"Attribute Wrong Type",
// 			fmt.Sprintf(`rel expected to be basetypes.StringValue, was: %T`, relAttribute))
// 	}

// 	if diags.HasError() {
// 		return NewLinksValueUnknown(), diags
// 	}

// 	return LinksValue{
// 		Href:  hrefVal,
// 		Rel:   relVal,
// 		state: attr.ValueStateKnown,
// 	}, diags
// }

// func NewLinksValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) LinksValue {
// 	object, diags := NewLinksValue(attributeTypes, attributes)

// 	if diags.HasError() {
// 		// This could potentially be added to the diag package.
// 		diagsStrings := make([]string, 0, len(diags))

// 		for _, diagnostic := range diags {
// 			diagsStrings = append(diagsStrings, fmt.Sprintf(
// 				"%s | %s | %s",
// 				diagnostic.Severity(),
// 				diagnostic.Summary(),
// 				diagnostic.Detail()))
// 		}

// 		panic("NewLinksValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
// 	}

// 	return object
// }

// func (t LinksType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
// 	if in.Type() == nil {
// 		return NewLinksValueNull(), nil
// 	}

// 	if !in.Type().Equal(t.TerraformType(ctx)) {
// 		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
// 	}

// 	if !in.IsKnown() {
// 		return NewLinksValueUnknown(), nil
// 	}

// 	if in.IsNull() {
// 		return NewLinksValueNull(), nil
// 	}

// 	attributes := map[string]attr.Value{}

// 	val := map[string]tftypes.Value{}

// 	err := in.As(&val)

// 	if err != nil {
// 		return nil, err
// 	}

// 	for k, v := range val {
// 		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

// 		if err != nil {
// 			return nil, err
// 		}

// 		attributes[k] = a
// 	}

// 	return NewLinksValueMust(LinksValue{}.AttributeTypes(ctx), attributes), nil
// }

// func (t LinksType) ValueType(ctx context.Context) attr.Value {
// 	return LinksValue{}
// }

// var _ basetypes.ObjectValuable = LinksValue{}

// type LinksValue struct {
// 	Href  basetypes.StringValue `tfsdk:"href"`
// 	Rel   basetypes.StringValue `tfsdk:"rel"`
// 	state attr.ValueState
// }

// func (v LinksValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
// 	attrTypes := make(map[string]tftypes.Type, 2)

// 	var val tftypes.Value
// 	var err error

// 	attrTypes["href"] = basetypes.StringType{}.TerraformType(ctx)
// 	attrTypes["rel"] = basetypes.StringType{}.TerraformType(ctx)

// 	objectType := tftypes.Object{AttributeTypes: attrTypes}

// 	switch v.state {
// 	case attr.ValueStateKnown:
// 		vals := make(map[string]tftypes.Value, 2)

// 		val, err = v.Href.ToTerraformValue(ctx)

// 		if err != nil {
// 			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
// 		}

// 		vals["href"] = val

// 		val, err = v.Rel.ToTerraformValue(ctx)

// 		if err != nil {
// 			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
// 		}

// 		vals["rel"] = val

// 		if err := tftypes.ValidateValue(objectType, vals); err != nil {
// 			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
// 		}

// 		return tftypes.NewValue(objectType, vals), nil
// 	case attr.ValueStateNull:
// 		return tftypes.NewValue(objectType, nil), nil
// 	case attr.ValueStateUnknown:
// 		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
// 	default:
// 		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
// 	}
// }

// func (v LinksValue) IsNull() bool {
// 	return v.state == attr.ValueStateNull
// }

// func (v LinksValue) IsUnknown() bool {
// 	return v.state == attr.ValueStateUnknown
// }

// func (v LinksValue) String() string {
// 	return "LinksValue"
// }

// func (v LinksValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
// 	var diags diag.Diagnostics

// 	objVal, diags := types.ObjectValue(
// 		map[string]attr.Type{
// 			"href": basetypes.StringType{},
// 			"rel":  basetypes.StringType{},
// 		},
// 		map[string]attr.Value{
// 			"href": v.Href,
// 			"rel":  v.Rel,
// 		})

// 	return objVal, diags
// }

// func (v LinksValue) Equal(o attr.Value) bool {
// 	other, ok := o.(LinksValue)

// 	if !ok {
// 		return false
// 	}

// 	if v.state != other.state {
// 		return false
// 	}

// 	if v.state != attr.ValueStateKnown {
// 		return true
// 	}

// 	if !v.Href.Equal(other.Href) {
// 		return false
// 	}

// 	if !v.Rel.Equal(other.Rel) {
// 		return false
// 	}

// 	return true
// }

// func (v LinksValue) Type(ctx context.Context) attr.Type {
// 	return LinksType{
// 		basetypes.ObjectType{
// 			AttrTypes: v.AttributeTypes(ctx),
// 		},
// 	}
// }

// func (v LinksValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
// 	return map[string]attr.Type{
// 		"href": basetypes.StringType{},
// 		"rel":  basetypes.StringType{},
// 	}
// }
