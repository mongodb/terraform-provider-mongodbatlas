package acc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var _ plancheck.PlanCheck = debugPlan{}

type debugPlan struct{}

func (e debugPlan) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	rd, err := json.Marshal(req.Plan)
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("error marshaling machine-readable plan output: %s", err))
	}
	tflog.Info(ctx, fmt.Sprintf("req.Plan - %s\n", string(rd)))
}

func DebugPlan() plancheck.PlanCheck {
	return debugPlan{}
}
