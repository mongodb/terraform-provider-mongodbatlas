package genconfigmapper

import (
	"log"

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
	Bool         *BoolAttribute
	Float64      *Float64Attribute
	Int64        *Int64Attribute
	List         *ListAttribute
	ListNested   *ListNestedAttribute
	Map          *resource.MapAttribute
	MapNested    *resource.MapNestedAttribute
	Number       *resource.NumberAttribute
	Object       *resource.ObjectAttribute
	Set          *resource.SetAttribute
	SetNested    *resource.SetNestedAttribute
	SingleNested *resource.SingleNestedAttribute
	String       *resource.StringAttribute

	Name               string
	Description        *string `json:"description,omitempty"`
	DeprecationMessage *string `json:"deprecation_message,omitempty"`
	Sensitive          *bool   `json:"sensitive,omitempty"`

	IsComputed *bool `json:"is_computed,omitempty"`
	IsOptional *bool `json:"is_optional,omitempty"`
	IsRequired *bool `json:"is_required,omitempty"`

	// TODO:
	// PlanModifiers PlanModifiers `json:"plan_modifiers,omitempty"`
	// Validators BoolValidators `json:"validators,omitempty"`

}

func temp() {
	listA := ListAttribute{
		Default: &ListDefault{
			Custom: &CustomDefault{
				Definition: "",
				Imports:    []string{""},
			},
		},
		ElementType: Bool,
	}
	log.Print(listA.ElementType)
}

type BoolAttribute struct {
	Default *bool `json:"default,omitempty"`
}
type Float64Attribute struct {
	Default *float64 `json:"default,omitempty"`
}
type Int64Attribute struct {
	Default *int64 `json:"default,omitempty"`
}
type ListAttribute struct {
	Default     *ListDefault `json:"default,omitempty"`
	ElementType ElemType     `json:"element_type"`
}
type ListNestedAttribute struct {
	Default      *ListDefault          `json:"default,omitempty"`
	NestedObject NestedAttributeObject `json:"nested_object"`
}
type ListDefault struct {
	Custom *CustomDefault `json:"custom,omitempty"`
}

type NestedAttributeObject struct {
	Attributes Attributes `json:"attributes,omitempty"`
	// TODO:
	// PlanModifiers schema.ObjectPlanModifiers `json:"plan_modifiers,omitempty"`
	// Validators schema.ObjectValidators `json:"validators,omitempty"`
}

type CustomDefault struct {
	Definition string
	Imports    []string
}

type ElemType int

const (
	Bool ElemType = iota
	Float64
	Int64
	List
	Map
	Number
	Object
	Set
	String
)
