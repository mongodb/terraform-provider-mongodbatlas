//nolint:gocritic
package flexcluster_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

var (
	projectID                    = "projectId"
	id                           = "id"
	createDate                   = "2021-08-17T17:00:00Z"
	mongoDBVersion               = "8.0"
	name                         = "myCluster"
	clusterType                  = "REPLICASET"
	stateName                    = "IDLE"
	versionReleaseSystem         = "LTS"
	terminationProtectionEnabled = true
	createDateTime, _            = conversion.StringToTime(createDate)
	providerName                 = "AWS"
	regionName                   = "us-east-1"
	backingProviderName          = "AWS"
	diskSizeGb                   = 100.0
	standardConnectionString     = "mongodb://localhost:27017"
	standardSrvConnectionString  = "mongodb+srv://localhost:27017"
	key1                         = "key1"
	value1                       = "value1"
)

func TestNewTFModel(t *testing.T) {
	testCases := map[string]tfModelTestCase{
		"Complete TF state": {
			expectedtfModel: &flexcluster.TFModel{
				ProjectId: types.StringValue(projectID),
				Id:        types.StringValue(id),
				Tags: types.MapValueMust(types.StringType, map[string]attr.Value{
					key1: types.StringValue(value1),
				}),
				ProviderSettings: flexcluster.TFProviderSettings{
					ProviderName:        types.StringValue(providerName),
					RegionName:          types.StringValue(regionName),
					BackingProviderName: types.StringValue(backingProviderName),
					DiskSizeGb:          types.Float64Value(diskSizeGb),
				},
				ConnectionStrings: flexcluster.TFConnectionStrings{
					Standard:    types.StringValue(standardConnectionString),
					StandardSrv: types.StringValue(standardSrvConnectionString),
				},
				CreateDate:           types.StringValue(createDate),
				MongoDbversion:       types.StringValue(mongoDBVersion),
				Name:                 types.StringValue(name),
				ClusterType:          types.StringValue(clusterType),
				StateName:            types.StringValue(stateName),
				VersionReleaseSystem: types.StringValue(versionReleaseSystem),
				BackupSettings: flexcluster.TFBackupSettings{
					Enabled: types.BoolValue(true),
				},
				TerminationProtectionEnabled: types.BoolValue(terminationProtectionEnabled),
			},
			SDKReq: &admin.FlexClusterDescription20250101{
				GroupId: &projectID,
				Id:      &id,
				Tags: &[]admin.ResourceTag{
					{
						Key:   key1,
						Value: value1,
					},
				},
				ProviderSettings: admin.FlexProviderSettings20250101{
					ProviderName:        &providerName,
					RegionName:          &regionName,
					BackingProviderName: &backingProviderName,
					DiskSizeGB:          &diskSizeGb,
				},
				ConnectionStrings: &admin.FlexConnectionStrings20250101{
					Standard:    &standardConnectionString,
					StandardSrv: &standardSrvConnectionString,
				},
				CreateDate:           &createDateTime,
				MongoDBVersion:       &mongoDBVersion,
				Name:                 &name,
				ClusterType:          &clusterType,
				StateName:            &stateName,
				VersionReleaseSystem: &versionReleaseSystem,
				BackupSettings: &admin.FlexBackupSettings20250101{
					Enabled: conversion.Pointer(true),
				},
				TerminationProtectionEnabled: &terminationProtectionEnabled,
			},
		},
		"Nil values": {
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
