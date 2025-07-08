package clouduserorgassignment_test

import (
	"context"
	"testing"

    "github.com/stretchr/testify/assert"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clouduserorgassignment"
	// "go.mongodb.org/atlas-sdk/v20231115003/admin" use latest version
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.CloudUserOrgAssignment
	expectedTFModel *clouduserorgassignment.TFModel
}

func TestCloudUserOrgAssignmentSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
		"Complete SDK response": {
			SDKResp: &admin.CloudUserOrgAssignment{
			},
			expectedTFModel: &clouduserorgassignment.TFModel{
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := clouduserorgassignment.NewTFModel(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}


type tfToSDKModelTestCase struct {
	tfModel        *clouduserorgassignment.TFModel
	expectedSDKReq *admin.CloudUserOrgAssignment
}

func TestCloudUserOrgAssignmentTFModelToSDK(t *testing.T) {
	testCases := map[string]tfToSDKModelTestCase{
		"Complete TF state": {
			tfModel: &clouduserorgassignment.TFModel{
			},
			expectedSDKReq: &admin.CloudUserOrgAssignment{
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := clouduserorgassignment.NewAtlasReq(context.Background(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}


