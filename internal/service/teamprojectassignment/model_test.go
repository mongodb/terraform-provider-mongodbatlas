package teamprojectassignment_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/teamprojectassignment"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
)

const (
	testProjectID = "project-123"
	testTeamID    = "team-123"
)

var (
	testProjectRoles = []string{"PROJECT_OWNER", "PROJECT_READ_ONLY", "PROJECT_MEMBER"}
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.TeamRole
	expectedTFModel *teamprojectassignment.TFModel
}

func TestTeamProjectAssignmentSDKToTFModel(t *testing.T) {
	ctx := t.Context()

	fullResp := &admin.TeamRole{
		TeamId:    admin.PtrString(testTeamID),
		RoleNames: &testProjectRoles,
	}

	expectedRoles, _ := types.SetValueFrom(ctx, types.StringType, testProjectRoles)
	expectedFullModel := &teamprojectassignment.TFModel{
		ProjectId: types.StringValue(testProjectID),
		TeamId:    types.StringValue(testTeamID),
		RoleNames: expectedRoles,
	}

	fullNilResp := &admin.TeamRole{
		TeamId:    admin.PtrString(""),
		RoleNames: nil,
	}

	expectedNilModel := &teamprojectassignment.TFModel{
		ProjectId: types.StringValue(testProjectID),
		TeamId:    types.StringValue(""),
		RoleNames: types.SetNull(types.StringType),
	}

	testCases := map[string]sdkToTFModelTestCase{
		"Complete SDK response": {
			SDKResp:         fullResp,
			expectedTFModel: expectedFullModel,
		},
		"nil SDK response": {
			SDKResp:         nil,
			expectedTFModel: nil,
		},
		"Empty SDK response": {
			SDKResp:         fullNilResp,
			expectedTFModel: expectedNilModel,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := teamprojectassignment.NewTFModel(t.Context(), tc.SDKResp, testProjectID)
			assert.False(t, diags.HasError(), "expected no diagnostics")
			assert.Equal(t, tc.expectedTFModel, resultModel, "TFModel did not match expected")
		})
	}
}

func TestNewAtlasReq(t *testing.T) {
	ctx := t.Context()

	roles, _ := types.SetValueFrom(ctx, types.StringType, testProjectRoles)
	testCases := map[string]struct {
		plan     *teamprojectassignment.TFModel
		expected *[]admin.TeamRole
	}{
		"Complete TF state": {
			plan: &teamprojectassignment.TFModel{
				ProjectId: types.StringValue(testProjectID),
				TeamId:    types.StringValue(testTeamID),
				RoleNames: roles,
			},
			expected: &[]admin.TeamRole{
				{
					TeamId:    admin.PtrString(testTeamID),
					RoleNames: &testProjectRoles,
				},
			},
		},
		"No roles": {
			plan: &teamprojectassignment.TFModel{
				ProjectId: types.StringValue(testProjectID),
				TeamId:    types.StringValue(testTeamID),
				RoleNames: types.SetNull(types.StringType),
			},
			expected: &[]admin.TeamRole{
				{
					TeamId:    admin.PtrString(testTeamID),
					RoleNames: &[]string{},
				},
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := teamprojectassignment.NewAtlasReq(ctx, tc.plan)
			assert.False(t, diags.HasError(), "expected no diagnostics")

			assert.Len(t, *apiReqResult, len(*tc.expected), "slice lengths don't match")

			for i := range *tc.expected {
				expectedItem := (*tc.expected)[i]
				actualItem := (*apiReqResult)[i]

				assert.Equal(t, *expectedItem.TeamId, *actualItem.TeamId, "TeamId values don't match")

				if expectedItem.RoleNames == nil {
					assert.Nil(t, actualItem.RoleNames, "expected RoleNames to be nil")
				} else {
					assert.Equal(t, *expectedItem.RoleNames, *actualItem.RoleNames, "RoleNames values don't match")
				}
			}
		})
	}
}

func TestNewAtlasUpdateReq(t *testing.T) {
	ctx := t.Context()

	roles, _ := types.SetValueFrom(ctx, types.StringType, testProjectRoles)

	testCases := map[string]struct {
		plan     *teamprojectassignment.TFModel
		expected *admin.TeamRole
	}{
		"Complete TF state": {
			plan: &teamprojectassignment.TFModel{
				ProjectId: types.StringValue(testProjectID),
				TeamId:    types.StringValue(testTeamID),
				RoleNames: roles,
			},
			expected: &admin.TeamRole{
				TeamId:    admin.PtrString(testTeamID),
				RoleNames: &testProjectRoles,
			},
		},
		"No roles": {
			plan: &teamprojectassignment.TFModel{
				ProjectId: types.StringValue(testProjectID),
				TeamId:    types.StringValue(testTeamID),
				RoleNames: types.SetNull(types.StringType),
			},
			expected: &admin.TeamRole{
				TeamId:    admin.PtrString(testTeamID),
				RoleNames: &[]string{},
			},
		},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := teamprojectassignment.NewAtlasUpdateReq(ctx, tc.plan)
			assert.False(t, diags.HasError(), "expected no diagnostics")

			assert.Equal(t, *tc.expected.TeamId, *apiReqResult.TeamId, "TeamId values don't match")

			if tc.expected.RoleNames == nil {
				assert.Nil(t, apiReqResult.RoleNames, "expected RoleNames to be nil")
			} else {
				assert.Equal(t, *tc.expected.RoleNames, *apiReqResult.RoleNames, "RoleNames values don't match")
			}
		})
	}
}
