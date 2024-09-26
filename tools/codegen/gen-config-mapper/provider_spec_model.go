package genconfigmapper

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

type CodeSpecification struct {
	Resources Resource
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
	Name               string
	Description        *string
	DeprecationMessage *string

	Sensitive  *bool
	IsComputed *bool
	IsOptional *bool
	IsRequired *bool

	Bool         *BoolAttribute
	Float64      *Float64Attribute
	Int64        *Int64Attribute
	List         *ListAttribute
	ListNested   *ListNestedAttribute
	Map          *MapAttribute
	MapNested    *MapNestedAttribute
	Number       *NumberAttribute
	Object       *ObjectAttribute
	Set          *SetAttribute
	SetNested    *SetNestedAttribute
	SingleNested *SingleNestedAttribute
	String       *StringAttribute
}

type BoolAttribute struct {
	Default *bool
}
type Float64Attribute struct {
	Default *float64
}
type Int64Attribute struct {
	Default *int64
}
type MapAttribute struct {
	Default     *CustomDefault
	ElementType ElemType
}
type MapNestedAttribute struct {
	Default      *CustomDefault
	NestedObject NestedAttributeObject
}
type NumberAttribute struct {
	Default *CustomDefault
}
type ObjectAttribute struct {
	Default *CustomDefault
}
type SetAttribute struct {
	Default     *CustomDefault
	ElementType ElemType
}
type SetNestedAttribute struct {
	Default      *CustomDefault
	NestedObject NestedAttributeObject
}
type SingleNestedAttribute struct {
	Default    *CustomDefault
	Attributes Attributes
}
type StringAttribute struct {
	Default *string
}
type ListAttribute struct {
	Default     *CustomDefault
	ElementType ElemType
}
type ListNestedAttribute struct {
	Default      *CustomDefault
	NestedObject NestedAttributeObject
}
type NestedAttributeObject struct {
	Attributes Attributes
}

type CustomDefault struct {
	Definition string
	Imports    []string
}
