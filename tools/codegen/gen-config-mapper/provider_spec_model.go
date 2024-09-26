package genconfigmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type CodeSpecification struct {
	DataSource       DataSource
	DataSourcePlural DataSourcePlural
	Resources        Resource
}

type DataSourcePlural struct {
	Schema *Schema
	Name   string
}

type DataSource struct {
	Schema *Schema
	Name   string
}

type Resource struct {
	Schema *Schema
	Name   string
}

type Schema struct {
	Description         *string
	MarkdownDescription *string
	DeprecationMessage  *string

	Attributes Attributes
}

type Attributes []Attribute

type Attribute struct {
	Bool         *resource.BoolAttribute
	Float64      *resource.Float64Attribute
	Int64        *resource.Int64Attribute
	List         *resource.ListAttribute
	ListNested   *resource.ListNestedAttribute
	Map          *resource.MapAttribute
	MapNested    *resource.MapNestedAttribute
	Number       *resource.NumberAttribute
	Object       *resource.ObjectAttribute
	Set          *resource.SetAttribute
	SetNested    *resource.SetNestedAttribute
	SingleNested *resource.SingleNestedAttribute
	String       *resource.StringAttribute

	Name string
}
