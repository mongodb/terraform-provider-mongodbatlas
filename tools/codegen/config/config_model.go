package config

type Config struct {
	Resources map[string]Resource `yaml:"resources"`
}

type Resource struct {
	Create        *APIOperation `yaml:"create"`
	Read          *APIOperation `yaml:"read"`
	Update        *APIOperation `yaml:"update"`
	Delete        *APIOperation `yaml:"delete"`
	VersionHeader string        `yaml:"version_header"`
	SchemaOptions SchemaOptions `yaml:"schema"`
}

type APIOperation struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

type SchemaOptions struct {
	Ignores   []string            `yaml:"ignores"`
	Aliases   map[string]string   `yaml:"aliases"`
	Overrides map[string]Override `yaml:"overrides"`
	Timeouts  []string            `yaml:"timeouts"`
}

type Override struct {
	Computability *Computability `yaml:"computability,omitempty"`
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
