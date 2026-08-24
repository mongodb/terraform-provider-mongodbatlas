package streamprocessor

// Vertical Autoscaling (Phase 1) for stream processors.
//
// The autoscaling config (`options.autoscaling`) and the read-only top-level
// `effectiveTier` field are still `@Hidden` in the Atlas API and are therefore
// not part of the officially generated SDK yet. During development they are
// provided by a hand-patched LOCAL copy of the versioned SDK wired in via the
// `replace` directive in go.mod (see .sdkstub/, per the ASP "Terraform
// Development" wiki's "Local validation in Cloud Dev" guidance). Once the fields
// are un-hidden and land in the official SDK, remove the stub + replace directive;
// no changes to this file should be required because it already targets the real
// admin.StreamsAutoscaling type.

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/atlas-sdk-go/admin"
)

type TFAutoscalingModel struct {
	MinTier types.String `tfsdk:"min_tier"`
	MaxTier types.String `tfsdk:"max_tier"`
}

var AutoscalingObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"min_tier": types.StringType,
	"max_tier": types.StringType,
}}

// newAutoscalingReq converts the TF autoscaling object into the SDK request type.
// Returns nil when the block is not configured (unset), which the API treats as
// "keep persisted config" on update and "no autoscaling" on create.
func newAutoscalingReq(ctx context.Context, autoscaling types.Object) (*admin.StreamsAutoscaling, diag.Diagnostics) {
	if autoscaling.IsNull() || autoscaling.IsUnknown() {
		return nil, nil
	}
	tfModel := &TFAutoscalingModel{}
	if diags := autoscaling.As(ctx, tfModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}
	// Block presence is Terraform's enablement signal. The API requires this
	// transport field for create/PATCH semantics, but it is not exposed in state.
	req := &admin.StreamsAutoscaling{Enabled: new(true)}
	if !tfModel.MinTier.IsNull() && !tfModel.MinTier.IsUnknown() {
		req.MinTier = tfModel.MinTier.ValueStringPointer()
	}
	if !tfModel.MaxTier.IsNull() && !tfModel.MaxTier.IsUnknown() {
		req.MaxTier = tfModel.MaxTier.ValueStringPointer()
	}
	return req, nil
}

// autoscalingFromPlanOptions extracts the autoscaling SDK request from the plan's
// `options` object, for endpoints (e.g. :startWith) that take autoscaling top-level.
// Returns nil when options or the autoscaling block is unset.
func autoscalingFromPlanOptions(ctx context.Context, options types.Object) (*admin.StreamsAutoscaling, diag.Diagnostics) {
	if options.IsNull() || options.IsUnknown() {
		return nil, nil
	}
	optionsModel := &TFOptionsModel{}
	if diags := options.As(ctx, optionsModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}
	return newAutoscalingReq(ctx, optionsModel.Autoscaling)
}

// resolveAutoscalingForUpdate applies the PATCH tri-state semantics: it returns the
// plan's autoscaling request when configured; an explicit disable ({enabled:false})
// when the block was removed from the plan but was present in prior state; or nil when
// autoscaling is absent from both, so the API preserves whatever is persisted.
func resolveAutoscalingForUpdate(ctx context.Context, plan, state *TFStreamProcessorRSModel) (*admin.StreamsAutoscaling, diag.Diagnostics) {
	planAutoscaling, diags := autoscalingFromPlanOptions(ctx, plan.Options)
	if diags.HasError() {
		return nil, diags
	}
	if planAutoscaling != nil {
		return planAutoscaling, nil
	}
	stateAutoscaling, diags := autoscalingFromPlanOptions(ctx, state.Options)
	if diags.HasError() {
		return nil, diags
	}
	if stateAutoscaling != nil {
		return &admin.StreamsAutoscaling{Enabled: new(false)}, nil
	}
	return nil, nil
}

// effectiveTierFromResp derives the read-only `effective_tier` from the API response.
// Falls back to the baseline `tier` when the API does not return an explicit
// effectiveTier (they are equal whenever autoscaling is disabled).
func effectiveTierFromResp(effectiveTier string, tier *string) types.String {
	if effectiveTier != "" {
		return types.StringValue(effectiveTier)
	}
	return types.StringPointerValue(tier)
}

// convertAutoscalingToTF converts the SDK response type into a TF object.
func convertAutoscalingToTF(ctx context.Context, autoscaling *admin.StreamsAutoscaling) (types.Object, diag.Diagnostics) {
	if autoscaling == nil {
		return types.ObjectNull(AutoscalingObjectType.AttributeTypes()), nil
	}
	tfModel := TFAutoscalingModel{
		MinTier: types.StringPointerValue(autoscaling.MinTier),
		MaxTier: types.StringPointerValue(autoscaling.MaxTier),
	}
	return types.ObjectValueFrom(ctx, AutoscalingObjectType.AttributeTypes(), tfModel)
}
