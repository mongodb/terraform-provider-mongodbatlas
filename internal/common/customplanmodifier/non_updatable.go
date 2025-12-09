package customplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	planmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// NonUpdatable returns a plan modifier that ensures that update operations fails when the attribute is changed.
// This is useful for attributes only supported in create and not in update.
// It shows a helpful error message helping the user to update their config to match the state.
// Never use a schema.Default for create only attributes, instead use WithXXXDefault, the default will lead to plan changes that are not expected after import.
// Implement CopyFromPlan if the attribute is not in the API Response.
func NonUpdatable() NonUpdatableModifier {
	return &nonUpdatableAttributePlanModifier{}
}

// Single interface allows customplanmodifier.NonUpdatable() to be used by all attribute types, simplifying code generation of auto-generated resources
type NonUpdatableModifier interface {
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

// Plan modifier that implements non-updatable behavior for multiple attribute types
type nonUpdatableAttributePlanModifier struct{}

func (d *nonUpdatableAttributePlanModifier) Description(ctx context.Context) string {
	return d.MarkdownDescription(ctx)
}

func (d *nonUpdatableAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Ensures that update operations fail when attempting to modify a non-updatable attribute."
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	validateNonUpdatable(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	validateNonUpdatable(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	validateNonUpdatable(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyFloat64(ctx context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	validateNonUpdatable(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyNumber(ctx context.Context, req planmodifier.NumberRequest, resp *planmodifier.NumberResponse) {
	validateNonUpdatable(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	validateNonUpdatable(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	validateNonUpdatable(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *nonUpdatableAttributePlanModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	validateNonUpdatable(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *nonUpdatableAttributePlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	validateNonUpdatable(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

// validateNonUpdatable checks if an attribute value has changed and adds an error if it has
func validateNonUpdatable(planValue, stateValue attr.Value, attrPath path.Path, diagnostics *diag.Diagnostics,
) {
	if !stateValue.IsNull() && !stateValue.Equal(planValue) {
		diagnostics.AddError(
			fmt.Sprintf("%s cannot be updated", attrPath),
			fmt.Sprintf("%s cannot be updated", attrPath),
		)
	}
}
