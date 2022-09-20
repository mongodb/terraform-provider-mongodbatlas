package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMongoDBAtlasServerlessInstance_byName(t *testing.T) {
	var (
		instanceName   = acctest.RandomWithPrefix("test-serverless-instance")
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		datasourceName = "data.mongodbatlas_serverless_instance.test_two"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckBetaFeatures(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasServerlessInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasServerlessInstanceDSConfig(projectID, instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "name"),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "state_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "create_date"),
					resource.TestCheckResourceAttrSet(datasourceName, "mongo_db_version"),
					resource.TestCheckResourceAttrSet(datasourceName, "continuous_backup_enabled"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasServerlessInstanceDSConfig(projectID, name string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_serverless_instance" "test_two" {
			name        = mongodbatlas_serverless_instance.test.name
			project_id  = mongodbatlas_serverless_instance.test.project_id
		}

	`, testAccMongoDBAtlasServerlessInstanceConfig(projectID, name))
}
