package employeeaccess_test

import (
	"context"
	"testing"

    "github.com/stretchr/testify/assert"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/employeeaccess"
	// "go.mongodb.org/atlas-sdk/v20231115003/admin" use latest version
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.EmployeeAccess
	expectedTFModel *employeeaccess.TFEmployeeAccessModel
}

func TestEmployeeAccessSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
		"Complete SDK response": {
			SDKResp: &admin.EmployeeAccess{
			},
			expectedTFModel: &employeeaccess.TFEmployeeAccessModel{
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := employeeaccess.NewTFEmployeeAccess(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}


type tfToSDKModelTestCase struct {
	tfModel        *employeeaccess.TFEmployeeAccessModel
	expectedSDKReq *admin.EmployeeAccess
}

func TestEmployeeAccessTFModelToSDK(t *testing.T) {
	testCases := map[string]tfToSDKModelTestCase{
		"Complete TF state": {
			tfModel: &employeeaccess.TFEmployeeAccessModel{
			},
			expectedSDKReq: &admin.EmployeeAccess{
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := employeeaccess.NewEmployeeAccessReq(context.Background(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}


