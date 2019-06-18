package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMongoDBAtlasProjects_basic(t *testing.T) {
	projectName := fmt.Sprintf("test-datasource-project-%s", acctest.RandString(10))
	orgID := "5b71ff2f96e82120d0aaec14"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceProjectsConfig(projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_projects.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_projects.test", "org_id"),
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
		resource "mongodbatlas_projects" "test" {
			name   = "%[1]s"
			org_id = "%[2]s"
		}
`, projectName, orgID)
}

func testAccMongoDBAtlasProjectsConfigWithDS(projectName, orgID string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_projects" "test" {
		}
	`, testAccMongoDBAtlasDataSourceProjectsConfig(projectName, orgID))
}
