package cloudbackupsnapshot_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
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
		projectName                    = acctest.RandomWithPrefix("test-acc")
		clusterName                    = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		description                    = "My description in my cluster"
		retentionInDays                = "4"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "retention_in_days", retentionInDays),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(dataSourceName, "description", description),
					resource.TestCheckResourceAttr(dataSourceName, "retention_in_days", retentionInDays),
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
		log.Printf("[DEBUG] cloudBackupSnapshot ID: %s", ids["snapshot_id"])
		requestParameters := &matlas.SnapshotReqPathParameters{
			SnapshotID:  ids["snapshot_id"],
			GroupID:     ids["project_id"],
			ClusterName: ids["cluster_name"],
		}
		_, _, err := acc.Conn().CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), requestParameters)
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
		requestParameters := &matlas.SnapshotReqPathParameters{
			SnapshotID:  ids["snapshot_id"],
			GroupID:     ids["project_id"],
			ClusterName: ids["cluster_name"],
		}
		res, _, _ := acc.Conn().CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), requestParameters)
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

func configBasic(orgID, projectName, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "backup_project" {
	org_id = %[1]q
	name   = %[2]q
}
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = mongodbatlas_project.backup_project.id
  name         = %[3]q
  disk_size_gb = 10

  // Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true //enable cloud backup snapshots
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_cluster.my_cluster.project_id
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = %[4]q
  retention_in_days = %[5]q
}

data "mongodbatlas_cloud_backup_snapshot" "test" {
	snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.snapshot_id
	project_id        = mongodbatlas_cluster.my_cluster.project_id
	cluster_name      = mongodbatlas_cluster.my_cluster.name
}

data "mongodbatlas_cloud_backup_snapshots" "test" {
	project_id        = mongodbatlas_cluster.my_cluster.project_id
	cluster_name      = mongodbatlas_cluster.my_cluster.name
}

data "mongodbatlas_cloud_backup_snapshots" "pagination" {
	project_id        = mongodbatlas_cluster.my_cluster.project_id
	cluster_name      = mongodbatlas_cluster.my_cluster.name
	page_num = 1
	items_per_page = 5
}



	`, orgID, projectName, clusterName, description, retentionInDays)
}
