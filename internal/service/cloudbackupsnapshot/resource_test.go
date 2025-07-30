package cloudbackupsnapshot_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName   = "mongodbatlas_cloud_backup_snapshot.test"
	dataSourceName = "data.mongodbatlas_cloud_backup_snapshot.test"
)

func TestAccBackupRSCloudBackupSnapshot_basic(t *testing.T) {
	var (
		dataSourcePluralSimpleName     = "data.mongodbatlas_cloud_backup_snapshots.test"
		dataSourcePluralPaginationName = "data.mongodbatlas_cloud_backup_snapshots.pagination"
		clusterInfo                    = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true})
		description                    = "My description in my cluster"
		retentionInDays                = "4"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(&clusterInfo, description, retentionInDays),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "type", "replicaSet"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "snapshot_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.Name),
					resource.TestCheckResourceAttr(resourceName, "replica_set_name", clusterInfo.Name),
					resource.TestCheckResourceAttr(resourceName, "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "retention_in_days", retentionInDays),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "type", "replicaSet"),
					resource.TestCheckResourceAttr(dataSourceName, "members.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "snapshot_ids.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_name", clusterInfo.Name),
					resource.TestCheckResourceAttr(dataSourceName, "replica_set_name", clusterInfo.Name),
					resource.TestCheckResourceAttr(dataSourceName, "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr(dataSourceName, "description", description),
					resource.TestCheckResourceAttrSet(dataSourcePluralSimpleName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralPaginationName, "results.#"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"retention_in_days"},
			},
		},
	})
}

func TestAccBackupRSCloudBackupSnapshot_sharded(t *testing.T) {
	var (
		projectID       = acc.ProjectIDExecution(t)
		clusterName     = acc.RandomClusterName()
		description     = "My description in my cluster"
		retentionInDays = "4"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configSharded(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "type", "shardedCluster"),
					resource.TestCheckResourceAttrWith(resourceName, "members.#", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttrWith(resourceName, "snapshot_ids.#", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "retention_in_days", retentionInDays),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "type", "shardedCluster"),
					resource.TestCheckResourceAttrWith(dataSourceName, "members.#", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttrWith(dataSourceName, "snapshot_ids.#", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttr(dataSourceName, "description", description)),
			},
		},
	})
}

func TestAccBackupRSCloudBackupSnapshot_timeouts(t *testing.T) {
	var (
		clusterInfo     = acc.GetClusterInfo(t, &acc.ClusterRequest{CloudBackup: true})
		description     = "Timeout test snapshot"
		retentionInDays = "1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy, // resource is deleted when creation times out
		Steps: []resource.TestStep{
			{
				Config:      configTimeout(&clusterInfo, description, retentionInDays),
				ExpectError: regexp.MustCompile("will run cleanup because delete_on_create_timeout is true"),
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		if ids["snapshot_id"] == "" {
			return fmt.Errorf("no ID is set")
		}
		if _, _, err := acc.ConnV2().CloudBackupsApi.GetReplicaSetBackup(context.Background(), ids["project_id"], ids["cluster_name"], ids["snapshot_id"]).Execute(); err != nil {
			return fmt.Errorf("cloudBackupSnapshot (%s) does not exist", rs.Primary.Attributes["snapshot_id"])
		}

		if _, _, err := acc.ConnV2().ClustersApi.GetCluster(context.Background(), ids["project_id"], ids["cluster_name"]).Execute(); err != nil {
			return fmt.Errorf("cluster (%s : %s) does not exist", ids["project_id"], ids["cluster_name"])
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	if acc.ExistingClusterUsed() {
		return nil
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_backup_snapshot" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		res, _, _ := acc.ConnV2().CloudBackupsApi.GetReplicaSetBackup(context.Background(), ids["project_id"], ids["cluster_name"], ids["snapshot_id"]).Execute()
		if res != nil {
			return fmt.Errorf("cloudBackupSnapshot (%s) still exists", rs.Primary.Attributes["snapshot_id"])
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"], rs.Primary.Attributes["snapshot_id"]), nil
	}
}

func configBasic(info *acc.ClusterInfo, description, retentionInDays string) string {
	return info.TerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_cloud_backup_snapshot" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			description       = %[3]q
			retention_in_days = %[4]q
		}

		data "mongodbatlas_cloud_backup_snapshot" "test" {
			snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.snapshot_id
			cluster_name     = %[1]s
			project_id       = %[2]q
		}

		data "mongodbatlas_cloud_backup_snapshots" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]q
		}

		data "mongodbatlas_cloud_backup_snapshots" "pagination" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			page_num = 1
			items_per_page = 5
		}
	`, info.TerraformNameRef, info.ProjectID, description, retentionInDays)
}

func configSharded(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "my_cluster" {
			project_id   = %[1]q
			name           = %[2]q
			cluster_type   = "SHARDED"
			backup_enabled = true
		
			replication_specs {
				num_shards = 3
		
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
		}

		resource "mongodbatlas_cloud_backup_snapshot" "test" {
			project_id        = mongodbatlas_advanced_cluster.my_cluster.project_id
			cluster_name      = mongodbatlas_advanced_cluster.my_cluster.name
			description       = %[3]q
			retention_in_days = %[4]q
		}

		data "mongodbatlas_cloud_backup_snapshot" "test" {
			snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.snapshot_id
			project_id        = mongodbatlas_advanced_cluster.my_cluster.project_id
			cluster_name      = mongodbatlas_advanced_cluster.my_cluster.name
		}

	`, projectID, clusterName, description, retentionInDays)
}

func configTimeout(info *acc.ClusterInfo, description, retentionInDays string) string {
	return info.TerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_cloud_backup_snapshot" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			description       = %[3]q
			retention_in_days = %[4]q
			delete_on_create_timeout = true # default value
			
			timeouts {
				create = "10s"
			}
		}

		data "mongodbatlas_cloud_backup_snapshot" "test" {
			snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.snapshot_id
			cluster_name     = %[1]s
			project_id       = %[2]q
		}

		data "mongodbatlas_cloud_backup_snapshots" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]q
		}

		data "mongodbatlas_cloud_backup_snapshots" "pagination" {
			cluster_name     = %[1]s
			project_id       = %[2]q
			page_num = 1
			items_per_page = 5
		}
	`, info.TerraformNameRef, info.ProjectID, description, retentionInDays)
}
