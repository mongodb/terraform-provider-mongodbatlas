package config

type Config struct {
	Resources map[string]Resource `yaml:"resources"`
}

type Resource struct {
	Create             *APIOperation `yaml:"create"`
	Read               *APIOperation `yaml:"read"`
	Update             *APIOperation `yaml:"update"`
	Delete             *APIOperation `yaml:"delete"`
	MoveState          *MoveState    `yaml:"move_state"`
	IDAttributes       []string      `yaml:"id_attributes"`
	DeprecationMessage *string       `yaml:"deprecation_message"`
	DataSources        *DataSources  `yaml:"datasources"` // when defined, data source(s) are generated with independent schema options
	VersionHeader      string        `yaml:"version_header"`
	SchemaOptions      SchemaOptions `yaml:"schema"`
}

// DataSources defines the configuration for generating data sources independently from resources
type DataSources struct {
	Read          *APIOperation `yaml:"read"`   // singular data source read operation
	List          *APIOperation `yaml:"list"`   // plural data source list operation
	SchemaOptions SchemaOptions `yaml:"schema"` // data source specific schema options (aliases, overrides, ignores)
}

type APIOperation struct {
	Wait              *Wait  `yaml:"wait"`
	Path              string `yaml:"path"`
	Method            string `yaml:"method"`
	StaticRequestBody string `yaml:"static_request_body"`
	SchemaIgnore      bool   `yaml:"schema_ignore"`
}

type Wait struct {
	StateProperty     string   `yaml:"state_property"` // defined in camel case as found in API response body, e.g. "stateName"
	PendingStates     []string `yaml:"pending_states"`
	TargetStates      []string `yaml:"target_states"`
	TimeoutSeconds    int      `yaml:"timeout_seconds"`
	MinTimeoutSeconds int      `yaml:"min_timeout_seconds"`
	DelaySeconds      int      `yaml:"delay_seconds"`
}

type MoveState struct {
	SourceResources []string `yaml:"source_resources"`
}

type SchemaOptions struct {
	Aliases   map[string]string   `yaml:"aliases"` // keys and values use camelCase (e.g., groupId: projectId, nestedObject.innerAttr: renamedAttr). Supports path params and request/response body fields via APIName preservation and apiname tag generation
	Overrides map[string]Override `yaml:"overrides"`
	Ignores   []string            `yaml:"ignores"`
}

type Override struct {
	Computability       *Computability `yaml:"computability,omitempty"`
	Sensitive           *bool          `yaml:"sensitive"`
	IncludeNullOnUpdate *bool          `yaml:"include_null_on_update"`
	Type                *Type          `yaml:"type"`
	Description         string         `yaml:"description"`
	PlanModifiers       []PlanModifier `yaml:"plan_modifiers"`
	Validators          []Validator    `yaml:"validators"`
}

type PlanModifier struct {
	Definition string   `yaml:"definition"`
	Imports    []string `yaml:"imports"`
}

type Validator struct {
	Definition string   `yaml:"definition"`
	Imports    []string `yaml:"imports"`
}

type Computability struct {
	Optional bool `yaml:"optional"`
	Computed bool `yaml:"computed"`
	Required bool `yaml:"required"`
}

type Type string

const (
	List Type = "list"
	Set  Type = "set"
)
