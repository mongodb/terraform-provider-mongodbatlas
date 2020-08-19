package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasCloudProviderSnapshot_basic(t *testing.T) {
	var (
		cloudProviderSnapshot = matlas.CloudProviderSnapshot{}
		resourceName          = "mongodbatlas_cloud_provider_snapshot.test"
		projectID             = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName           = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		description           = "My description in my cluster"
		retentionInDays       = "4"
	)

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
	var (
		resourceName    = "mongodbatlas_cloud_provider_snapshot.test"
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName     = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		description     = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays = "5"
	)

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

		ids := decodeStateID(rs.Primary.ID)

		if ids["snapshot_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] cloudProviderSnapshot ID: %s", ids["snapshot_id"])

		requestParameters := &matlas.SnapshotReqPathParameters{
			SnapshotID:  ids["snapshot_id"],
			GroupID:     ids["project_id"],
			ClusterName: ids["cluster_name"],
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

		ids := decodeStateID(rs.Primary.ID)

		requestParameters := &matlas.SnapshotReqPathParameters{
			SnapshotID:  ids["snapshot_id"],
			GroupID:     ids["project_id"],
			ClusterName: ids["cluster_name"],
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
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		log.Printf("%s-%s-%s", ids["project_id"], ids["cluster_name"], ids["snapshot_id"])

		return fmt.Sprintf("%s-%s-%s", ids["project_id"], ids["cluster_name"], ids["snapshot_id"]), nil
	}
}

func testAccMongoDBAtlasCloudProviderSnapshotConfig(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 10

			// Provider Settings "block"
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

func TestResourceMongoDBAtlasCloudProviderSnapshot_snapshotID(t *testing.T) {
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
