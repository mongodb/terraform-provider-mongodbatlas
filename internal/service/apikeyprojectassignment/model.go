package apikeyprojectassignment

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
)

func NewTFModel(ctx context.Context, apiKeys []admin.ApiKeyUserDetails, projectID, apiKeyID string) (*TFModel, diag.Diagnostics) {
	for _, apiKey := range apiKeys {
		if apiKey.GetId() != apiKeyID {
			continue
		}
		return apiKeyUserDetailsToTFModel(ctx, &apiKey, projectID)
	}
	return nil, diag.Diagnostics{diag.NewErrorDiagnostic("API key not found", fmt.Sprintf("API key %s not found", apiKeyID))}
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
		Roles:     roleNames,
		ApiKeyId:  types.StringValue(apiKey.GetId()),
		ProjectId: types.StringValue(projectID),
	}, diags
}

func NewAtlasCreateReq(ctx context.Context, plan *TFModel) (*[]admin.UserAccessRoleAssignment, diag.Diagnostics) {
	roleNames := conversion.TypesSetToString(ctx, plan.Roles)
	return &[]admin.UserAccessRoleAssignment{
		{
			Roles:  &roleNames,
			UserId: plan.ApiKeyId.ValueStringPointer(),
		},
	}, nil
}

func NewAtlasUpdateReq(ctx context.Context, plan *TFModel) (*admin.UpdateAtlasProjectApiKey, diag.Diagnostics) {
	roleNames := conversion.TypesSetToString(ctx, plan.Roles)
	return &admin.UpdateAtlasProjectApiKey{
		Roles: &roleNames,
	}, nil
}

func NewTFModelDSP(ctx context.Context, projectID string, apiKeys []admin.ApiKeyUserDetails) (*TFModelDSP, diag.Diagnostics) {
	results := make([]TFModel, 0, len(apiKeys))
	for _, apiKey := range apiKeys {
		model, diags := apiKeyUserDetailsToTFModel(ctx, &apiKey, projectID)
		if diags.HasError() {
			return nil, diags
		}
		results = append(results, *model)
	}

	return &TFModelDSP{
		ProjectId: types.StringValue(projectID),
		Results:   results,
	}, nil
}
