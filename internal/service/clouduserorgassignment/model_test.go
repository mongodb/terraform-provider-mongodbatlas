package clouduserorgassignment_test

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clouduserorgassignment"
)

const (
	testUserID              = "user-123"
	testUsername            = "jdoe"
	testFirstName           = "John"
	testLastName            = "Doe"
	testCountry             = "CA"
	testMobile              = "+1555123456"
	testInviter             = "admin"
	testOrgMembershipStatus = "ACTIVE"

	testOrgRoleOwner      = "ORG_OWNER"
	testOrgRoleMember     = "ORG_MEMBER"
	testProjectRoleOwner  = "PROJECT_OWNER"
	testProjectRoleRead   = "PROJECT_READ_ONLY"
	testProjectRoleMember = "PROJECT_MEMBER"

	testProjectID1 = "project1"
	testProjectID2 = "project2"

	testOrgID = "org-123"
)

var (
	when          = time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	testCreatedAt = when.Format(time.RFC3339)

	testLastAuth = when.Add(-2 * time.Hour).Format(time.RFC3339)

	testTeamIDs  = []string{"teamA", "teamB"}
	testOrgRoles = []string{"owner", "readWrite"}

	testOrgRolesMultiple   = []string{testOrgRoleOwner, testOrgRoleMember}
	testProjectRolesSingle = []string{testProjectRoleOwner}
)

func createRolesObject(ctx context.Context, orgRoles []string, projectAssignments []clouduserorgassignment.TFRolesProjectRoleAssignmentsModel) types.Object {
	orgRolesSet, _ := types.SetValueFrom(ctx, types.StringType, orgRoles)
	var praList types.List
	if len(projectAssignments) == 0 {
		praList = types.ListNull(clouduserorgassignment.ProjectRoleAssignmentsAttrType.ElemType.(types.ObjectType))
	} else {
		praList, _ = types.ListValueFrom(ctx,
			clouduserorgassignment.ProjectRoleAssignmentsAttrType.ElemType.(types.ObjectType),
			projectAssignments,
		)
	}
	obj, _ := types.ObjectValue(
		clouduserorgassignment.RolesObjectAttrTypes,
		map[string]attr.Value{
			"org_roles":                orgRolesSet,
			"project_role_assignments": praList,
		},
	)
	return obj
}

type sdkToTFModelTestCase struct {
	SDKResp         *admin.OrgUserResponse
	expectedTFModel *clouduserorgassignment.TFModel
}

func TestNewTFModel_SDKToTFModel(t *testing.T) {
	ctx := t.Context()

	fullResp := &admin.OrgUserResponse{
		Id:                  testUserID,
		Username:            testUsername,
		FirstName:           admin.PtrString(testFirstName),
		LastName:            admin.PtrString(testLastName),
		Country:             admin.PtrString(testCountry),
		MobileNumber:        admin.PtrString(testMobile),
		OrgMembershipStatus: testOrgMembershipStatus,
		CreatedAt:           admin.PtrTime(when),
		LastAuth:            admin.PtrTime(when.Add(-2 * time.Hour)),
		TeamIds:             &testTeamIDs,
		Roles: admin.OrgUserRolesResponse{
			OrgRoles: &testOrgRoles,
		},
	}

	orgRolesSet, _ := types.SetValueFrom(ctx, types.StringType, testOrgRoles)
	expectedRoles, _ := types.ObjectValue(
		clouduserorgassignment.RolesObjectAttrTypes,
		map[string]attr.Value{
			"org_roles":                orgRolesSet,
			"project_role_assignments": types.ListNull(clouduserorgassignment.ProjectRoleAssignmentsAttrType),
		},
	)

	expectedTeams, _ := types.SetValueFrom(ctx, types.StringType, testTeamIDs)

	expectedFullModel := &clouduserorgassignment.TFModel{
		UserId:              types.StringValue(testUserID),
		Username:            types.StringValue(testUsername),
		FirstName:           types.StringValue(testFirstName),
		LastName:            types.StringValue(testLastName),
		Country:             types.StringValue(testCountry),
		MobileNumber:        types.StringValue(testMobile),
		InviterUsername:     types.StringNull(),
		OrgMembershipStatus: types.StringValue(testOrgMembershipStatus),
		CreatedAt:           types.StringValue(testCreatedAt),
		InvitationCreatedAt: types.StringNull(),
		InvitationExpiresAt: types.StringNull(),
		LastAuth:            types.StringValue(testLastAuth),
		Roles:               expectedRoles,
		TeamIds:             expectedTeams,
		OrgId:               types.StringValue(testOrgID),
	}

	testCases := map[string]sdkToTFModelTestCase{
		"nil SDK response": {
			SDKResp:         nil,
			expectedTFModel: nil,
		},
		"fully populated SDK response": {
			SDKResp:         fullResp,
			expectedTFModel: expectedFullModel,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotModel, diags := clouduserorgassignment.NewTFModel(ctx, tc.SDKResp, testOrgID)
			assert.False(t, diags.HasError(), "expected no diagnostics")
			assert.Equal(t, tc.expectedTFModel, gotModel, "TFModel did not match expected")
		})
	}
}

func TestNewOrgUserReq(t *testing.T) {
	ctx := t.Context()

	singleOrgRole := []string{"owner"}
	projectAssignment := clouduserorgassignment.TFRolesProjectRoleAssignmentsModel{
		ProjectId:    types.StringValue(testProjectID1),
		ProjectRoles: types.SetValueMust(types.StringType, []attr.Value{types.StringValue(testProjectRoleOwner)}),
	}

	testCases := map[string]struct {
		plan     *clouduserorgassignment.TFModel
		expected *admin.OrgUserRequest
	}{
		"with org roles": {
			plan: &clouduserorgassignment.TFModel{
				Username: types.StringValue("bob"),
				Roles:    createRolesObject(ctx, singleOrgRole, nil),
			},
			expected: &admin.OrgUserRequest{
				Username: "bob",
				Roles:    admin.OrgUserRolesRequest{OrgRoles: singleOrgRole},
			},
		},
		"with both org roles and project role assignments": {
			plan: &clouduserorgassignment.TFModel{
				Username: types.StringValue("alice"),
				Roles:    createRolesObject(ctx, testOrgRolesMultiple, []clouduserorgassignment.TFRolesProjectRoleAssignmentsModel{projectAssignment}),
			},
			expected: &admin.OrgUserRequest{
				Username: "alice",
				Roles:    admin.OrgUserRolesRequest{OrgRoles: testOrgRolesMultiple},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, diags := clouduserorgassignment.NewOrgUserReq(ctx, tc.plan)
			assert.False(t, diags.HasError(), "expected no diagnostics")
			assert.Equal(t, tc.expected, req)
		})
	}
}

func TestNewAtlasUpdateReq(t *testing.T) {
	ctx := t.Context()

	singleOrgRole := []string{"owner"}

	testCases := map[string]struct {
		plan     *clouduserorgassignment.TFModel
		expected *admin.OrgUserUpdateRequest
	}{
		"null roles": {
			plan: &clouduserorgassignment.TFModel{
				Roles: types.ObjectNull(clouduserorgassignment.RolesObjectAttrTypes),
			},
			expected: &admin.OrgUserUpdateRequest{
				Roles: &admin.OrgUserRolesRequest{OrgRoles: nil},
			},
		},
		"with org roles": {
			plan: &clouduserorgassignment.TFModel{
				Roles: createRolesObject(ctx, singleOrgRole, nil),
			},
			expected: &admin.OrgUserUpdateRequest{
				Roles: &admin.OrgUserRolesRequest{OrgRoles: singleOrgRole},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, diags := clouduserorgassignment.NewAtlasUpdateReq(ctx, tc.plan)
			assert.False(t, diags.HasError(), "expected no diagnostics")
			assert.Equal(t, tc.expected, req)
		})
	}
}

func TestNewTFRoles(t *testing.T) {
	ctx := t.Context()

	testCases := map[string]struct {
		roles          *admin.OrgUserRolesResponse
		expectedObject types.Object
	}{
		"nil roles": {
			roles:          nil,
			expectedObject: types.ObjectNull(clouduserorgassignment.RolesObjectAttrTypes),
		},

		"roles with both roles": {
			roles: &admin.OrgUserRolesResponse{
				OrgRoles: &testOrgRolesMultiple,
				GroupRoleAssignments: &[]admin.GroupRoleAssignment{
					{
						GroupId:    admin.PtrString(testProjectID1),
						GroupRoles: &testProjectRolesSingle,
					},
				},
			},
			expectedObject: createRolesObject(ctx, testOrgRolesMultiple, []clouduserorgassignment.TFRolesProjectRoleAssignmentsModel{
				{
					ProjectId:    types.StringValue(testProjectID1),
					ProjectRoles: types.SetValueMust(types.StringType, []attr.Value{types.StringValue(testProjectRoleOwner)}),
				},
			}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			obj, diags := clouduserorgassignment.NewTFRoles(ctx, tc.roles)
			assert.False(t, diags.HasError(), "unexpected diagnostics")
			assert.Equal(t, tc.expectedObject, obj, "created roles object did not match expected")
		})
	}
}
