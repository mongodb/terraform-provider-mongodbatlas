package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *rs) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		1: {
			PriorSchema:   resourceSchemaV1(ctx),
			StateUpgrader: stateUpgraderFromV1,
		},
	}
}

func stateUpgraderFromV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	resp.Diagnostics.AddError("UpgradeState not implemented", "UpgradeState not implemented yet")
}
