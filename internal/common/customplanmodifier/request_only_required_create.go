package customplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequestOnlyRequiredOnCreate returns a plan modifier that fails planning if the value is
// missing (null/unknown) during create (i.e., when state is null), but allows omission on read/import.
func RequestOnlyRequiredOnCreate() RequestOnlyRequiredOnCreateModifier {
	return &requestOnlyRequiredOnCreateAttributePlanModifier{}
}

// Single interface so the modifier can be applied to any attribute type.
type RequestOnlyRequiredOnCreateModifier interface {
	planmodifier.String
	planmodifier.Bool
	planmodifier.Int64
	planmodifier.Float64
	planmodifier.Number
	planmodifier.List
	planmodifier.Map
	planmodifier.Set
	planmodifier.Object
}

type requestOnlyRequiredOnCreateAttributePlanModifier struct{}

func (m *requestOnlyRequiredOnCreateAttributePlanModifier) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m *requestOnlyRequiredOnCreateAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Ensures that create operations fail when attempting to create a resource with a missing required attribute."
}

func (m *requestOnlyRequiredOnCreateAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if isCreate(&req.State) && (req.PlanValue.IsNull() || req.PlanValue.IsUnknown()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s is required when creating this resource", req.Path),
			fmt.Sprintf("Provide a value for %s during resource creation.", req.Path),
		)
	}
}

func (m *requestOnlyRequiredOnCreateAttributePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if isCreate(&req.State) && (req.PlanValue.IsNull() || req.PlanValue.IsUnknown()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s is required when creating this resource", req.Path),
			fmt.Sprintf("Provide a value for %s during resource creation.", req.Path),
		)
	}
}

func (m *requestOnlyRequiredOnCreateAttributePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if isCreate(&req.State) && (req.PlanValue.IsNull() || req.PlanValue.IsUnknown()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s is required when creating this resource", req.Path),
			fmt.Sprintf("Provide a value for %s during resource creation.", req.Path),
		)
	}
}

func (m *requestOnlyRequiredOnCreateAttributePlanModifier) PlanModifyFloat64(ctx context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	if isCreate(&req.State) && (req.PlanValue.IsNull() || req.PlanValue.IsUnknown()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s is required when creating this resource", req.Path),
			fmt.Sprintf("Provide a value for %s during resource creation.", req.Path),
		)
	}
}

func (m *requestOnlyRequiredOnCreateAttributePlanModifier) PlanModifyNumber(ctx context.Context, req planmodifier.NumberRequest, resp *planmodifier.NumberResponse) {
	if isCreate(&req.State) && (req.PlanValue.IsNull() || req.PlanValue.IsUnknown()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s is required when creating this resource", req.Path),
			fmt.Sprintf("Provide a value for %s during resource creation.", req.Path),
		)
	}
}

func (m *requestOnlyRequiredOnCreateAttributePlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if isCreate(&req.State) && (req.PlanValue.IsNull() || req.PlanValue.IsUnknown()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s is required when creating this resource", req.Path),
			fmt.Sprintf("Provide a value for %s during resource creation.", req.Path),
		)
	}
}

func (m *requestOnlyRequiredOnCreateAttributePlanModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if isCreate(&req.State) && (req.PlanValue.IsNull() || req.PlanValue.IsUnknown()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s is required when creating this resource", req.Path),
			fmt.Sprintf("Provide a value for %s during resource creation.", req.Path),
		)
	}
}

func (m *requestOnlyRequiredOnCreateAttributePlanModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if isCreate(&req.State) && (req.PlanValue.IsNull() || req.PlanValue.IsUnknown()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s is required when creating this resource", req.Path),
			fmt.Sprintf("Provide a value for %s during resource creation.", req.Path),
		)
	}
}

func (m *requestOnlyRequiredOnCreateAttributePlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if isCreate(&req.State) && (req.PlanValue.IsNull() || req.PlanValue.IsUnknown()) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("%s is required when creating this resource", req.Path),
			fmt.Sprintf("Provide a value for %s during resource creation.", req.Path),
		)
	}
}
