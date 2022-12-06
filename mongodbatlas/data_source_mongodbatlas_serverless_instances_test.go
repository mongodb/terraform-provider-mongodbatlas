package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccClusterDSServerlessInstances_basic(t *testing.T) {
	var (
		clusterName    = acctest.RandomWithPrefix("test-acc-serverless")
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		datasourceName = "data.mongodbatlas_serverless_instances.data_serverless"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckBetaFeatures(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasServerlessInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasServerlessInstancesDSConfig(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.state_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.continuous_backup_enabled"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.termination_protection_enabled"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasServerlessInstancesDSConfig(projectID, name string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_serverless_instances" "data_serverless" {
			project_id         = mongodbatlas_serverless_instance.test.project_id
		}
	`, testAccMongoDBAtlasServerlessInstanceConfig(projectID, name, true))
}
