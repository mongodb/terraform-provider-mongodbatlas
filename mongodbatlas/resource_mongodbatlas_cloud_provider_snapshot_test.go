package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasCloudProviderSnapshot_basic(t *testing.T) {
	var cloudProviderSnapshot = matlas.CloudProviderSnapshot{}

	resourceName := "mongodbatlas_cloud_provider_snapshot.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	description := fmt.Sprintf("My description in %s", clusterName)
	retentionInDays := "1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCloudProviderSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotConfig(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotExists(resourceName, &cloudProviderSnapshot),
					testAccCheckMongoDBAtlasCloudProviderSnapshotAttributes(&cloudProviderSnapshot, description),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "retention_in_days", retentionInDays),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCloudProviderSnapshot_importBasic(t *testing.T) {

	resourceName := "mongodbatlas_cloud_provider_snapshot.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	description := fmt.Sprintf("My description in %s", clusterName)
	retentionInDays := "1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCloudProviderSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotConfig(projectID, clusterName, description, retentionInDays),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasCloudProviderSnapshotImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"retention_in_days"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotExists(resourceName string, cloudProviderSnapshot *matlas.CloudProviderSnapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.Attributes["snapshot_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] cloudProviderSnapshot ID: %s", rs.Primary.Attributes["snapshot_id"])

		requestParameters := &matlas.SnapshotReqPathParameters{
			SnapshotID:  rs.Primary.Attributes["snapshot_id"],
			GroupID:     rs.Primary.Attributes["project_id"],
			ClusterName: rs.Primary.Attributes["cluster_name"],
		}

		res, _, err := conn.CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), requestParameters)
		if err == nil {
			*cloudProviderSnapshot = *res
			return nil
		}
		return fmt.Errorf("cloudProviderSnapshot (%s) does not exist", rs.Primary.Attributes["snapshot_id"])
	}
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotAttributes(cloudProviderSnapshot *matlas.CloudProviderSnapshot, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloudProviderSnapshot.Description != description {
			return fmt.Errorf("bad cloudProviderSnapshot description: %s", cloudProviderSnapshot.Description)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_provider_snapshot" {
			continue
		}

		requestParameters := &matlas.SnapshotReqPathParameters{
			SnapshotID:  rs.Primary.Attributes["snapshot_id"],
			GroupID:     rs.Primary.Attributes["project_id"],
			ClusterName: rs.Primary.Attributes["cluster_name"],
		}

		res, _, _ := conn.CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), requestParameters)

		if res != nil {
			return fmt.Errorf("cloudProviderSnapshot (%s) still exists", rs.Primary.Attributes["snapshot_id"])
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"], rs.Primary.Attributes["snapshot_id"]), nil
	}
}

func testAccMongoDBAtlasCloudProviderSnapshotConfig(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 5

			//Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true //enable cloud provider snapshots
			provider_disk_iops          = 100
			provider_encrypt_ebs_volume = false
		}

		resource "mongodbatlas_cloud_provider_snapshot" "test" {
			project_id        = mongodbatlas_cluster.my_cluster.project_id
			cluster_name      = mongodbatlas_cluster.my_cluster.name
			description       = "%s"
			retention_in_days = %s
		}
	`, projectID, clusterName, description, retentionInDays)
}
