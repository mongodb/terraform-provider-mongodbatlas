package searchdeployment_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/searchdeployment"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.ApiSearchDeploymentResponse
	expectedTFModel *searchdeployment.TFSearchDeploymentRSModel
	name            string
	clusterName     string
}

const (
	dummyDeploymentID = "111111111111111111111111"
	dummyProjectID    = "222222222222222222222222"
	stateName         = "IDLE"
	clusterName       = "Cluster0"
	instanceSize      = "S20_HIGHCPU_NVME"
	nodeCount         = 2
)

func TestSearchDeploymentSDKToTFModel(t *testing.T) {
	testCases := []sdkToTFModelTestCase{
		{
			name:        "Complete SDK response",
			clusterName: clusterName,
			SDKResp: &admin.ApiSearchDeploymentResponse{
				Id:        admin.PtrString(dummyDeploymentID),
				GroupId:   admin.PtrString(dummyProjectID),
				StateName: admin.PtrString(stateName),
				Specs: &[]admin.ApiSearchDeploymentSpec{
					{
						InstanceSize: instanceSize,
						NodeCount:    nodeCount,
					},
				},
			},
			expectedTFModel: &searchdeployment.TFSearchDeploymentRSModel{
				ID:          types.StringValue(dummyDeploymentID),
				ClusterName: types.StringValue(clusterName),
				ProjectID:   types.StringValue(dummyProjectID),
				StateName:   types.StringValue(stateName),
				Specs:       tfSpecsList(t, instanceSize, nodeCount),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := searchdeployment.NewTFSearchDeployment(context.Background(), tc.clusterName, tc.SDKResp, nil, false)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}

func TestSearchDeploymentTFModelToSDK(t *testing.T) {
	testCases := []struct {
		name           string
		tfModel        *searchdeployment.TFSearchDeploymentRSModel
		expectedSDKReq admin.ApiSearchDeploymentRequest
	}{
		{
			name: "Complete TF state",
			tfModel: &searchdeployment.TFSearchDeploymentRSModel{
				ID:          types.StringValue(dummyDeploymentID),
				ClusterName: types.StringValue(clusterName),
				ProjectID:   types.StringValue(dummyProjectID),
				StateName:   types.StringValue(stateName),
				Specs:       tfSpecsList(t, instanceSize, nodeCount),
			},
			expectedSDKReq: admin.ApiSearchDeploymentRequest{
				Specs: []admin.ApiSearchDeploymentSpec{
					{
						InstanceSize: instanceSize,
						NodeCount:    nodeCount,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apiReqResult := searchdeployment.NewSearchDeploymentReq(context.Background(), tc.tfModel)
			if !reflect.DeepEqual(apiReqResult, tc.expectedSDKReq) {
				t.Errorf("created sdk model did not match expected output")
			}
		})
	}
}

func tfSpecsList(t *testing.T, instanceSize string, nodeCount int64) basetypes.ListValue {
	t.Helper()
	tfSpecsList, diags := types.ListValueFrom(context.Background(), searchdeployment.SpecObjectType, []searchdeployment.TFSearchNodeSpecModel{
		{
			InstanceSize: types.StringValue(instanceSize),
			NodeCount:    types.Int64Value(nodeCount),
		},
	})
	if diags.HasError() {
		t.Errorf("failed to create terraform spec lists model: %s", diags.Errors()[0].Summary())
	}
	return tfSpecsList
}
