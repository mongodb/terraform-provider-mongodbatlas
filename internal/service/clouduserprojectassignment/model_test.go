package clouduserprojectassignment_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clouduserprojectassignment"
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
	testOrgMembershipStatus = "ACTIVE"
	testInviterUsername     = ""

	testProjectRoleOwner  = "PROJECT_OWNER"
	testProjectRoleMember = "PROJECT_MEMBER"

	testProjectID = "project-123"
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
		"Empty SDK response": {
			SDKResp: &admin.GroupUserResponse{
				Id:                  "",
				Username:            "",
				FirstName:           nil,
				LastName:            nil,
				Country:             nil,
				MobileNumber:        nil,
				OrgMembershipStatus: "",
				CreatedAt:           nil,
				LastAuth:            nil,
				InvitationCreatedAt: nil,
				InvitationExpiresAt: nil,
				InviterUsername:     nil,
				Roles:               nil,
			},
			expectedTFModel: &clouduserprojectassignment.TFModel{
				UserId:              types.StringValue(""),
				Username:            types.StringValue(""),
				ProjectId:           types.StringValue(testProjectID),
				FirstName:           types.StringNull(),
				LastName:            types.StringNull(),
				Country:             types.StringNull(),
				MobileNumber:        types.StringNull(),
				OrgMembershipStatus: types.StringValue(""),
				CreatedAt:           types.StringNull(),
				LastAuth:            types.StringNull(),
				InvitationCreatedAt: types.StringNull(),
				InvitationExpiresAt: types.StringNull(),
				InviterUsername:     types.StringNull(),
				Roles:               types.SetNull(types.StringType),
			},
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
		"Nil model": {
			plan: &clouduserprojectassignment.TFModel{
				Username: types.StringNull(),
				Roles:    types.SetNull(types.StringType),
			},
			expected: &admin.GroupUserRequest{
				Username: "",
				Roles:    []string{},
			},
		},
		"Empty model": {
			plan: &clouduserprojectassignment.TFModel{
				Username: types.StringValue(""),
				Roles:    types.SetValueMust(types.StringType, []attr.Value{}),
			},
			expected: &admin.GroupUserRequest{
				Username: "",
				Roles:    []string{},
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

func TestNewAtlasUpdateReq(t *testing.T) {
	ctx := t.Context()

	type args struct {
		stateRoles []string
		planRoles  []string
	}
	tests := []struct {
		name            string
		args            args
		wantAddRoles    []string
		wantRemoveRoles []string
	}{
		{
			name: "add and remove roles",
			args: args{
				stateRoles: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"},
				planRoles:  []string{"GROUP_OWNER", "GROUP_DATA_ACCESS_READ_ONLY"},
			},
			wantAddRoles:    []string{"GROUP_OWNER"},
			wantRemoveRoles: []string{"GROUP_READ_ONLY"},
		},
		{
			name: "no changes",
			args: args{
				stateRoles: []string{"GROUP_OWNER"},
				planRoles:  []string{"GROUP_OWNER"},
			},
			wantAddRoles:    []string{},
			wantRemoveRoles: []string{},
		},
		{
			name: "all roles removed",
			args: args{
				stateRoles: []string{"GROUP_OWNER"},
				planRoles:  []string{},
			},
			wantAddRoles:    []string{},
			wantRemoveRoles: []string{"GROUP_OWNER"},
		},
		{
			name: "all roles added",
			args: args{
				stateRoles: []string{},
				planRoles:  []string{"GROUP_OWNER"},
			},
			wantAddRoles:    []string{"GROUP_OWNER"},
			wantRemoveRoles: []string{},
		},
		{
			name: "nil roles",
			args: args{
				stateRoles: nil,
				planRoles:  []string{},
			},
			wantAddRoles:    []string{},
			wantRemoveRoles: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planRoles, _ := types.SetValueFrom(ctx, types.StringType, tt.args.planRoles)

			state := tt.args.stateRoles
			plan := &clouduserprojectassignment.TFModel{Roles: planRoles}

			addReqs, removeReqs, diags := clouduserprojectassignment.NewAtlasUpdateReq(ctx, plan, state)
			assert.False(t, diags.HasError(), "expected no diagnostics")

			var gotAddRoles, gotRemoveRoles []string
			for _, r := range addReqs {
				gotAddRoles = append(gotAddRoles, r.GroupRole)
			}
			for _, r := range removeReqs {
				gotRemoveRoles = append(gotRemoveRoles, r.GroupRole)
			}

			assert.ElementsMatch(t, tt.wantAddRoles, gotAddRoles, "add roles mismatch")
			assert.ElementsMatch(t, tt.wantRemoveRoles, gotRemoveRoles, "remove roles mismatch")
		})
	}
}
