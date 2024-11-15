package conversion

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func PatchPayloadHasChangesTpf[T any, R any](ctx context.Context, diags *diag.Diagnostics, state, plan *T, converter func(ctx context.Context, input *T, diags *diag.Diagnostics) *R, reqPatch *R) bool {
	stateReq := converter(ctx, state, diags)
	planReq := converter(ctx, plan, diags)
	if diags.HasError() {
		return false
	}
	noChanges, err := PatchPayloadNoChanges(stateReq, planReq, reqPatch)
	if err != nil {
		diags.AddError(fmt.Sprintf("error creating patch payload %T", reqPatch), err.Error())
		return false
	}
	return !noChanges
}
