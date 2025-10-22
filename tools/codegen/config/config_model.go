package config

type Config struct {
	Resources map[string]Resource `yaml:"resources"`
}

type Resource struct {
	Create        *APIOperation `yaml:"create"`
	Read          *APIOperation `yaml:"read"`
	Update        *APIOperation `yaml:"update"`
	Delete        *APIOperation `yaml:"delete"`
	VersionHeader string        `yaml:"version_header"` // when not defined latest version defined in API Spec of the resource is used
	SchemaOptions SchemaOptions `yaml:"schema"`
}

type APIOperation struct {
	Wait              *Wait  `yaml:"wait"`
	Path              string `yaml:"path"`
	Method            string `yaml:"method"`
	StaticRequestBody string `yaml:"static_request_body"` // use at the moment for Delete when it's done with a PATCH or PUT and needs to send a static request body.
}

type Wait struct {
	StateProperty     string   `yaml:"state_property"` // defined in camel case as found in API response body, e.g. "stateName"
	PendingStates     []string `yaml:"pending_states"`
	TargetStates      []string `yaml:"target_states"`
	TimeoutSeconds    int      `yaml:"timeout_seconds"`
	MinTimeoutSeconds int      `yaml:"min_timeout_seconds"`
	DelaySeconds      int      `yaml:"delay_seconds"`
}

type SchemaOptions struct {
	UseCustomNestedTypes *bool               `yaml:"use_custom_nested_types"` // Tmp flag to disable custom nested types usage until typing is supported for all nested attributes. Defaults to true. - CLOUDP-352973
	Ignores              []string            `yaml:"ignores"`
	Aliases              map[string]string   `yaml:"aliases"` // only supports modifying path param names, full alias support is not yet implemented in conversion logic for request/response bodies
	Overrides            map[string]Override `yaml:"overrides"`
	Timeouts             []string            `yaml:"timeouts"`
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
