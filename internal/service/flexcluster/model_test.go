//nolint:gocritic
package flexcluster_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240805004/admin"
)

type sdkToTFModelTestCase struct {
	input           *admin.FlexClusterDescription20250101
	expectedTFModel *flexcluster.TFModel
}

func TestFlexClusterSDKToTFModel(t *testing.T) {
	time, _ := admin.StringToTime("2021-08-17T17:00:00Z")
	testCases := map[string]sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
		"Complete SDK response": {
			input: &admin.FlexClusterDescription20250101{
				ClusterType:                  conversion.StringPtr("REPLICASET"),
				CreateDate:                   conversion.Pointer(time),
				GroupId:                      conversion.StringPtr("5f3b6d1b0b8e4f0e7b8f6b1b"),
				Id:                           conversion.StringPtr("5f3b6d1b0b8e4f0e7b8f6b1c"),
				MongoDBVersion:               conversion.StringPtr("8.0"),
				Name:                         conversion.StringPtr("myCluster"),
				StateName:                    conversion.StringPtr("IDLE"),
				TerminationProtectionEnabled: conversion.Pointer(true),
				VersionReleaseSystem:         conversion.StringPtr("LTS"),
			},
			expectedTFModel: &flexcluster.TFModel{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := flexcluster.NewTFModel(context.Background(), tc.input)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

type tfModelTestCase struct {
	expectedtfModel *flexcluster.TFModel
	SDKReq          *admin.FlexClusterDescription20250101
}

func TestNewTFModel(t *testing.T) {
	testCases := map[string]tfModelTestCase{
		"Complete TF state": {
			expectedtfModel: &flexcluster.TFModel{
				ProjectId:                    types.StringNull(),
				Id:                           types.StringNull(),
				Tags:                         types.Map{},
				ProviderSettings:             flexcluster.TFProviderSettings{},
				ConnectionStrings:            flexcluster.TFConnectionStrings{},
				CreateDate:                   types.StringNull(),
				MongoDbversion:               types.StringNull(),
				Name:                         types.StringNull(),
				ClusterType:                  types.StringNull(),
				StateName:                    types.StringNull(),
				VersionReleaseSystem:         types.StringNull(),
				BackupSettings:               flexcluster.TFBackupSettings{},
				TerminationProtectionEnabled: types.BoolNull(),
			},
			SDKReq: &admin.FlexClusterDescription20250101{
				GroupId:                      nil,
				Id:                           nil,
				Tags:                         &[]admin.ResourceTag{},
				ProviderSettings:             admin.FlexProviderSettings20250101{},
				ConnectionStrings:            &admin.FlexConnectionStrings20250101{},
				CreateDate:                   nil,
				MongoDBVersion:               nil,
				Name:                         nil,
				ClusterType:                  nil,
				StateName:                    nil,
				VersionReleaseSystem:         nil,
				BackupSettings:               &admin.FlexBackupSettings20250101{},
				TerminationProtectionEnabled: nil,
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			tfModel, diags := flexcluster.NewTFModel(context.Background(), tc.SDKReq)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedtfModel, tfModel, "created TF model did not match expected output")
		})
	}
}

func TestNewAtlasCreateReq(t *testing.T) {
	// testCases := map[string]tfToSDKModelTestCase{
	// 	"Complete TF state": {
	// 		tfModel:        &flexcluster.TFModel{},
	// 		expectedSDKReq: &admin.FlexClusterDescriptionCreate20250101{},
	// 	},
	// }

	// for testName, tc := range testCases {
	// 	t.Run(testName, func(t *testing.T) {
	// 		apiReqResult, diags := flexcluster.NewAtlasCreateReq(context.Background(), tc.tfModel)
	// 		if diags.HasError() {
	// 			t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
	// 		}
	// 		assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
	// 	})
	// }
}
