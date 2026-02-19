package validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	schemavalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// VariantDefinition describes which sibling attributes are allowed and required when a particular discriminator value is active.
// Required attributes must also be present in the Allowed list.
type VariantDefinition struct {
	Allowed  []string
	Required []string
}

type DiscriminatorDefinition struct {
	Mapping map[string]VariantDefinition
}

// ValidateDiscriminator returns a config-phase string validator that checks sibling
// attribute presence/absence based on the active discriminator value.
//   - If the discriminator value is unknown, null, or not found in Mapping all checks are skipped.
//   - Required attributes for the active variant must be non-null (unknown is accepted as "set").
//   - Type-specific attributes from other variants must be null.
//   - Note: Unset Optional+Computed attributes are null in the config, they only become unknown later during PlanResourceChange
func ValidateDiscriminator(def DiscriminatorDefinition) schemavalidator.String {
	return discriminatorValidator{def: def}
}

type discriminatorValidator struct {
	def DiscriminatorDefinition
}

func (v discriminatorValidator) Description(_ context.Context) string {
	return "validates sibling attributes based on the selected discriminator value"
}

func (v discriminatorValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v discriminatorValidator) ValidateString(_ context.Context, req schemavalidator.StringRequest, resp *schemavalidator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	discriminatorValue := req.ConfigValue.ValueString()
	variant, ok := v.def.Mapping[discriminatorValue]
	if !ok {
		return
	}

	parentPath := req.Path.ParentPath()
	discriminatorName := lastPathStepName(req.Path)

	allTypeSpecific := allTypeSpecificAttrs(v.def)
	activeAllowed := toSet(variant.Allowed)

	for _, name := range variant.Required {
		siblingPath := parentPath.AtName(name)
		if isNullInConfig(req.Config, siblingPath) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing Required Attribute",
				fmt.Sprintf("Attribute %q must be set when %s is %q", name, discriminatorName, discriminatorValue),
			)
		}
	}

	for name := range allTypeSpecific {
		if activeAllowed[name] {
			continue
		}
		siblingPath := parentPath.AtName(name)
		if !isNullInConfig(req.Config, siblingPath) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Attribute Combination",
				fmt.Sprintf("Attribute %q is not allowed when %s is %q", name, discriminatorName, discriminatorValue),
			)
		}
	}
}

func allTypeSpecificAttrs(def DiscriminatorDefinition) map[string]bool {
	result := make(map[string]bool)
	for _, variant := range def.Mapping {
		for _, name := range variant.Allowed {
			result[name] = true
		}
	}
	return result
}

func toSet(names []string) map[string]bool {
	result := make(map[string]bool, len(names))
	for _, name := range names {
		result[name] = true
	}
	return result
}

func lastPathStepName(p path.Path) string {
	steps := p.Steps()
	if len(steps) == 0 {
		return ""
	}
	if nameStep, ok := steps[len(steps)-1].(path.PathStepAttributeName); ok {
		return string(nameStep)
	}
	return ""
}

func isNullInConfig(config tfsdk.Config, attrPath path.Path) bool {
	val, ok := resolveConfigValue(config, attrPath)
	if !ok {
		return true
	}
	return val.IsNull()
}

func resolveConfigValue(config tfsdk.Config, attrPath path.Path) (tftypes.Value, bool) {
	tfPath := pathToTFTypesPath(attrPath)
	rawVal, remaining, err := tftypes.WalkAttributePath(config.Raw, tfPath)
	if err != nil || len(remaining.Steps()) > 0 {
		return tftypes.Value{}, false
	}
	val, ok := rawVal.(tftypes.Value)
	return val, ok
}

func pathToTFTypesPath(p path.Path) *tftypes.AttributePath {
	result := tftypes.NewAttributePath()
	for _, step := range p.Steps() {
		switch s := step.(type) {
		case path.PathStepAttributeName:
			result = result.WithAttributeName(string(s))
		case path.PathStepElementKeyInt:
			result = result.WithElementKeyInt(int(s))
		case path.PathStepElementKeyString:
			result = result.WithElementKeyString(string(s))
		}
	}
	return result
}
