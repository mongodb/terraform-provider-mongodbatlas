package customplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// POC (unsetting Optional+Computed sizing/autoscaling attributes on advanced_cluster):
//
// These attribute plan modifiers make it possible to REMOVE an Optional+Computed attribute from the
// configuration and have that removal reflected in the plan (and, paired with keepUnknown in the
// resource-level ModifyPlan and ForceUpdateAttr in the update payload, actually sent to the API as an
// unset) — instead of the current silent no-op where the prior value is kept.
//
// Mechanism: when the attribute is absent from config (ConfigValue is null) but the prior state holds
// a concrete value, the planned value is forced to unknown ("known after apply"). Keeping the field
// Optional+Computed (rather than switching to Optional-only) avoids "provider produced inconsistent
// result after apply", because Computed permits the server-returned value to be accepted into state.
//
// IMPORTANT LIMITATION: Terraform cannot distinguish "removed from config" from "never set in config"
// — both present as a null config value with a populated state. So these modifiers fire on ANY
// omission, which introduces a perpetual "known after apply" for users who simply leave the attribute
// unset. This is the fundamental Optional+Computed indistinguishability tradeoff; see the PR description.

const keepUnknownOnRemovalDesc = "POC: forces the planned value to unknown when the attribute is omitted from config but present in state, so the removal is detectable and can be unset via the API."

func KeepUnknownOnRemovalString() planmodifier.String   { return keepUnknownOnRemovalString{} }
func KeepUnknownOnRemovalBool() planmodifier.Bool       { return keepUnknownOnRemovalBool{} }
func KeepUnknownOnRemovalInt64() planmodifier.Int64     { return keepUnknownOnRemovalInt64{} }
func KeepUnknownOnRemovalFloat64() planmodifier.Float64 { return keepUnknownOnRemovalFloat64{} }

type keepUnknownOnRemovalString struct{}

func (m keepUnknownOnRemovalString) Description(context.Context) string {
	return keepUnknownOnRemovalDesc
}
func (m keepUnknownOnRemovalString) MarkdownDescription(context.Context) string {
	return keepUnknownOnRemovalDesc
}
func (m keepUnknownOnRemovalString) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// StateValue must be a real, non-empty value; empty string is treated as effectively unset.
	if req.ConfigValue.IsNull() && IsKnown(req.StateValue) && req.StateValue.ValueString() != "" {
		resp.PlanValue = types.StringUnknown()
	}
}

type keepUnknownOnRemovalBool struct{}

func (m keepUnknownOnRemovalBool) Description(context.Context) string {
	return keepUnknownOnRemovalDesc
}
func (m keepUnknownOnRemovalBool) MarkdownDescription(context.Context) string {
	return keepUnknownOnRemovalDesc
}
func (m keepUnknownOnRemovalBool) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.ConfigValue.IsNull() && IsKnown(req.StateValue) {
		resp.PlanValue = types.BoolUnknown()
	}
}

type keepUnknownOnRemovalInt64 struct{}

func (m keepUnknownOnRemovalInt64) Description(context.Context) string {
	return keepUnknownOnRemovalDesc
}
func (m keepUnknownOnRemovalInt64) MarkdownDescription(context.Context) string {
	return keepUnknownOnRemovalDesc
}
func (m keepUnknownOnRemovalInt64) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.ConfigValue.IsNull() && IsKnown(req.StateValue) {
		resp.PlanValue = types.Int64Unknown()
	}
}

type keepUnknownOnRemovalFloat64 struct{}

func (m keepUnknownOnRemovalFloat64) Description(context.Context) string {
	return keepUnknownOnRemovalDesc
}
func (m keepUnknownOnRemovalFloat64) MarkdownDescription(context.Context) string {
	return keepUnknownOnRemovalDesc
}
func (m keepUnknownOnRemovalFloat64) PlanModifyFloat64(_ context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	if req.ConfigValue.IsNull() && IsKnown(req.StateValue) {
		resp.PlanValue = types.Float64Unknown()
	}
}
