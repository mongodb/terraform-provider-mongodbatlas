package validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	schemavalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func (v discriminatorValidator) ValidateString(ctx context.Context, req schemavalidator.StringRequest, resp *schemavalidator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	discriminatorValue := req.ConfigValue.ValueString()
	variant, ok := v.def.Mapping[discriminatorValue]
	if !ok {
		tflog.Warn(ctx, "Discriminator validation skipped: no mapping found for discriminator value",
			map[string]any{"path": req.Path.String(), "value": discriminatorValue})
		return
	}

	discriminatorName := lastPathStepName(req.Path)

	siblingAttrs, ok := parentObjectAttrs(ctx, req.Config.Raw, req.Path)
	if !ok {
		tflog.Warn(ctx, "Discriminator validation skipped: unable to resolve parent object attributes",
			map[string]any{"path": req.Path.String()})
		return
	}

	allTypeSpecific := allTypeSpecificAttrs(v.def)
	activeAllowed := toSet(variant.Allowed)

	for _, name := range variant.Required {
		val, exists := siblingAttrs[name]
		if !exists || val.IsNull() {
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
		val, exists := siblingAttrs[name]
		if exists && !val.IsNull() {
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

// parentObjectAttrs walks raw to the parent of attrPath using ApplyTerraform5AttributePathStep
// and returns the parent object's attributes as a map keyed by attribute name.
// If any step in the path cannot be resolved, it returns nil, false.
func parentObjectAttrs(ctx context.Context, raw tftypes.Value, attrPath path.Path) (map[string]tftypes.Value, bool) {
	current := raw
	for _, step := range attrPath.ParentPath().Steps() {
		tfStep, ok := toTFTypesStep(ctx, step)
		if !ok {
			return nil, false
		}
		rawVal, err := current.ApplyTerraform5AttributePathStep(tfStep)
		if err != nil {
			return nil, false
		}
		val, ok := rawVal.(tftypes.Value)
		if !ok {
			return nil, false
		}
		current = val
	}
	var attrs map[string]tftypes.Value
	if err := current.As(&attrs); err != nil {
		return nil, false
	}
	return attrs, true
}

func toTFTypesStep(ctx context.Context, step path.PathStep) (tftypes.AttributePathStep, bool) {
	switch s := step.(type) {
	case path.PathStepAttributeName:
		return tftypes.AttributeName(string(s)), true
	case path.PathStepElementKeyInt:
		return tftypes.ElementKeyInt(int64(s)), true
	case path.PathStepElementKeyString:
		return tftypes.ElementKeyString(string(s)), true
	case path.PathStepElementKeyValue:
		tfVal, err := s.ToTerraformValue(ctx)
		if err != nil {
			return nil, false
		}
		return tftypes.ElementKeyValue(tfVal), true
	default:
		return nil, false
	}
}
