package clouduserorgassignment

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewTFModel(ctx context.Context, apiResp *admin.OrgUserResponse) (*TFModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var rolesObj types.Object
	var rolesDiags diag.Diagnostics

	if apiResp == nil {
		return nil, diags
	}

	rolesObj, rolesDiags = NewTFRoles(ctx, &apiResp.Roles)
	diags.Append(rolesDiags...)

	var teamIds types.Set
	if apiResp.TeamIds != nil {
		teamIds, _ = types.SetValueFrom(ctx, types.StringType, *apiResp.TeamIds)
	} else {
		teamIds = types.SetNull(types.StringType)
	}
	return &TFModel{
		Country:             types.StringValue(apiResp.GetCountry()),
		CreatedAt:           types.StringValue(conversion.TimeToString(apiResp.GetCreatedAt())),
		FirstName:           types.StringValue(apiResp.GetFirstName()),
		UserId:              types.StringValue(apiResp.GetId()),
		InvitationCreatedAt: types.StringValue(conversion.TimeToString(apiResp.GetInvitationCreatedAt())),
		InvitationExpiresAt: types.StringValue(conversion.TimeToString(apiResp.GetInvitationExpiresAt())),
		InviterUsername:     types.StringValue(apiResp.GetInviterUsername()),
		LastAuth:            types.StringValue(conversion.TimeToString(apiResp.GetLastAuth())),
		LastName:            types.StringValue(apiResp.GetLastName()),
		MobileNumber:        types.StringValue(apiResp.GetMobileNumber()),
		// OrgId:               types.StringNull(), // Not returned by API, must be set elsewhere
		OrgMembershipStatus: types.StringValue(apiResp.GetOrgMembershipStatus()),
		Roles:               rolesObj,
		TeamIds:             teamIds,
		Username:            types.StringValue(apiResp.GetUsername()),
	}, diags
}

func NewTFRoles(ctx context.Context, roles *admin.OrgUserRolesResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if roles == nil {
		return types.ObjectNull(RolesObjectAttrTypes), diags
	}
	orgRoles, _ := types.SetValueFrom(ctx, types.StringType, roles.GetOrgRoles())
	praList := NewTFProjectRoleAssignments(ctx, roles.GroupRoleAssignments)
	rolesObj, _ := types.ObjectValue(
		RolesObjectAttrTypes,
		map[string]attr.Value{
			"org_roles":                orgRoles,
			"project_role_assignments": praList,
		},
	)
	return rolesObj, diags
}

func NewTFProjectRoleAssignments(ctx context.Context, groupRoleAssignments *[]admin.GroupRoleAssignment) types.List {
	var projectRoleAssignments []TFRolesProjectRoleAssignmentsModel
	if groupRoleAssignments != nil {
		for _, pra := range *groupRoleAssignments {
			projectId := types.StringNull()
			if pra.GroupId != nil {
				projectId = types.StringValue(*pra.GroupId)
			}
			projectRoles := types.SetNull(types.StringType)
			if pra.GroupRoles != nil {
				projectRoles, _ = types.SetValueFrom(ctx, types.StringType, *pra.GroupRoles)
			}
			projectRoleAssignments = append(projectRoleAssignments, TFRolesProjectRoleAssignmentsModel{
				ProjectId:    projectId,
				ProjectRoles: projectRoles,
			})
		}
	}
	praList, _ := types.ListValueFrom(ctx, ProjectRoleAssignmentsAttrType.ElemType.(types.ObjectType), projectRoleAssignments)
	return praList
}

func NewOrgUserReq(ctx context.Context, plan *TFModel) (*admin.OrgUserRequest, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	roles, rolesDiags := NewOrgUserRolesRequest(ctx, plan.Roles)
	diags.Append(rolesDiags...)
	return &admin.OrgUserRequest{
		Roles:    *roles,
		Username: plan.Username.ValueString(),
	}, diags
}

func NewAtlasUpdateReq(ctx context.Context, plan *TFModel) (*admin.OrgUserUpdateRequest, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	roles, rolesDiags := NewOrgUserRolesRequest(ctx, plan.Roles)
	diags.Append(rolesDiags...)

	return &admin.OrgUserUpdateRequest{
		Roles: roles,
	}, diags
}

func NewOrgUserRolesRequest(ctx context.Context, rolesObj types.Object) (*admin.OrgUserRolesRequest, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if rolesObj.IsNull() || rolesObj.IsUnknown() {
		return &admin.OrgUserRolesRequest{}, diags
	}
	var rolesModel TFRolesModel
	diags.Append(rolesObj.As(ctx, &rolesModel, basetypes.ObjectAsOptions{})...)
	var orgRoles []string
	if !rolesModel.OrgRoles.IsNull() && !rolesModel.OrgRoles.IsUnknown() {
		rolesModel.OrgRoles.ElementsAs(ctx, &orgRoles, false)
	}
	// project_role_assignments is computed/read-only, do not send in request
	return &admin.OrgUserRolesRequest{
		OrgRoles: orgRoles,
	}, diags
}
