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
	Create        APIOperation
	Read          APIOperation
	Update        APIOperation
	Delete        APIOperation
	VersionHeader string
}

type APIOperation struct {
	Wait       *Wait
	HTTPMethod string
	Path       string
}

type Wait struct {
	StateProperty     string
	PendingStates     []string
	TargetStates      []string
	TimeoutSeconds    int
	MinTimeoutSeconds int
	DelaySeconds      int
}

type Schema struct {
	Description        *string
	DeprecationMessage *string

	Attributes Attributes
}

type Attributes []Attribute

// Add this field to the Attribute struct
// Usage AttributeUsage
type Attribute struct {
	SetNested                *SetNestedAttribute
	String                   *StringAttribute
	Float64                  *Float64Attribute
	List                     *ListAttribute
	Bool                     *BoolAttribute
	ListNested               *ListNestedAttribute
	Map                      *MapAttribute
	MapNested                *MapNestedAttribute
	Int64                    *Int64Attribute
	Number                   *NumberAttribute
	Set                      *SetAttribute
	SingleNested             *SingleNestedAttribute
	Timeouts                 *TimeoutsAttribute
	Description              *string
	DeprecationMessage       *string
	Name                     stringcase.SnakeCaseString
	ComputedOptionalRequired ComputedOptionalRequired
	ReqBodyUsage             AttributeReqBodyUsage
	Sensitive                bool
}

type AttributeReqBodyUsage int

const (
	AllRequestBodies = iota // by default attribute is sent in request bodies
	OmitInUpdateBody
	OmitAlways // this covers computed-only attributes and attributes which are only used for path/query params
)

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
