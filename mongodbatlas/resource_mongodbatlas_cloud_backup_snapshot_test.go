package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasCloudBackupSnapshot_basic(t *testing.T) {
	var (
		cloudBackupSnapshot = matlas.CloudProviderSnapshot{}
		resourceName        = "mongodbatlas_cloud_backup_snapshot.test"
		projectID           = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName         = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		description         = "My description in my cluster"
		retentionInDays     = "4"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupSnapshotConfig(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupSnapshotExists(resourceName, &cloudBackupSnapshot),
					testAccCheckMongoDBAtlasCloudBackupSnapshotAttributes(&cloudBackupSnapshot, description),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "retention_in_days", retentionInDays),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasCloudBackupSnapshotImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"retention_in_days"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasCloudBackupSnapshotExists(resourceName string, cloudBackupSnapshot *matlas.CloudProviderSnapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		if ids["snapshot_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] cloudBackupSnapshot ID: %s", ids["snapshot_id"])

		requestParameters := &matlas.SnapshotReqPathParameters{
			SnapshotID:  ids["snapshot_id"],
			GroupID:     ids["project_id"],
			ClusterName: ids["cluster_name"],
		}

		res, _, err := conn.CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), requestParameters)
		if err == nil {
			*cloudBackupSnapshot = *res
			return nil
		}

		return fmt.Errorf("cloudBackupSnapshot (%s) does not exist", rs.Primary.Attributes["snapshot_id"])
	}
}

func testAccCheckMongoDBAtlasCloudBackupSnapshotAttributes(cloudBackupSnapshot *matlas.CloudProviderSnapshot, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloudBackupSnapshot.Description != description {
			return fmt.Errorf("bad cloudBackupSnapshot description: %s", cloudBackupSnapshot.Description)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasCloudBackupSnapshotDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_backup_snapshot" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		requestParameters := &matlas.SnapshotReqPathParameters{
			SnapshotID:  ids["snapshot_id"],
			GroupID:     ids["project_id"],
			ClusterName: ids["cluster_name"],
		}

		res, _, _ := conn.CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), requestParameters)

		if res != nil {
			return fmt.Errorf("cloudBackupSnapshot (%s) still exists", rs.Primary.Attributes["snapshot_id"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasCloudBackupSnapshotImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"], rs.Primary.Attributes["snapshot_id"]), nil
	}
}

func testAccMongoDBAtlasCloudBackupSnapshotConfig(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = %[1]q
  name         = %[2]q
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
  description       = %[3]q
  retention_in_days = %[4]q
}
	`, projectID, clusterName, description, retentionInDays)
}

func TestResourceMongoDBAtlasCloudBackupSnapshot_snapshotID(t *testing.T) {
	got, err := splitSnapshotImportID("5cf5a45a9ccf6400e60981b6-projectname-environment-mongo-global-cluster-5cf5a45a9ccf6400e60981b7")
	if err != nil {
		t.Errorf("splitSnapshotImportID returned error(%s), expected error=nil", err)
	}

	expected := &matlas.SnapshotReqPathParameters{
		GroupID:     "5cf5a45a9ccf6400e60981b6",
		ClusterName: "projectname-environment-mongo-global-cluster",
		SnapshotID:  "5cf5a45a9ccf6400e60981b7",
	}

	if diff := deep.Equal(expected, got); diff != nil {
		t.Errorf("Bad splitSnapshotImportID return \n got = %#v\nwant = %#v \ndiff = %#v", expected, *got, diff)
	}

	if _, err := splitSnapshotImportID("5cf5a45a9ccf6400e60981b6projectname-environment-mongo-global-cluster5cf5a45a9ccf6400e60981b7"); err == nil {
		t.Error("splitSnapshotImportID expected to have error")
	}
}
