package streamworkspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// FailoverRegionsWriteOnce requires replacement when failover_regions is already configured
// and being changed. The backend treats failover_regions as write-once: it can be set via PATCH
// exactly once, but cannot be modified or removed after that.
type FailoverRegionsWriteOnce struct{}

func (m FailoverRegionsWriteOnce) Description(_ context.Context) string {
	return "Requires replacement when failover_regions is already configured and being changed."
}

func (m FailoverRegionsWriteOnce) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m FailoverRegionsWriteOnce) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if len(req.StateValue.Elements()) == 0 {
		return
	}
	if req.PlanValue.IsUnknown() {
		return
	}
	if req.PlanValue.Equal(req.StateValue) {
		return
	}
	resp.RequiresReplace = true
}
