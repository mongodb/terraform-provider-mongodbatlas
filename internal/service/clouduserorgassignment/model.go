package clouduserorgassignment

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewTFModel(ctx context.Context, apiResp *admin.OrgUserResponse, orgID string) (*TFModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var rolesObj types.Object
	var rolesDiags diag.Diagnostics

	if apiResp == nil {
		return nil, diags
	}

	rolesObj, rolesDiags = NewTFRoles(ctx, &apiResp.Roles)
	diags.Append(rolesDiags...)

	teamIDs := conversion.TFSetValueOrNull(ctx, apiResp.TeamIds, types.StringType)

	return &TFModel{
		OrgId:               types.StringValue(orgID),
		Country:             types.StringPointerValue(apiResp.Country),
		CreatedAt:           types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.CreatedAt)),
		FirstName:           types.StringPointerValue(apiResp.FirstName),
		UserId:              types.StringValue(apiResp.GetId()),
		InvitationCreatedAt: types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.InvitationCreatedAt)),
		InvitationExpiresAt: types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.InvitationExpiresAt)),
		InviterUsername:     types.StringPointerValue(apiResp.InviterUsername),
		LastAuth:            types.StringPointerValue(conversion.TimePtrToStringPtr(apiResp.LastAuth)),
		LastName:            types.StringPointerValue(apiResp.LastName),
		MobileNumber:        types.StringPointerValue(apiResp.MobileNumber),
		OrgMembershipStatus: types.StringValue(apiResp.GetOrgMembershipStatus()),
		Roles:               rolesObj,
		TeamIds:             teamIDs,
		Username:            types.StringValue(apiResp.GetUsername()),
	}, diags
}

func NewTFRoles(ctx context.Context, roles *admin.OrgUserRolesResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if roles == nil {
		return types.ObjectNull(RolesObjectAttrTypes), diags
	}
	orgRoles := conversion.TFSetValueOrNull(ctx, roles.OrgRoles, types.StringType)
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
	if groupRoleAssignments == nil {
		return types.ListNull(ProjectRoleAssignmentsAttrType)
	}

	var projectRoleAssignments []TFRolesProjectRoleAssignmentsModel

	for _, pra := range *groupRoleAssignments {
		projectID := types.StringPointerValue(pra.GroupId)
		projectRoles := conversion.TFSetValueOrNull(ctx, pra.GroupRoles, types.StringType)

		projectRoleAssignments = append(projectRoleAssignments, TFRolesProjectRoleAssignmentsModel{
			ProjectId:    projectID,
			ProjectRoles: projectRoles,
		})
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
		return &admin.OrgUserRolesRequest{
			OrgRoles: nil,
		}, diags
	}
	var rolesModel TFRolesModel
	diags.Append(rolesObj.As(ctx, &rolesModel, basetypes.ObjectAsOptions{})...)
	var orgRoles []string
	if !rolesModel.OrgRoles.IsNull() && !rolesModel.OrgRoles.IsUnknown() {
		rolesModel.OrgRoles.ElementsAs(ctx, &orgRoles, false)
	} else {
		orgRoles = nil
	}

	return &admin.OrgUserRolesRequest{
		OrgRoles: orgRoles,
	}, diags
}
