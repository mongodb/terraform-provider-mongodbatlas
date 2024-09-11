package resourcepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"
)

func NewTFResourcePolicyModel(ctx context.Context, apiResp *admin.ApiAtlasResourcePolicy) (*TFResourcePolicyModel, diag.Diagnostics) {
	createdByUser, diags := NewUserMetadataObjectType(ctx, apiResp.CreatedByUser)
	if diags.HasError() {
		return nil, diags
	}
	lastUpdatedByUser, diags := NewUserMetadataObjectType(ctx, apiResp.LastUpdatedByUser)
	if diags.HasError() {
		return nil, diags
	}
	policies, diags := NewPolicyObjectType(ctx, apiResp.Policies)
	if diags.HasError() {
		return nil, diags
	}
	return &TFResourcePolicyModel{
		CreatedByUser:     *createdByUser,
		CreatedDate:       types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.CreatedDate)),
		ID:                types.StringPointerValue(apiResp.Id),
		LastUpdatedByUser: *lastUpdatedByUser,
		LastUpdatedDate:   types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.LastUpdatedDate)),
		Name:              types.StringPointerValue(apiResp.Name),
		OrgID:             types.StringPointerValue(apiResp.OrgId),
		Policies:          *policies,
		Version:           types.StringPointerValue(apiResp.Version),
	}, nil
}

func NewUserMetadataObjectType(ctx context.Context, apiResp *admin.ApiAtlasUserMetadata) (*types.Object, diag.Diagnostics) {
	if apiResp == nil {
		empty := types.ObjectNull(UserMetadataObjectType.AttrTypes)
		return &empty, nil
	}
	tfModel := TFUserMetadataModel{
		ID:   types.StringPointerValue(apiResp.Id),
		Name: types.StringPointerValue(apiResp.Name),
	}
	objType, diags := types.ObjectValueFrom(ctx, UserMetadataObjectType.AttrTypes, tfModel)
	if diags.HasError() {
		return nil, diags
	}
	return &objType, nil
}

func NewPolicyObjectType(ctx context.Context, apiResp *[]admin.ApiAtlasPolicy) (*types.List, diag.Diagnostics) {
	if apiResp == nil {
		empty := types.ListNull(PolicyObjectType)
		return &empty, nil
	}
	tfModels := make([]TFPolicyModel, len(*apiResp))
	for i, item := range *apiResp {
		tfModels[i] = TFPolicyModel{
			Body: types.StringPointerValue(item.Body),
			ID:   types.StringPointerValue(item.Id),
		}
	}
	listType, diags := types.ListValueFrom(ctx, PolicyObjectType, tfModels)
	if diags.HasError() {
		return nil, diags
	}
	return &listType, nil
}

func NewResourcePolicyReq(ctx context.Context, plan *TFResourcePolicyModel) (*admin.ApiAtlasResourcePolicy, diag.Diagnostics) {
	return &admin.ApiAtlasResourcePolicy{}, nil
}
