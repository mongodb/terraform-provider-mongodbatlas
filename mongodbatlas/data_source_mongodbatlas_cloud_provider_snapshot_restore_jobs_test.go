package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobs_basic(t *testing.T) {
	projectID := "5cf5a45a9ccf6400e60981b6"
	clusterName := "MyCluster"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobsConfigWithDS(projectID, clusterName),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobsConfigWithDS(projectID, clusterName string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_cloud_provider_snapshot_restore_jobs" "test" {
			project_id     = "%s"
			cluster_name = "%s"
		}`, projectID, clusterName)
}
