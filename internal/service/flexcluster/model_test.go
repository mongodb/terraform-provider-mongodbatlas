//nolint:gocritic
package flexcluster_test

// import (
// )

// type sdkToTFModelTestCase struct {
// 	SDKResp         *admin.FlexCluster
// 	expectedTFModel *flexcluster.TFModel
// }

// func TestFlexClusterSDKToTFModel(t *testing.T) {
// 	testCases := map[string]sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
// 		"Complete SDK response": {
// 			SDKResp:         &admin.FlexCluster{},
// 			expectedTFModel: &flexcluster.TFModel{},
// 		},
// 	}

// 	for testName, tc := range testCases {
// 		t.Run(testName, func(t *testing.T) {
// 			resultModel, diags := flexcluster.NewTFModel(context.Background(), tc.SDKResp)
// 			if diags.HasError() {
// 				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
// 			}
// 			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
// 		})
// 	}
// }

// type tfToSDKModelTestCase struct {
// 	tfModel        *flexcluster.TFModel
// 	expectedSDKReq *admin.FlexCluster
// }

// func TestFlexClusterTFModelToSDK(t *testing.T) {
// 	testCases := map[string]tfToSDKModelTestCase{
// 		"Complete TF state": {
// 			tfModel:        &flexcluster.TFModel{},
// 			expectedSDKReq: &admin.FlexCluster{},
// 		},
// 	}

// 	for testName, tc := range testCases {
// 		t.Run(testName, func(t *testing.T) {
// 			apiReqResult, diags := flexcluster.NewAtlasReq(context.Background(), tc.tfModel)
// 			if diags.HasError() {
// 				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
// 			}
// 			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
// 		})
// 	}
// }
