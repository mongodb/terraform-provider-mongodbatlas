package flexcluster_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
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
	connectionStringsObject, _   = flexcluster.ConvertConnectionStringsToTF(context.Background(), &admin.FlexConnectionStrings20241113{
		Standard:    &standardConnectionString,
		StandardSrv: &standardSrvConnectionString,
	})
	backupSettingsObject, _ = flexcluster.ConvertBackupSettingsToTF(context.Background(), &admin.FlexBackupSettings20241113{
		Enabled: conversion.Pointer(true),
	})
	providerSettingsObject, _ = flexcluster.ConvertProviderSettingsToTF(context.Background(), admin.FlexProviderSettings20241113{
		ProviderName:        &providerName,
		RegionName:          &regionName,
		BackingProviderName: &backingProviderName,
		DiskSizeGB:          &diskSizeGb,
	})
)

type NewTFModelTestCase struct {
	expectedTFModel *flexcluster.TFModel
	input           *admin.FlexClusterDescription20241113
}

type NewTFModelDSPTestCase struct {
	expectedTFModelDSP *flexcluster.TFModelDSP
	input              []admin.FlexClusterDescription20241113
}

type NewAtlasCreateReqTestCase struct {
	input          *flexcluster.TFModel
	expectedSDKReq *admin.FlexClusterDescriptionCreate20241113
}

type NewAtlasUpdateReqTestCase struct {
	input          *flexcluster.TFModel
	expectedSDKReq *admin.FlexClusterDescriptionUpdate20241113
}

func TestNewTFModel(t *testing.T) {
	providerSettingsTF := &flexcluster.TFProviderSettings{
		ProviderName:        types.StringNull(),
		RegionName:          types.StringNull(),
		BackingProviderName: types.StringNull(),
		DiskSizeGb:          types.Float64Null(),
	}
	nilProviderSettingsObject, _ := types.ObjectValueFrom(context.Background(), flexcluster.ProviderSettingsType.AttributeTypes(), providerSettingsTF)
	testCases := map[string]NewTFModelTestCase{
		"Complete TF state": {
			expectedTFModel: &flexcluster.TFModel{
				ProjectId: types.StringValue(projectID),
				Id:        types.StringValue(id),
				Tags: types.MapValueMust(types.StringType, map[string]attr.Value{
					key1: types.StringValue(value1),
				}),
				ProviderSettings:             *providerSettingsObject,
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
			input: &admin.FlexClusterDescription20241113{
				GroupId: &projectID,
				Id:      &id,
				Tags: &[]admin.ResourceTag{
					{
						Key:   key1,
						Value: value1,
					},
				},
				ProviderSettings: admin.FlexProviderSettings20241113{
					ProviderName:        &providerName,
					RegionName:          &regionName,
					BackingProviderName: &backingProviderName,
					DiskSizeGB:          &diskSizeGb,
				},
				ConnectionStrings: &admin.FlexConnectionStrings20241113{
					Standard:    &standardConnectionString,
					StandardSrv: &standardSrvConnectionString,
				},
				CreateDate:           &createDateTime,
				MongoDBVersion:       &mongoDBVersion,
				Name:                 &name,
				ClusterType:          &clusterType,
				StateName:            &stateName,
				VersionReleaseSystem: &versionReleaseSystem,
				BackupSettings: &admin.FlexBackupSettings20241113{
					Enabled: conversion.Pointer(true),
				},
				TerminationProtectionEnabled: &terminationProtectionEnabled,
			},
		},
		"Nil values": {
			expectedTFModel: &flexcluster.TFModel{
				ProjectId:                    types.StringNull(),
				Id:                           types.StringNull(),
				Tags:                         types.MapValueMust(types.StringType, map[string]attr.Value{}),
				ProviderSettings:             nilProviderSettingsObject,
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
			input: &admin.FlexClusterDescription20241113{
				GroupId:                      nil,
				Id:                           nil,
				Tags:                         &[]admin.ResourceTag{},
				ProviderSettings:             admin.FlexProviderSettings20241113{},
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

func TestNewTFModelDSP(t *testing.T) {
	testCases := map[string]NewTFModelDSPTestCase{
		"Complete TF state": {
			expectedTFModelDSP: &flexcluster.TFModelDSP{
				ProjectId: types.StringValue(projectID),
				Results: []flexcluster.TFModel{
					{
						ProjectId: types.StringValue(projectID),
						Id:        types.StringValue(id),
						Tags: types.MapValueMust(types.StringType, map[string]attr.Value{
							key1: types.StringValue(value1),
						}),
						ProviderSettings:             *providerSettingsObject,
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
					{
						ProjectId: types.StringValue(projectID),
						Id:        types.StringValue("id-2"),
						Tags: types.MapValueMust(types.StringType, map[string]attr.Value{
							key1: types.StringValue(value1),
						}),
						ProviderSettings:             *providerSettingsObject,
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
				},
			},
			input: []admin.FlexClusterDescription20241113{
				{
					GroupId: &projectID,
					Id:      &id,
					Tags: &[]admin.ResourceTag{
						{
							Key:   key1,
							Value: value1,
						},
					},
					ProviderSettings: admin.FlexProviderSettings20241113{
						ProviderName:        &providerName,
						RegionName:          &regionName,
						BackingProviderName: &backingProviderName,
						DiskSizeGB:          &diskSizeGb,
					},
					ConnectionStrings: &admin.FlexConnectionStrings20241113{
						Standard:    &standardConnectionString,
						StandardSrv: &standardSrvConnectionString,
					},
					CreateDate:           &createDateTime,
					MongoDBVersion:       &mongoDBVersion,
					Name:                 &name,
					ClusterType:          &clusterType,
					StateName:            &stateName,
					VersionReleaseSystem: &versionReleaseSystem,
					BackupSettings: &admin.FlexBackupSettings20241113{
						Enabled: conversion.Pointer(true),
					},
					TerminationProtectionEnabled: &terminationProtectionEnabled,
				},
				{
					GroupId: &projectID,
					Id:      conversion.StringPtr("id-2"),
					Tags: &[]admin.ResourceTag{
						{
							Key:   key1,
							Value: value1,
						},
					},
					ProviderSettings: admin.FlexProviderSettings20241113{
						ProviderName:        &providerName,
						RegionName:          &regionName,
						BackingProviderName: &backingProviderName,
						DiskSizeGB:          &diskSizeGb,
					},
					ConnectionStrings: &admin.FlexConnectionStrings20241113{
						Standard:    &standardConnectionString,
						StandardSrv: &standardSrvConnectionString,
					},
					CreateDate:           &createDateTime,
					MongoDBVersion:       &mongoDBVersion,
					Name:                 &name,
					ClusterType:          &clusterType,
					StateName:            &stateName,
					VersionReleaseSystem: &versionReleaseSystem,
					BackupSettings: &admin.FlexBackupSettings20241113{
						Enabled: conversion.Pointer(true),
					},
					TerminationProtectionEnabled: &terminationProtectionEnabled,
				},
			},
		},
		"No Flex Clusters": {
			expectedTFModelDSP: &flexcluster.TFModelDSP{
				ProjectId: types.StringValue(projectID),
				Results:   []flexcluster.TFModel{},
			},
			input: []admin.FlexClusterDescription20241113{},
		},
	}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			tfModelDSP, diags := flexcluster.NewTFModelDSP(context.Background(), projectID, tc.input)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModelDSP, tfModelDSP, "created TF model DSP did not match expected output")
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
				ProviderSettings:             *providerSettingsObject,
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
			expectedSDKReq: &admin.FlexClusterDescriptionCreate20241113{
				Name: name,
				Tags: &[]admin.ResourceTag{
					{
						Key:   key1,
						Value: value1,
					},
				},
				ProviderSettings: admin.FlexProviderSettingsCreate20241113{
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
				ProviderSettings:             *providerSettingsObject,
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
			expectedSDKReq: &admin.FlexClusterDescriptionUpdate20241113{
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
