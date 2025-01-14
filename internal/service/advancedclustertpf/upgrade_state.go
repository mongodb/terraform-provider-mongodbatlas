package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *rs) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		1: {StateUpgrader: stateUpgraderFromV1},
	}
}

func stateUpgraderFromV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	diags := &resp.Diagnostics
	projectID, name := getProjectIDClusterNameFromRawState(diags, req.RawState)
	if diags.HasError() {
		return
	}
	setStateResponse(ctx, &resp.Diagnostics, &resp.State, projectID, name)
}
