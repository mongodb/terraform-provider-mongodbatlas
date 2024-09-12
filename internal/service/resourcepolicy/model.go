package resourcepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"
)

func NewTFResourcePolicyModel(ctx context.Context, input *admin.ApiAtlasResourcePolicy) (*TFResourcePolicyModel, diag.Diagnostics) {
	createdByUser, diags := NewUserMetadataObjectType(ctx, input.CreatedByUser)
	if diags.HasError() {
		return nil, diags
	}
	lastUpdatedByUser, diags := NewUserMetadataObjectType(ctx, input.LastUpdatedByUser)
	if diags.HasError() {
		return nil, diags
	}
	policies, diags := NewPolicyObjectType(ctx, input.Policies)
	if diags.HasError() {
		return nil, diags
	}
	return &TFResourcePolicyModel{
		CreatedByUser:     *createdByUser,
		CreatedDate:       types.StringPointerValue(conversion.TimePtrToStringPtr(input.CreatedDate)),
		ID:                types.StringPointerValue(input.Id),
		LastUpdatedByUser: *lastUpdatedByUser,
		LastUpdatedDate:   types.StringPointerValue(conversion.TimePtrToStringPtr(input.LastUpdatedDate)),
		Name:              types.StringPointerValue(input.Name),
		OrgID:             types.StringPointerValue(input.OrgId),
		Policies:          *policies,
		Version:           types.StringPointerValue(input.Version),
	}, nil
}

func NewUserMetadataObjectType(ctx context.Context, input *admin.ApiAtlasUserMetadata) (*types.Object, diag.Diagnostics) {
	var nilPointer *admin.ApiAtlasUserMetadata
	if input == nilPointer {
		empty := types.ObjectNull(UserMetadataObjectType.AttrTypes)
		return &empty, nil
	}
	tfModel := TFUserMetadataModel{
		ID:   types.StringPointerValue(input.Id),
		Name: types.StringPointerValue(input.Name),
	}
	objType, diags := types.ObjectValueFrom(ctx, UserMetadataObjectType.AttrTypes, tfModel)
	if diags.HasError() {
		return nil, diags
	}
	return &objType, nil
}

func NewPolicyObjectType(ctx context.Context, input *[]admin.ApiAtlasPolicy) (*types.List, diag.Diagnostics) {
	var nilPointer *[]admin.ApiAtlasPolicy
	if input == nilPointer {
		empty := types.ListNull(PolicyObjectType)
		return &empty, nil
	}
	tfModels := make([]TFPolicyModel, len(*input))
	for i, item := range *input {
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

func NewTFPoliciesModelToSDK(ctx context.Context, input types.List) (*[]admin.ApiAtlasPolicyCreate, diag.Diagnostics) {
	var tfPolicies []TFPolicyModel
	if diags := input.ElementsAs(ctx, &tfPolicies, false); diags.HasError() {
		return nil, diags
	}
	apiModels := make([]admin.ApiAtlasPolicyCreate, len(tfPolicies))
	for i, item := range tfPolicies {
		apiModels[i] = admin.ApiAtlasPolicyCreate{
			Body: item.Body.ValueStringPointer(),
		}
	}
	return &apiModels, nil
}
