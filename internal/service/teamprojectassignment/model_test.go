package teamprojectassignment_test

import (
	"context"
	"testing"

    "github.com/stretchr/testify/assert"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/teamprojectassignment"
	// "go.mongodb.org/atlas-sdk/v20231115003/admin" use latest version
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.TeamProjectAssignment
	expectedTFModel *teamprojectassignment.TFModel
}

func TestTeamProjectAssignmentSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
		"Complete SDK response": {
			SDKResp: &admin.TeamProjectAssignment{
			},
			expectedTFModel: &teamprojectassignment.TFModel{
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := teamprojectassignment.NewTFModel(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}


type tfToSDKModelTestCase struct {
	tfModel        *teamprojectassignment.TFModel
	expectedSDKReq *admin.TeamProjectAssignment
}

func TestTeamProjectAssignmentTFModelToSDK(t *testing.T) {
	testCases := map[string]tfToSDKModelTestCase{
		"Complete TF state": {
			tfModel: &teamprojectassignment.TFModel{
			},
			expectedSDKReq: &admin.TeamProjectAssignment{
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := teamprojectassignment.NewAtlasReq(context.Background(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}


