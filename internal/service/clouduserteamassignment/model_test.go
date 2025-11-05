package clouduserteamassignment_test

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clouduserteamassignment"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
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
	testInviterUsername     = ""

	testOrgRoleOwner      = "ORG_OWNER"
	testOrgRoleMember     = "ORG_MEMBER"
	testProjectRoleOwner  = "PROJECT_OWNER"
	testProjectRoleRead   = "PROJECT_READ_ONLY"
	testProjectRoleMember = "PROJECT_MEMBER"

	testTeamID1    = "team1"
	testTeamID2    = "team2"
	testProjectID1 = "project-123"
	testOrgID      = "org-123"
)

var (
	when                    = time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	testCreatedAt           = when.Format(time.RFC3339)
	testInvitationCreatedAt = when.Add(-24 * time.Hour).Format(time.RFC3339)
	testInvitationExpiresAt = when.Add(24 * time.Hour).Format(time.RFC3339)
	testLastAuth            = when.Add(-2 * time.Hour).Format(time.RFC3339)

	testTeamIDs  = []string{"team1", "team2"}
	testOrgRoles = []string{"owner", "readWrite"}
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.OrgUserResponse
	expectedTFModel *clouduserteamassignment.TFUserTeamAssignmentModel
}

func TestUserTeamAssignmentSDKToTFModel(t *testing.T) {
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
		InvitationCreatedAt: admin.PtrTime(when.Add(-24 * time.Hour)),
		InvitationExpiresAt: admin.PtrTime(when.Add(24 * time.Hour)),
		InviterUsername:     admin.PtrString(testInviterUsername),
		TeamIds:             &testTeamIDs,
		Roles: admin.OrgUserRolesResponse{
			OrgRoles: &testOrgRoles,
		},
	}

	orgRolesSet, _ := types.SetValueFrom(ctx, types.StringType, testOrgRoles)
	expectedRoles, _ := types.ObjectValue(clouduserteamassignment.RolesObjectAttrTypes, map[string]attr.Value{
		"org_roles":                orgRolesSet,
		"project_role_assignments": types.ListNull(clouduserteamassignment.ProjectRoleAssignmentsAttrType),
	})
	expectedTeams, _ := types.SetValueFrom(ctx, types.StringType, testTeamIDs)
	expectedFullModel := &clouduserteamassignment.TFUserTeamAssignmentModel{
		UserId:              types.StringValue(testUserID),
		Username:            types.StringValue(testUsername),
		FirstName:           types.StringValue(testFirstName),
		LastName:            types.StringValue(testLastName),
		Country:             types.StringValue(testCountry),
		MobileNumber:        types.StringValue(testMobile),
		OrgMembershipStatus: types.StringValue(testOrgMembershipStatus),
		CreatedAt:           types.StringValue(testCreatedAt),
		LastAuth:            types.StringValue(testLastAuth),
		InvitationCreatedAt: types.StringValue(testInvitationCreatedAt),
		InvitationExpiresAt: types.StringValue(testInvitationExpiresAt),
		InviterUsername:     types.StringValue(testInviterUsername),
		OrgId:               types.StringNull(),
		TeamId:              types.StringNull(),
		Roles:               expectedRoles,
		TeamIds:             expectedTeams,
	}

	testCases := map[string]sdkToTFModelTestCase{
		"nil SDK response": {
			SDKResp:         nil,
			expectedTFModel: nil,
		},
		"Complete SDK response": {
			SDKResp:         fullResp,
			expectedTFModel: expectedFullModel,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := clouduserteamassignment.NewTFUserTeamAssignmentModel(t.Context(), tc.SDKResp)
			assert.False(t, diags.HasError(), "expected no diagnostics")
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

func createRolesObject(ctx context.Context, orgRoles []string, projectAssignments []clouduserteamassignment.TFProjectRoleAssignmentsModel) types.Object {
	orgRolesSet, _ := types.SetValueFrom(ctx, types.StringType, orgRoles)

	var projectRoleAssignmentsList types.List
	if len(projectAssignments) == 0 {
		projectRoleAssignmentsList = types.ListNull(clouduserteamassignment.ProjectRoleAssignmentsAttrType.ElemType.(types.ObjectType))
	} else {
		projectRoleAssignmentsList, _ = types.ListValueFrom(ctx, clouduserteamassignment.ProjectRoleAssignmentsAttrType.ElemType.(types.ObjectType), projectAssignments)
	}

	obj, _ := types.ObjectValue(
		clouduserteamassignment.RolesObjectAttrTypes,
		map[string]attr.Value{
			"org_roles":                orgRolesSet,
			"project_role_assignments": projectRoleAssignmentsList,
		},
	)
	return obj
}

func TestNewUserTeamAssignmentReq(t *testing.T) {
	ctx := t.Context()
	projectAssignment := clouduserteamassignment.TFProjectRoleAssignmentsModel{
		ProjectId:    types.StringValue(testProjectID1),
		ProjectRoles: types.SetValueMust(types.StringType, []attr.Value{types.StringValue(testProjectRoleOwner)}),
	}
	teams, _ := types.SetValueFrom(ctx, types.StringType, testTeamIDs)
	testCases := map[string]struct {
		plan     *clouduserteamassignment.TFUserTeamAssignmentModel
		expected *admin.AddOrRemoveUserFromTeam
	}{
		"Complete model": {
			plan: &clouduserteamassignment.TFUserTeamAssignmentModel{
				OrgId:               types.StringValue(testOrgID),
				TeamId:              types.StringValue(testTeamID1),
				UserId:              types.StringValue(testUserID),
				Username:            types.StringValue(testUsername),
				OrgMembershipStatus: types.StringValue(testOrgMembershipStatus),
				Roles: createRolesObject(ctx, testOrgRoles, []clouduserteamassignment.TFProjectRoleAssignmentsModel{
					projectAssignment,
				}),
				TeamIds:             teams,
				InvitationCreatedAt: types.StringValue(testInvitationCreatedAt),
				InvitationExpiresAt: types.StringValue(testInvitationExpiresAt),
				InviterUsername:     types.StringValue(testInviterUsername),
				Country:             types.StringValue(testCountry),
				FirstName:           types.StringValue(testFirstName),
				LastName:            types.StringValue(testLastName),
				CreatedAt:           types.StringValue(testCreatedAt),
				LastAuth:            types.StringValue(testLastAuth),
				MobileNumber:        types.StringValue(testMobile),
			},
			expected: &admin.AddOrRemoveUserFromTeam{
				Id: testUserID,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, diags := clouduserteamassignment.NewUserTeamAssignmentReq(ctx, tc.plan)
			assert.False(t, diags.HasError(), "expected no diagnostics")
			assert.Equal(t, tc.expected, req)
		})
	}
}
