package conversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func PatchPayloadTpf[TFModel any, SDKRequest any](ctx context.Context, diags *diag.Diagnostics, state, plan *TFModel, converter func(ctx context.Context, input *TFModel, diags *diag.Diagnostics) *SDKRequest) *SDKRequest {
	stateReq := converter(ctx, state, diags)
	planReq := converter(ctx, plan, diags)
	if diags.HasError() {
		return nil
	}
	reqPatch := new(SDKRequest)
	noChanges, err := PatchPayloadNoChanges(stateReq, planReq, reqPatch)
	if err != nil {
		diags.AddError(fmt.Sprintf("error creating patch payload %T", reqPatch), err.Error())
		return nil
	}
	if noChanges {
		return nil
	}
	return reqPatch
}
