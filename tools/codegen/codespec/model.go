package codespec

import "github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/stringcase"

type ElemType int

const (
	Bool ElemType = iota
	Float64
	Int64
	Number
	String
	Unknown
)

type Model struct {
	Resources []Resource
}

type Resource struct {
	Schema     *Schema
	Name       stringcase.SnakeCaseString
	Operations APIOperations
}

type APIOperations struct {
	CreatePath    string
	ReadPath      string
	UpdatePath    string
	DeletePath    string
	VersionHeader string
}
type Schema struct {
	Description        *string
	DeprecationMessage *string

	Attributes Attributes
}

type Attributes []Attribute

type Attribute struct {
	List      *ListAttribute
	SetNested *SetNestedAttribute

	Float64 *Float64Attribute
	String  *StringAttribute

	Bool         *BoolAttribute
	ListNested   *ListNestedAttribute
	Map          *MapAttribute
	MapNested    *MapNestedAttribute
	Number       *NumberAttribute
	Set          *SetAttribute
	Int64        *Int64Attribute
	SingleNested *SingleNestedAttribute
	Timeouts     *TimeoutsAttribute

	Description              *string
	Name                     stringcase.SnakeCaseString
	DeprecationMessage       *string
	Sensitive                *bool
	ComputedOptionalRequired ComputedOptionalRequired
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
type SetAttribute struct {
	Default     *CustomDefault
	ElementType ElemType
}
type SetNestedAttribute struct {
	Default      *CustomDefault
	NestedObject NestedAttributeObject
}
type SingleNestedAttribute struct {
	Default      *CustomDefault
	NestedObject NestedAttributeObject
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

type TimeoutsAttribute struct {
	ConfigurableTimeouts []Operation
}

type Operation int

const (
	Create Operation = iota
	Update
	Read
	Delete
)

type CustomDefault struct {
	Definition string
	Imports    []string
}
