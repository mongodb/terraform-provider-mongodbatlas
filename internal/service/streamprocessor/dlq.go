package streamprocessor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
)

func dlqFromPlanOptions(ctx context.Context, options types.Object) (*admin.StreamsDLQ, diag.Diagnostics) {
	if options.IsNull() || options.IsUnknown() {
		return nil, nil
	}
	optionsModel := &TFOptionsModel{}
	if diags := options.As(ctx, optionsModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}
	return newDlqReq(ctx, optionsModel.Dlq)
}

// resolveDlqForUpdate applies the PATCH tri-state semantics for DLQ. SPM uses
// an empty DLQ object as the explicit clear signal; omission preserves the
// existing DLQ.
func resolveDlqForUpdate(ctx context.Context, plan, state *TFStreamProcessorRSModel) (*admin.StreamsDLQ, bool, diag.Diagnostics) {
	planDLQ, diags := dlqFromPlanOptions(ctx, plan.Options)
	if diags.HasError() {
		return nil, false, diags
	}
	if planDLQ != nil {
		return planDLQ, false, nil
	}
	stateDLQ, diags := dlqFromPlanOptions(ctx, state.Options)
	if diags.HasError() {
		return nil, false, diags
	}
	if stateDLQ != nil {
		return nil, true, nil
	}
	return nil, false, nil
}
