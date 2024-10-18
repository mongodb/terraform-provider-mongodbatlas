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
	connectionStringsObject, _   = flexcluster.ConvertConnectionStringsToTF(context.Background(), &admin.FlexConnectionStrings20250101{
		Standard:    &standardConnectionString,
		StandardSrv: &standardSrvConnectionString,
	})
	backupSettingsObject, _ = flexcluster.ConvertBackupSettingsToTF(context.Background(), &admin.FlexBackupSettings20250101{
		Enabled: conversion.Pointer(true),
	})
)

type NewTFModelTestCase struct {
	input           *admin.FlexClusterDescription20250101
	expectedTFModel *flexcluster.TFModel
}

type NewAtlasCreateReqTestCase struct {
	input          *flexcluster.TFModel
	expectedSDKReq *admin.FlexClusterDescriptionCreate20250101
}

type NewAtlasUpdateReqTestCase struct {
	input          *flexcluster.TFModel
	expectedSDKReq *admin.FlexClusterDescription20250101
}

func TestNewTFModel(t *testing.T) {
	testCases := map[string]NewTFModelTestCase{
		"Complete TF state": {
			expectedTFModel: &flexcluster.TFModel{
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
				ConnectionStrings:            *connectionStringsObject,
				CreateDate:                   types.StringValue(createDate),
				MongoDbversion:               types.StringValue(mongoDBVersion),
				Name:                         types.StringValue(name),
				ClusterType:                  types.StringValue(clusterType),
				StateName:                    types.StringValue(stateName),
				VersionReleaseSystem:         types.StringValue(versionReleaseSystem),
				BackupSettings:               *backupSettingsObject,
				TerminationProtectionEnabled: types.BoolValue(terminationProtectionEnabled),
			},
			input: &admin.FlexClusterDescription20250101{
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
			expectedTFModel: &flexcluster.TFModel{
				ProjectId:                    types.StringNull(),
				Id:                           types.StringNull(),
				Tags:                         types.Map{},
				ProviderSettings:             flexcluster.TFProviderSettings{},
				ConnectionStrings:            types.ObjectNull(flexcluster.ConnectionStringsType.AttrTypes),
				CreateDate:                   types.StringNull(),
				MongoDbversion:               types.StringNull(),
				Name:                         types.StringNull(),
				ClusterType:                  types.StringNull(),
				StateName:                    types.StringNull(),
				VersionReleaseSystem:         types.StringNull(),
				BackupSettings:               types.ObjectNull(flexcluster.BackupSettingsType.AttrTypes),
				TerminationProtectionEnabled: types.BoolNull(),
			},
			input: &admin.FlexClusterDescription20250101{
				GroupId:                      nil,
				Id:                           nil,
				Tags:                         &[]admin.ResourceTag{},
				ProviderSettings:             admin.FlexProviderSettings20250101{},
				ConnectionStrings:            nil,
				CreateDate:                   nil,
				MongoDBVersion:               nil,
				Name:                         nil,
				ClusterType:                  nil,
				StateName:                    nil,
				VersionReleaseSystem:         nil,
				BackupSettings:               nil,
				TerminationProtectionEnabled: nil,
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			tfModel, diags := flexcluster.NewTFModel(context.Background(), tc.input)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, tfModel, "created TF model did not match expected output")
		})
	}
}

func TestNewAtlasCreateReq(t *testing.T) {
	testCases := map[string]NewAtlasCreateReqTestCase{
		"Complete TF state": {
			input: &flexcluster.TFModel{
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
				ConnectionStrings:            *connectionStringsObject,
				CreateDate:                   types.StringValue(createDate),
				MongoDbversion:               types.StringValue(mongoDBVersion),
				Name:                         types.StringValue(name),
				ClusterType:                  types.StringValue(clusterType),
				StateName:                    types.StringValue(stateName),
				VersionReleaseSystem:         types.StringValue(versionReleaseSystem),
				BackupSettings:               *backupSettingsObject,
				TerminationProtectionEnabled: types.BoolValue(terminationProtectionEnabled),
			},
			expectedSDKReq: &admin.FlexClusterDescriptionCreate20250101{
				Name: name,
				Tags: &[]admin.ResourceTag{
					{
						Key:   key1,
						Value: value1,
					},
				},
				ProviderSettings: admin.FlexProviderSettingsCreate20250101{
					RegionName:          regionName,
					BackingProviderName: backingProviderName,
				},
				TerminationProtectionEnabled: &terminationProtectionEnabled,
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := flexcluster.NewAtlasCreateReq(context.Background(), tc.input)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}

func TestNewAtlasUpdateReq(t *testing.T) {
	testCases := map[string]NewAtlasUpdateReqTestCase{
		"Complete TF state": {
			input: &flexcluster.TFModel{
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
				ConnectionStrings:            *connectionStringsObject,
				CreateDate:                   types.StringValue(createDate),
				MongoDbversion:               types.StringValue(mongoDBVersion),
				Name:                         types.StringValue(name),
				ClusterType:                  types.StringValue(clusterType),
				StateName:                    types.StringValue(stateName),
				VersionReleaseSystem:         types.StringValue(versionReleaseSystem),
				BackupSettings:               *backupSettingsObject,
				TerminationProtectionEnabled: types.BoolValue(terminationProtectionEnabled),
			},
			expectedSDKReq: &admin.FlexClusterDescription20250101{
				Tags: &[]admin.ResourceTag{
					{
						Key:   key1,
						Value: value1,
					},
				},
				TerminationProtectionEnabled: &terminationProtectionEnabled,
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := flexcluster.NewAtlasUpdateReq(context.Background(), tc.input)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}
