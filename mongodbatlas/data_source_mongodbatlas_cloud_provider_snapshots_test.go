package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMongoDBAtlasCloudProviderSnapshots_basic(t *testing.T) {
	projectID := "5d0f1f73cf09a29120e173cf"
	clusterName := "MyClusterTest"
	description := "SomeDescription"
	retentionInDays := "1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudProviderSnapshotsConfig(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "project_id"),
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

func testAccMongoDBAtlasDataSourceCloudProviderSnapshotsConfig(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cloud_provider_snapshot" "test" {
			project_id          = "%s"
			cluster_name      = "%s"
			description       = "%s"
			retention_in_days = %s
		}
	`, projectID, clusterName, description, retentionInDays)
}

func testAccMongoDBAtlasCloudProviderSnapshotsConfigWithDS() string {
	return `
		data "mongodbatlas_cloud_provider_snapshots" "test" {
			project_id     = "${mongodbatlas_cloud_provider_snapshot.test.project_id}"
			cluster_name = "${mongodbatlas_cloud_provider_snapshot.test.cluster_name}"
		}`
}
