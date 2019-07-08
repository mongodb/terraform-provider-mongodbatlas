package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasCloudProviderSnapshot_basic(t *testing.T) {
	var cloudProviderSnapshot matlas.CloudProviderSnapshot

	resourceName := "data.mongodbatlas_cloud_provider_snapshot.test"
	groupID := "5d0f1f73cf09a29120e173cf"
	clusterName := "MyClusterTest"
	description := "SomeDescription"
	retentionInDays := "1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudProviderSnapshotConfig(groupID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotExists("mongodbatlas_cloud_provider_snapshot.test", &cloudProviderSnapshot),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "group_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "description"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "retention_in_days"),
				),
			},
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotConfigWithDS(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotExists(resourceName, &cloudProviderSnapshot),
					resource.TestCheckResourceAttrSet(resourceName, "snapshot_id"),
					resource.TestCheckResourceAttrSet(resourceName, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_name"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceCloudProviderSnapshotConfig(groupID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cloud_provider_snapshot" "test" {
		group_id          = "%s"
		cluster_name      = "%s"
		description       = "%s"
		retention_in_days = %s
	}
`, groupID, clusterName, description, retentionInDays)
}

func testAccMongoDBAtlasCloudProviderSnapshotConfigWithDS() string {
	return `
		data "mongodbatlas_cloud_provider_snapshot" "test" {
			snapshot_id  = "5d1285acd5ec13b6c2d1726a"
			group_id     = "${mongodbatlas_cloud_provider_snapshot.test.group_id}"
			cluster_name = "${mongodbatlas_cloud_provider_snapshot.test.cluster_name}"
		}`
}
