package advancedcluster_test

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	mockadmin20240530 "go.mongodb.org/atlas-sdk/v20240530005/mockadmin"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

const (
	resourceName           = "mongodbatlas_advanced_cluster.test"
	dataSourceName         = "data.mongodbatlas_advanced_cluster.test"
	dataSourcePluralName   = "data.mongodbatlas_advanced_clusters.test"
	dataSourcesTFOldSchema = `
	data "mongodbatlas_advanced_cluster" "test" {
		project_id = mongodbatlas_advanced_cluster.test.project_id
		name 	     = mongodbatlas_advanced_cluster.test.name
		depends_on = [mongodbatlas_advanced_cluster.test]
	}

	data "mongodbatlas_advanced_clusters" "test" {
		project_id = mongodbatlas_advanced_cluster.test.project_id
		depends_on = [mongodbatlas_advanced_cluster.test]
	}`
	dataSourcesTFNewSchema = `
	data "mongodbatlas_advanced_cluster" "test" {
		project_id = mongodbatlas_advanced_cluster.test.project_id
		name 	     = mongodbatlas_advanced_cluster.test.name
		use_replication_spec_per_shard = true
		depends_on = [mongodbatlas_advanced_cluster.test]
	}
			
	data "mongodbatlas_advanced_clusters" "test" {
		use_replication_spec_per_shard = true
		project_id = mongodbatlas_advanced_cluster.test.project_id
		depends_on = [mongodbatlas_advanced_cluster.test]
	}`
	freeInstanceSize   = "M0"
	sharedInstanceSize = "M2"
)

var (
	configServerManagementModeFixedToDedicated = "FIXED_TO_DEDICATED"
	configServerManagementModeAtlasManaged     = "ATLAS_MANAGED"
	mockConfig                                 = unit.MockConfigAdvancedClusterTPF
)

func TestGetReplicationSpecAttributesFromOldAPI(t *testing.T) {
	var (
		projectID   = "11111"
		clusterName = "testCluster"
		ID          = "111111"
		numShard    = 2
		zoneName    = "ZoneName managed by Terraform"
	)

	testCases := map[string]struct {
		mockCluster    *admin20240530.AdvancedClusterDescription
		mockResponse   *http.Response
		mockError      error
		expectedResult map[string]advancedcluster.OldShardConfigMeta
		expectedError  error
	}{
		"Error in the API call": {
			mockCluster:    &admin20240530.AdvancedClusterDescription{},
			mockResponse:   &http.Response{StatusCode: 400},
			mockError:      errGeneric,
			expectedError:  errGeneric,
			expectedResult: nil,
		},
		"Successful": {
			mockCluster: &admin20240530.AdvancedClusterDescription{
				ReplicationSpecs: &[]admin20240530.ReplicationSpec{
					{
						NumShards: &numShard,
						Id:        &ID,
						ZoneName:  &zoneName,
					},
				},
			},
			mockResponse:  &http.Response{},
			mockError:     nil,
			expectedError: nil,
			expectedResult: map[string]advancedcluster.OldShardConfigMeta{
				zoneName: {ID: ID, NumShard: numShard},
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			testObject := mockadmin20240530.NewClustersApi(t)

			testObject.EXPECT().GetCluster(mock.Anything, mock.Anything, mock.Anything).Return(admin20240530.GetClusterApiRequest{ApiService: testObject}).Once()
			testObject.EXPECT().GetClusterExecute(mock.Anything).Return(tc.mockCluster, tc.mockResponse, tc.mockError).Once()

			result, err := advancedcluster.GetReplicationSpecAttributesFromOldAPI(t.Context(), projectID, clusterName, testObject)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func testAccAdvancedClusterFlexUpgrade(t *testing.T, instanceSize string, includeDedicated bool) resource.TestCase {
	t.Helper()
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 1)
	defaultZoneName := "Zone 1" // Uses backend default as in existing tests

	steps := []resource.TestStep{
		{
			Config: configTenant(t, true, projectID, clusterName, defaultZoneName, instanceSize),
			Check:  checkTenant(true, projectID, clusterName),
		},
		{
			Config: configFlexCluster(t, projectID, clusterName, "AWS", "US_EAST_1", defaultZoneName, false),
			Check:  checkFlexClusterConfig(projectID, clusterName, "AWS", "US_EAST_1", false),
		},
	}
	if includeDedicated {
		steps = append(steps, resource.TestStep{
			Config: acc.ConvertAdvancedClusterToPreviewProviderV2(t, true, acc.ConfigBasicDedicated(projectID, clusterName, defaultZoneName)),
			Check:  checksBasicDedicated(projectID, clusterName),
		})
	}

	return resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps:                    steps,
	}
}

func TestAccAdvancedCluster_basicTenant_flexUpgrade_dedicatedUpgrade(t *testing.T) {
	resource.Test(t, testAccAdvancedClusterFlexUpgrade(t, freeInstanceSize, true))
}

func TestAccAdvancedCluster_sharedTier_flexUpgrade(t *testing.T) {
	resource.Test(t, testAccAdvancedClusterFlexUpgrade(t, sharedInstanceSize, false))
}
func TestAccMockableAdvancedCluster_tenantUpgrade(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 1)
		defaultZoneName        = "Zone 1" // Uses backend default to avoid non-empty plan, see CLOUDP-294339
	)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, &resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: acc.ConvertAdvancedClusterToPreviewProviderV2(t, true, configTenant(t, true, projectID, clusterName, defaultZoneName, freeInstanceSize)),
				Check:  checkTenant(true, projectID, clusterName),
			},
			{
				Config: acc.ConvertAdvancedClusterToPreviewProviderV2(t, true, acc.ConfigBasicDedicated(projectID, clusterName, defaultZoneName)),
				Check:  checksBasicDedicated(projectID, clusterName),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_replicaSetAWSProvider(t *testing.T) {
	resource.ParallelTest(t, replicaSetAWSProviderTestCase(t, true))
}

func replicaSetAWSProviderTestCase(t *testing.T, usePreviewProvider bool) resource.TestCase {
	t.Helper()
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 6)
	)

	return resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAWSProvider(t, usePreviewProvider, ReplicaSetAWSConfig{
					ProjectID:          projectID,
					ClusterName:        clusterName,
					ClusterType:        "REPLICASET",
					DiskSizeGB:         60,
					NodeCountElectable: 3,
					WithAnalyticsSpecs: true,
				}),
				Check: checkReplicaSetAWSProvider(usePreviewProvider, projectID, clusterName, 60, 3, true, true),
			},
			// empty plan when analytics block is removed
			acc.TestStepCheckEmptyPlan(configAWSProvider(t, usePreviewProvider, ReplicaSetAWSConfig{
				ProjectID:          projectID,
				ClusterName:        clusterName,
				ClusterType:        "REPLICASET",
				DiskSizeGB:         60,
				NodeCountElectable: 3,
				WithAnalyticsSpecs: false,
			})),
			{
				Config: configAWSProvider(t, usePreviewProvider, ReplicaSetAWSConfig{
					ProjectID:          projectID,
					ClusterName:        clusterName,
					ClusterType:        "REPLICASET",
					DiskSizeGB:         50,
					NodeCountElectable: 5,
					WithAnalyticsSpecs: false, // other update made after removed analytics block, computed value is expected to be the same
				}),
				Check: checkReplicaSetAWSProvider(usePreviewProvider, projectID, clusterName, 50, 5, true, true),
			},
			{ // testing transition from replica set to sharded cluster
				Config: configAWSProvider(t, usePreviewProvider, ReplicaSetAWSConfig{
					ProjectID:          projectID,
					ClusterName:        clusterName,
					ClusterType:        "SHARDED",
					DiskSizeGB:         50,
					NodeCountElectable: 5,
					WithAnalyticsSpecs: false,
				}),
				Check: checkReplicaSetAWSProvider(usePreviewProvider, projectID, clusterName, 50, 5, true, true),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs", "retain_backups_enabled"),
		},
	}
}

func TestAccClusterAdvancedCluster_replicaSetMultiCloud(t *testing.T) {
	resource.ParallelTest(t, replicaSetMultiCloudTestCase(t, true))
}

func replicaSetMultiCloudTestCase(t *testing.T, usePreviewProvider bool) resource.TestCase {
	t.Helper()
	var (
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName        = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
	)

	return resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicaSetMultiCloud(t, usePreviewProvider, orgID, projectName, clusterName),
				Check:  checkReplicaSetMultiCloud(usePreviewProvider, clusterName, 3),
			},
			{
				Config: configReplicaSetMultiCloud(t, usePreviewProvider, orgID, projectName, clusterNameUpdated),
				Check:  checkReplicaSetMultiCloud(usePreviewProvider, clusterNameUpdated, 3),
			},
			acc.TestStepImportCluster(resourceName),
		},
	}
}

func TestAccClusterAdvancedCluster_singleShardedMultiCloud(t *testing.T) {
	resource.ParallelTest(t, singleShardedMultiCloudTestCase(t, true))
}

func singleShardedMultiCloudTestCase(t *testing.T, usePreviewProvider bool) resource.TestCase {
	t.Helper()
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 7)
		clusterNameUpdated     = acc.RandomClusterName()
	)

	return resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedOldSchemaMultiCloud(t, usePreviewProvider, projectID, clusterName, 1, "M10", nil),
				Check:  checkShardedOldSchemaMultiCloud(usePreviewProvider, clusterName, 1, "M10", true, nil),
			},
			{
				Config: configShardedOldSchemaMultiCloud(t, usePreviewProvider, projectID, clusterNameUpdated, 1, "M10", nil),
				Check:  checkShardedOldSchemaMultiCloud(usePreviewProvider, clusterNameUpdated, 1, "M10", true, nil),
			},
			acc.TestStepImportCluster(resourceName),
		},
	}
}

func TestAccClusterAdvancedCluster_unpausedToPaused(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 4)
		instanceSize           = "M10"
		anotherInstanceSize    = "M20"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configSingleProviderPaused(t, true, projectID, clusterName, false, instanceSize),
				Check:  checkSingleProviderPaused(true, clusterName, false),
			},
			{
				Config: configSingleProviderPaused(t, true, projectID, clusterName, true, instanceSize), // only pause to avoid `OPERATION_INVALID_MEMBER_REPLICATION_LAG`, more info in HELP-72502
				Check:  checkSingleProviderPaused(true, clusterName, true),
			},
			{
				Config:      configSingleProviderPaused(t, true, projectID, clusterName, true, anotherInstanceSize),
				ExpectError: regexp.MustCompile("CANNOT_UPDATE_PAUSED_CLUSTER"),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_pausedToUnpaused(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 4)
		instanceSize           = "M10"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configSingleProviderPaused(t, true, projectID, clusterName, true, instanceSize),
				Check:  checkSingleProviderPaused(true, clusterName, true),
			},
			{
				Config: configSingleProviderPaused(t, true, projectID, clusterName, false, instanceSize),
				Check:  checkSingleProviderPaused(true, clusterName, false),
			},
			{
				Config:      configSingleProviderPaused(t, true, projectID, clusterName, true, instanceSize),
				ExpectError: regexp.MustCompile("CANNOT_PAUSE_RECENTLY_RESUMED_CLUSTER"),
			},
			{
				Config: configSingleProviderPaused(t, true, projectID, clusterName, false, instanceSize),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_advancedConfig_oldMongoDBVersion(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 4)

		processArgs20240530 = &admin20240530.ClusterDescriptionProcessArgs{
			DefaultReadConcern:               conversion.StringPtr("available"),
			DefaultWriteConcern:              conversion.StringPtr("1"),
			FailIndexKeyTooLong:              conversion.Pointer(false),
			JavascriptEnabled:                conversion.Pointer(true),
			MinimumEnabledTlsProtocol:        conversion.StringPtr("TLS1_2"),
			NoTableScan:                      conversion.Pointer(false),
			OplogSizeMB:                      conversion.Pointer(1000),
			SampleRefreshIntervalBIConnector: conversion.Pointer(310),
			SampleSizeBIConnector:            conversion.Pointer(110),
			TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
		}
		processArgs = &admin.ClusterDescriptionProcessArgs20240805{
			ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: conversion.IntPtr(-1), // this will not be set in the TF configuration
			DefaultMaxTimeMS: conversion.IntPtr(65),
		}

		processArgsCipherConfig = &admin.ClusterDescriptionProcessArgs20240805{
			TlsCipherConfigMode:            conversion.StringPtr("CUSTOM"),
			CustomOpensslCipherConfigTls12: &[]string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"},
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configAdvanced(t, true, projectID, clusterName, "6.0", processArgs20240530, processArgs),
				ExpectError: regexp.MustCompile(advancedcluster.ErrorDefaultMaxTimeMinVersion),
			},
			{
				Config: configAdvanced(t, true, projectID, clusterName, "6.0", processArgs20240530, processArgsCipherConfig),
				Check:  checkAdvanced(true, clusterName, "TLS1_2", processArgsCipherConfig),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_advancedConfig(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 4)
		clusterNameUpdated     = acc.RandomClusterName()
		processArgs20240530    = &admin20240530.ClusterDescriptionProcessArgs{
			DefaultReadConcern:               conversion.StringPtr("available"),
			DefaultWriteConcern:              conversion.StringPtr("1"),
			FailIndexKeyTooLong:              conversion.Pointer(false),
			JavascriptEnabled:                conversion.Pointer(true),
			MinimumEnabledTlsProtocol:        conversion.StringPtr("TLS1_2"),
			NoTableScan:                      conversion.Pointer(false),
			OplogSizeMB:                      conversion.Pointer(1000),
			SampleRefreshIntervalBIConnector: conversion.Pointer(310),
			SampleSizeBIConnector:            conversion.Pointer(110),
			TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
		}
		processArgs = &admin.ClusterDescriptionProcessArgs20240805{
			ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: conversion.IntPtr(-1), // this will not be set in the TF configuration
			TlsCipherConfigMode: conversion.StringPtr("DEFAULT"),
		}

		processArgs20240530Updated = &admin20240530.ClusterDescriptionProcessArgs{
			DefaultReadConcern:               conversion.StringPtr("available"),
			DefaultWriteConcern:              conversion.StringPtr("0"),
			FailIndexKeyTooLong:              conversion.Pointer(false),
			JavascriptEnabled:                conversion.Pointer(true),
			MinimumEnabledTlsProtocol:        conversion.StringPtr("TLS1_2"),
			NoTableScan:                      conversion.Pointer(false),
			OplogSizeMB:                      conversion.Pointer(1000),
			SampleRefreshIntervalBIConnector: conversion.Pointer(310),
			SampleSizeBIConnector:            conversion.Pointer(110),
			TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
		}
		processArgsUpdated = &admin.ClusterDescriptionProcessArgs20240805{
			DefaultMaxTimeMS: conversion.IntPtr(65),
			ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: conversion.IntPtr(100),
			TlsCipherConfigMode:            conversion.StringPtr("CUSTOM"),
			CustomOpensslCipherConfigTls12: &[]string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"},
		}
		processArgsUpdatedCipherConfig = &admin.ClusterDescriptionProcessArgs20240805{
			DefaultMaxTimeMS: conversion.IntPtr(65),
			ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: conversion.IntPtr(100),
			TlsCipherConfigMode: conversion.StringPtr("DEFAULT"), // To unset TlsCipherConfigMode, user needs to set this to DEFAULT
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvanced(t, true, projectID, clusterName, "", processArgs20240530, processArgs),
				Check:  checkAdvanced(true, clusterName, "TLS1_2", processArgs),
			},
			{
				Config: configAdvanced(t, true, projectID, clusterNameUpdated, "", processArgs20240530Updated, processArgsUpdated),
				Check:  checkAdvanced(true, clusterNameUpdated, "TLS1_2", processArgsUpdated),
			},
			{
				Config: configAdvanced(t, true, projectID, clusterNameUpdated, "", processArgs20240530Updated, processArgsUpdatedCipherConfig),
				Check:  checkAdvanced(true, clusterNameUpdated, "TLS1_2", processArgsUpdatedCipherConfig),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_defaultWrite(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 4)
		clusterNameUpdated     = acc.RandomClusterName()
		processArgs            = &admin20240530.ClusterDescriptionProcessArgs{
			DefaultReadConcern:               conversion.StringPtr("available"),
			DefaultWriteConcern:              conversion.StringPtr("1"),
			JavascriptEnabled:                conversion.Pointer(true),
			MinimumEnabledTlsProtocol:        conversion.StringPtr("TLS1_2"),
			NoTableScan:                      conversion.Pointer(false),
			OplogSizeMB:                      conversion.Pointer(1000),
			SampleRefreshIntervalBIConnector: conversion.Pointer(310),
			SampleSizeBIConnector:            conversion.Pointer(110),
		}
		processArgsUpdated = &admin20240530.ClusterDescriptionProcessArgs{
			DefaultReadConcern:               conversion.StringPtr("available"),
			DefaultWriteConcern:              conversion.StringPtr("majority"),
			JavascriptEnabled:                conversion.Pointer(true),
			MinimumEnabledTlsProtocol:        conversion.StringPtr("TLS1_2"),
			NoTableScan:                      conversion.Pointer(false),
			OplogSizeMB:                      conversion.Pointer(1000),
			SampleRefreshIntervalBIConnector: conversion.Pointer(310),
			SampleSizeBIConnector:            conversion.Pointer(110),
			TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvancedDefaultWrite(t, true, projectID, clusterName, processArgs),
				Check:  checkAdvancedDefaultWrite(true, clusterName, "1", "TLS1_2"),
			},
			{
				Config: configAdvancedDefaultWrite(t, true, projectID, clusterNameUpdated, processArgsUpdated),
				Check:  checkAdvancedDefaultWrite(true, clusterNameUpdated, "majority", "TLS1_2"),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedClusterConfig_replicationSpecsAutoScaling(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 4)
		autoScaling            = &admin.AdvancedAutoScalingSettings{
			Compute: &admin.AdvancedComputeAutoScaling{Enabled: conversion.Pointer(false), MaxInstanceSize: conversion.StringPtr("")},
			DiskGB:  &admin.DiskGBAutoScaling{Enabled: conversion.Pointer(true)},
		}
		autoScalingUpdated = &admin.AdvancedAutoScalingSettings{
			Compute: &admin.AdvancedComputeAutoScaling{Enabled: conversion.Pointer(true), MaxInstanceSize: conversion.StringPtr("M20")},
			DiskGB:  &admin.DiskGBAutoScaling{Enabled: conversion.Pointer(true)},
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicationSpecsAutoScaling(t, true, projectID, clusterName, autoScaling, "M10", 10, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckExistsCluster(resourceName),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "name", clusterName),
					acc.TestCheckResourceAttrSetPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.#"),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "false"),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "advanced_configuration.0.oplog_min_retention_hours", "5.5"),
				),
			},
			{
				Config: configReplicationSpecsAutoScaling(t, true, projectID, clusterName, autoScalingUpdated, "M20", 20, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckExistsCluster(resourceName),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "name", clusterName),
					acc.TestCheckResourceAttrSetPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.#"),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "true"),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M10"), // modified instance size in config is ignored
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb", "10"),   // modified disk size gb in config is ignored
				),
			},
			// empty plan when auto_scaling block is removed (also aligns instance_size/disk_size_gb to values in state)
			acc.TestStepCheckEmptyPlan(configReplicationSpecsAutoScaling(t, true, projectID, clusterName, nil, "M10", 10, 1)),
			{
				Config: configReplicationSpecsAutoScaling(t, true, projectID, clusterName, nil, "M10", 10, 2), // other change after autoscaling block removed, preserves previous state
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckExistsCluster(resourceName),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "name", clusterName),
					acc.TestCheckResourceAttrSetPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.#"),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "true"), // autoscaling value is preserved
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.node_count", "2"),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.0.electable_specs.0.instance_size", "M10"),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb", "10"),
				),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedClusterConfig_replicationSpecsAnalyticsAutoScaling(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 4)
		clusterNameUpdated     = acc.RandomClusterName()
		autoScaling            = &admin.AdvancedAutoScalingSettings{
			Compute: &admin.AdvancedComputeAutoScaling{Enabled: conversion.Pointer(false), MaxInstanceSize: conversion.StringPtr("")},
			DiskGB:  &admin.DiskGBAutoScaling{Enabled: conversion.Pointer(true)},
		}
		autoScalingUpdated = &admin.AdvancedAutoScalingSettings{
			Compute: &admin.AdvancedComputeAutoScaling{Enabled: conversion.Pointer(true), MaxInstanceSize: conversion.StringPtr("M20")},
			DiskGB:  &admin.DiskGBAutoScaling{Enabled: conversion.Pointer(true)},
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicationSpecsAnalyticsAutoScaling(t, true, projectID, clusterName, autoScaling, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckExistsCluster(resourceName),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "name", clusterName),
					acc.TestCheckResourceAttrSetPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.#"),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "false"),
				),
			},
			{
				Config: configReplicationSpecsAnalyticsAutoScaling(t, true, projectID, clusterNameUpdated, autoScalingUpdated, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckExistsCluster(resourceName),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "name", clusterNameUpdated),
					acc.TestCheckResourceAttrSetPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.#"),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "true"),
				),
			},
			// empty plan when analytics_auto_scaling block is removed
			acc.TestStepCheckEmptyPlan(configReplicationSpecsAnalyticsAutoScaling(t, true, projectID, clusterNameUpdated, nil, 1)),
			{
				Config: configReplicationSpecsAnalyticsAutoScaling(t, true, projectID, clusterNameUpdated, nil, 2), // other changes after analytics_auto_scaling block removed, preserves previous state
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckExistsCluster(resourceName),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "name", clusterNameUpdated),
					acc.TestCheckResourceAttrSetPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.#"),
					acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "true"),
				),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedClusterConfig_singleShardedTransitionToOldSchemaExpectsError(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 9)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGeoShardedOldSchema(t, true, projectID, clusterName, 1, 1, false),
				Check:  checkGeoShardedOldSchema(true, clusterName, 1, 1, true, true),
			},
			acc.TestStepImportCluster(resourceName),
			{
				Config:      configGeoShardedOldSchema(t, true, projectID, clusterName, 1, 2, false),
				ExpectError: regexp.MustCompile(advancedcluster.ErrorOperationNotPermitted),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_withTags(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to check correctly plural data source in the different test steps
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "tags"),
				Check:  checkKeyValueBlocks(true, true, "tags"),
			},
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "tags", acc.ClusterTagsMap1, acc.ClusterTagsMap2),
				Check:  checkKeyValueBlocks(true, true, "tags", acc.ClusterTagsMap1, acc.ClusterTagsMap2),
			},
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "tags", acc.ClusterTagsMap3),
				Check:  checkKeyValueBlocks(true, true, "tags", acc.ClusterTagsMap3),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_withLabels(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to check correctly plural data source in the different test steps
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "labels"),
				Check:  checkKeyValueBlocks(true, true, "labels"),
			},
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "labels", acc.ClusterLabelsMap1, acc.ClusterLabelsMap2),
				Check:  checkKeyValueBlocks(true, true, "labels", acc.ClusterLabelsMap1, acc.ClusterLabelsMap2),
			},
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "labels", acc.ClusterLabelsMap3),
				Check:  checkKeyValueBlocks(true, true, "labels", acc.ClusterLabelsMap3),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_withLabelIgnored(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "labels", acc.ClusterLabelsMapIgnored),
				ExpectError: regexp.MustCompile(advancedclustertpf.ErrLegacyIgnoreLabel.Error()),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_selfManagedSharding(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 6)
		checks                 = []resource.TestCheckFunc{
			acc.CheckExistsCluster(resourceName),
			acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "global_cluster_self_managed_sharding", "true"),
			acc.TestCheckResourceAttrPreviewProviderV2(true, dataSourceName, "global_cluster_self_managed_sharding", "true"),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGeoShardedOldSchema(t, true, projectID, clusterName, 1, 1, true),
				Check: resource.ComposeAggregateTestCheckFunc(checks...,
				),
			},
			acc.TestStepImportCluster(resourceName),
			{
				Config:      configGeoShardedOldSchema(t, true, projectID, clusterName, 1, 1, false),
				ExpectError: regexp.MustCompile("CANNOT_MODIFY_GLOBAL_CLUSTER_MANAGEMENT_SETTING"),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_selfManagedShardingIncorrectType(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 1)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configIncorrectTypeGobalClusterSelfManagedSharding(t, true, projectID, clusterName),
				ExpectError: regexp.MustCompile("CANNOT_SET_SELF_MANAGED_SHARDING_FOR_NON_GLOBAL_CLUSTER"),
			},
		},
	})
}

func TestAccMockableAdvancedCluster_symmetricShardedOldSchema(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 12)
	)

	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedOldSchemaMultiCloud(t, true, projectID, clusterName, 2, "M10", &configServerManagementModeFixedToDedicated),
				Check:  checkShardedOldSchemaMultiCloud(true, clusterName, 2, "M10", false, &configServerManagementModeFixedToDedicated),
			},
			{
				Config: configShardedOldSchemaMultiCloud(t, true, projectID, clusterName, 2, "M20", &configServerManagementModeAtlasManaged),
				Check:  checkShardedOldSchemaMultiCloud(true, clusterName, 2, "M20", false, &configServerManagementModeAtlasManaged),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs"), // Import with old schema will NOT use `num_shards`
		},
	})
}

func TestAccClusterAdvancedClusterConfig_symmetricGeoShardedOldSchema(t *testing.T) {
	resource.ParallelTest(t, symmetricGeoShardedOldSchemaTestCase(t, true))
}

func symmetricGeoShardedOldSchemaTestCase(t *testing.T, usePreviewProvider bool) resource.TestCase {
	t.Helper()
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 18)

	return resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGeoShardedOldSchema(t, usePreviewProvider, projectID, clusterName, 2, 2, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkGeoShardedOldSchema(usePreviewProvider, clusterName, 2, 2, true, false),
					acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER")),
			},
			{
				Config: configGeoShardedOldSchema(t, usePreviewProvider, projectID, clusterName, 3, 3, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkGeoShardedOldSchema(usePreviewProvider, clusterName, 3, 3, true, false),
					acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER")),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs"), // Import with old schema will NOT use `num_shards`
		},
	}
}

func TestAccMockableAdvancedCluster_symmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 6)

	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedOldSchemaDiskSizeGBElectableLevel(t, true, projectID, clusterName, 50),
				Check:  checkShardedOldSchemaDiskSizeGBElectableLevel(true, 50),
			},
			{
				Config: configShardedOldSchemaDiskSizeGBElectableLevel(t, true, projectID, clusterName, 55),
				Check:  checkShardedOldSchemaDiskSizeGBElectableLevel(true, 55),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs"), // Import with old schema will NOT use `num_shards`
		},
	})
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedNewSchemaToAsymmetricAddingRemovingShard(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedNewSchema(t, true, orgID, projectName, clusterName, 50, "M10", "M10", nil, nil, false, false),
				Check:  checkShardedNewSchema(true, 50, "M10", "M10", nil, nil, false, false),
			},
			{
				Config: configShardedNewSchema(t, true, orgID, projectName, clusterName, 55, "M10", "M20", nil, nil, true, false), // add middle replication spec and transition to asymmetric
				Check:  checkShardedNewSchema(true, 55, "M10", "M20", nil, nil, true, true),
			},
			{
				Config: configShardedNewSchema(t, true, orgID, projectName, clusterName, 55, "M10", "M20", nil, nil, false, false), // removes middle replication spec
				Check:  checkShardedNewSchema(true, 55, "M10", "M20", nil, nil, true, false),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedClusterConfig_asymmetricShardedNewSchema(t *testing.T) {
	resource.ParallelTest(t, asymmetricShardedNewSchemaTestCase(t, true))
}

func asymmetricShardedNewSchemaTestCase(t *testing.T, usePreviewProvider bool) resource.TestCase {
	t.Helper()
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	return resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedNewSchema(t, usePreviewProvider, orgID, projectName, clusterName, 50, "M30", "M40", admin.PtrInt(2000), admin.PtrInt(2500), false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkShardedNewSchema(usePreviewProvider, 50, "M30", "M40", admin.PtrInt(2000), admin.PtrInt(2500), true, false),
					resource.TestCheckResourceAttr("data.mongodbatlas_advanced_clusters.test-replication-specs-per-shard-false", "results.#", "0"),
					acc.CheckIndependentShardScalingMode(resourceName, clusterName, "SHARD")),
			},
			acc.TestStepImportCluster(resourceName),
		},
	}
}

func TestAccClusterAdvancedClusterConfig_asymmetricShardedNewSchemaInconsistentDisk(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configShardedNewSchema(t, true, orgID, projectName, clusterName, 50, "M30", "M40", admin.PtrInt(2000), admin.PtrInt(2500), false, true),
				ExpectError: regexp.MustCompile("DISK_SIZE_GB_INCONSISTENT"), // API Error when disk size is not consistent across all shards
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_asymmetricGeoShardedNewSchemaAddingRemovingShard(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 9)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGeoShardedNewSchema(t, true, projectID, clusterName, false),
				Check:  checkGeoShardedNewSchema(true, false),
			},
			{
				Config: configGeoShardedNewSchema(t, true, projectID, clusterName, true),
				Check:  checkGeoShardedNewSchema(true, true),
			},
			{
				Config: configGeoShardedNewSchema(t, true, projectID, clusterName, false),
				Check:  checkGeoShardedNewSchema(true, false),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedClusterConfig_shardedTransitionFromOldToNewSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkShardedTransitionOldToNewSchema(true, false),
					acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER")),
			},
			{
				Config: configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, false),
				Check:  checkShardedTransitionOldToNewSchema(true, true),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedClusterConfig_geoShardedTransitionFromOldToNewSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, false),
				Check:  checkGeoShardedTransitionOldToNewSchema(true, false),
			},
			{
				Config: configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true),
				Check:  checkGeoShardedTransitionOldToNewSchema(true, true),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccAdvancedCluster_replicaSetScalingStrategyAndRedactClientLogData(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(t, true, orgID, projectName, clusterName, "WORKLOAD_TYPE", true),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "WORKLOAD_TYPE", true),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(t, true, orgID, projectName, clusterName, "SEQUENTIAL", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "SEQUENTIAL", false),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(t, true, orgID, projectName, clusterName, "NODE_TYPE", true),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "NODE_TYPE", true),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(t, true, orgID, projectName, clusterName, "NODE_TYPE", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "NODE_TYPE", false),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccAdvancedCluster_replicaSetScalingStrategyAndRedactClientLogDataOldSchema(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(t, true, orgID, projectName, clusterName, "WORKLOAD_TYPE", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "WORKLOAD_TYPE", false),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(t, true, orgID, projectName, clusterName, "SEQUENTIAL", true),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "SEQUENTIAL", true),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(t, true, orgID, projectName, clusterName, "NODE_TYPE", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "NODE_TYPE", false),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs"), // Import with old schema will NOT use `num_shards`
		},
	})
}

// TestAccClusterAdvancedCluster_priorityOldSchema will be able to be simplied or deleted in CLOUDP-275825
func TestAccClusterAdvancedCluster_priorityOldSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 6)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configPriority(t, true, projectID, clusterName, true, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
			{
				Config: configPriority(t, true, projectID, clusterName, true, false),
				Check:  acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.#", "2"),
			},
			{
				Config:      configPriority(t, true, projectID, clusterName, true, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
			// Extra step added to allow deletion, otherwise we get `Error running post-test destroy` since validation of TF fails
			{
				Config: configPriority(t, true, projectID, clusterName, true, false),
				Check:  acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.#", "2"),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs"), // Import with old schema will NOT use `num_shards`
		},
	})
}

// TestAccClusterAdvancedCluster_priorityNewSchema will be able to be simplied or deleted in CLOUDP-275825
func TestAccClusterAdvancedCluster_priorityNewSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 3)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configPriority(t, true, projectID, clusterName, false, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
			{
				Config: configPriority(t, true, projectID, clusterName, false, false),
				Check:  acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.#", "2"),
			},
			{
				Config:      configPriority(t, true, projectID, clusterName, false, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
			// Extra step added to allow deletion, otherwise we get `Error running post-test destroy` since validation of TF fails
			{
				Config: configPriority(t, true, projectID, clusterName, false, false),
				Check:  acc.TestCheckResourceAttrPreviewProviderV2(true, resourceName, "replication_specs.0.region_configs.#", "2"),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_biConnectorConfig(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 4)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configBiConnectorConfig(t, true, projectID, clusterName, false),
				Check:  checkTenantBiConnectorConfig(true, projectID, clusterName, false),
			},
			{
				Config: configBiConnectorConfig(t, true, projectID, clusterName, true),
				Check:  checkTenantBiConnectorConfig(true, projectID, clusterName, true),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_pinnedFCVWithVersionUpgradeAndDowngrade(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // Using single project to assert plural data source
		clusterName = acc.RandomClusterName()
	)

	now := time.Now()
	// Time 7 days from now, truncated to the beginning of the day
	sevenDaysFromNow := now.AddDate(0, 0, 7).Truncate(24 * time.Hour)
	firstExpirationDate := conversion.TimeToString(sevenDaysFromNow)
	// Time 8 days from now
	eightDaysFromNow := sevenDaysFromNow.AddDate(0, 0, 1)
	updatedExpirationDate := conversion.TimeToString(eightDaysFromNow)
	invalidDateFormat := "invalid"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configFCVPinning(t, orgID, projectName, clusterName, nil, "7.0"),
				Check:  acc.CheckFCVPinningConfig(true, resourceName, dataSourceName, dataSourcePluralName, 7, nil, nil),
			},
			{ // pins fcv
				Config: configFCVPinning(t, orgID, projectName, clusterName, &firstExpirationDate, "7.0"),
				Check:  acc.CheckFCVPinningConfig(true, resourceName, dataSourceName, dataSourcePluralName, 7, admin.PtrString(firstExpirationDate), admin.PtrInt(7)),
			},
			{ // using incorrect format
				Config:      configFCVPinning(t, orgID, projectName, clusterName, &invalidDateFormat, "7.0"),
				ExpectError: regexp.MustCompile("expiration_date format is incorrect: " + invalidDateFormat),
			},
			{ // updates expiration date of fcv
				Config: configFCVPinning(t, orgID, projectName, clusterName, &updatedExpirationDate, "7.0"),
				Check:  acc.CheckFCVPinningConfig(true, resourceName, dataSourceName, dataSourcePluralName, 7, admin.PtrString(updatedExpirationDate), admin.PtrInt(7)),
			},
			{ // upgrade mongodb version with fcv pinned
				Config: configFCVPinning(t, orgID, projectName, clusterName, &updatedExpirationDate, "8.0"),
				Check:  acc.CheckFCVPinningConfig(true, resourceName, dataSourceName, dataSourcePluralName, 8, admin.PtrString(updatedExpirationDate), admin.PtrInt(7)),
			},
			{ // downgrade mongodb version with fcv pinned
				Config: configFCVPinning(t, orgID, projectName, clusterName, &updatedExpirationDate, "7.0"),
				Check:  acc.CheckFCVPinningConfig(true, resourceName, dataSourceName, dataSourcePluralName, 7, admin.PtrString(updatedExpirationDate), admin.PtrInt(7)),
			},
			{ // unpins fcv
				Config: configFCVPinning(t, orgID, projectName, clusterName, nil, "7.0"),
				Check:  acc.CheckFCVPinningConfig(true, resourceName, dataSourceName, dataSourcePluralName, 7, nil, nil),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccAdvancedCluster_oldToNewSchemaWithAutoscalingEnabled(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, false, true),
				Check:  acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER"),
			},
			{
				Config: configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, true),
				Check:  acc.CheckIndependentShardScalingMode(resourceName, clusterName, "SHARD"),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccAdvancedCluster_oldToNewSchemaWithAutoscalingDisabledToEnabled(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, false, false),
				Check:  acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER"),
			},
			{
				Config: configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, false),
				Check:  acc.CheckIndependentShardScalingMode(resourceName, clusterName, "CLUSTER"),
			},
			{
				Config: configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, true),
				Check:  acc.CheckIndependentShardScalingMode(resourceName, clusterName, "SHARD"),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccMockableAdvancedCluster_replicasetAdvConfigUpdate(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 3)
		checksMap              = map[string]string{
			"state_name": "IDLE",
		}
		checksSet = []string{
			"replication_specs.0.container_id.AWS:US_EAST_1",
			"mongo_db_major_version",
		}
		timeoutCheck   = resource.TestCheckResourceAttr(resourceName, "timeouts.create", "6000s") // timeouts.create is not set on data sources
		tagsLabelsMap  = map[string]string{"key": "env", "value": "test"}
		tagsCheck      = checkKeyValueBlocks(true, false, "tags", tagsLabelsMap)
		labelsCheck    = checkKeyValueBlocks(true, false, "labels", tagsLabelsMap)
		checks         = checkAggr(true, checksSet, checksMap, timeoutCheck)
		afterUpdateMap = map[string]string{
			"state_name":                    "IDLE",
			"backup_enabled":                "true",
			"bi_connector_config.0.enabled": "true",
			"pit_enabled":                   "true",
			"redact_client_log_data":        "true",
			"replica_set_scaling_strategy":  "NODE_TYPE",
			"root_cert_type":                "ISRGROOTX1",
			"version_release_system":        "CONTINUOUS",
			"advanced_configuration.0.change_stream_options_pre_and_post_images_expire_after_seconds": "100",
			"advanced_configuration.0.default_read_concern":                                           "available",
			"advanced_configuration.0.default_write_concern":                                          "majority",
			"advanced_configuration.0.javascript_enabled":                                             "true",
			"advanced_configuration.0.minimum_enabled_tls_protocol":                                   "TLS1_2",
			"advanced_configuration.0.no_table_scan":                                                  "true",
			"advanced_configuration.0.sample_refresh_interval_bi_connector":                           "310",
			"advanced_configuration.0.sample_size_bi_connector":                                       "110",
			"advanced_configuration.0.transaction_lifetime_limit_seconds":                             "300",
			"advanced_configuration.0.tls_cipher_config_mode":                                         "CUSTOM",
			"advanced_configuration.0.custom_openssl_cipher_config_tls12.#":                           "1",
			"advanced_configuration.0.default_max_time_ms":                                            "65",
		}
		checksUpdate = checkAggr(true, checksSet, afterUpdateMap, timeoutCheck, tagsCheck, labelsCheck)
		fullUpdate   = `
	backup_enabled = true
	bi_connector_config {
		enabled = true
	}
	labels {
		key   = "env"
		value = "test"
	}
	tags {
		key   = "env"
		value = "test"
	}
	pit_enabled = true
	redact_client_log_data = true
	replica_set_scaling_strategy = "NODE_TYPE"
	root_cert_type = "ISRGROOTX1"
	version_release_system = "CONTINUOUS"
	
	advanced_configuration {
		change_stream_options_pre_and_post_images_expire_after_seconds = 100
		default_read_concern                                           = "available"
		default_write_concern                                          = "majority"
		javascript_enabled                                             = true
		minimum_enabled_tls_protocol                                   = "TLS1_2" # This cluster does not support TLS1.0 or TLS1.1. If you must use old TLS versions contact MongoDB support
		no_table_scan                                                  = true
		sample_refresh_interval_bi_connector                           = 310
		sample_size_bi_connector                                       = 110
		transaction_lifetime_limit_seconds                             = 300
		custom_openssl_cipher_config_tls12							   = ["TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"]
		tls_cipher_config_mode               						   = "CUSTOM"
		default_max_time_ms											   = 65
	}
`
	)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, &resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasicReplicaset(t, projectID, clusterName, "", ""),
				Check:  checks,
			},
			{
				Config: configBasicReplicaset(t, projectID, clusterName, fullUpdate, ""),
				Check:  checksUpdate,
			},
			{
				Config: configBasicReplicaset(t, projectID, clusterName, "", ""),
				Check:  checks,
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccMockableAdvancedCluster_shardedAddAnalyticsAndAutoScaling(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 8)
		checksMap              = map[string]string{
			"state_name": "IDLE",
			"project_id": projectID,
			"name":       clusterName,
		}
		checksUpdatedMap = map[string]string{
			"replication_specs.0.region_configs.0.auto_scaling.0.disk_gb_enabled":    "true",
			"replication_specs.0.region_configs.0.electable_specs.0.instance_size":   "M30",
			"replication_specs.0.region_configs.0.analytics_specs.0.instance_size":   "M30",
			"replication_specs.0.region_configs.0.analytics_specs.0.node_count":      "1",
			"replication_specs.0.region_configs.0.analytics_specs.0.disk_iops":       "2000",
			"replication_specs.0.region_configs.0.analytics_specs.0.ebs_volume_type": "PROVISIONED",
			"replication_specs.1.region_configs.0.analytics_specs.0.instance_size":   "M30",
			"replication_specs.1.region_configs.0.analytics_specs.0.node_count":      "1",
			"replication_specs.1.region_configs.0.analytics_specs.0.ebs_volume_type": "PROVISIONED",
			"replication_specs.1.region_configs.0.analytics_specs.0.disk_iops":       "1000",
		}
		checksUpdated = checkAggr(true, nil, checksUpdatedMap)
	)
	if config.PreviewProviderV2AdvancedCluster() { // SDKv2 don't set "computed" specs in the state
		checksMap["replication_specs.0.region_configs.0.electable_specs.0.instance_size"] = "M30"
		checksMap["replication_specs.0.region_configs.0.analytics_specs.0.node_count"] = "0"
	}
	checks := checkAggr(true, nil, checksMap)
	checksMap["replication_specs.0.region_configs.0.analytics_specs.0.node_count"] = "1" // analytics_specs is kept even if it's removed from the config
	checksAfter := checkAggr(true, nil, checksMap)
	unit.CaptureOrMockTestCaseAndRun(t, mockConfig, &resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configSharded(t, projectID, clusterName, false),
				Check:  checks,
			},
			{
				Config: configSharded(t, projectID, clusterName, true),
				Check:  checksUpdated,
			},
			{
				Config: configSharded(t, projectID, clusterName, false),
				Check:  checksAfter,
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccAdvancedCluster_removeBlocksFromConfig(t *testing.T) {
	if !config.PreviewProviderV2AdvancedCluster() { // SDKv2 don't set "computed" specs in the state
		t.Skip("This test is not applicable for SDKv2")
	}
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 15)
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBlocks(t, projectID, clusterName, "M10", true),
				Check:  checkBlocks("M10"),
			},
			// removing blocks generates an empty plan
			acc.TestStepCheckEmptyPlan(configBlocks(t, projectID, clusterName, "M10", false)),
			{
				Config: configBlocks(t, projectID, clusterName, "M20", false), // applying a change after removing blocks preserves previous state
				Check:  checkBlocks("M20"),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccAdvancedCluster_createTimeoutWithDeleteOnCreateReplicaset(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 3)
		configCall             = func(t *testing.T, timeoutSection string) string {
			t.Helper()
			return configBasicReplicaset(t, projectID, clusterName, "", timeoutSection)
		}
		waitOnClusterDeleteDone = func() {
			diags := &diag.Diagnostics{}
			clusterResp, _ := advancedclustertpf.GetClusterDetails(t.Context(), diags, projectID, clusterName, acc.MongoDBClient, false)
			if clusterResp == nil {
				t.Fatalf("cluster %s not found in %s", clusterName, projectID)
			}
			advancedclustertpf.AwaitChanges(t.Context(), acc.MongoDBClient, &advancedclustertpf.ClusterWaitParams{
				ProjectID:   projectID,
				ClusterName: clusterName,
				Timeout:     60 * time.Second,
				IsDelete:    true,
			}, "waiting for cluster to be deleted after cleanup in create timeout", diags)
			time.Sleep(1 * time.Minute) // decrease the chance of `CONTAINER_WAITING_FOR_FAST_RECORD_CLEAN_UP`: "A transient error occurred. Please try again in a minute or use a different name"
		}
	)
	resource.ParallelTest(t, *createCleanupTest(t, configCall, waitOnClusterDeleteDone, true))
}

func TestAccAdvancedCluster_createTimeoutWithDeleteOnCreateFlex(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 1)
		configCall             = func(t *testing.T, timeoutSection string) string {
			t.Helper()
			return acc.ConvertAdvancedClusterToPreviewProviderV2(t, true, fmt.Sprintf(`
			resource "mongodbatlas_advanced_cluster" "test" {
				project_id   = %[1]q
				name         = %[2]q
				cluster_type = "REPLICASET"
				replication_specs {
					region_configs {
						provider_name = "FLEX"
						backing_provider_name = "AWS"
						region_name = "US_EAST_1"
						priority      = 7
					}
				}
				%[3]s
			}`, projectID, clusterName, timeoutSection))
		}
		waitOnClusterDeleteDone = func() {
			err := flexcluster.WaitStateTransitionDelete(t.Context(), &admin.GetFlexClusterApiParams{
				GroupId: projectID,
				Name:    clusterName,
			}, acc.ConnV2().FlexClustersApi)
			require.NoError(t, err)
			time.Sleep(1 * time.Minute) // decrease the chance of `CONTAINER_WAITING_FOR_FAST_RECORD_CLEAN_UP`: "A transient error occurred. Please try again in a minute or use a different name"
		}
	)
	resource.ParallelTest(t, *createCleanupTest(t, configCall, waitOnClusterDeleteDone, false))
}

func createCleanupTest(t *testing.T, configCall func(t *testing.T, timeoutSection string) string, waitOnClusterDeleteDone func(), isUpdateSupported bool) *resource.TestCase {
	t.Helper()
	var (
		timeoutsStrShort = `
			timeouts {
				create = "2s"
			}
			delete_on_create_timeout = true
		`
		timeoutsStrLong      = strings.ReplaceAll(timeoutsStrShort, "2s", "6000s")
		timeoutsStrLongFalse = strings.ReplaceAll(timeoutsStrLong, "true", "false")
	)
	steps := []resource.TestStep{
		{
			Config:      configCall(t, timeoutsStrShort),
			ExpectError: regexp.MustCompile("context deadline exceeded"),
		},
		// OK create should keep the delete_on_create_timeout flag and should be no cleanup
		{
			PreConfig: waitOnClusterDeleteDone,
			Config:    configCall(t, timeoutsStrLong),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceName, "delete_on_create_timeout", "true"),
			),
		},
		acc.TestStepImportCluster(resourceName),
	}
	if isUpdateSupported {
		steps = append(steps,
			// Switch delete_on_create_timeout to false
			resource.TestStep{
				Config: configCall(t, timeoutsStrLongFalse),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "delete_on_create_timeout", "false"),
				),
			},
		)
		deleteOnCreateTimeoutRemoved := configCall(t, "")
		if config.PreviewProviderV2AdvancedCluster() {
			steps = append(steps,
				resource.TestStep{
					Config: deleteOnCreateTimeoutRemoved,
					Check:  resource.TestCheckNoResourceAttr(resourceName, "delete_on_create_timeout"),
				})
		} else {
			// removing an optional false value has no affect in SDKv2, as false==null and no-plan-change
			steps = append(steps, acc.TestStepCheckEmptyPlan(deleteOnCreateTimeoutRemoved))
		}
		steps = append(steps, acc.TestStepImportCluster(resourceName))
	}
	return &resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps:                    steps,
	}
}

func configBasicReplicaset(t *testing.T, projectID, clusterName, extra, timeoutStr string) string {
	t.Helper()
	if timeoutStr == "" {
		timeoutStr = `
			timeouts {
				create = "6000s"
			}`
	}
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, true, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			%[4]s		
			project_id = %[1]q
			name = %[2]q
			cluster_type = "REPLICASET"
			replication_specs {
				region_configs {
					priority        = 7
					provider_name = "AWS"
					region_name     = "US_EAST_1"
					auto_scaling {
						compute_scale_down_enabled = false
						compute_enabled = false
						disk_gb_enabled = true
					}
					electable_specs {
						node_count = 3
						instance_size = "M10"
						disk_size_gb = 10
					}
				}
			}
			%[3]s
		}
	`, projectID, clusterName, extra, timeoutStr)) + dataSourcesTFNewSchema
}

func configSharded(t *testing.T, projectID, clusterName string, withUpdate bool) string {
	t.Helper()
	var autoScaling, analyticsSpecs string
	if withUpdate {
		autoScaling = `
			auto_scaling {
				disk_gb_enabled = true
			}`
		analyticsSpecs = `
			analytics_specs {
				instance_size   = "M30"
				node_count      = 1
				ebs_volume_type = "PROVISIONED"
				disk_iops       = 2000
			}`
	}
	// SDK v2 Implementation receives many warnings, one of them: `.replication_specs[1].region_configs[0].analytics_specs[0].disk_iops: was cty.NumberIntVal(2000), but now cty.NumberIntVal(1000)`
	// Therefore, in TPF we are forced to set the value that will be returned by the API (1000)
	// The rule is: For any replication spec, the `(analytics|electable|read_only)_spec.disk_iops` must be the same across all region_configs
	// The API raises no errors, but the response reflects this rule
	analyticsSpecsForSpec2 := strings.ReplaceAll(analyticsSpecs, "2000", "1000")
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, true, fmt.Sprintf(`
	resource "mongodbatlas_advanced_cluster" "test" {
		project_id   = %[1]q
		name         = %[2]q
		cluster_type = "SHARDED"

		replication_specs { # shard 1
			region_configs {
				electable_specs {
					instance_size   = "M30"
					disk_iops       = 2000
					node_count      = 3
					ebs_volume_type = "PROVISIONED"
					}
				%[3]s
				%[4]s
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
			}
		}
		replication_specs { # shard 2
			region_configs {
				electable_specs {
					instance_size   = "M30"
					ebs_volume_type = "PROVISIONED"
					disk_iops       = 1000
					node_count      = 3
				}
				%[3]s
				%[5]s
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
			}
		}
	}
	`, projectID, clusterName, autoScaling, analyticsSpecs, analyticsSpecsForSpec2)) + dataSourcesTFNewSchema
}

func configBlocks(t *testing.T, projectID, clusterName, instanceSize string, defineBlocks bool) string {
	t.Helper()
	var extraConfig0, extraConfig1, electableSpecs0 string
	autoScalingBlocks := `
		auto_scaling {
			disk_gb_enabled            = true
			compute_enabled            = true
			compute_min_instance_size  = "M10"
			compute_max_instance_size  = "M30"
			compute_scale_down_enabled = true
		}
		analytics_auto_scaling {
			disk_gb_enabled            = true
			compute_enabled            = true
			compute_min_instance_size  = "M10"
			compute_max_instance_size  = "M30"
			compute_scale_down_enabled = true
		}
	`
	if defineBlocks {
		electableSpecs0 = `
			electable_specs {
				instance_size   = "M10"
				node_count      = 5
			}
		`
		// read only + autoscaling blocks
		extraConfig0 = `
			read_only_specs {
				instance_size = "M10"
				node_count    = 2
			}
		` + autoScalingBlocks
		// read only + analytics + autoscaling blocks
		extraConfig1 = `
			read_only_specs {
				instance_size = "M10"
				node_count    = 1
			}
			analytics_specs {
				instance_size = "M10"
				node_count    = 4
			}
		` + autoScalingBlocks
	}
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, true, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "GEOSHARDED"

			replication_specs { 
				zone_name = "Zone 1"
				region_configs {
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
					%[6]s
					%[4]s
				}
			}

			replication_specs { 
				zone_name = "Zone 2"
				region_configs {
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
					electable_specs {
						instance_size   = %[3]q
						node_count      = 3
					}
					%[5]s
				}
				region_configs { // region with no electable specs
					provider_name = "AWS"
					priority      = 0
					region_name   = "US_EAST_1"
					%[4]s
				}
			}
		}
	`, projectID, clusterName, instanceSize, extraConfig0, extraConfig1, electableSpecs0))
}

func checkBlocks(instanceSize string) resource.TestCheckFunc {
	checksMap := map[string]string{
		"replication_specs.0.region_configs.0.electable_specs.0.instance_size": "M10",
		"replication_specs.0.region_configs.0.electable_specs.0.node_count":    "5",
		"replication_specs.0.region_configs.0.read_only_specs.0.instance_size": "M10",
		"replication_specs.0.region_configs.0.read_only_specs.0.node_count":    "2",
		"replication_specs.0.region_configs.0.analytics_specs.0.node_count":    "0",

		"replication_specs.1.region_configs.0.electable_specs.0.instance_size": instanceSize,
		"replication_specs.1.region_configs.0.electable_specs.0.node_count":    "3",
		"replication_specs.1.region_configs.0.read_only_specs.0.instance_size": instanceSize,
		"replication_specs.1.region_configs.0.read_only_specs.0.node_count":    "1",
		"replication_specs.1.region_configs.0.analytics_specs.0.instance_size": "M10",
		"replication_specs.1.region_configs.0.analytics_specs.0.node_count":    "4",

		"replication_specs.1.region_configs.1.read_only_specs.0.instance_size": instanceSize,
		"replication_specs.1.region_configs.1.read_only_specs.0.node_count":    "2",
	}
	for repSpecsIdx := range 2 {
		for _, block := range []string{"auto_scaling", "analytics_auto_scaling"} {
			checksMap[fmt.Sprintf("replication_specs.%d.region_configs.0.%s.disk_gb_enabled", repSpecsIdx, block)] = "true"
			checksMap[fmt.Sprintf("replication_specs.%d.region_configs.0.%s.compute_enabled", repSpecsIdx, block)] = "true"
			checksMap[fmt.Sprintf("replication_specs.%d.region_configs.0.%s.compute_scale_down_enabled", repSpecsIdx, block)] = "true"
			checksMap[fmt.Sprintf("replication_specs.%d.region_configs.0.%s.compute_min_instance_size", repSpecsIdx, block)] = "M10"
			checksMap[fmt.Sprintf("replication_specs.%d.region_configs.0.%s.compute_max_instance_size", repSpecsIdx, block)] = "M30"
		}
	}
	return resource.ComposeAggregateTestCheckFunc(acc.AddAttrChecksPreviewProviderV2(true, resourceName, nil, checksMap)...)
}

func checkAggr(usePreviewProvider bool, attrsSet []string, attrsMap map[string]string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	extraChecks := extra
	extraChecks = append(extraChecks, acc.CheckExistsCluster(resourceName))
	return acc.CheckRSAndDSPreviewProviderV2(usePreviewProvider, resourceName, admin.PtrString(dataSourceName), nil, attrsSet, attrsMap, extraChecks...)
}

func configTenant(t *testing.T, usePreviewProvider bool, projectID, name, zoneName, instanceSize string) string {
	t.Helper()
	zoneNameLine := ""
	if zoneName != "" {
		zoneNameLine = fmt.Sprintf("zone_name = %q", zoneName)
	}
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = %[4]q
					}
					provider_name         = "TENANT"
					backing_provider_name = "AWS"
					region_name           = "US_EAST_1"
					priority              = 7
				}
				%[3]s
			}
		}
	`, projectID, name, zoneNameLine, instanceSize)) + dataSourcesTFNewSchema
}

func checkTenant(usePreviewProvider bool, projectID, name string) resource.TestCheckFunc {
	pluralChecks := acc.AddAttrSetChecksPreviewProviderV2(usePreviewProvider, dataSourcePluralName, nil,
		[]string{"results.#", "results.0.replication_specs.#", "results.0.name", "results.0.termination_protection_enabled", "results.0.global_cluster_self_managed_sharding"}...)
	return checkAggr(usePreviewProvider,
		[]string{"replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"},
		map[string]string{
			"project_id":                           projectID,
			"name":                                 name,
			"termination_protection_enabled":       "false",
			"global_cluster_self_managed_sharding": "false"},
		pluralChecks...)
}

func checksBasicDedicated(projectID, name string) resource.TestCheckFunc {
	originalChecks := checkTenant(true, projectID, name)
	checkMap := map[string]string{
		"replication_specs.0.region_configs.0.electable_specs.0.node_count":    "3",
		"replication_specs.0.region_configs.0.electable_specs.0.instance_size": "M10",
		"replication_specs.0.region_configs.0.provider_name":                   "AWS",
	}
	return checkAggr(true, nil, checkMap, originalChecks)
}

func configWithKeyValueBlocks(t *testing.T, usePreviewProvider bool, orgID, projectName, clusterName, blockName string, blocks ...map[string]string) string {
	t.Helper()
	var extraConfig string
	for _, block := range blocks {
		extraConfig += fmt.Sprintf(`
			%[1]s {
				key   = %[2]q
				value = %[3]q
			}
		`, blockName, block["key"], block["value"])
	}

	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}
		
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q
			cluster_type = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
				}
			}

			%[4]s
		}
	`, orgID, projectName, clusterName, extraConfig)) + dataSourcesTFNewSchema
}

func checkKeyValueBlocks(usePreviewProvider, includeDataSources bool, blockName string, blocks ...map[string]string) resource.TestCheckFunc {
	if config.PreviewProviderV2AdvancedCluster() {
		return checkKeyValueBlocksPreviewProviderV2(usePreviewProvider, includeDataSources, blockName, blocks...)
	}
	const pluralPrefix = "results.0."
	lenStr := strconv.Itoa(len(blocks))
	keyHash := blockName + ".#"
	keyStar := blockName + ".*"
	checks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, resourceName, keyHash, lenStr),
	}
	if includeDataSources {
		checks = append(checks,
			acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, dataSourceName, keyHash, lenStr),
			acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, dataSourcePluralName, pluralPrefix+keyHash, lenStr))
	}
	for _, block := range blocks {
		checks = append(checks,
			acc.TestCheckTypeSetElemNestedAttrsPreviewProviderV2(usePreviewProvider, resourceName, keyStar, block),
		)
		if includeDataSources {
			checks = append(checks,
				acc.TestCheckTypeSetElemNestedAttrsPreviewProviderV2(usePreviewProvider, dataSourceName, keyStar, block),
				acc.TestCheckTypeSetElemNestedAttrsPreviewProviderV2(usePreviewProvider, dataSourcePluralName, pluralPrefix+keyStar, block))
		}
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkKeyValueBlocksPreviewProviderV2(usePreviewProvider, includeDataSources bool, blockName string, blocks ...map[string]string) resource.TestCheckFunc {
	const pluralPrefix = "results.0."
	lenStr := strconv.Itoa(len(blocks))
	keyPct := blockName + ".%"
	checks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, resourceName, keyPct, lenStr),
	}
	if includeDataSources {
		checks = append(checks,
			acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, dataSourceName, keyPct, lenStr),
			acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, dataSourcePluralName, pluralPrefix+keyPct, lenStr))
	}
	for _, block := range blocks {
		key := blockName + "." + block["key"]
		value := block["value"]
		checks = append(checks,
			acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, resourceName, key, value),
		)
		if includeDataSources {
			checks = append(checks,
				acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, dataSourceName, key, value),
				acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, dataSourcePluralName, pluralPrefix+key, value))
		}
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

type ReplicaSetAWSConfig struct {
	ProjectID          string
	ClusterName        string
	ClusterType        string
	DiskSizeGB         int
	NodeCountElectable int
	WithAnalyticsSpecs bool
}

func configAWSProvider(t *testing.T, usePreviewProvider bool, configInfo ReplicaSetAWSConfig) string {
	t.Helper()
	analyticsSpecs := ""
	if configInfo.WithAnalyticsSpecs {
		analyticsSpecs = `
		analytics_specs {
			instance_size = "M10"
			node_count    = 1
		}`
	}
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = %[3]q
			retain_backups_enabled = "true"
			disk_size_gb = %[4]d

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = %[5]d
					}
					%[6]s
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
		}
	`, configInfo.ProjectID, configInfo.ClusterName, configInfo.ClusterType, configInfo.DiskSizeGB, configInfo.NodeCountElectable, analyticsSpecs)) + dataSourcesTFOldSchema
}

func checkReplicaSetAWSProvider(usePreviewProvider bool, projectID, name string, diskSizeGB, nodeCountElectable int, checkDiskSizeGBInnerLevel, checkExternalID bool) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, resourceName, "retain_backups_enabled", "true"),
	}
	additionalChecks = append(additionalChecks,
		acc.TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider, resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider, dataSourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)))

	if checkDiskSizeGBInnerLevel {
		additionalChecks = append(additionalChecks,
			checkAggr(usePreviewProvider, []string{}, map[string]string{
				"replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
				"replication_specs.0.region_configs.0.analytics_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
			}),
		)
	}

	if checkExternalID {
		additionalChecks = append(additionalChecks, acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, resourceName, "replication_specs.0.external_id"))
	}

	return checkAggr(usePreviewProvider,
		[]string{"replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"},
		map[string]string{
			"project_id":   projectID,
			"disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.electable_specs.0.node_count": fmt.Sprintf("%d", nodeCountElectable),
			"replication_specs.0.region_configs.0.analytics_specs.0.node_count": "1",
			"name": name},
		additionalChecks...,
	)
}

func configIncorrectTypeGobalClusterSelfManagedSharding(t *testing.T, usePreviewProvider bool, projectID, name string) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q

			cluster_type = "REPLICASET"
			global_cluster_self_managed_sharding = true # invalid, can only by used with GEOSHARDED clusters

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
		}
	`, projectID, name))
}

func configReplicaSetMultiCloud(t *testing.T, usePreviewProvider bool, orgID, projectName, name string) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q
			cluster_type = "REPLICASET"
			retain_backups_enabled = false

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}

				region_configs {
					read_only_specs {
						instance_size = "M10"
						node_count    = 2
					}
					provider_name = "GCP"
					priority      = 0
					region_name   = "US_EAST_4"
				}

				region_configs {
					read_only_specs {
						instance_size = "M10"
						node_count    = 2
					}
					provider_name = "GCP"
					priority      = 0
					region_name   = "NORTH_AMERICA_NORTHEAST_1"
				}
			}
		}
	`, orgID, projectName, name)) + dataSourcesTFNewSchema
}

func checkReplicaSetMultiCloud(usePreviewProvider bool, name string, regionConfigs int) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, resourceName, "retain_backups_enabled", "false"),
		acc.TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider, resourceName, "replication_specs.0.region_configs.#", acc.JSONEquals(strconv.Itoa(regionConfigs))),
		acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, resourceName, "replication_specs.0.external_id"),
		acc.TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider, dataSourceName, "replication_specs.0.region_configs.#", acc.JSONEquals(strconv.Itoa(regionConfigs))),
		acc.TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider, dataSourcePluralName, "results.0.replication_specs.0.region_configs.#", acc.JSONEquals(strconv.Itoa(regionConfigs))),
		acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, dataSourcePluralName, "results.#"),
		acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, dataSourcePluralName, "results.0.replication_specs.#"),
		acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, dataSourcePluralName, "results.0.name"),
	}
	return checkAggr(usePreviewProvider,
		[]string{"project_id", "replication_specs.#", "replication_specs.0.id"},
		map[string]string{
			"name": name},
		additionalChecks...,
	)
}

func configShardedOldSchemaMultiCloud(t *testing.T, usePreviewProvider bool, projectID, name string, numShards int, analyticsSize string, configServerManagementMode *string) string {
	t.Helper()
	var rootConfig string
	if configServerManagementMode != nil {
		// valid values: FIXED_TO_DEDICATED or ATLAS_MANAGED (default)
		// only valid for Major version 8 and later
		// cluster must be SHARDED
		rootConfig = fmt.Sprintf(`
		  mongo_db_major_version = "8"
		  config_server_management_mode = %[1]q
		`, *configServerManagementMode)
	}
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "SHARDED"
			%[5]s

			replication_specs {
				num_shards = %[3]d
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = %[4]q
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 2
					}
					provider_name = "AZURE"
					priority      = 6
					region_name   = "US_EAST_2"
				}
			}
		}
	`, projectID, name, numShards, analyticsSize, rootConfig)) + dataSourcesTFOldSchema
}

func checkShardedOldSchemaMultiCloud(usePreviewProvider bool, name string, numShards int, analyticsSize string, verifyExternalID bool, configServerManagementMode *string) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider, resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider, resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider, resourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider, dataSourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider, dataSourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithPreviewProviderV2(usePreviewProvider, dataSourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
	}
	if verifyExternalID {
		additionalChecks = append(
			additionalChecks,
			acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, resourceName, "replication_specs.0.external_id"))
	}
	if configServerManagementMode != nil {
		additionalChecks = append(additionalChecks,
			acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, resourceName, "config_server_management_mode", *configServerManagementMode),
			acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, resourceName, "config_server_type"),
			acc.TestCheckResourceAttrPreviewProviderV2(usePreviewProvider, dataSourceName, "config_server_management_mode", *configServerManagementMode),
			acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, dataSourceName, "config_server_type"),
		)
	}

	return checkAggr(usePreviewProvider,
		[]string{"project_id", "replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"},
		map[string]string{
			"name":                           name,
			"replication_specs.0.num_shards": strconv.Itoa(numShards),
			"replication_specs.0.region_configs.0.analytics_specs.0.instance_size": analyticsSize,
		},
		additionalChecks...)
}

func configSingleProviderPaused(t *testing.T, usePreviewProvider bool, projectID, clusterName string, paused bool, instanceSize string) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			paused       = %[3]t
			cluster_type = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = %[4]q
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
		}
`, projectID, clusterName, paused, instanceSize)) + dataSourcesTFNewSchema
}

func checkSingleProviderPaused(usePreviewProvider bool, name string, paused bool) resource.TestCheckFunc {
	return checkAggr(usePreviewProvider,
		[]string{"project_id", "replication_specs.#", "replication_specs.0.region_configs.#"},
		map[string]string{
			"name":   name,
			"paused": strconv.FormatBool(paused)})
}

func configAdvanced(t *testing.T, usePreviewProvider bool, projectID, clusterName, mongoDBMajorVersion string, p20240530 *admin20240530.ClusterDescriptionProcessArgs, p *admin.ClusterDescriptionProcessArgs20240805) string {
	t.Helper()
	changeStreamOptionsStr := ""
	defaultMaxTimeStr := ""
	tlsCipherConfigModeStr := ""
	customOpensslCipherConfigTLS12Str := ""
	mongoDBMajorVersionStr := ""

	if p != nil {
		if p.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds != nil && p.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds != conversion.IntPtr(-1) {
			changeStreamOptionsStr = fmt.Sprintf(`change_stream_options_pre_and_post_images_expire_after_seconds = %[1]d`, *p.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds)
		}
		if p.DefaultMaxTimeMS != nil {
			defaultMaxTimeStr = fmt.Sprintf(`default_max_time_ms = %[1]d`, *p.DefaultMaxTimeMS)
		}
		if p.TlsCipherConfigMode != nil {
			tlsCipherConfigModeStr = fmt.Sprintf(`tls_cipher_config_mode = %[1]q`, *p.TlsCipherConfigMode)
			if p.CustomOpensslCipherConfigTls12 != nil && len(*p.CustomOpensslCipherConfigTls12) > 0 {
				customOpensslCipherConfigTLS12Str = fmt.Sprintf(
					`custom_openssl_cipher_config_tls12 = [%s]`,
					acc.JoinQuotedStrings(*p.CustomOpensslCipherConfigTls12),
				)
			}
		}
	}
	if mongoDBMajorVersion != "" {
		mongoDBMajorVersionStr = fmt.Sprintf(`mongo_db_major_version = %[1]q`, mongoDBMajorVersion)
	}

	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id             = %[1]q
			name                   = %[2]q
			cluster_type           = "REPLICASET"
			%[13]s

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}

			advanced_configuration  {
				fail_index_key_too_long              = %[3]t
				javascript_enabled                   = %[4]t
				minimum_enabled_tls_protocol         = %[5]q
				no_table_scan                        = %[6]t
				oplog_size_mb                        = %[7]d
				sample_size_bi_connector			 = %[8]d
				sample_refresh_interval_bi_connector = %[9]d
			    transaction_lifetime_limit_seconds   = %[10]d
			    %[11]s
				%[12]s
				%[14]s
				%[15]s
			}
		}
	`, projectID, clusterName,
		p20240530.GetFailIndexKeyTooLong(), p20240530.GetJavascriptEnabled(), p20240530.GetMinimumEnabledTlsProtocol(), p20240530.GetNoTableScan(),
		p20240530.GetOplogSizeMB(), p20240530.GetSampleSizeBIConnector(), p20240530.GetSampleRefreshIntervalBIConnector(), p20240530.GetTransactionLifetimeLimitSeconds(),
		changeStreamOptionsStr, defaultMaxTimeStr, mongoDBMajorVersionStr, tlsCipherConfigModeStr, customOpensslCipherConfigTLS12Str)) + dataSourcesTFNewSchema
}

func checkAdvanced(usePreviewProvider bool, name, tls string, processArgs *admin.ClusterDescriptionProcessArgs20240805) resource.TestCheckFunc {
	advancedConfig := map[string]string{
		"name": name,
		"advanced_configuration.0.minimum_enabled_tls_protocol":         tls,
		"advanced_configuration.0.fail_index_key_too_long":              "false",
		"advanced_configuration.0.javascript_enabled":                   "true",
		"advanced_configuration.0.no_table_scan":                        "false",
		"advanced_configuration.0.oplog_size_mb":                        "1000",
		"advanced_configuration.0.sample_refresh_interval_bi_connector": "310",
		"advanced_configuration.0.sample_size_bi_connector":             "110",
		"advanced_configuration.0.transaction_lifetime_limit_seconds":   "300",
	}

	if processArgs.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds != nil {
		advancedConfig["advanced_configuration.0.change_stream_options_pre_and_post_images_expire_after_seconds"] = strconv.Itoa(*processArgs.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds)
	}

	if processArgs.DefaultMaxTimeMS != nil {
		advancedConfig["advanced_configuration.0.default_max_time_ms"] = strconv.Itoa(*processArgs.DefaultMaxTimeMS)
	}

	if processArgs.TlsCipherConfigMode != nil && processArgs.CustomOpensslCipherConfigTls12 != nil {
		advancedConfig["advanced_configuration.0.tls_cipher_config_mode"] = "CUSTOM"
		advancedConfig["advanced_configuration.0.custom_openssl_cipher_config_tls12.#"] = strconv.Itoa(len(*processArgs.CustomOpensslCipherConfigTls12))
	} else {
		advancedConfig["advanced_configuration.0.tls_cipher_config_mode"] = "DEFAULT"
	}

	pluralChecks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, dataSourcePluralName, "results.#"),
		acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, dataSourcePluralName, "results.0.replication_specs.#"),
		acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, dataSourcePluralName, "results.0.name"),
	}

	return checkAggr(usePreviewProvider,
		[]string{"project_id", "replication_specs.#", "replication_specs.0.region_configs.#"},
		advancedConfig,
		pluralChecks...,
	)
}

func configAdvancedDefaultWrite(t *testing.T, usePreviewProvider bool, projectID, clusterName string, p *admin20240530.ClusterDescriptionProcessArgs) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id             = %[1]q
			name                   = %[2]q
			cluster_type           = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}

			advanced_configuration  {
				javascript_enabled                   = %[3]t
				minimum_enabled_tls_protocol         = %[4]q
				no_table_scan                        = %[5]t
				oplog_size_mb                        = %[6]d
				sample_size_bi_connector			 = %[7]d
				sample_refresh_interval_bi_connector = %[8]d
				default_read_concern                 = %[9]q
				default_write_concern                = %[10]q
			}
		}
	`, projectID, clusterName, p.GetJavascriptEnabled(), p.GetMinimumEnabledTlsProtocol(), p.GetNoTableScan(),
		p.GetOplogSizeMB(), p.GetSampleSizeBIConnector(), p.GetSampleRefreshIntervalBIConnector(), p.GetDefaultReadConcern(), p.GetDefaultWriteConcern())) + dataSourcesTFNewSchema
}

func checkAdvancedDefaultWrite(usePreviewProvider bool, name, writeConcern, tls string) resource.TestCheckFunc {
	pluralChecks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, dataSourcePluralName, "results.#"),
		acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, dataSourcePluralName, "results.0.replication_specs.#"),
		acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, dataSourcePluralName, "results.0.name"),
	}
	return checkAggr(usePreviewProvider,
		[]string{"project_id", "replication_specs.#", "replication_specs.0.region_configs.#"},
		map[string]string{
			"name": name,
			"advanced_configuration.0.minimum_enabled_tls_protocol":         tls,
			"advanced_configuration.0.default_write_concern":                writeConcern,
			"advanced_configuration.0.default_read_concern":                 "available",
			"advanced_configuration.0.fail_index_key_too_long":              "false",
			"advanced_configuration.0.javascript_enabled":                   "true",
			"advanced_configuration.0.no_table_scan":                        "false",
			"advanced_configuration.0.oplog_size_mb":                        "1000",
			"advanced_configuration.0.sample_refresh_interval_bi_connector": "310",
			"advanced_configuration.0.sample_size_bi_connector":             "110",
			"advanced_configuration.0.tls_cipher_config_mode":               "DEFAULT"},
		pluralChecks...)
}

func configReplicationSpecsAutoScaling(t *testing.T, usePreviewProvider bool, projectID, clusterName string, autoScalingSettings *admin.AdvancedAutoScalingSettings, elecInstanceSize string, elecDiskSizeGB, analyticsNodeCount int) string {
	t.Helper()
	lifecycleIgnoreChanges := ""
	autoScalingCompute := autoScalingSettings.GetCompute()
	if autoScalingCompute.GetEnabled() {
		lifecycleIgnoreChanges = `
		lifecycle {
			ignore_changes = [
				replication_specs.0.region_configs.0.electable_specs.0.instance_size,
				replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb
			]
        }`
	}

	autoScalingBlock := ""
	if autoScalingSettings != nil {
		autoScalingBlock = fmt.Sprintf(`auto_scaling {
			compute_enabled = %t
			disk_gb_enabled = %t
			compute_max_instance_size = %q
		}`, autoScalingSettings.Compute.GetEnabled(), autoScalingSettings.DiskGB.GetEnabled(), autoScalingSettings.Compute.GetMaxInstanceSize())
	}

	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id             = %[1]q
			name                   = %[2]q
			cluster_type           = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = %[3]q
						disk_size_gb = %[4]d
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = %[5]d
					}
					%[6]s
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
			advanced_configuration  {
				oplog_min_retention_hours = 5.5
			}
			%[7]s
		}		
	`, projectID, clusterName, elecInstanceSize, elecDiskSizeGB, analyticsNodeCount, autoScalingBlock, lifecycleIgnoreChanges))
}

func configReplicationSpecsAnalyticsAutoScaling(t *testing.T, usePreviewProvider bool, projectID, clusterName string, analyticsAutoScalingSettings *admin.AdvancedAutoScalingSettings, analyticsNodeCount int) string {
	t.Helper()

	analyticsAutoScalingBlock := ""
	if analyticsAutoScalingSettings != nil {
		analyticsAutoScalingBlock = fmt.Sprintf(`
				analytics_auto_scaling {
					compute_enabled = %t
					disk_gb_enabled = %t
					compute_max_instance_size = %q
				}`, analyticsAutoScalingSettings.Compute.GetEnabled(), analyticsAutoScalingSettings.DiskGB.GetEnabled(), analyticsAutoScalingSettings.Compute.GetMaxInstanceSize())
	}

	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id             = %[1]q
			name                   = %[2]q
			cluster_type           = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = %[3]d
					}
					%[4]s
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
		}
	`, projectID, clusterName, analyticsNodeCount, analyticsAutoScalingBlock))
}

func configGeoShardedOldSchema(t *testing.T, usePreviewProvider bool, projectID, name string, numShardsFirstZone, numShardsSecondZone int, selfManagedSharding bool) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = %[1]q
			name = %[2]q
			backup_enabled = false
			mongo_db_major_version = "7.0"
			cluster_type   = "GEOSHARDED"
			global_cluster_self_managed_sharding = %[5]t
			disk_size_gb  = 60

			replication_specs {
				zone_name  = "zone n1"
				num_shards = %[3]d

				region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				analytics_specs {
					instance_size = "M10"
					node_count    = 0
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				}
			}

			replication_specs {
				zone_name  = "zone n2"
				num_shards = %[4]d

				region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				analytics_specs {
					instance_size = "M10"
					node_count    = 0
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
				}
			}
		}

	`, projectID, name, numShardsFirstZone, numShardsSecondZone, selfManagedSharding)) + dataSourcesTFOldSchema
}

func checkGeoShardedOldSchema(usePreviewProvider bool, name string, numShardsFirstZone, numShardsSecondZone int, isLatestProviderVersion, verifyExternalID bool) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{}

	if verifyExternalID {
		additionalChecks = append(additionalChecks, acc.TestCheckResourceAttrSetPreviewProviderV2(usePreviewProvider, resourceName, "replication_specs.0.external_id"))
	}

	if isLatestProviderVersion { // checks that will not apply if doing migration test with older version
		additionalChecks = append(additionalChecks, checkAggr(usePreviewProvider,
			[]string{"replication_specs.0.zone_id", "replication_specs.0.zone_id"},
			map[string]string{
				"replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb": "60",
				"replication_specs.0.region_configs.0.analytics_specs.0.disk_size_gb": "60",
			}))
	}

	return checkAggr(usePreviewProvider,
		[]string{"project_id", "replication_specs.0.id", "replication_specs.1.id"},
		map[string]string{
			"name":                           name,
			"disk_size_gb":                   "60",
			"replication_specs.0.num_shards": strconv.Itoa(numShardsFirstZone),
			"replication_specs.1.num_shards": strconv.Itoa(numShardsSecondZone),
		},
		additionalChecks...,
	)
}

func configShardedOldSchemaDiskSizeGBElectableLevel(t *testing.T, usePreviewProvider bool, projectID, name string, diskSizeGB int) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = %[1]q
			name = %[2]q
			backup_enabled = false
			mongo_db_major_version = "7.0"
			cluster_type   = "SHARDED"

			replication_specs {
				num_shards = 2

				region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
					disk_size_gb  = %[3]d
				}
				analytics_specs {
					instance_size = "M10"
					node_count    = 0
					disk_size_gb  = %[3]d
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				}
			}
		}
	`, projectID, name, diskSizeGB)) + dataSourcesTFOldSchema
}

func checkShardedOldSchemaDiskSizeGBElectableLevel(usePreviewProvider bool, diskSizeGB int) resource.TestCheckFunc {
	return checkAggr(usePreviewProvider,
		[]string{},
		map[string]string{
			"replication_specs.0.num_shards": "2",
			"disk_size_gb":                   fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.analytics_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
		})
}

func configShardedNewSchema(t *testing.T, usePreviewProvider bool, orgID, projectName, name string, diskSizeGB int, firstInstanceSize, lastInstanceSize string, firstDiskIOPS, lastDiskIOPS *int, includeMiddleSpec, increaseDiskSizeShard2 bool) string {
	t.Helper()
	var thirdReplicationSpec string
	var diskSizeGBShard2 = diskSizeGB
	if increaseDiskSizeShard2 {
		diskSizeGBShard2 = diskSizeGB + 10
	}
	if includeMiddleSpec {
		thirdReplicationSpec = fmt.Sprintf(`
			replication_specs {
				region_configs {
					electable_specs {
						instance_size = %[1]q
						node_count    = 3
						disk_size_gb  = %[2]d
					}
					analytics_specs {
						instance_size = %[1]q
						node_count    = 1
						disk_size_gb  = %[2]d
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
			}
		`, firstInstanceSize, diskSizeGB)
	}
	var firstDiskIOPSAttrs string
	if firstDiskIOPS != nil {
		firstDiskIOPSAttrs = fmt.Sprintf(`
			disk_iops = %d
			ebs_volume_type = "PROVISIONED"
		`, *firstDiskIOPS)
	}
	var lastDiskIOPSAttrs string
	if lastDiskIOPS != nil {
		lastDiskIOPSAttrs = fmt.Sprintf(`
			disk_iops = %d
			ebs_volume_type = "PROVISIONED"
		`, *lastDiskIOPS)
	}
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			cluster_type   = "SHARDED"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = %[4]q
						node_count    = 3
						disk_size_gb  = %[9]d
						%[6]s
					}
					analytics_specs {
						instance_size = %[4]q
						node_count    = 1
						disk_size_gb  = %[9]d
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
			}

			%[8]s

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = %[5]q
						node_count    = 3
						disk_size_gb  = %[10]d
						%[7]s
					}
					analytics_specs {
						instance_size = %[5]q
						node_count    = 1
						disk_size_gb  = %[10]d
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
			use_replication_spec_per_shard = true
		}

		data "mongodbatlas_advanced_clusters" "test-replication-specs-per-shard-false" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			use_replication_spec_per_shard = false
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			use_replication_spec_per_shard = true
		}
	`, orgID, projectName, name, firstInstanceSize, lastInstanceSize, firstDiskIOPSAttrs, lastDiskIOPSAttrs, thirdReplicationSpec, diskSizeGB, diskSizeGBShard2))
}

func checkShardedNewSchema(usePreviewProvider bool, diskSizeGB int, firstInstanceSize, lastInstanceSize string, firstDiskIops, lastDiskIops *int, isAsymmetricCluster, includeMiddleSpec bool) resource.TestCheckFunc {
	amtOfReplicationSpecs := 2
	if includeMiddleSpec {
		amtOfReplicationSpecs = 3
	}

	lastSpecIndex := 1
	if includeMiddleSpec {
		lastSpecIndex = 2
	}

	clusterChecks := map[string]string{
		"disk_size_gb":        fmt.Sprintf("%d", diskSizeGB),
		"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
		"replication_specs.0.region_configs.0.electable_specs.0.instance_size":                              firstInstanceSize,
		fmt.Sprintf("replication_specs.%d.region_configs.0.electable_specs.0.instance_size", lastSpecIndex): lastInstanceSize,
		"replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb":                               fmt.Sprintf("%d", diskSizeGB),
		fmt.Sprintf("replication_specs.%d.region_configs.0.electable_specs.0.disk_size_gb", lastSpecIndex):  fmt.Sprintf("%d", diskSizeGB),
		"replication_specs.0.region_configs.0.analytics_specs.0.disk_size_gb":                               fmt.Sprintf("%d", diskSizeGB),
		fmt.Sprintf("replication_specs.%d.region_configs.0.analytics_specs.0.disk_size_gb", lastSpecIndex):  fmt.Sprintf("%d", diskSizeGB),
	}
	if firstDiskIops != nil {
		clusterChecks["replication_specs.0.region_configs.0.electable_specs.0.disk_iops"] = fmt.Sprintf("%d", *firstDiskIops)
	}
	if lastDiskIops != nil {
		clusterChecks[fmt.Sprintf("replication_specs.%d.region_configs.0.electable_specs.0.disk_iops", lastSpecIndex)] = fmt.Sprintf("%d", *lastDiskIops)
	}

	// plural data source checks
	pluralChecks := acc.AddAttrSetChecksPreviewProviderV2(usePreviewProvider, dataSourcePluralName, nil,
		[]string{"results.#", "results.0.replication_specs.#", "results.0.replication_specs.0.region_configs.#", "results.0.name", "results.0.termination_protection_enabled", "results.0.global_cluster_self_managed_sharding"}...)

	pluralChecks = acc.AddAttrChecksPrefixPreviewProviderV2(usePreviewProvider, dataSourcePluralName, pluralChecks, clusterChecks, "results.0")
	if isAsymmetricCluster {
		pluralChecks = append(pluralChecks, checkAggr(usePreviewProvider, []string{}, map[string]string{
			"replication_specs.0.id": "",
			"replication_specs.1.id": "",
		}))
		pluralChecks = acc.AddAttrChecksPreviewProviderV2(usePreviewProvider, dataSourcePluralName, pluralChecks, map[string]string{
			"results.0.replication_specs.0.id": "",
			"results.0.replication_specs.1.id": "",
		})
	} else {
		pluralChecks = append(pluralChecks, checkAggr(usePreviewProvider, []string{"replication_specs.0.id", "replication_specs.1.id"}, map[string]string{}))
		pluralChecks = acc.AddAttrSetChecksPreviewProviderV2(usePreviewProvider, dataSourcePluralName, pluralChecks, "results.0.replication_specs.0.id", "results.0.replication_specs.1.id")
	}
	return checkAggr(usePreviewProvider,
		[]string{"replication_specs.0.external_id", "replication_specs.0.zone_id", "replication_specs.1.external_id", "replication_specs.1.zone_id"},
		clusterChecks,
		pluralChecks...,
	)
}

func configGeoShardedNewSchema(t *testing.T, usePreviewProvider bool, projectID, name string, includeThirdShardInFirstZone bool) string {
	t.Helper()
	var thirdReplicationSpec string
	if includeThirdShardInFirstZone {
		thirdReplicationSpec = `
			replication_specs {
				zone_name  = "zone n1"
				region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				}
			}
		`
	}
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = %[1]q
			name = %[2]q
			backup_enabled = false
			mongo_db_major_version = "7.0"
			cluster_type   = "GEOSHARDED"
			replication_specs {
				zone_name  = "zone n1"
				region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				}
			}
			%[3]s
			replication_specs {
				zone_name  = "zone n2"
				region_configs {
				electable_specs {
					instance_size = "M20"
					node_count    = 3
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
				}
			}
    	}
	`, projectID, name, thirdReplicationSpec)) + dataSourcesTFNewSchema
}

func checkGeoShardedNewSchema(usePreviewProvider, includeThirdShardInFirstZone bool) resource.TestCheckFunc {
	var amtOfReplicationSpecs int
	if includeThirdShardInFirstZone {
		amtOfReplicationSpecs = 3
	} else {
		amtOfReplicationSpecs = 2
	}
	clusterChecks := map[string]string{
		"replication_specs.#":                fmt.Sprintf("%d", amtOfReplicationSpecs),
		"replication_specs.0.container_id.%": "1",
		"replication_specs.1.container_id.%": "1",
	}
	return checkAggr(usePreviewProvider, []string{}, clusterChecks)
}

func configShardedTransitionOldToNewSchema(t *testing.T, usePreviewProvider bool, projectID, name string, useNewSchema, autoscaling bool) string {
	t.Helper()
	var numShardsStr string
	if !useNewSchema {
		numShardsStr = `num_shards = 2`
	}
	var autoscalingStr string
	if autoscaling {
		autoscalingStr = `auto_scaling {
			compute_enabled = true
			disk_gb_enabled = true
			compute_max_instance_size = "M20"
		}`
	}
	replicationSpec := fmt.Sprintf(`
		replication_specs {
			%[1]s
			region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				analytics_specs {
					instance_size = "M10"
					node_count    = 1
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
				%[2]s
			}
		}
	`, numShardsStr, autoscalingStr)

	var replicationSpecs string
	if useNewSchema {
		replicationSpecs = fmt.Sprintf(`
			%[1]s
			%[1]s
		`, replicationSpec)
	} else {
		replicationSpecs = replicationSpec
	}

	var dataSources = dataSourcesTFOldSchema
	if useNewSchema {
		dataSources = dataSourcesTFNewSchema
	}

	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = %[1]q
			name = %[2]q
			backup_enabled = false
			cluster_type   = "SHARDED"

			%[3]s
		}

	`, projectID, name, replicationSpecs)) + dataSources
}

func checkShardedTransitionOldToNewSchema(usePreviewProvider, useNewSchema bool) resource.TestCheckFunc {
	var amtOfReplicationSpecs int
	if useNewSchema {
		amtOfReplicationSpecs = 2
	} else {
		amtOfReplicationSpecs = 1
	}
	var checksForNewSchema []resource.TestCheckFunc
	if useNewSchema {
		checksForNewSchema = []resource.TestCheckFunc{
			checkAggr(usePreviewProvider, []string{"replication_specs.1.id", "replication_specs.0.external_id", "replication_specs.1.external_id"},
				map[string]string{
					"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
					"replication_specs.1.region_configs.0.electable_specs.0.instance_size": "M10",
					"replication_specs.1.region_configs.0.analytics_specs.0.instance_size": "M10",
				}),
		}
	}

	return checkAggr(usePreviewProvider,
		[]string{"replication_specs.0.id"},
		map[string]string{
			"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
			"replication_specs.0.region_configs.0.electable_specs.0.instance_size": "M10",
			"replication_specs.0.region_configs.0.analytics_specs.0.instance_size": "M10",
		},
		checksForNewSchema...,
	)
}

func configGeoShardedTransitionOldToNewSchema(t *testing.T, usePreviewProvider bool, projectID, name string, useNewSchema bool) string {
	t.Helper()
	var numShardsStr string
	if !useNewSchema {
		numShardsStr = `num_shards = 2`
	}
	replicationSpec := `
		replication_specs {
			%[1]s
			region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				analytics_specs {
					instance_size = "M10"
					node_count    = 1
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = %[2]q
			}
			zone_name = %[3]q
		}
	`

	var replicationSpecs string
	if !useNewSchema {
		replicationSpecs = fmt.Sprintf(`
			%[1]s
			%[2]s
		`, fmt.Sprintf(replicationSpec, numShardsStr, "US_EAST_1", "zone 1"), fmt.Sprintf(replicationSpec, numShardsStr, "EU_WEST_1", "zone 2"))
	} else {
		replicationSpecs = fmt.Sprintf(`
			%[1]s
			%[2]s
			%[3]s
			%[4]s
		`, fmt.Sprintf(replicationSpec, numShardsStr, "US_EAST_1", "zone 1"), fmt.Sprintf(replicationSpec, numShardsStr, "US_EAST_1", "zone 1"),
			fmt.Sprintf(replicationSpec, numShardsStr, "EU_WEST_1", "zone 2"), fmt.Sprintf(replicationSpec, numShardsStr, "EU_WEST_1", "zone 2"))
	}

	var dataSources = dataSourcesTFOldSchema
	if useNewSchema {
		dataSources = dataSourcesTFNewSchema
	}

	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = %[1]q
			name = %[2]q
			backup_enabled = false
			cluster_type   = "GEOSHARDED"

			%[3]s
		}
	`, projectID, name, replicationSpecs)) + dataSources
}

func checkGeoShardedTransitionOldToNewSchema(usePreviewProvider, useNewSchema bool) resource.TestCheckFunc {
	if useNewSchema {
		return checkAggr(usePreviewProvider,
			[]string{"replication_specs.0.id", "replication_specs.1.id", "replication_specs.2.id", "replication_specs.3.id",
				"replication_specs.0.external_id", "replication_specs.1.external_id", "replication_specs.2.external_id", "replication_specs.3.external_id",
			},
			map[string]string{
				"replication_specs.#":           "4",
				"replication_specs.0.zone_name": "zone 1",
				"replication_specs.1.zone_name": "zone 1",
				"replication_specs.2.zone_name": "zone 2",
				"replication_specs.3.zone_name": "zone 2",
			},
		)
	}
	return checkAggr(usePreviewProvider,
		[]string{"replication_specs.0.id", "replication_specs.1.id"},
		map[string]string{
			"replication_specs.#":           "2",
			"replication_specs.0.zone_name": "zone 1",
			"replication_specs.1.zone_name": "zone 2",
		},
	)
}

func configReplicaSetScalingStrategyAndRedactClientLogData(t *testing.T, usePreviewProvider bool, orgID, projectName, name, replicaSetScalingStrategy string, redactClientLogData bool) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			cluster_type   = "SHARDED"
			replica_set_scaling_strategy = %[4]q
			redact_client_log_data = %[5]t

			replication_specs {
				region_configs {
					electable_specs {
						instance_size ="M10"
						node_count    = 3
						disk_size_gb  = 10
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
						disk_size_gb  = 10
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
			}
		}
	`, orgID, projectName, name, replicaSetScalingStrategy, redactClientLogData)) + dataSourcesTFNewSchema
}

func configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(t *testing.T, usePreviewProvider bool, orgID, projectName, name, replicaSetScalingStrategy string, redactClientLogData bool) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			cluster_type   = "SHARDED"
			replica_set_scaling_strategy = %[4]q
			redact_client_log_data = %[5]t

			replication_specs {
				num_shards = 2
				region_configs {
					electable_specs {
						instance_size ="M10"
						node_count    = 3
						disk_size_gb  = 10
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
						disk_size_gb  = 10
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
			}
		}
	`, orgID, projectName, name, replicaSetScalingStrategy, redactClientLogData)) + dataSourcesTFOldSchema
}

func checkReplicaSetScalingStrategyAndRedactClientLogData(usePreviewProvider bool, replicaSetScalingStrategy string, redactClientLogData bool) resource.TestCheckFunc {
	clusterChecks := map[string]string{
		"replica_set_scaling_strategy": replicaSetScalingStrategy,
		"redact_client_log_data":       strconv.FormatBool(redactClientLogData),
	}

	// plural data source checks
	pluralChecks := acc.AddAttrSetChecksPreviewProviderV2(usePreviewProvider, dataSourcePluralName, nil,
		[]string{"results.#", "results.0.replica_set_scaling_strategy", "results.0.redact_client_log_data"}...)

	return checkAggr(usePreviewProvider,
		[]string{},
		clusterChecks,
		pluralChecks...,
	)
}

func configPriority(t *testing.T, usePreviewProvider bool, projectID, clusterName string, oldSchema, swapPriorities bool) string {
	t.Helper()
	const (
		config7 = `
			region_configs {
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				electable_specs {
					node_count    = 2
					instance_size = "M10"
				}
			}
		`
		config6 = `
			region_configs {
				provider_name = "AWS"
				priority      = 6
				region_name   = "US_WEST_2"
				electable_specs {
					node_count    = 1
					instance_size = "M10"
				}
			}
		`
	)
	strType, strNumShards, strConfigs := "REPLICASET", "", config7+config6
	if oldSchema {
		strType = "SHARDED"
		strNumShards = "num_shards = 2"
	}
	if swapPriorities {
		strConfigs = config6 + config7
	}

	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type   = %[3]q
			backup_enabled = false
			
			replication_specs {
 					%[4]s
 					%[5]s
			}
		}
	`, projectID, clusterName, strType, strNumShards, strConfigs))
}

func configBiConnectorConfig(t *testing.T, usePreviewProvider bool, projectID, name string, enabled bool) string {
	t.Helper()
	additionalConfig := `
		bi_connector_config {
			enabled = false
		}	
	`
	if enabled {
		additionalConfig = `
			bi_connector_config {
				enabled         = true
				read_preference = "secondary"
			}	
		`
	}

	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, usePreviewProvider, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}

			%[3]s
		}
	`, projectID, name, additionalConfig)) + dataSourcesTFOldSchema
}

func checkTenantBiConnectorConfig(usePreviewProvider bool, projectID, name string, enabled bool) resource.TestCheckFunc {
	attrsMap := map[string]string{
		"project_id": projectID,
		"name":       name,
	}
	if enabled {
		attrsMap["bi_connector_config.0.enabled"] = "true"
		attrsMap["bi_connector_config.0.read_preference"] = "secondary"
	} else {
		attrsMap["bi_connector_config.0.enabled"] = "false"
	}
	return checkAggr(usePreviewProvider, nil, attrsMap)
}

func configFCVPinning(t *testing.T, orgID, projectName, clusterName string, pinningExpirationDate *string, mongoDBMajorVersion string) string {
	t.Helper()
	var pinnedFCVAttr string
	if pinningExpirationDate != nil {
		pinnedFCVAttr = fmt.Sprintf(`
		pinned_fcv {
    		expiration_date = %q
  		}
		`, *pinningExpirationDate)
	}

	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, true, fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id = %[1]q
			name   = %[2]q
		}
		
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = mongodbatlas_project.test.id
			name         = %[3]q

			cluster_type = "REPLICASET"

			mongo_db_major_version = %[4]q

			%[5]s

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
		}

	`, orgID, projectName, clusterName, mongoDBMajorVersion, pinnedFCVAttr)) + dataSourcesTFNewSchema
}

func configFlexCluster(t *testing.T, projectID, clusterName, providerName, region, zoneName string, withTags bool) string {
	t.Helper()
	zoneNameLine := ""
	if zoneName != "" {
		zoneNameLine = fmt.Sprintf("zone_name = %q", zoneName)
	}
	tags := ""
	if withTags {
		tags = `
			tags {
				key = "testKey"
				value = "testValue"
			}`
	}
	return acc.ConvertAdvancedClusterToPreviewProviderV2(t, true, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"
			replication_specs {
				region_configs {
					provider_name = "FLEX"
					backing_provider_name = %[3]q
					region_name = %[4]q
					priority      = 7
				}
				%[5]s
			}
			%[6]s
			termination_protection_enabled = false
		}
	`, projectID, clusterName, providerName, region, zoneNameLine, tags)+dataSourcesTFOldSchema+
		strings.ReplaceAll(acc.FlexDataSource, "mongodbatlas_flex_cluster.", "mongodbatlas_advanced_cluster."))
}

func TestAccClusterFlexCluster_basic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyFlexCluster,
		Steps: []resource.TestStep{
			{
				Config: configFlexCluster(t, projectID, clusterName, "AWS", "US_EAST_1", "", false),
				Check:  checkFlexClusterConfig(projectID, clusterName, "AWS", "US_EAST_1", false),
			},
			{
				Config: configFlexCluster(t, projectID, clusterName, "AWS", "US_EAST_1", "", true),
				Check:  checkFlexClusterConfig(projectID, clusterName, "AWS", "US_EAST_1", true),
			},
			acc.TestStepImportCluster(resourceName),
			{
				Config:      configFlexCluster(t, projectID, clusterName, "AWS", "US_EAST_2", "", true),
				ExpectError: regexp.MustCompile("flex cluster update is not supported except for tags and termination_protection_enabled fields"),
			},
		},
	})
}

func checkFlexClusterConfig(projectID, clusterName, providerName, region string, tagsCheck bool) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{acc.CheckExistsFlexCluster()}
	attrMapAdvCluster := map[string]string{
		"name":                                 clusterName,
		"cluster_type":                         "REPLICASET",
		"termination_protection_enabled":       "false",
		"replication_specs.#":                  "1",
		"replication_specs.0.region_configs.#": "1",
		"replication_specs.0.region_configs.0.provider_name":         "FLEX",
		"replication_specs.0.region_configs.0.backing_provider_name": providerName,
		"replication_specs.0.region_configs.0.region_name":           region,
	}
	attrSetAdvCluster := []string{
		"backup_enabled",
		"connection_strings.0.standard",
		"connection_strings.0.standard_srv",
		"create_date",
		"mongo_db_version",
		"state_name",
		"version_release_system",
	}
	attrMapFlex := map[string]string{
		"project_id":                     projectID,
		"name":                           clusterName,
		"termination_protection_enabled": "false",
	}
	attrSetFlex := []string{
		"backup_settings.enabled",
		"cluster_type",
		"connection_strings.standard",
		"create_date",
		"id",
		"mongo_db_version",
		"state_name",
		"version_release_system",
		"provider_settings.provider_name",
	}
	if tagsCheck {
		attrMapFlex["tags.testKey"] = "testValue"
		tagsMap := map[string]string{"key": "testKey", "value": "testValue"}
		tagsCheck := checkKeyValueBlocks(true, true, "tags", tagsMap)
		checks = append(checks, tagsCheck)
	}
	pluralMap := map[string]string{
		"project_id": projectID,
		"results.#":  "1",
	}
	checks = acc.AddAttrChecks(acc.FlexDataSourceName, checks, attrMapFlex)
	checks = acc.AddAttrSetChecks(acc.FlexDataSourceName, checks, attrSetFlex...)
	checks = acc.AddAttrChecks(acc.FlexDataSourcePluralName, checks, pluralMap)
	checks = acc.AddAttrChecksPrefix(acc.FlexDataSourcePluralName, checks, attrMapFlex, "results.0")
	checks = acc.AddAttrSetChecksPrefix(acc.FlexDataSourcePluralName, checks, attrSetFlex, "results.0")
	checks = acc.AddAttrChecks(dataSourcePluralName, checks, pluralMap)
	ds := conversion.StringPtr(dataSourceName)
	dsp := conversion.StringPtr(dataSourcePluralName)
	return acc.CheckRSAndDSPreviewProviderV2(true, resourceName, ds, dsp, attrSetAdvCluster, attrMapAdvCluster, checks...)
}
