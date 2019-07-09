package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasCloudProviderSnapshot_basic(t *testing.T) {
	var cloudProviderSnapshot = matlas.CloudProviderSnapshot{
		RetentionInDays: 1,
	}

	resourceName := "mongodbatlas_cloud_provider_snapshot.test"

	projectID := "5d0f1f73cf09a29120e173cf"
	clusterName := "MyClusterTest"
	description := "SomeDescription"
	retentionInDays := "1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCloudProviderSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotConfig(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotExists(resourceName, &cloudProviderSnapshot),
					testAccCheckMongoDBAtlasCloudProviderSnapshotAttributes(&cloudProviderSnapshot, retentionInDays),
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

	projectID := "5d0f1f73cf09a29120e173cf"
	clusterName := "MyClusterTest"
	description := "SomeDescription"
	retentionInDays := "1"

	resource.Test(t, resource.TestCase{
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
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] cloudProviderSnapshot ID: %s", rs.Primary.ID)

		requestParameters := &matlas.SnapshotReqPathParameters{
			SnapshotID:  rs.Primary.ID,
			GroupID:     rs.Primary.Attributes["project_id"],
			ClusterName: rs.Primary.Attributes["cluster_name"],
		}

		_, _, err := conn.CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), requestParameters)
		if err == nil {
			return nil
		}

		return fmt.Errorf("cloudProviderSnapshot (%s) does not exist", rs.Primary.ID)
	}
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotAttributes(cloudProviderSnapshot *matlas.CloudProviderSnapshot, retentionInDays string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if strconv.Itoa(cloudProviderSnapshot.RetentionInDays) != retentionInDays {
			return fmt.Errorf("bad cloudProviderSnapshot retentionInDays: %s", strconv.Itoa(cloudProviderSnapshot.RetentionInDays))
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
			SnapshotID:  rs.Primary.ID,
			GroupID:     rs.Primary.Attributes["project_id"],
			ClusterName: rs.Primary.Attributes["cluster_name"],
		}

		_, err := conn.CloudProviderSnapshots.Delete(context.Background(), requestParameters)

		if err != nil {
			return fmt.Errorf("cloudProviderSnapshot (%s) still exists", rs.Primary.ID)
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
		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"], rs.Primary.ID), nil
	}
}

func testAccMongoDBAtlasCloudProviderSnapshotConfig(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cloud_provider_snapshot" "test" {
			project_id        = "%s"
			cluster_name      = "%s"
			description       = "%s"
			retention_in_days = %s
		}
	`, projectID, clusterName, description, retentionInDays)
}
