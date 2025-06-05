package config

import "fmt"

type APISpec struct {
	IsDefault *bool  `yaml:"is_default,omitempty"`
	Name      string `yaml:"name"`
	URL       string `yaml:"url"`
}
type Config struct {
	Resources map[string]Resource `yaml:"resources"`
	APISpecs  []APISpec           `yaml:"api_specs,omitempty"`
}

func (c *Config) APISpecsNames() []string {
	specs := make([]string, 0, len(c.APISpecs))
	for _, spec := range c.APISpecs {
		specs = append(specs, spec.Name)
	}
	return specs
}

func (c *Config) DefaultAPISpecName() string {
	defaultSpecs := []string{}
	for _, spec := range c.APISpecs {
		if spec.IsDefault != nil && *spec.IsDefault {
			defaultSpecs = append(defaultSpecs, spec.Name)
		}
	}
	if len(defaultSpecs) == 0 {
		panic("No default API spec defined in the configuration")
	}
	if len(defaultSpecs) > 1 {
		panic(fmt.Sprintf("Multiple default API specs defined in the configuration, please define only one default API spec %v", defaultSpecs))
	}
	return defaultSpecs[0]
}

func (c *Config) GetAPISpecName(name *string) string {
	if name == nil {
		return c.DefaultAPISpecName()
	}
	for _, spec := range c.APISpecs {
		if spec.Name == *name {
			return spec.Name
		}
	}
	panic(fmt.Sprintf("API spec with name %s not found in the configuration", *name))
}

type Resource struct {
	Create        *APIOperation `yaml:"create"`
	Read          *APIOperation `yaml:"read"`
	Update        *APIOperation `yaml:"update"`
	Delete        *APIOperation `yaml:"delete"`
	APISpec       *string       `yaml:"api_spec,omitempty"`
	VersionHeader string        `yaml:"version_header"`
	SchemaOptions SchemaOptions `yaml:"schema"`
}

type APIOperation struct {
	Wait   *Wait  `yaml:"wait"`
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
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
	Ignores   []string            `yaml:"ignores"`
	Aliases   map[string]string   `yaml:"aliases"` // only supports modifying path param names, full alias support is not yet implemented in conversion logic for request/response bodies
	Overrides map[string]Override `yaml:"overrides"`
	Timeouts  []string            `yaml:"timeouts"`
}

type Override struct {
	Computability *Computability `yaml:"computability,omitempty"`
	Sensitive     *bool          `yaml:"sensitive"`
	Description   string         `yaml:"description"`
	PlanModifiers []PlanModifier `yaml:"plan_modifiers"`
	Validators    []Validator    `yaml:"validators"`
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
