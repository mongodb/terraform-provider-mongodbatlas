package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasProject_basic(t *testing.T) {
	var project matlas.Project

	resourceName := "data.mongodbatlas_project.test"
	projectName := fmt.Sprintf("test-datasource-project-%s", acctest.RandString(10))
	orgID := "5b71ff2f96e82120d0aaec14"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceProjectConfig(projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists("mongodbatlas_projects.test", &project),
					resource.TestCheckResourceAttrSet("mongodbatlas_projects.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_projects.test", "org_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectConfigWithDS(projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceProjectConfig(projectName, orgID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_projects" "test" {
			name   = "%[1]s"
			org_id = "%[2]s"
		}
`, projectName, orgID)
}

func testAccMongoDBAtlasProjectConfigWithDS(projectName, orgID string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_project" "test" {
			name = "${mongodbatlas_projects.test.name}"
		}
	`, testAccMongoDBAtlasDataSourceProjectConfig(projectName, orgID))
}
