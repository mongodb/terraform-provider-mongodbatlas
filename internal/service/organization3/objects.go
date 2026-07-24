package organization3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"go.mongodb.org/atlas-sdk/v20250312020/admin"
)

var secretMetadataObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"secret_id":  types.StringType,
	"created_at": types.StringType,
	"expires_at": types.StringType,
}}

var rotationObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"expires_after_hours":        types.Int64Type,
	"rotate_before_expiry_hours": types.Int64Type,
	"secret_version":             types.Int64Type,
	"current_secret":             secretMetadataObjectType,
	"old_secret":                 secretMetadataObjectType,
}}

func RotationFromObject(ctx context.Context, object types.Object) (TFClientSecretRotationModel, diag.Diagnostics) {
	if object.IsNull() {
		return TFClientSecretRotationModel{}, nil
	}
	var rotation TFClientSecretRotationModel
	diags := object.As(ctx, &rotation, basetypes.ObjectAsOptions{})
	return rotation, diags
}

func rotationToObject(ctx context.Context, rotation *TFClientSecretRotationModel) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(ctx, rotationObjectType.AttrTypes, rotation)
}

func SecretMetadataFromObject(ctx context.Context, object types.Object) (TFSecretMetadataModel, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return TFSecretMetadataModel{}, nil
	}
	var metadata TFSecretMetadataModel
	diags := object.As(ctx, &metadata, basetypes.ObjectAsOptions{})
	return metadata, diags
}

func secretMetadataToObject(ctx context.Context, metadata *TFSecretMetadataModel) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(ctx, secretMetadataObjectType.AttrTypes, metadata)
}

func secretMetadataNullObject(ctx context.Context) (types.Object, diag.Diagnostics) {
	return types.ObjectNull(secretMetadataObjectType.AttrTypes), nil
}

func secretMetadataFromAPI(secret *admin.ServiceAccountSecret) TFSecretMetadataModel {
	return TFSecretMetadataModel{
		SecretID:  types.StringValue(secret.GetId()),
		CreatedAt: types.StringValue(formatRFC3339(secret.GetCreatedAt())),
		ExpiresAt: types.StringValue(formatRFC3339(secret.GetExpiresAt())),
	}
}

func mergeSecretMetadataFromAPI(stateMeta *TFSecretMetadataModel, apiSecret *admin.ServiceAccountSecret) TFSecretMetadataModel {
	merged := secretMetadataFromAPI(apiSecret)
	if stateMeta != nil && !stateMeta.SecretID.IsNull() && stateMeta.SecretID.ValueString() != "" {
		merged.SecretID = stateMeta.SecretID
	}
	return merged
}

func findSecretByID(secrets []admin.ServiceAccountSecret, secretID string) (*admin.ServiceAccountSecret, bool) {
	for i := range secrets {
		if secrets[i].GetId() == secretID {
			return &secrets[i], true
		}
	}
	return nil, false
}

// SecretMetadataObjectTypeForTest exposes secretMetadataObjectType for unit tests.
func SecretMetadataObjectTypeForTest() types.ObjectType {
	return secretMetadataObjectType
}

// markPlanForRotation shapes the Terraform plan for an upcoming rotation: known secret_version and old_secret,
// unknown current_secret and client_secret so apply creates the next Atlas secret.
func markPlanForRotation(
	ctx context.Context,
	plan *TFModel,
	planRotation, stateRotation *TFClientSecretRotationModel,
	targetVersion int64,
) diag.Diagnostics {
	var diags diag.Diagnostics
	planRotation.SecretVersion = types.Int64Value(targetVersion)
	promotedOld, promoteDiags := secretMetadataForRotationPromotion(ctx, planRotation, stateRotation)
	diags.Append(promoteDiags...)
	if diags.HasError() {
		return diags
	}
	oldSecretObject, oldDiags := secretMetadataToObject(ctx, &promotedOld)
	diags.Append(oldDiags...)
	if diags.HasError() {
		return diags
	}
	planRotation.OldSecret = oldSecretObject
	planRotation.CurrentSecret = types.ObjectUnknown(secretMetadataObjectType.AttrTypes)
	plan.ClientSecret = types.StringUnknown()
	rotationObject, objectDiags := rotationToObject(ctx, planRotation)
	diags.Append(objectDiags...)
	if diags.HasError() {
		return diags
	}
	plan.ClientSecretRotation = rotationObject
	return diags
}

// MarkPlanForRotationForTest exposes markPlanForRotation for unit tests.
func MarkPlanForRotationForTest(
	ctx context.Context,
	plan *TFModel,
	planRotation, stateRotation *TFClientSecretRotationModel,
	targetVersion int64,
) diag.Diagnostics {
	return markPlanForRotation(ctx, plan, planRotation, stateRotation, targetVersion)
}

func secretMetadataForRotationPromotion(
	ctx context.Context,
	planRotation, stateRotation *TFClientSecretRotationModel,
) (TFSecretMetadataModel, diag.Diagnostics) {
	for _, object := range []types.Object{stateRotation.CurrentSecret, planRotation.CurrentSecret} {
		meta, diags := SecretMetadataFromObject(ctx, object)
		if diags.HasError() {
			return TFSecretMetadataModel{}, diags
		}
		if !meta.SecretID.IsNull() && meta.SecretID.ValueString() != "" {
			return meta, nil
		}
	}
	var diags diag.Diagnostics
	diags.AddError(
		"Missing current_secret metadata",
		"current_secret.secret_id must be present in state to promote to old_secret during rotation; restore current_secret and old_secret in Terraform state if a prior apply left them null",
	)
	return TFSecretMetadataModel{}, diags
}

func validateRotationStateForWrite(ctx context.Context, rotationObject types.Object) diag.Diagnostics {
	if rotationObject.IsNull() || rotationObject.IsUnknown() {
		return nil
	}
	rotation, diags := RotationFromObject(ctx, rotationObject)
	if diags.HasError() {
		return diags
	}
	secretVersion := int64(0)
	if !rotation.SecretVersion.IsNull() && !rotation.SecretVersion.IsUnknown() {
		secretVersion = rotation.SecretVersion.ValueInt64()
	}
	if secretVersion < 1 {
		return nil
	}
	currentMeta, currentDiags := SecretMetadataFromObject(ctx, rotation.CurrentSecret)
	diags.Append(currentDiags...)
	if currentMeta.SecretID.IsNull() || currentMeta.SecretID.ValueString() == "" {
		diags.AddError(
			"Invalid rotation state",
			"Refusing to write state with an empty current_secret; restore current_secret metadata in Terraform state before applying again",
		)
	}
	if secretVersion >= 2 {
		oldMeta, oldDiags := SecretMetadataFromObject(ctx, rotation.OldSecret)
		diags.Append(oldDiags...)
		if oldMeta.SecretID.IsNull() || oldMeta.SecretID.ValueString() == "" {
			diags.AddError(
				"Invalid rotation state",
				"Refusing to write state with an empty old_secret at secret_version >= 2; restore old_secret metadata in Terraform state before applying again",
			)
		}
	}
	return diags
}

func rotationWithDefaults(rotation *TFClientSecretRotationModel, expiresAfterHours int64) {
	if rotation.RotateBeforeExpiryHours.IsNull() || rotation.RotateBeforeExpiryHours.IsUnknown() {
		rotation.RotateBeforeExpiryHours = types.Int64Value(effectiveRotateBeforeExpiryHours(expiresAfterHours, types.Int64Null()))
	}
}
