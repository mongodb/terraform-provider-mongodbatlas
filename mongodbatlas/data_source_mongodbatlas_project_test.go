package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasProject_basic(t *testing.T) {
	var project matlas.Project

	resourceName := "data.mongodbatlas_project.test"
	projectName := fmt.Sprintf("test-datasource-project-%s", acctest.RandString(10))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceProjectConfig(projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists("mongodbatlas_project.test", &project),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
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
		resource "mongodbatlas_project" "test" {
			name   = "%[1]s"
			org_id = "%[2]s"
		}
`, projectName, orgID)
}

func testAccMongoDBAtlasProjectConfigWithDS(projectName, orgID string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_project" "test" {
			name = "${mongodbatlas_project.test.name}"
		}
	`, testAccMongoDBAtlasDataSourceProjectConfig(projectName, orgID))
}
