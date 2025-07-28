package clouduserprojectassignment_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clouduserprojectassignment"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
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

	testProjectRoleOwner  = "PROJECT_OWNER"
	testProjectRoleRead   = "PROJECT_READ_ONLY"
	testProjectRoleMember = "PROJECT_MEMBER"

	testProjectID = "project-123"
	testOrgID     = "org-123"
)

var (
	when                    = time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	testCreatedAt           = when.Format(time.RFC3339)
	testInvitationCreatedAt = when.Add(-24 * time.Hour).Format(time.RFC3339)
	testInvitationExpiresAt = when.Add(24 * time.Hour).Format(time.RFC3339)
	testLastAuth            = when.Add(-2 * time.Hour).Format(time.RFC3339)

	testProjectRoles = []string{testProjectRoleMember, testProjectRoleOwner}
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.GroupUserResponse
	expectedTFModel *clouduserprojectassignment.TFModel
}

func TestCloudUserProjectAssignmentSDKToTFModel(t *testing.T) {
	ctx := t.Context()

	fullResp := &admin.GroupUserResponse{
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
		Roles:               testProjectRoles,
	}

	expectedRoles, _ := types.SetValueFrom(ctx, types.StringType, testProjectRoles)

	expectedFullModel := &clouduserprojectassignment.TFModel{
		UserId:              types.StringValue(testUserID),
		Username:            types.StringValue(testUsername),
		ProjectId:           types.StringValue(testProjectID),
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
		Roles:               expectedRoles,
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
			resultModel, diags := clouduserprojectassignment.NewTFModel(t.Context(), testProjectID, tc.SDKResp)
			assert.False(t, diags.HasError(), "expected no diagnostics")
			assert.Equal(t, tc.expectedTFModel, resultModel, "TFModel did not match expected")
		})
	}
}

func TestNewProjectUserRequest(t *testing.T) {
	ctx := t.Context()
	expectedRoles, _ := types.SetValueFrom(ctx, types.StringType, testProjectRoles)

	testCases := map[string]struct {
		plan     *clouduserprojectassignment.TFModel
		expected *admin.GroupUserRequest
	}{
		"Complete model": {
			plan: &clouduserprojectassignment.TFModel{
				UserId:              types.StringValue(testUserID),
				Username:            types.StringValue(testUsername),
				ProjectId:           types.StringValue(testProjectID),
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
				Roles:               expectedRoles,
			},
			expected: &admin.GroupUserRequest{
				Username: testUsername,
				Roles:    testProjectRoles,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, diags := clouduserprojectassignment.NewProjectUserReq(ctx, tc.plan)
			assert.False(t, diags.HasError(), "expected no diagnostics")
			assert.Equal(t, tc.expected, req)
		})
	}
}
