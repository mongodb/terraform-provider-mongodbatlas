package conversion

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SingleListTFToSDK[TFModel, SDKRequest any](ctx context.Context, diags *diag.Diagnostics, input *types.List, fnTransform func(tf TFModel) SDKRequest) SDKRequest {
	var resp SDKRequest
	if input == nil || input.IsUnknown() || input.IsNull() || len(input.Elements()) == 0 {
		return resp
	}
	elements := make([]TFModel, len(input.Elements()))
	diags.Append(input.ElementsAs(ctx, &elements, false)...)
	if diags.HasError() {
		return resp
	}
	return fnTransform(elements[0])
}
