package resourcepolicy

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewTFModel(ctx context.Context, input *admin.ApiAtlasResourcePolicy) (*TFModel, diag.Diagnostics) {
	diags := &diag.Diagnostics{}
	createdByUser := NewUserMetadataObjectType(ctx, input.CreatedByUser, diags)
	lastUpdatedByUser := NewUserMetadataObjectType(ctx, input.LastUpdatedByUser, diags)
	policies := NewTFPolicies(ctx, input.Policies, diags)
	if diags.HasError() {
		return nil, *diags
	}
	return &TFModel{
		CreatedByUser:     createdByUser,
		CreatedDate:       types.StringPointerValue(conversion.TimePtrToStringPtr(input.CreatedDate)),
		ID:                types.StringPointerValue(input.Id),
		LastUpdatedByUser: lastUpdatedByUser,
		LastUpdatedDate:   types.StringPointerValue(conversion.TimePtrToStringPtr(input.LastUpdatedDate)),
		Name:              types.StringPointerValue(input.Name),
		Description:       types.StringPointerValue(input.Description),
		OrgID:             types.StringPointerValue(input.OrgId),
		Policies:          policies,
		Version:           types.StringPointerValue(input.Version),
	}, nil
}

func NewUserMetadataObjectType(ctx context.Context, input *admin.ApiAtlasUserMetadata, diags *diag.Diagnostics) types.Object {
	var nilPointer *admin.ApiAtlasUserMetadata
	if input == nilPointer {
		return types.ObjectNull(UserMetadataObjectType.AttrTypes)
	}
	tfModel := TFUserMetadataModel{
		ID:   types.StringPointerValue(input.Id),
		Name: types.StringPointerValue(input.Name),
	}
	objType, diagsLocal := types.ObjectValueFrom(ctx, UserMetadataObjectType.AttrTypes, tfModel)
	diags.Append(diagsLocal...)
	return objType
}

func NewTFPolicies(ctx context.Context, input *[]admin.ApiAtlasPolicy, diags *diag.Diagnostics) []TFPolicyModel {
	var nilPointer *[]admin.ApiAtlasPolicy
	if input == nilPointer {
		return []TFPolicyModel{}
	}
	tfModels := make([]TFPolicyModel, len(*input))
	for i, item := range *input {
		tfModels[i] = TFPolicyModel{
			Body: types.StringPointerValue(item.Body),
			ID:   types.StringPointerValue(item.Id),
		}
	}
	return tfModels
}

func NewAdminPolicies(ctx context.Context, input []TFPolicyModel) []admin.ApiAtlasPolicyCreate {
	apiModels := make([]admin.ApiAtlasPolicyCreate, len(input))
	for i, item := range input {
		apiModels[i] = admin.ApiAtlasPolicyCreate{
			Body: item.Body.ValueString(),
		}
	}
	return apiModels
}

func NewTFModelDSP(ctx context.Context, orgID string, input []admin.ApiAtlasResourcePolicy) (*TFModelDSP, diag.Diagnostics) {
	diags := &diag.Diagnostics{}
	tfModels := make([]TFModel, len(input))
	for i, item := range input {
		tfModel, diagsLocal := NewTFModel(ctx, &item)
		diags.Append(diagsLocal...)
		tfModels[i] = *tfModel
	}
	if diags.HasError() {
		return nil, *diags
	}
	return &TFModelDSP{
		ResourcePolicies: tfModels,
		Results:          tfModels,
		OrgID:            types.StringValue(orgID),
	}, *diags
}
