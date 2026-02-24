package codespec

import (
	"fmt"
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

var ElementTypeToSchemaString = map[ElemType]string{
	Bool:           "types.BoolType",
	Float64:        "types.Float64Type",
	Int64:          "types.Int64Type",
	Number:         "types.NumberType",
	String:         "types.StringType",
	CustomTypeJSON: CustomTypeJSONVar.Schema,
}

var ElementTypeToModelString = map[ElemType]string{
	Bool:           "types.Bool",
	Float64:        "types.Float64",
	Int64:          "types.Int64",
	Number:         "types.Number",
	String:         "types.String",
	CustomTypeJSON: CustomTypeJSONVar.Model,
}

type Model struct {
	Resources []Resource
}

type Resource struct {
	Schema       *Schema       `yaml:"schema,omitempty"`
	Operations   APIOperations `yaml:"operations"`
	MoveState    *MoveState    `yaml:"move_state,omitempty"`
	DataSources  *DataSources  `yaml:"data_sources,omitempty"`
	Name         string        `yaml:"name"`
	PackageName  string        `yaml:"packageName"`
	IDAttributes []string      `yaml:"id_attributes,omitempty"`
}

// DataSources holds the data source configuration within a resource
type DataSources struct {
	Schema     *DataSourceSchema `yaml:"schema,omitempty"`
	Operations APIOperations     `yaml:"operations"` // only Read and List operations
}

// DataSourceSchema holds schema information specific to data sources
type DataSourceSchema struct {
	SingularDSDescription *string     `yaml:"singular_ds_description,omitempty"`
	SingularDSAttributes  *Attributes `yaml:"singular_ds_attributes,omitempty"`
	PluralDSDescription   *string     `yaml:"plural_ds_description,omitempty"`
	PluralDSAttributes    *Attributes `yaml:"plural_ds_attributes,omitempty"`
	DeprecationMessage    *string     `yaml:"deprecation_message,omitempty"`
}

type APIOperations struct {
	Delete        *APIOperation `yaml:"delete,omitempty"`
	Create        *APIOperation `yaml:"create,omitempty"` // optional to support datasource-only API resources
	Read          *APIOperation `yaml:"read,omitempty"`   // optional to support datasource-only API resources
	List          *APIOperation `yaml:"list,omitempty"`   // for plural data sources
	Update        *APIOperation `yaml:"update,omitempty"`
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

type MoveState struct {
	SourceResources []string `yaml:"source_resources"`
}

type Schema struct {
	Description        *string        `yaml:"description,omitempty"`
	DeprecationMessage *string        `yaml:"deprecation_message,omitempty"`
	Discriminator      *Discriminator `yaml:"discriminator,omitempty"`

	Attributes Attributes `yaml:"attributes"`
}

// DiscriminatorAttrName pairs the original API property name with the Terraform schema name.
// TFSchemaName will have aliasing configurations applied.
type DiscriminatorAttrName struct {
	APIName      string `yaml:"api_name"`
	TFSchemaName string `yaml:"tf_schema_name"`
}

type Discriminator struct {
	Mapping        map[string]DiscriminatorType `yaml:"mapping"`
	PropertyName   DiscriminatorAttrName        `yaml:"property_name"`
	SkipValidation bool                         `yaml:"skip_validation,omitempty"`
}

type DiscriminatorType struct {
	// Allowed enumerates every sub-type specific attributes valid for this discriminator value.
	Allowed []DiscriminatorAttrName `yaml:"allowed"`
	// Required is a subset of Allowed that the user must explicitly set in their configuration.
	Required []DiscriminatorAttrName `yaml:"required,omitempty"`
}

type Attributes []Attribute

type Attribute struct {
	Set                         *SetAttribute            `yaml:"set,omitempty"`
	String                      *StringAttribute         `yaml:"string,omitempty"`
	Float64                     *Float64Attribute        `yaml:"float64,omitempty"`
	List                        *ListAttribute           `yaml:"list,omitempty"`
	Bool                        *BoolAttribute           `yaml:"bool,omitempty"`
	ListNested                  *ListNestedAttribute     `yaml:"list_nested,omitempty"`
	Map                         *MapAttribute            `yaml:"map,omitempty"`
	MapNested                   *MapNestedAttribute      `yaml:"map_nested,omitempty"`
	Number                      *NumberAttribute         `yaml:"number,omitempty"`
	Int64                       *Int64Attribute          `yaml:"int64,omitempty"`
	Timeouts                    *TimeoutsAttribute       `yaml:"timeouts,omitempty"`
	SingleNested                *SingleNestedAttribute   `yaml:"single_nested,omitempty"`
	SetNested                   *SetNestedAttribute      `yaml:"set_nested,omitempty"`
	Description                 *string                  `yaml:"description,omitempty"`
	DeprecationMessage          *string                  `yaml:"deprecation_message,omitempty"`
	CustomType                  *CustomType              `yaml:"custom_type,omitempty"`
	ComputedOptionalRequired    ComputedOptionalRequired `yaml:"computed_optional_required"`
	TFSchemaName                string                   `yaml:"tf_schema_name"`
	TFModelName                 string                   `yaml:"tf_model_name"`
	APIName                     string                   `yaml:"api_name,omitempty"` // original API property name (camelCase), used for apiname tag when different from Uncapitalize(TFModelName)
	ReqBodyUsage                AttributeReqBodyUsage    `yaml:"req_body_usage"`
	Sensitive                   bool                     `yaml:"sensitive"`
	CreateOnly                  bool                     `yaml:"create_only"` // leveraged for defining plan modifier which avoids updates on this attribute
	PresentInAnyResponse        bool                     `yaml:"present_in_any_response"`
	RequestOnlyRequiredOnCreate bool                     `yaml:"request_only_required_on_create"` // Flags API property only present in create request body as required. These properties are intentionally modified to optional attributes to preserve creation validation, but allows omitting value on imports.
	ListTypeAsMap               bool                     `yaml:"list_type_as_map,omitempty"`      // Flags API property to be defined as a Map type while API defines as list of key-value pairs (used for tags and labels).
	SkipStateListMerge          bool                     `yaml:"skip_state_list_merge,omitempty"` // When true, nested list elements are not merged with state during unmarshal.
	ImmutableComputed           bool                     `yaml:"immutable_computed,omitempty"`    // When true, adds UseStateForUnknown plan modifier for computed attributes.
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
	AllRequestBodies        AttributeReqBodyUsage = "all_request_bodies" // by default attribute is sent in request bodies
	OmitInUpdateBody        AttributeReqBodyUsage = "omit_in_update_body"
	SendNullAsNullOnUpdate  AttributeReqBodyUsage = "send_null_as_null_on_update"  // attributes with null value are sent as null in update request body
	SendNullAsEmptyOnUpdate AttributeReqBodyUsage = "send_null_as_empty_on_update" // attributes with null value are sent as empty in update request body (collections only)
	OmitAlways              AttributeReqBodyUsage = "omit_always"                  // this covers computed-only attributes and attributes which are only used for path/query params
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
	Discriminator *Discriminator `yaml:"discriminator,omitempty"`
	Attributes    Attributes     `yaml:"attributes"`
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
	Name    string            `yaml:"name,omitempty"` // Nested object name without the "TF" & "Model" prefix and suffix. Used for type overrides.
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
		Name:    name,
		Model:   fmt.Sprintf("customtypes.ObjectValue[TF%sModel]", name),
		Schema:  fmt.Sprintf("customtypes.NewObjectType[TF%sModel](ctx)", name),
	}
}

func NewCustomListType(elemType ElemType) *CustomType {
	elemTypeStr := ElementTypeToModelString[elemType]
	return &CustomType{
		Package: CustomTypesPkg,
		Model:   fmt.Sprintf("customtypes.ListValue[%s]", elemTypeStr),
		Schema:  fmt.Sprintf("customtypes.NewListType[%s](ctx)", elemTypeStr),
	}
}

func NewCustomNestedListType(name string) *CustomType {
	return &CustomType{
		Package: CustomTypesPkg,
		Name:    name,
		Model:   fmt.Sprintf("customtypes.NestedListValue[TF%sModel]", name),
		Schema:  fmt.Sprintf("customtypes.NewNestedListType[TF%sModel](ctx)", name),
	}
}

func NewCustomSetType(elemType ElemType) *CustomType {
	elemTypeStr := ElementTypeToModelString[elemType]
	return &CustomType{
		Package: CustomTypesPkg,
		Model:   fmt.Sprintf("customtypes.SetValue[%s]", elemTypeStr),
		Schema:  fmt.Sprintf("customtypes.NewSetType[%s](ctx)", elemTypeStr),
	}
}

func NewCustomNestedSetType(name string) *CustomType {
	return &CustomType{
		Package: CustomTypesPkg,
		Name:    name,
		Model:   fmt.Sprintf("customtypes.NestedSetValue[TF%sModel]", name),
		Schema:  fmt.Sprintf("customtypes.NewNestedSetType[TF%sModel](ctx)", name),
	}
}

func NewCustomMapType(elemType ElemType) *CustomType {
	elemTypeStr := ElementTypeToModelString[elemType]
	return &CustomType{
		Package: CustomTypesPkg,
		Model:   fmt.Sprintf("customtypes.MapValue[%s]", elemTypeStr),
		Schema:  fmt.Sprintf("customtypes.NewMapType[%s](ctx)", elemTypeStr),
	}
}

func NewCustomNestedMapType(name string) *CustomType {
	return &CustomType{
		Package: CustomTypesPkg,
		Model:   fmt.Sprintf("customtypes.NestedMapValue[TF%sModel]", name),
		Schema:  fmt.Sprintf("customtypes.NewNestedMapType[TF%sModel](ctx)", name),
	}
}
