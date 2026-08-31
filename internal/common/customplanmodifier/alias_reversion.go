package customplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AliasReversion prevents a resource from reverting from canonicalAttribute to
// deprecatedAttribute once state contains only the canonical attribute. Both
// arguments are root-level string attribute names.
func AliasReversion(canonicalAttribute, deprecatedAttribute string) planmodifier.String {
	return aliasReversion{
		canonicalPath:  path.Root(canonicalAttribute),
		deprecatedPath: path.Root(deprecatedAttribute),
	}
}

type aliasReversion struct {
	canonicalPath  path.Path
	deprecatedPath path.Path
}

func (m aliasReversion) Description(_ context.Context) string {
	return fmt.Sprintf("Prevents reverting from %s to %s.", m.canonicalPath, m.deprecatedPath)
}

func (m aliasReversion) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m aliasReversion) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var planCanonical, planDeprecated, stateCanonical, stateDeprecated types.String

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, m.canonicalPath, &planCanonical)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, m.deprecatedPath, &planDeprecated)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, m.canonicalPath, &stateCanonical)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, m.deprecatedPath, &stateDeprecated)...)
	if resp.Diagnostics.HasError() || !IsAliasReversion(stateCanonical, stateDeprecated, planCanonical, planDeprecated) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		m.deprecatedPath,
		"Cannot revert deprecated attribute alias",
		fmt.Sprintf("This resource already uses %s in state. Use %s instead of %s.", m.canonicalPath, m.canonicalPath, m.deprecatedPath),
	)
}

// IsAliasReversion reports whether state containing only the canonical alias is
// being changed back to the deprecated alias. A state that contains both aliases
// is treated as legacy state and is not subject to the one-way restriction.
func IsAliasReversion(stateCanonical, stateDeprecated, planCanonical, planDeprecated types.String) bool {
	return IsKnown(stateCanonical) &&
		stateDeprecated.IsNull() &&
		planCanonical.IsNull() &&
		IsKnown(planDeprecated)
}
