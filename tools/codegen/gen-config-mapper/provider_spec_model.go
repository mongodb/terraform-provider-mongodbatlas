package genconfigmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type CodeSpecification struct {
	DataSource       DataSource       `json:"datasource,omitempty"`
	DataSourcePlural DataSourcePlural `json:"datasourceplural,omitempty"`
	Resources        Resource         `json:"resource,omitempty"`
}

type DataSourcePlural struct {
	Schema *Schema `json:"schema,omitempty"`
	Name   string  `json:"name"`
}

type DataSource struct {
	Schema *Schema `json:"schema,omitempty"`
	Name   string  `json:"name"`
}

type Resource struct {
	Schema *Schema `json:"schema,omitempty"`
	Name   string  `json:"name"`
}

type Schema struct {
	Description         *string `json:"description,omitempty"`
	MarkdownDescription *string `json:"markdown_description,omitempty"`
	DeprecationMessage  *string `json:"deprecation_message,omitempty"`

	// Blocks              Blocks  `json:"blocks,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`
}

type Attributes []Attribute

type Attribute struct {
	Bool         *resource.BoolAttribute         `json:"bool,omitempty"`
	Float64      *resource.Float64Attribute      `json:"float64,omitempty"`
	Int64        *resource.Int64Attribute        `json:"int64,omitempty"`
	List         *resource.ListAttribute         `json:"list,omitempty"`
	ListNested   *resource.ListNestedAttribute   `json:"list_nested,omitempty"`
	Map          *resource.MapAttribute          `json:"map,omitempty"`
	MapNested    *resource.MapNestedAttribute    `json:"map_nested,omitempty"`
	Number       *resource.NumberAttribute       `json:"number,omitempty"`
	Object       *resource.ObjectAttribute       `json:"object,omitempty"`
	Set          *resource.SetAttribute          `json:"set,omitempty"`
	SetNested    *resource.SetNestedAttribute    `json:"set_nested,omitempty"`
	SingleNested *resource.SingleNestedAttribute `json:"single_nested,omitempty"`
	String       *resource.StringAttribute       `json:"string,omitempty"`

	Name string `json:"name"`
}
