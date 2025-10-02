package customplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	planmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// CreateOnlyAttributePlanModifier returns a plan modifier that ensures that update operations fails when the attribute is changed.
// This is useful for attributes only supported in create and not in update.
// It shows a helpful error message helping the user to update their config to match the state.
// Never use a schema.Default for create only attributes, instead use WithXXXDefault, the default will lead to plan changes that are not expected after import.
// Implement CopyFromPlan if the attribute is not in the API Response.
func CreateOnlyAttributePlanModifier() CreateOnlyModifier {
	return &createOnlyAttributePlanModifier{}
}

type CreateOnlyModifier interface {
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

// Plan modifier that implements create-only behavior for multiple attribute types
type createOnlyAttributePlanModifier struct{}

func (d *createOnlyAttributePlanModifier) Description(ctx context.Context) string {
	return d.MarkdownDescription(ctx)
}

func (d *createOnlyAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Ensures that update operations fail when attempting to modify a create-only attribute."
}

func (d *createOnlyAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	validateCreateOnly(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *createOnlyAttributePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	validateCreateOnly(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *createOnlyAttributePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	validateCreateOnly(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *createOnlyAttributePlanModifier) PlanModifyFloat64(ctx context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	validateCreateOnly(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *createOnlyAttributePlanModifier) PlanModifyNumber(ctx context.Context, req planmodifier.NumberRequest, resp *planmodifier.NumberResponse) {
	validateCreateOnly(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *createOnlyAttributePlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	validateCreateOnly(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *createOnlyAttributePlanModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	validateCreateOnly(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *createOnlyAttributePlanModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	validateCreateOnly(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

func (d *createOnlyAttributePlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	validateCreateOnly(req.PlanValue, req.StateValue, req.Path, &resp.Diagnostics)
}

// validateCreateOnly checks if an attribute value has changed and adds an error if it has
func validateCreateOnly(planValue, stateValue attr.Value, attrPath path.Path, diagnostics *diag.Diagnostics,
) {
	if !stateValue.IsNull() && !stateValue.Equal(planValue) {
		diagnostics.AddError(
			fmt.Sprintf("%s cannot be updated", attrPath),
			fmt.Sprintf("%s cannot be updated", attrPath),
		)
	}
}
