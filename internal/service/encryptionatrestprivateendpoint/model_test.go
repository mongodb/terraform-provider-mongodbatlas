package encryptionatrestprivateendpoint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/encryptionatrestprivateendpoint"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

const (
	testCloudProvider    = "AZURE"
	testErrMsg           = "error occurred"
	testID               = "6666666999999adsfsgdg"
	testRegionName       = "US_EAST"
	testStatus           = "PENDING_ACCEPTANCE"
	testPEConnectionName = "mongodb-atlas-US_EAST-666666666067bd1e20a8bf14"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.EARPrivateEndpoint
	expectedTFModel *encryptionatrestprivateendpoint.TFEarPrivateEndpointModel
}

func TestEncryptionAtRestPrivateEndpointSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{
		"Complete SDK response": {
			SDKResp: &admin.EARPrivateEndpoint{
				CloudProvider:                 admin.PtrString(testCloudProvider),
				ErrorMessage:                  admin.PtrString(""),
				Id:                            admin.PtrString(testID),
				RegionName:                    admin.PtrString(testRegionName),
				Status:                        admin.PtrString(testStatus),
				PrivateEndpointConnectionName: admin.PtrString(testPEConnectionName),
			},
			expectedTFModel: &encryptionatrestprivateendpoint.TFEarPrivateEndpointModel{
				CloudProvider:                 types.StringValue(testCloudProvider),
				ErrorMessage:                  types.StringNull(),
				ID:                            types.StringValue(testID),
				RegionName:                    types.StringValue(testRegionName),
				Status:                        types.StringValue(testStatus),
				PrivateEndpointConnectionName: types.StringValue(testPEConnectionName),
			},
		},
		"Complete SDK response with error message": {
			SDKResp: &admin.EARPrivateEndpoint{
				CloudProvider:                 admin.PtrString(testCloudProvider),
				ErrorMessage:                  admin.PtrString(testErrMsg),
				Id:                            admin.PtrString(testID),
				RegionName:                    admin.PtrString(testRegionName),
				Status:                        admin.PtrString(testStatus),
				PrivateEndpointConnectionName: admin.PtrString(testPEConnectionName),
			},
			expectedTFModel: &encryptionatrestprivateendpoint.TFEarPrivateEndpointModel{
				CloudProvider:                 types.StringValue(testCloudProvider),
				ErrorMessage:                  types.StringValue(testErrMsg),
				ID:                            types.StringValue(testID),
				RegionName:                    types.StringValue(testRegionName),
				Status:                        types.StringValue(testStatus),
				PrivateEndpointConnectionName: types.StringValue(testPEConnectionName),
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel := encryptionatrestprivateendpoint.NewTFEarPrivateEndpoint(tc.SDKResp)
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

type tfToSDKModelTestCase struct {
	tfModel        *encryptionatrestprivateendpoint.TFEarPrivateEndpointModel
	expectedSDKReq *admin.EARPrivateEndpoint
}

func TestEncryptionAtRestPrivateEndpointTFModelToSDK(t *testing.T) {
	testCases := map[string]tfToSDKModelTestCase{
		"Complete TF state": {
			tfModel: &encryptionatrestprivateendpoint.TFEarPrivateEndpointModel{
				CloudProvider:                 types.StringValue(testCloudProvider),
				ErrorMessage:                  types.StringNull(),
				ID:                            types.StringValue(testID),
				RegionName:                    types.StringValue(testRegionName),
				Status:                        types.StringValue(testStatus),
				PrivateEndpointConnectionName: types.StringValue(testPEConnectionName)},
			expectedSDKReq: &admin.EARPrivateEndpoint{
				CloudProvider:                 admin.PtrString(testCloudProvider),
				ErrorMessage:                  nil,
				Id:                            admin.PtrString(testID),
				RegionName:                    admin.PtrString(testRegionName),
				Status:                        admin.PtrString(testStatus),
				PrivateEndpointConnectionName: admin.PtrString(testPEConnectionName)},
		},
		"Complete TF state with error message": {
			tfModel: &encryptionatrestprivateendpoint.TFEarPrivateEndpointModel{
				CloudProvider:                 types.StringValue(testCloudProvider),
				ErrorMessage:                  types.StringValue(testErrMsg),
				ID:                            types.StringValue(testID),
				RegionName:                    types.StringValue(testRegionName),
				Status:                        types.StringValue(testStatus),
				PrivateEndpointConnectionName: types.StringValue(testPEConnectionName)},
			expectedSDKReq: &admin.EARPrivateEndpoint{
				CloudProvider:                 admin.PtrString(testCloudProvider),
				ErrorMessage:                  admin.PtrString(testErrMsg),
				Id:                            admin.PtrString(testID),
				RegionName:                    admin.PtrString(testRegionName),
				Status:                        admin.PtrString(testStatus),
				PrivateEndpointConnectionName: admin.PtrString(testPEConnectionName)},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult := encryptionatrestprivateendpoint.NewEarPrivateEndpointReq(tc.tfModel)
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}
