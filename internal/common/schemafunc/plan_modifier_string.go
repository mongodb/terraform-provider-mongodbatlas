package schemafunc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func PlanModifyStringUpdateOnly() UpdateOnlyString {
	return UpdateOnlyString{}
}

type UpdateOnlyString struct{}

func (u UpdateOnlyString) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if (req.PlanValue.IsUnknown() || !req.PlanValue.IsNull()) && req.State.Raw.IsNull() {
		errMsg := fmt.Sprintf("Update only attribute set on create: %s", req.Path)
		resp.Diagnostics.AddError(errMsg, errMsg)
	}
}

func (u UpdateOnlyString) Description(ctx context.Context) string {
	return u.MarkdownDescription(ctx)
}

func (u UpdateOnlyString) MarkdownDescription(ctx context.Context) string {
	return "Checks the attribute is never set on create"
}
