package employeeaccessgrant_test

import (
	"context"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/employeeaccessgrant"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.EmployeeAccessGrant
	expectedTFModel *employeeaccessgrant.TFEmployeeAccessModel
}

func TestEmployeeAccessSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
		"Complete SDK response": {
			SDKResp:         &admin.EmployeeAccessGrant{},
			expectedTFModel: &employeeaccessgrant.TFEmployeeAccessModel{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := employeeaccessgrant.NewTFEmployeeAccess(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

type tfToSDKModelTestCase struct {
	tfModel        *employeeaccessgrant.TFEmployeeAccessModel
	expectedSDKReq *admin.EmployeeAccessGrant
}

func TestEmployeeAccessTFModelToSDK(t *testing.T) {
	testCases := map[string]tfToSDKModelTestCase{
		"Complete TF state": {
			tfModel:        &employeeaccessgrant.TFEmployeeAccessModel{},
			expectedSDKReq: &admin.EmployeeAccessGrant{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := employeeaccessgrant.NewEmployeeAccessReq(context.Background(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}
