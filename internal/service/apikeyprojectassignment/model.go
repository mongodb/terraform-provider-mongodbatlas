package apikeyprojectassignment

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312003/admin"
)

func NewTFModel(ctx context.Context, apiResp *admin.PaginatedApiApiUser, apiKeyID, projectID string) (*TFModel, diag.Diagnostics) {
	apiKeyUserDetails := apiResp.GetResults()

	for _, apiKey := range apiKeyUserDetails {
		if apiKey.GetId() != apiKeyID {
			continue
		}

		return apiKeyUserDetailsToTFModel(ctx, &apiKey, projectID)
	}

	return &TFModel{
		ProjectId: types.StringValue(projectID),
		ApiKeyId:  types.StringValue(apiKeyID),
	}, nil
}

func NewTFModelPatch(ctx context.Context, apiKey *admin.ApiKeyUserDetails, projectID string) (*TFModel, diag.Diagnostics) {
	return apiKeyUserDetailsToTFModel(ctx, apiKey, projectID)
}

func apiKeyUserDetailsToTFModel(ctx context.Context, apiKey *admin.ApiKeyUserDetails, projectID string) (*TFModel, diag.Diagnostics) {
	// filter for the specific project roles
	projectRoles := make([]string, 0, len(apiKey.GetRoles()))
	for _, role := range apiKey.GetRoles() {
		if !strings.HasPrefix(role.GetRoleName(), "ORG_") && role.GetGroupId() == projectID {
			projectRoles = append(projectRoles, role.GetRoleName())
		}
	}

	roleNames, diags := types.SetValueFrom(ctx, types.StringType, projectRoles)
	if diags.HasError() {
		return nil, diags
	}

	return &TFModel{
		RoleNames: roleNames,
		ApiKeyId:  types.StringValue(apiKey.GetId()),
		ProjectId: types.StringValue(projectID),
	}, diags
}

func NewAtlasCreateReq(ctx context.Context, plan *TFModel) (*[]admin.UserAccessRoleAssignment, diag.Diagnostics) {
	roleNames := conversion.TypesSetToString(ctx, plan.RoleNames)
	return &[]admin.UserAccessRoleAssignment{
		{
			Roles:  &roleNames,
			UserId: plan.ApiKeyId.ValueStringPointer(),
		},
	}, nil
}

func NewAtlasUpdateReq(ctx context.Context, plan *TFModel) (*admin.UpdateAtlasProjectApiKey, diag.Diagnostics) {
	roleNames := conversion.TypesSetToString(ctx, plan.RoleNames)
	return &admin.UpdateAtlasProjectApiKey{
		Roles: &roleNames,
	}, nil
}
