package streamprivatelinkendpoint_test

import (
	"context"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprivatelinkendpoint"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20241113002/admin"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.StreamsPrivateLinkConnection
	expectedTFModel *streamprivatelinkendpoint.TFModel
}

func TestStreamPrivatelinkEndpointSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
		"Complete SDK response": {
			SDKResp:         &admin.StreamsPrivateLinkConnection{},
			expectedTFModel: &streamprivatelinkendpoint.TFModel{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := streamprivatelinkendpoint.NewTFModel(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

type tfToSDKModelTestCase struct {
	tfModel        *streamprivatelinkendpoint.TFModel
	expectedSDKReq *admin.StreamsPrivateLinkConnection
}

func TestStreamPrivatelinkEndpointTFModelToSDK(t *testing.T) {
	testCases := map[string]tfToSDKModelTestCase{
		"Complete TF state": {
			tfModel:        &streamprivatelinkendpoint.TFModel{},
			expectedSDKReq: &admin.StreamsPrivateLinkConnection{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := streamprivatelinkendpoint.NewAtlasReq(context.Background(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}
