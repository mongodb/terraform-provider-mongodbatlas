package acc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var standardClusterResource = `
resource "mongodbatlas_advanced_cluster" "cluster_info" {
  backup_enabled = false
  cluster_type   = "REPLICASET"
  name           = "my-name"
  pit_enabled    = false
  project_id     = "project"

  replication_specs = [{
    region_configs = [{
      auto_scaling = {
        disk_gb_enabled = false
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority        = 7
      provider_name   = "AWS"
      read_only_specs = {}
      region_name     = "US_WEST_2"
    }]
    zone_name = "Zone 1"
  }]

}
`
var overrideClusterResource = `
resource "mongodbatlas_advanced_cluster" "cluster_info" {
  project_id             = mongodbatlas_project.test.id
  backup_enabled         = true
  cluster_type           = "GEOSHARDED"
  mongo_db_major_version = "6.0"
  name                   = "my-name"
  pit_enabled            = true
  retain_backups_enabled = true

  advanced_configuration = {
    oplog_min_retention_hours = 8
  }

  replication_specs = [{
    region_configs = [{
      auto_scaling = {
        disk_gb_enabled = false
      }
      electable_specs = {
        ebs_volume_type = "STANDARD"
        instance_size   = "M30"
        node_count      = 30
      }
      priority        = 7
      provider_name   = "AZURE"
      read_only_specs = {}
      region_name     = "MY_REGION_1"
    }]
    zone_name = "Zone X"
  }]

}
`

var dependsOnClusterResource = `
resource "mongodbatlas_advanced_cluster" "cluster_info" {
  backup_enabled = false
  cluster_type   = "REPLICASET"
  name           = "my-name"
  pit_enabled    = false
  project_id     = "project"

  replication_specs = [{
    region_configs = [{
      auto_scaling = {
        disk_gb_enabled = false
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority        = 7
      provider_name   = "AWS"
      read_only_specs = {}
      region_name     = "US_WEST_2"
    }]
    zone_name = "Zone 1"
  }]

  depends_on = [mongodbatlas_project.project_execution]
}
`
var dependsOnMultiResource = `
resource "mongodbatlas_advanced_cluster" "cluster_info" {
  backup_enabled = false
  cluster_type   = "REPLICASET"
  name           = "my-name"
  pit_enabled    = false
  project_id     = "project"

  replication_specs = [{
    region_configs = [{
      auto_scaling = {
        disk_gb_enabled = false
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority        = 7
      provider_name   = "AWS"
      read_only_specs = {}
      region_name     = "US_WEST_2"
    }]
    zone_name = "Zone 1"
  }]

  depends_on = [mongodbatlas_private_endpoint_regional_mode.atlasrm, mongodbatlas_privatelink_endpoint_service.atlasple]
}
`
var twoReplicationSpecs = `
resource "mongodbatlas_advanced_cluster" "cluster_info" {
  backup_enabled = false
  cluster_type   = "REPLICASET"
  name           = "my-name"
  pit_enabled    = false
  project_id     = "project"

  replication_specs = [{
    region_configs = [{
      auto_scaling = {
        disk_gb_enabled = false
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority        = 7
      provider_name   = "AWS"
      read_only_specs = {}
      region_name     = "US_WEST_1"
    }]
    zone_name = "Zone 1"
    }, {
    region_configs = [{
      auto_scaling = {
        disk_gb_enabled = false
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority        = 7
      provider_name   = "AWS"
      read_only_specs = {}
      region_name     = "EU_WEST_2"
    }]
    zone_name = "Zone 2"
  }]

}
`
var twoRegionConfigs = `
resource "mongodbatlas_advanced_cluster" "cluster_info" {
  backup_enabled = false
  cluster_type   = "REPLICASET"
  name           = "my-name"
  pit_enabled    = false
  project_id     = "project"

  replication_specs = [{
    region_configs = [{
      auto_scaling = {
        disk_gb_enabled = false
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority        = 7
      provider_name   = "AWS"
      read_only_specs = {}
      region_name     = "US_WEST_1"
      }, {
      auto_scaling = {
        disk_gb_enabled = false
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority        = 7
      provider_name   = "AWS"
      read_only_specs = {}
      region_name     = "EU_WEST_1"
    }]
    zone_name = "Zone 1"
  }]

}
`

var autoScalingDiskEnabled = `
resource "mongodbatlas_advanced_cluster" "cluster_info" {
  backup_enabled = false
  cluster_type   = "REPLICASET"
  name           = "my-name"
  pit_enabled    = false
  project_id     = "project"

  replication_specs = [{
    region_configs = [{
      auto_scaling = {
        disk_gb_enabled = true
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
      priority        = 7
      provider_name   = "AWS"
      read_only_specs = {}
      region_name     = "US_WEST_2"
    }]
    zone_name = "Zone 1"
  }]
  tags = {
    ArchiveTest = "true"
    Owner       = "test"
  }

}
`
var readOnlyAndPriority = `
resource "mongodbatlas_advanced_cluster" "cluster_info" {
  backup_enabled = false
  cluster_type   = "REPLICASET"
  name           = "my-name"
  pit_enabled    = false
  project_id     = "project"

  replication_specs = [{
    region_configs = [{
      auto_scaling = {
        disk_gb_enabled = false
      }
      electable_specs = {
        instance_size = "M10"
        node_count    = 5
      }
      priority      = 5
      provider_name = "AWS"
      read_only_specs = {
        instance_size = "M10"
        node_count    = 1
      }
      region_name = "US_EAST_1"
    }]
    zone_name = "Zone 1"
  }]

}
`

func Test_ClusterResourceHcl(t *testing.T) {
	var (
		clusterName = "my-name"
		testCases   = map[string]struct {
			expected string
			req      acc.ClusterRequest
		}{
			"defaults": {
				standardClusterResource,
				acc.ClusterRequest{ClusterName: clusterName},
			},
			"dependsOn": {
				dependsOnClusterResource,
				acc.ClusterRequest{ClusterName: clusterName, ResourceDependencyName: "mongodbatlas_project.project_execution"},
			},
			"dependsOnMulti": {
				dependsOnMultiResource,
				acc.ClusterRequest{ClusterName: clusterName, ResourceDependencyName: "mongodbatlas_private_endpoint_regional_mode.atlasrm, mongodbatlas_privatelink_endpoint_service.atlasple"},
			},
			"twoReplicationSpecs": {
				twoReplicationSpecs,
				acc.ClusterRequest{ClusterName: clusterName, ReplicationSpecs: []acc.ReplicationSpecRequest{
					{Region: "US_WEST_1", ZoneName: "Zone 1"},
					{Region: "EU_WEST_2", ZoneName: "Zone 2"},
				}},
			},
			"overrideClusterResource": {
				overrideClusterResource,
				acc.ClusterRequest{
					ProjectID:            "mongodbatlas_project.test.id",
					ClusterName:          clusterName,
					Geosharded:           true,
					CloudBackup:          true,
					MongoDBMajorVersion:  "6.0",
					RetainBackupsEnabled: true,
					ReplicationSpecs: []acc.ReplicationSpecRequest{
						{Region: "MY_REGION_1", ZoneName: "Zone X", InstanceSize: "M30", NodeCount: 30, ProviderName: constant.AZURE, EbsVolumeType: "STANDARD"},
					},
					PitEnabled: true,
					AdvancedConfiguration: map[string]any{
						acc.ClusterAdvConfigOplogMinRetentionHours: 8,
					},
				},
			},
			"twoRegionConfigs": {
				twoRegionConfigs,
				acc.ClusterRequest{ClusterName: clusterName, ReplicationSpecs: []acc.ReplicationSpecRequest{
					{
						Region:             "US_WEST_1",
						InstanceSize:       "M10",
						NodeCount:          3,
						ExtraRegionConfigs: []acc.ReplicationSpecRequest{{Region: "EU_WEST_1", InstanceSize: "M10", NodeCount: 3, ProviderName: constant.AWS}},
					},
				},
				},
			},
			"autoScalingDiskEnabled": {
				autoScalingDiskEnabled,
				acc.ClusterRequest{ClusterName: clusterName, Tags: map[string]string{
					"ArchiveTest": "true", "Owner": "test",
				}, ReplicationSpecs: []acc.ReplicationSpecRequest{
					{AutoScalingDiskGbEnabled: true},
				}},
			},
			"readOnlyAndPriority": {
				readOnlyAndPriority,
				acc.ClusterRequest{
					ClusterName: clusterName,
					ReplicationSpecs: []acc.ReplicationSpecRequest{
						{Priority: 5, NodeCount: 5, Region: "US_EAST_1", NodeCountReadOnly: 1},
					}},
			},
		}
	)
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req := tc.req
			if req.ProjectID == "" {
				req.ProjectID = "project"
			}
			config, actualClusterName, actualResourceName, err := acc.ClusterResourceHcl(&req)
			require.NoError(t, err)
			assert.Equal(t, "mongodbatlas_advanced_cluster.cluster_info", actualResourceName)
			assert.Equal(t, clusterName, actualClusterName)
			assert.Equal(t, tc.expected, config)
		})
	}
}

var expectedDatasource = `
data "mongodbatlas_advanced_cluster" "cluster_info" {
  name       = "my-datasource-cluster"
  project_id = "datasource-project"
}
`

func Test_ClusterDatasourceHcl(t *testing.T) {
	expectedClusterName := "my-datasource-cluster"
	config, clusterName, resourceName, err := acc.ClusterDatasourceHcl(&acc.ClusterRequest{
		ClusterName: expectedClusterName,
		ProjectID:   "datasource-project",
	})
	require.NoError(t, err)
	assert.Equal(t, "data.mongodbatlas_advanced_cluster.cluster_info", resourceName)
	assert.Equal(t, expectedClusterName, clusterName)
	assert.Equal(t, expectedDatasource, config)
}
