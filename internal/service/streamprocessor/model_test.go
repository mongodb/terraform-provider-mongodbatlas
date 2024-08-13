package streamprocessor_test

import (
	"context"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprocessor"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.StreamsProcessor
	expectedTFModel *streamprocessor.TFStreamProcessorRSModel
}

func TestStreamProcessorSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
		"Complete SDK response": {
			SDKResp:         &admin.StreamsProcessor{},
			expectedTFModel: &streamprocessor.TFStreamProcessorRSModel{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := streamprocessor.NewTFStreamProcessor(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

type tfToSDKModelTestCase struct {
	tfModel        *streamprocessor.TFStreamProcessorRSModel
	expectedSDKReq *admin.StreamsProcessor
}

func TestStreamProcessorTFModelToSDK(t *testing.T) {
	testCases := map[string]tfToSDKModelTestCase{
		"Complete TF state": {
			tfModel:        &streamprocessor.TFStreamProcessorRSModel{},
			expectedSDKReq: &admin.StreamsProcessor{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := streamprocessor.NewStreamProcessorReq(context.Background(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}
