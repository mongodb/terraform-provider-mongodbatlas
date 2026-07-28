package organization3

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	defaultExpiresAfterHours = int64(720)
	// Default overlap window: rotate when half of expires_after_hours remains.
	defaultRotateBeforeExpiryDivisor int64 = 2
)

type TFModel struct {
	Name                 types.String `tfsdk:"name"`
	OrgOwnerID           types.String `tfsdk:"org_owner_id"`
	OrgID                types.String `tfsdk:"org_id"`
	ClientID             types.String `tfsdk:"client_id"`
	ClientSecret         types.String `tfsdk:"client_secret"`
	ClientSecretRotation types.Object `tfsdk:"client_secret_rotation"`
}

type TFClientSecretRotationModel struct {
	CurrentSecret           types.Object `tfsdk:"current_secret"`
	OldSecret               types.Object `tfsdk:"old_secret"`
	ExpiresAfterHours       types.Int64  `tfsdk:"expires_after_hours"`
	RotateBeforeExpiryHours types.Int64  `tfsdk:"rotate_before_expiry_hours"`
	SecretVersion           types.Int64  `tfsdk:"secret_version"`
}

type TFSecretMetadataModel struct {
	SecretID  types.String `tfsdk:"secret_id"`
	CreatedAt types.String `tfsdk:"created_at"`
	ExpiresAt types.String `tfsdk:"expires_at"`
}

func effectiveRotateBeforeExpiryHours(expiresAfterHours int64, configured types.Int64) int64 {
	if !configured.IsNull() && !configured.IsUnknown() {
		return configured.ValueInt64()
	}
	return expiresAfterHours / defaultRotateBeforeExpiryDivisor
}

func RenewalDue(now, expiresAt time.Time, rotateBeforeExpiryHours int64) bool {
	renewAt := expiresAt.Add(-time.Duration(rotateBeforeExpiryHours) * time.Hour)
	return !now.Before(renewAt) || !now.Before(expiresAt)
}

func shouldDeleteOldSecret(stateVersion int64, oldSecret *TFSecretMetadataModel) bool {
	if oldSecret == nil {
		return false
	}
	return stateVersion >= 2 && !oldSecret.SecretID.IsNull() && oldSecret.SecretID.ValueString() != ""
}

// rotationDecision is whether plan or apply should rotate and the target secret_version.
// Forced when plan secret_version exceeds state; otherwise scheduled from state current_secret.expires_at and RenewalDue.
type rotationDecision struct {
	TargetVersion int64
	ShouldRotate  bool
}

// rotationPlanDecision evaluates rotationDecision for ModifyPlan and Update version resolution.
func rotationPlanDecision(
	ctx context.Context,
	planRotation, stateRotation *TFClientSecretRotationModel,
	stateVersion int64,
	now time.Time,
	diags *diag.Diagnostics,
) rotationDecision {
	targetVersion, shouldRotate := rotationTargetVersion(ctx, planRotation, stateRotation, stateVersion, now, diags)
	return rotationDecision{TargetVersion: targetVersion, ShouldRotate: shouldRotate}
}

func rotationTargetVersion(
	ctx context.Context,
	planRotation, stateRotation *TFClientSecretRotationModel,
	stateVersion int64,
	now time.Time,
	diags *diag.Diagnostics,
) (int64, bool) {
	if !planRotation.SecretVersion.IsUnknown() && !planRotation.SecretVersion.IsNull() {
		configVersion := planRotation.SecretVersion.ValueInt64()
		if configVersion > stateVersion {
			return configVersion, true
		}
	}

	currentSecret, currentDiags := SecretMetadataFromObject(ctx, stateRotation.CurrentSecret)
	diags.Append(currentDiags...)
	if diags.HasError() {
		return 0, false
	}
	if currentSecret.ExpiresAt.IsNull() || currentSecret.ExpiresAt.ValueString() == "" {
		diags.AddError(
			"Missing current_secret.expires_at",
			"State must include current_secret.expires_at from Atlas (create, read, or refresh) before rotation can be scheduled",
		)
		return 0, false
	}
	expiresAt, err := time.Parse(time.RFC3339, currentSecret.ExpiresAt.ValueString())
	if err != nil {
		diags.AddError(
			"Invalid current_secret.expires_at",
			fmt.Sprintf("Could not parse %q as RFC3339: %s", currentSecret.ExpiresAt.ValueString(), err),
		)
		return 0, false
	}
	expiresAfter := rotationPolicyExpiresAfterHours(planRotation, stateRotation)
	rotateBefore := effectiveRotateBeforeExpiryHours(
		expiresAfter,
		rotationPolicyRotateBeforeExpiryHours(planRotation, stateRotation),
	)
	if !RenewalDue(now, expiresAt, rotateBefore) {
		return 0, false
	}
	return stateVersion + 1, true
}

func rotationPolicyExpiresAfterHours(planRotation, stateRotation *TFClientSecretRotationModel) int64 {
	if !planRotation.ExpiresAfterHours.IsNull() && !planRotation.ExpiresAfterHours.IsUnknown() {
		return planRotation.ExpiresAfterHours.ValueInt64()
	}
	return stateRotation.ExpiresAfterHours.ValueInt64()
}

func rotationPolicyRotateBeforeExpiryHours(planRotation, stateRotation *TFClientSecretRotationModel) types.Int64 {
	if !planRotation.RotateBeforeExpiryHours.IsNull() && !planRotation.RotateBeforeExpiryHours.IsUnknown() {
		return planRotation.RotateBeforeExpiryHours
	}
	return stateRotation.RotateBeforeExpiryHours
}

func formatRFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// EffectiveRotateBeforeExpiryHoursForTest exposes effectiveRotateBeforeExpiryHours for unit tests.
func EffectiveRotateBeforeExpiryHoursForTest(expiresAfterHours int64, configured types.Int64) int64 {
	return effectiveRotateBeforeExpiryHours(expiresAfterHours, configured)
}

// ShouldDeleteOldSecretForTest exposes shouldDeleteOldSecret for unit tests.
func ShouldDeleteOldSecretForTest(stateVersion int64, oldSecret *TFSecretMetadataModel) bool {
	return shouldDeleteOldSecret(stateVersion, oldSecret)
}

// RotationPlanDecisionForTest exposes rotationPlanDecision for unit tests.
func RotationPlanDecisionForTest(
	ctx context.Context,
	planRotation, stateRotation *TFClientSecretRotationModel,
	stateVersion int64,
	now time.Time,
	diags *diag.Diagnostics,
) (int64, bool) {
	d := rotationPlanDecision(ctx, planRotation, stateRotation, stateVersion, now, diags)
	return d.TargetVersion, d.ShouldRotate
}
