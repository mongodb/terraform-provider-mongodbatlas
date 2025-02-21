package customplanmodifier

import (
	"context"
	"fmt"

	planmodifier "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
)

func InstanceSizeStringAttributePlanModifier() planmodifier.String {
	return &instanceSizeStringAttributePlanModifier{}
}

type instanceSizeStringAttributePlanModifier struct {
}

func (d *instanceSizeStringAttributePlanModifier) Description(ctx context.Context) string {
	return d.MarkdownDescription(ctx)
}

func (d *instanceSizeStringAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Ensures that a deprecation warning is displayed when instance_size is M2 or M5."
}

func (d *instanceSizeStringAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	planAttributeValue := req.PlanValue
	stateAttributeValue := req.StateValue

	if stateAttributeValue.ValueString() == "M2" || stateAttributeValue.ValueString() == "M5" ||
		planAttributeValue.ValueString() == "M2" || planAttributeValue.ValueString() == "M5" {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf(constant.DeprecationSharedTier, constant.ServerlessSharedEOLDate, "https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/serverless-shared-migration-guide"),
			req.Path.String(),
		)
		return
	}
}
