package resourcepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"
)

func NewTFResourcePolicy(ctx context.Context, apiResp *admin.ApiAtlasResourcePolicy) (*TFResourcePolicyModel, diag.Diagnostics) {
	createdBy, diags := newUserMetadataObjectType(ctx, apiResp.CreatedByUser)
	if diags.HasError() {
		return nil, diags
	}
	lastUpdatedBy, diags := newUserMetadataObjectType(ctx, apiResp.LastUpdatedByUser)
	if diags.HasError() {
		return nil, diags
	}
	policies, diags := newPoliciesListType(ctx, apiResp.Policies)
	if diags.HasError() {
		return nil, diags
	}
	return &TFResourcePolicyModel{
		CreatedByUser:     *createdBy,
		CreatedDate:       types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.CreatedDate)),
		ID:                types.StringValue(apiResp.GetId()),
		LastUpdatedByUser: *lastUpdatedBy,
		LastUpdatedDate:   types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.LastUpdatedDate)),
		Name:              types.StringPointerValue(apiResp.Name),
		OrgID:             types.StringValue(apiResp.GetOrgId()),
		Policies:          *policies,
		Version:           types.StringValue(apiResp.GetVersion()),
	}, nil
}

func newUserMetadataObjectType(ctx context.Context, userResp *admin.ApiAtlasUserMetadata) (*types.Object, diag.Diagnostics) {
	if userResp == nil {
		empty := types.ObjectNull(UserMetadataObjectType.AttrTypes)
		return &empty, nil
	}
	tfModel := TFUserMetadataModel{
		ID:   types.StringValue(userResp.GetId()),
		Name: types.StringValue(userResp.GetName()),
	}
	objType, diags := types.ObjectValueFrom(ctx, UserMetadataObjectType.AttrTypes, tfModel)
	if diags.HasError() {
		return nil, diags
	}
	return &objType, nil
}

func newPoliciesListType(ctx context.Context, apiResp *[]admin.ApiAtlasPolicy) (*types.List, diag.Diagnostics) {
	if apiResp == nil {
		empty := types.ListNull(PolicyObjectType)
		return &empty, nil
	}
	var tfList []TFPolicyModel
	for _, policy := range *apiResp {
		tfPolicy := TFPolicyModel{
			Body: types.StringValue(policy.GetBody()),
			ID:   types.StringValue(policy.GetId()),
		}
		tfList = append(tfList, tfPolicy)
	}
	listType, diags := types.ListValueFrom(ctx, PolicyObjectType, tfList)
	if diags.HasError() {
		return nil, diags
	}
	return &listType, nil
}

func NewResourcePolicyReq(ctx context.Context, plan *TFResourcePolicyModel) (*admin.ApiAtlasResourcePolicy, diag.Diagnostics) {
	return &admin.ApiAtlasResourcePolicy{}, nil
}
