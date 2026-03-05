package schema

import "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"

type Int64AttrGenerator struct {
	intModel codespec.Int64Attribute
	attr     codespec.Attribute
}

func (i *Int64AttrGenerator) AttributeCode() (CodeStatement, error) {
	return commonAttrStructure(&i.attr, &TypeSpecificStatements{
		AttrType:         "schema.Int64Attribute",
		PlanModifierType: "planmodifier.Int64",
		ValidatorType:    "validator.Int64",
	})
}

type Float64AttrGenerator struct {
	floatModel codespec.Float64Attribute
	attr       codespec.Attribute
}

func (f *Float64AttrGenerator) AttributeCode() (CodeStatement, error) {
	return commonAttrStructure(&f.attr, &TypeSpecificStatements{
		AttrType:         "schema.Float64Attribute",
		PlanModifierType: "planmodifier.Float64",
		ValidatorType:    "validator.Float64",
	})
}

type StringAttrGenerator struct {
	stringModel   codespec.StringAttribute
	discriminator *codespec.Discriminator
	attr          codespec.Attribute
}

func (s *StringAttrGenerator) AttributeCode() (CodeStatement, error) {
	var validators []CodeStatement
	if s.discriminator != nil && !s.discriminator.SkipValidation {
		validators = append(validators, discriminatorValidatorProperty(s.discriminator))
	}
	return commonAttrStructure(&s.attr, &TypeSpecificStatements{
		AttrType:         "schema.StringAttribute",
		PlanModifierType: "planmodifier.String",
		ValidatorType:    "validator.String",
		Validators:       validators,
	})
}

type BoolAttrGenerator struct {
	boolModel codespec.BoolAttribute
	attr      codespec.Attribute
}

func (s *BoolAttrGenerator) AttributeCode() (CodeStatement, error) {
	return commonAttrStructure(&s.attr, &TypeSpecificStatements{
		AttrType:         "schema.BoolAttribute",
		PlanModifierType: "planmodifier.Bool",
		ValidatorType:    "validator.Bool",
	})
}

type NumberAttrGenerator struct {
	numberModel codespec.NumberAttribute
	attr        codespec.Attribute
}

func (s *NumberAttrGenerator) AttributeCode() (CodeStatement, error) {
	return commonAttrStructure(&s.attr, &TypeSpecificStatements{
		AttrType:         "schema.NumberAttribute",
		PlanModifierType: "planmodifier.Number",
		ValidatorType:    "validator.Number",
	})
}

type ListAttrGenerator struct {
	listModel codespec.ListAttribute
	attr      codespec.Attribute
}

func (l *ListAttrGenerator) AttributeCode() (CodeStatement, error) {
	return commonAttrStructure(&l.attr, &TypeSpecificStatements{
		AttrType:         "schema.ListAttribute",
		PlanModifierType: "planmodifier.List",
		ValidatorType:    "validator.List",
		Properties:       []CodeStatement{ElementTypeProperty(l.listModel.ElementType)},
	})
}

type MapAttrGenerator struct {
	mapModel codespec.MapAttribute
	attr     codespec.Attribute
}

func (m *MapAttrGenerator) AttributeCode() (CodeStatement, error) {
	return commonAttrStructure(&m.attr, &TypeSpecificStatements{
		AttrType:         "schema.MapAttribute",
		PlanModifierType: "planmodifier.Map",
		ValidatorType:    "validator.Map",
		Properties:       []CodeStatement{ElementTypeProperty(m.mapModel.ElementType)},
	})
}

type SetAttrGenerator struct {
	setModel codespec.SetAttribute
	attr     codespec.Attribute
}

func (s *SetAttrGenerator) AttributeCode() (CodeStatement, error) {
	return commonAttrStructure(&s.attr, &TypeSpecificStatements{
		AttrType:         "schema.SetAttribute",
		PlanModifierType: "planmodifier.Set",
		ValidatorType:    "validator.Set",
		Properties:       []CodeStatement{ElementTypeProperty(s.setModel.ElementType)},
	})
}
