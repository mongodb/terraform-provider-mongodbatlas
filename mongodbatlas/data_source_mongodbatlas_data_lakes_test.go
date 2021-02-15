package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMongoDBAtlasDataLakes_basic(t *testing.T) {
	resourceName := "data.mongodbatlas_data_lakes.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataLakesDataSourceConfig(projectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataLakesDataSourceConfig(projectID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_data_lakes" "test" {
			project_id = "%s"
		}
	`, projectID)
}
