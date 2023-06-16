package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBackupGenericDSDataLakes_basic(t *testing.T) {
	resourceName := "data.mongodbatlas_data_lakes.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataLakesDataSourceConfig(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataLakesDataSourceConfig(orgID, projectName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "backup_project" {
			name   = %[2]q
			org_id = %[1]q
		}
		data "mongodbatlas_data_lakes" "test" {
			project_id = mongodbatlas_project.backup_project.id
		}
	`, orgID, projectName)
}
