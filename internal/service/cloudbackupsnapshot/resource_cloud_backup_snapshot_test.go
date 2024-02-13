package cloudbackupsnapshot_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName = "mongodbatlas_cloud_backup_snapshot.test"
)

func TestAccBackupRSCloudBackupSnapshot_basic(t *testing.T) {
	var (
		dataSourceName                 = "data.mongodbatlas_cloud_backup_snapshot.test"
		dataSourcePluralSimpleName     = "data.mongodbatlas_cloud_backup_snapshots.test"
		dataSourcePluralPaginationName = "data.mongodbatlas_cloud_backup_snapshots.pagination"
		orgID                          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterInfo                    = acc.GetClusterInfo(orgID, true)
		description                    = "My description in my cluster"
		retentionInDays                = "4"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(&clusterInfo, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterInfo.ClusterName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "retention_in_days", retentionInDays),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_name", clusterInfo.ClusterName),
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
		_, _, err := acc.ConnV2().CloudBackupsApi.GetReplicaSetBackup(context.Background(), ids["project_id"], ids["cluster_name"], ids["snapshot_id"]).Execute()
		if err == nil {
			return nil
		}

		return fmt.Errorf("cloudBackupSnapshot (%s) does not exist", rs.Primary.Attributes["snapshot_id"])
	}
}

func checkDestroy(s *terraform.State) error {
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
	return info.ClusterTerraformStr + fmt.Sprintf(`
		resource "mongodbatlas_cloud_backup_snapshot" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]s
			description       = %[3]q
			retention_in_days = %[4]q
		}

		data "mongodbatlas_cloud_backup_snapshot" "test" {
			snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.snapshot_id
			cluster_name     = %[1]s
			project_id       = %[2]s
		}

		data "mongodbatlas_cloud_backup_snapshots" "test" {
			cluster_name     = %[1]s
			project_id       = %[2]s
		}

		data "mongodbatlas_cloud_backup_snapshots" "pagination" {
			cluster_name     = %[1]s
			project_id       = %[2]s
			page_num = 1
			items_per_page = 5
		}
	`, info.ClusterNameStr, info.ProjectIDStr, description, retentionInDays)
}
