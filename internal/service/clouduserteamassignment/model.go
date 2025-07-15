package clouduserteamassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

func NewTFUserTeamAssignmentModel(ctx context.Context, orgID, teamID string, apiResp *admin.OrgUserResponse) (*TFUserTeamAssignmentModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var rolesObj types.Object
	var rolesDiags diag.Diagnostics

	if apiResp == nil {
		diags.AddError("Invalid data", "The API response for the user team assignment is nil and cannot be processed.")
		return nil, diags
	}

	rolesObj, rolesDiags = NewTFRolesModel(ctx, &apiResp.Roles)
	diags.Append(rolesDiags...)

	userTeamAssignment := TFUserTeamAssignmentModel{
		OrgId:               types.StringValue(orgID),
		TeamId:              types.StringValue(teamID),
		UserId:              types.StringValue(apiResp.GetId()),
		Username:            types.StringValue(apiResp.GetUsername()),
		OrgMembershipStatus: types.StringValue(apiResp.GetOrgMembershipStatus()),
		Roles:               rolesObj,
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

func NewTFRolesModel(ctx context.Context, roles *admin.OrgUserRolesResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if roles == nil {
		return types.ObjectNull(RolesObjectAttrTypes), diags
	}

	var orgRoles types.Set
	if roles.OrgRoles == nil || len(*roles.OrgRoles) == 0 {
		orgRoles = types.SetNull(types.StringType)
	} else {
		orgRoles, _ = types.SetValueFrom(ctx, types.StringType, *roles.OrgRoles)
	}

	projectRoleAssignmentsList := NewTFProjectRoleAssignments(ctx, roles.GroupRoleAssignments)

	rolesObj, _ := types.ObjectValue(
		RolesObjectAttrTypes,
		map[string]attr.Value{
			"project_role_assignments": projectRoleAssignmentsList,
			"org_roles":                orgRoles,
		},
	)

	return rolesObj, diags
}

func NewTFProjectRoleAssignments(ctx context.Context, groupRoleAssignments *[]admin.GroupRoleAssignment) types.List {
	if groupRoleAssignments == nil {
		return types.ListNull(ProjectRoleAssignmentsAttrType)
	}

	var projectRoleAssignments []TFProjectRoleAssignmentsModel

	for _, pra := range *groupRoleAssignments {
		projectID := types.StringPointerValue(pra.GroupId)
		var projectRoles types.Set
		if pra.GroupRoles == nil || len(*pra.GroupRoles) == 0 {
			projectRoles = types.SetNull(types.StringType)
		} else {
			projectRoles, _ = types.SetValueFrom(ctx, types.StringType, pra.GroupRoles)
		}
		projectRoleAssignments = append(projectRoleAssignments, TFProjectRoleAssignmentsModel{
			ProjectID:    projectID,
			ProjectRoles: projectRoles,
		})
	}

	praList, _ := types.ListValueFrom(ctx, ProjectRoleAssignmentsAttrType.ElemType.(types.ObjectType), projectRoleAssignments)
	return praList
}

func NewUserTeamAssignmentReq(ctx context.Context, plan *TFUserTeamAssignmentModel) (*admin.AddOrRemoveUserFromTeam, diag.Diagnostics) {
	addOrRemoveUserFromTeam := admin.AddOrRemoveUserFromTeam{
		Id: plan.UserId.ValueString(),
	}
	return &addOrRemoveUserFromTeam, nil
}
