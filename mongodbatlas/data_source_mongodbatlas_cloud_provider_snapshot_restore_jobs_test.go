package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobs_basic(t *testing.T) {
	groupID := "5cf5a45a9ccf6400e60981b6"
	clusterName := "MyCluster"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobsConfigWithDS(groupID, clusterName),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobsConfigWithDS(groupID, clusterName string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_cloud_provider_snapshot_restore_jobs" "test" {
			group_id     = "%s"
			cluster_name = "%s"
		}`, groupID, clusterName)
}
