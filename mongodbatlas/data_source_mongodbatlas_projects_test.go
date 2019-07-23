package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMongoDBAtlasProjects_basic(t *testing.T) {
	projectName := fmt.Sprintf("test-datasource-project-%s", acctest.RandString(10))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceProjectsConfig(projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectsConfigWithDS(projectName, orgID),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceProjectsConfig(projectName, orgID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%[1]s"
			org_id = "%[2]s"
		}
	`, projectName, orgID)
}

func testAccMongoDBAtlasProjectsConfigWithDS(projectName, orgID string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_projects" "test" {}
	`, testAccMongoDBAtlasDataSourceProjectsConfig(projectName, orgID))
}
