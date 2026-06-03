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
	"secret_id":    types.StringType,
	"created_at":   types.StringType,
	"expires_at":   types.StringType,
	"last_used_at": types.StringType,
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
	if object.IsNull() {
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
	model := TFSecretMetadataModel{
		SecretID:  types.StringValue(secret.GetId()),
		CreatedAt: types.StringValue(formatRFC3339(secret.GetCreatedAt())),
		ExpiresAt: types.StringValue(formatRFC3339(secret.GetExpiresAt())),
	}
	if secret.HasLastUsedAt() {
		model.LastUsedAt = types.StringValue(formatRFC3339(secret.GetLastUsedAt()))
	} else {
		model.LastUsedAt = types.StringNull()
	}
	return model
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

func rotationWithDefaults(rotation *TFClientSecretRotationModel, expiresAfterHours int64) {
	if rotation.RotateBeforeExpiryHours.IsNull() || rotation.RotateBeforeExpiryHours.IsUnknown() {
		rotation.RotateBeforeExpiryHours = types.Int64Value(effectiveRotateBeforeExpiryHours(expiresAfterHours, types.Int64Null()))
	}
}
