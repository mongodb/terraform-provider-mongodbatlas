package update

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func PatchPayloadCluster[TFModel any, SDKRequest any](ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, converter func(ctx context.Context, input *TFModel, diags *diag.Diagnostics) *SDKRequest) *SDKRequest {
	stateReq := converter(ctx, state, diags)
	planReq := converter(ctx, plan, diags)
	if diags.HasError() {
		return nil
	}
	req, err := PatchPayload(stateReq, planReq)
	if err != nil {
		diags.AddError(fmt.Sprintf("error creating patch payload %T", req), err.Error())
		return nil
	}
	return req
}
