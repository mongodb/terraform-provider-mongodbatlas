package clouduserteamassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

func NewTFUserTeamAssignmentModel(ctx context.Context, apiResp *admin.OrgUserResponse) (*TFUserTeamAssignmentModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var rolesObj types.Object
	var rolesDiags diag.Diagnostics

	if apiResp == nil {
		return nil, diags
	}

	rolesObj, rolesDiags = NewTFRolesModel(ctx, &apiResp.Roles)
	diags.Append(rolesDiags...)

	teamIDs := conversion.TFSetValueOrNull(ctx, apiResp.TeamIds, types.StringType)

	userTeamAssignment := TFUserTeamAssignmentModel{
		UserId:              types.StringValue(apiResp.GetId()),
		Username:            types.StringValue(apiResp.GetUsername()),
		OrgMembershipStatus: types.StringValue(apiResp.GetOrgMembershipStatus()),
		Roles:               rolesObj,
		TeamIds:             teamIDs,
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

	return &userTeamAssignment, nil
}

func NewTFRolesModel(ctx context.Context, roles *admin.OrgUserRolesResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if roles == nil {
		return types.ObjectNull(RolesObjectAttrTypes), diags
	}

	orgRoles := conversion.TFSetValueOrNull(ctx, roles.OrgRoles, types.StringType)

	projectRoleAssignmentsSet := NewTFProjectRoleAssignments(ctx, roles.GroupRoleAssignments)

	rolesObj, _ := types.ObjectValue(
		RolesObjectAttrTypes,
		map[string]attr.Value{
			"org_roles":                orgRoles,
			"project_role_assignments": projectRoleAssignmentsSet,
		},
	)
	return rolesObj, diags
}

func NewTFProjectRoleAssignments(ctx context.Context, groupRoleAssignments *[]admin.GroupRoleAssignment) types.Set {
	if groupRoleAssignments == nil {
		return types.SetNull(ProjectRoleAssignmentsAttrType)
	}

	var projectRoleAssignments []TFProjectRoleAssignmentsModel

	for _, pra := range *groupRoleAssignments {
		projectID := types.StringPointerValue(pra.GroupId)
		projectRoles := conversion.TFSetValueOrNull(ctx, pra.GroupRoles, types.StringType)

		projectRoleAssignments = append(projectRoleAssignments, TFProjectRoleAssignmentsModel{
			ProjectId:    projectID,
			ProjectRoles: projectRoles,
		})
	}

	praSet, _ := types.SetValueFrom(ctx, ProjectRoleAssignmentsAttrType.ElemType.(types.ObjectType), projectRoleAssignments)
	return praSet
}

func NewUserTeamAssignmentReq(ctx context.Context, plan *TFUserTeamAssignmentModel) (*admin.AddOrRemoveUserFromTeam, diag.Diagnostics) {
	addOrRemoveUserFromTeam := admin.AddOrRemoveUserFromTeam{
		Id: plan.UserId.ValueString(),
	}
	return &addOrRemoveUserFromTeam, nil
}
