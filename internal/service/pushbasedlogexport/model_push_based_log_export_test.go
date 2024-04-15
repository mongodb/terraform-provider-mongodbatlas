package pushbasedlogexport_test

import (
	"context"
	"testing"

    "github.com/stretchr/testify/assert"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/pushbasedlogexport"
	// "go.mongodb.org/atlas-sdk/v20231115003/admin" use latest version
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.PushBasedLogExport
	expectedTFModel *pushbasedlogexport.TFPushBasedLogExportModel
	name            string
}

func TestPushBasedLogExportSDKToTFModel(t *testing.T) {
	testCases := []sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
		{
			name: "Complete SDK response",
			SDKResp: &admin.PushBasedLogExport{
			},
			expectedTFModel: &pushbasedlogexport.TFPushBasedLogExportModel{
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := pushbasedlogexport.NewTFPushBasedLogExport(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			if !assert.Equal(t, resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}


type tfToSDKModelTestCase struct {
	name           string
	tfModel        *pushbasedlogexport.TFPushBasedLogExportModel
	expectedSDKReq *admin.PushBasedLogExport
}

func TestPushBasedLogExportTFModelToSDK(t *testing.T) {
	testCases := []tfToSDKModelTestCase{
		{
			name: "Complete TF state",
			tfModel: &pushbasedlogexport.TFPushBasedLogExportModel{
			},
			expectedSDKReq: &admin.PushBasedLogExport{
			},
		},
	}

	for _, tc := range testCases {
		apiReqResult, diags := pushbasedlogexport.NewPushBasedLogExportReq(context.Background(), tc.tfModel)
		if diags.HasError() {
			t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
		}
		if !assert.Equal(t, apiReqResult, tc.expectedSDKReq) {
			t.Errorf("created sdk model did not match expected output")
		}
	}
}


