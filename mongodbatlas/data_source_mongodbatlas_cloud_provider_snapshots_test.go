package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMongoDBAtlasCloudProviderSnapshots_basic(t *testing.T) {
	groupID := "5d0f1f73cf09a29120e173cf"
	clusterName := "MyClusterTest"
	description := "SomeDescription"
	retentionInDays := "1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudProviderSnapshotsConfig(groupID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "group_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "description"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "retention_in_days"),
				),
			},
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotsConfigWithDS(),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceCloudProviderSnapshotsConfig(groupID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cloud_provider_snapshot" "test" {
			group_id          = "%s"
			cluster_name      = "%s"
			description       = "%s"
			retention_in_days = %s
		}
	`, groupID, clusterName, description, retentionInDays)
}

func testAccMongoDBAtlasCloudProviderSnapshotsConfigWithDS() string {
	return `
		data "mongodbatlas_cloud_provider_snapshots" "test" {
			group_id     = "${mongodbatlas_cloud_provider_snapshot.test.group_id}"
			cluster_name = "${mongodbatlas_cloud_provider_snapshot.test.cluster_name}"
		}`
}
