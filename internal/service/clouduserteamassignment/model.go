package clouduserteamassignment

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

// TODO: `ctx` parameter and `diags` return value can be removed if tf schema has no complex data types (e.g., schema.ListAttribute, schema.SetAttribute)
func NewTFUserTeamAssignmentModel(ctx context.Context, orgID, teamID string, apiResp *admin.OrgUserResponse) (*TFUserTeamAssignmentModel, diag.Diagnostics) {
	// complexAttr, diagnostics := types.ListValueFrom(ctx, InnerObjectType, newTFComplexAttrModel(apiResp.ComplexAttr))
	// if diagnostics.HasError() {
	// 	return nil, diagnostics
	// }
	var diags diag.Diagnostics
	if apiResp == nil {
		diags.AddError("Invalid data", "The API response for the user team assignment is nil and cannot be processed.")
		return nil, diags
	}
	rolesModel, diags := NewTFRolesModel(ctx, &apiResp.Roles)
	if diags.HasError() {
		return nil, diags
	}
	userTeamAssignment := TFUserTeamAssignmentModel{
		OrgID:               types.StringValue(orgID),
		TeamID:              types.StringValue(teamID),
		UserID:              types.StringValue(apiResp.Id),
		Username:            types.StringValue(apiResp.Username),
		OrgMembershipStatus: types.StringValue(apiResp.OrgMembershipStatus),
		Roles:               rolesModel,
		InvitationCreatedAt: types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.InvitationCreatedAt)),
		InvitationExpiresAt: types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.InvitationExpiresAt)),
		InviterUsername:     types.StringPointerValue(apiResp.InviterUsername),
		Country:             types.StringPointerValue(apiResp.Country),
		FirstName:           types.StringPointerValue(apiResp.FirstName),
		LastName:            types.StringPointerValue(apiResp.LastName),
		CreatedAt:           types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.CreatedAt)),
		LastAuth:            types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.LastAuth)),
		MobileNumber:        types.StringPointerValue(apiResp.MobileNumber),
	}

	userTeamAssignment.TeamIDs = types.SetNull(types.StringType)
	if apiResp.TeamIds != nil {
		userTeamAssignment.TeamIDs, diags = types.SetValueFrom(ctx, types.StringType, apiResp.TeamIds)
		if diags.HasError() {
			return nil, diags
		}
	}

	return &userTeamAssignment, nil
}

func NewTFRolesModel(ctx context.Context, apiResp *admin.OrgUserRolesResponse) (*TFRolesModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	projectRoleAssignments := make([]*TFProjectRoleAssignmentsModel, len(*apiResp.GroupRoleAssignments))
	for i, roleAssignment := range *apiResp.GroupRoleAssignments {
		projectRoleAssignments[i] = &TFProjectRoleAssignmentsModel{
			ProjectID:    types.StringValue(roleAssignment.GetGroupId()),
			ProjectRoles: types.SetNull(types.StringType),
		}
		if roleAssignment.GetGroupRoles() != nil {
			projectRoles, diags := types.SetValueFrom(ctx, types.StringType, roleAssignment.GetGroupRoles())
			if diags.HasError() {
				return nil, diags
			}
			projectRoleAssignments[i].ProjectRoles = projectRoles
		}
	}

	orgRoles, _ := types.SetValueFrom(ctx, types.StringType, *apiResp.OrgRoles)

	return &TFRolesModel{
		ProjectRoleAssignments: projectRoleAssignments,
		OrgRoles:               orgRoles,
	}, diags
}

func NewCloudUserTeamAssignmentReq(ctx context.Context, plan *TFUserTeamAssignmentModel) (*admin.AddOrRemoveUserFromTeam, diag.Diagnostics) {
	addOrRemoveUserFromTeam := admin.AddOrRemoveUserFromTeam{
		Id: plan.UserID.ValueString(),
	}
	return &addOrRemoveUserFromTeam, nil
}
