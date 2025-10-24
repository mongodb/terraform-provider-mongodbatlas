package codespec

import (
	"fmt"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/stringcase"
)

type ElemType int

const (
	Bool ElemType = iota
	Float64
	Int64
	Number
	String
	CustomTypeJSON
	Unknown
)

type Model struct {
	Resources []Resource
}

type Resource struct {
	Schema     *Schema                    `yaml:"schema,omitempty"`
	Operations APIOperations              `yaml:"operations"`
	Name       stringcase.SnakeCaseString `yaml:"name"`
}

type APIOperations struct {
	Delete        *APIOperation `yaml:"delete,omitempty"`
	Create        APIOperation  `yaml:"create"`
	Read          APIOperation  `yaml:"read"`
	Update        APIOperation  `yaml:"update"`
	VersionHeader string        `yaml:"version_header"`
}

type APIOperation struct {
	Wait              *Wait  `yaml:"wait,omitempty"`
	HTTPMethod        string `yaml:"http_method"`
	Path              string `yaml:"path"`
	StaticRequestBody string `yaml:"static_request_body,omitempty"`
}

type Wait struct {
	StateProperty     string   `yaml:"state_property"`
	PendingStates     []string `yaml:"pending_states"`
	TargetStates      []string `yaml:"target_states"`
	TimeoutSeconds    int      `yaml:"timeout_seconds"`
	MinTimeoutSeconds int      `yaml:"min_timeout_seconds"`
	DelaySeconds      int      `yaml:"delay_seconds"`
}

type Schema struct {
	Description        *string `yaml:"description,omitempty"`
	DeprecationMessage *string `yaml:"deprecation_message,omitempty"`

	Attributes Attributes `yaml:"attributes"`
}

type Attributes []Attribute

// Add this field to the Attribute struct
// Usage AttributeUsage
type Attribute struct {
	Set                      *SetAttribute              `yaml:"set,omitempty"`
	String                   *StringAttribute           `yaml:"string,omitempty"`
	Float64                  *Float64Attribute          `yaml:"float64,omitempty"`
	List                     *ListAttribute             `yaml:"list,omitempty"`
	Bool                     *BoolAttribute             `yaml:"bool,omitempty"`
	ListNested               *ListNestedAttribute       `yaml:"list_nested,omitempty"`
	Map                      *MapAttribute              `yaml:"map,omitempty"`
	MapNested                *MapNestedAttribute        `yaml:"map_nested,omitempty"`
	Number                   *NumberAttribute           `yaml:"number,omitempty"`
	Int64                    *Int64Attribute            `yaml:"int64,omitempty"`
	Timeouts                 *TimeoutsAttribute         `yaml:"timeouts,omitempty"`
	SingleNested             *SingleNestedAttribute     `yaml:"single_nested,omitempty"`
	SetNested                *SetNestedAttribute        `yaml:"set_nested,omitempty"`
	Description              *string                    `yaml:"description,omitempty"`
	DeprecationMessage       *string                    `yaml:"deprecation_message,omitempty"`
	CustomType               *CustomType                `yaml:"custom_type,omitempty"`
	ComputedOptionalRequired ComputedOptionalRequired   `yaml:"computed_optional_required"`
	TFSchemaName             stringcase.SnakeCaseString `yaml:"tf_schema_name"`
	TFModelName              string                     `yaml:"tf_model_name"`
	ReqBodyUsage             AttributeReqBodyUsage      `yaml:"req_body_usage"`
	Sensitive                bool                       `yaml:"sensitive"`
	CreateOnly               bool                       `yaml:"create_only"` // leveraged for defining plan modifier which avoids updates on this attribute
}

type ComputedOptionalRequired string

const (
	Computed         ComputedOptionalRequired = "computed"
	ComputedOptional ComputedOptionalRequired = "computed_optional"
	Optional         ComputedOptionalRequired = "optional"
	Required         ComputedOptionalRequired = "required"
)

type AttributeReqBodyUsage string

const (
	AllRequestBodies    AttributeReqBodyUsage = "all_request_bodies" // by default attribute is sent in request bodies
	OmitInUpdateBody    AttributeReqBodyUsage = "omit_in_update_body"
	IncludeNullOnUpdate AttributeReqBodyUsage = "include_null_on_update" // attributes that always must be sent in update request body even if null
	OmitAlways          AttributeReqBodyUsage = "omit_always"            // this covers computed-only attributes and attributes which are only used for path/query params
)

type BoolAttribute struct {
	Default *bool `yaml:"default,omitempty"`
}
type Float64Attribute struct {
	Default *float64 `yaml:"default,omitempty"`
}
type Int64Attribute struct {
	Default *int64 `yaml:"default,omitempty"`
}
type MapAttribute struct {
	Default     *CustomDefault `yaml:"default,omitempty"`
	ElementType ElemType       `yaml:"element_type"`
}
type MapNestedAttribute struct {
	Default      *CustomDefault        `yaml:"default,omitempty"`
	NestedObject NestedAttributeObject `yaml:"nested_object"`
}
type NumberAttribute struct {
	Default *CustomDefault `yaml:"default,omitempty"`
}
type SetAttribute struct {
	Default     *CustomDefault `yaml:"default,omitempty"`
	ElementType ElemType       `yaml:"element_type"`
}
type SetNestedAttribute struct {
	Default      *CustomDefault        `yaml:"default,omitempty"`
	NestedObject NestedAttributeObject `yaml:"nested_object"`
}
type SingleNestedAttribute struct {
	Default      *CustomDefault        `yaml:"default,omitempty"`
	NestedObject NestedAttributeObject `yaml:"nested_object"`
}
type StringAttribute struct {
	Default *string `yaml:"default,omitempty"`
}
type ListAttribute struct {
	Default     *CustomDefault `yaml:"default,omitempty"`
	ElementType ElemType       `yaml:"element_type"`
}
type ListNestedAttribute struct {
	Default      *CustomDefault        `yaml:"default,omitempty"`
	NestedObject NestedAttributeObject `yaml:"nested_object"`
}
type NestedAttributeObject struct {
	Attributes Attributes `yaml:"attributes"`
}

type TimeoutsAttribute struct {
	ConfigurableTimeouts []Operation `yaml:"configurable_timeouts"`
}

type Operation string

const (
	Create Operation = "create"
	Update Operation = "update"
	Read   Operation = "read"
	Delete Operation = "delete"
)

type CustomDefault struct {
	Definition string   `yaml:"definition"`
	Imports    []string `yaml:"imports"`
}

type CustomTypePackage string

const (
	JSONTypesPkg   CustomTypePackage = "jsontypes"
	CustomTypesPkg CustomTypePackage = "customtypes"
)

type CustomType struct {
	Package CustomTypePackage `yaml:"package"`
	Model   string            `yaml:"model"`
	Schema  string            `yaml:"schema"`
}

var CustomTypeJSONVar = CustomType{
	Package: JSONTypesPkg,
	Model:   "jsontypes.Normalized",
	Schema:  "jsontypes.NormalizedType{}",
}

func NewCustomObjectType(name string) *CustomType {
	return &CustomType{
		Package: CustomTypesPkg,
		Model:   fmt.Sprintf("customtypes.ObjectValue[TF%sModel]", name),
		Schema:  fmt.Sprintf("customtypes.NewObjectType[TF%sModel](ctx)", name),
	}
}

func NewCustomNestedListType(name string) *CustomType {
	return &CustomType{
		Package: CustomTypesPkg,
		Model:   fmt.Sprintf("customtypes.NestedListValue[TF%sModel]", name),
		Schema:  fmt.Sprintf("customtypes.NewNestedListType[TF%sModel](ctx)", name),
	}
}
