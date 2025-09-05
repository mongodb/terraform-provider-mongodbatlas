package clouduserprojectassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

func NewTFModel(ctx context.Context, projectID string, apiResp *admin.GroupUserResponse) (*TFModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if apiResp == nil {
		return nil, diags
	}

	roles := conversion.TFSetValueOrNull(ctx, &apiResp.Roles, types.StringType)

	return &TFModel{
		Country:             types.StringPointerValue(apiResp.Country),
		CreatedAt:           types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.CreatedAt)),
		FirstName:           types.StringPointerValue(apiResp.FirstName),
		ProjectId:           types.StringValue(projectID),
		UserId:              types.StringValue(apiResp.Id),
		InvitationCreatedAt: types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.InvitationCreatedAt)),
		InvitationExpiresAt: types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.InvitationExpiresAt)),
		InviterUsername:     types.StringPointerValue(apiResp.InviterUsername),
		LastAuth:            types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.LastAuth)),
		LastName:            types.StringPointerValue(apiResp.LastName),
		MobileNumber:        types.StringPointerValue(apiResp.MobileNumber),
		OrgMembershipStatus: types.StringValue(apiResp.GetOrgMembershipStatus()),
		Roles:               roles,
		Username:            types.StringValue(apiResp.GetUsername()),
	}, diags
}

func NewProjectUserReq(ctx context.Context, plan *TFModel) (*admin.GroupUserRequest, diag.Diagnostics) {
	roleNames := []string{}
	if !plan.Roles.IsNull() && !plan.Roles.IsUnknown() {
		roleNames = conversion.TypesSetToString(ctx, plan.Roles)
	}

	addProjectUserReq := admin.GroupUserRequest{
		Username: plan.Username.ValueString(),
		Roles:    roleNames,
	}
	return &addProjectUserReq, nil
}

func NewAtlasUpdateReq(ctx context.Context, plan *TFModel, currentRoles []string) (addRequests, removeRequests []*admin.AddOrRemoveGroupRole, diags diag.Diagnostics) {
	var desiredRoles []string
	if !plan.Roles.IsNull() && !plan.Roles.IsUnknown() {
		desiredRoles = conversion.TypesSetToString(ctx, plan.Roles)
	}

	rolesToAdd, rolesToRemove := diffRoles(currentRoles, desiredRoles)

	addRequests = make([]*admin.AddOrRemoveGroupRole, 0, len(rolesToAdd))
	for _, role := range rolesToAdd {
		addRequests = append(addRequests, &admin.AddOrRemoveGroupRole{
			GroupRole: role,
		})
	}

	removeRequests = make([]*admin.AddOrRemoveGroupRole, 0, len(rolesToRemove))
	for _, role := range rolesToRemove {
		removeRequests = append(removeRequests, &admin.AddOrRemoveGroupRole{
			GroupRole: role,
		})
	}

	return addRequests, removeRequests, nil
}

func diffRoles(oldRoles, newRoles []string) (toAdd, toRemove []string) {
	oldRolesMap := make(map[string]bool, len(oldRoles))

	for _, role := range oldRoles {
		oldRolesMap[role] = true
	}

	for _, role := range newRoles {
		if oldRolesMap[role] {
			delete(oldRolesMap, role)
		} else {
			toAdd = append(toAdd, role)
		}
	}

	for role := range oldRolesMap {
		toRemove = append(toRemove, role)
	}

	return toAdd, toRemove
}
