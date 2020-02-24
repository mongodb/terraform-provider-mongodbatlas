package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasProject_byID(t *testing.T) {

	projectName := fmt.Sprintf("test-datasource-project-%s", acctest.RandString(10))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfigWithDSByID(projectName, orgID,
					[]*matlas.ProjectTeam{
						{
							TeamID:    "5e0fa8c99ccf641c722fe683",
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamID:    "5e1dd7b4f2a30ba80a70cd3a",
							RoleNames: []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
				),
			},
		},
	})
}

func TestAccDataSourceMongoDBAtlasProject_byName(t *testing.T) {

	projectName := fmt.Sprintf("test-datasource-project-%s", acctest.RandString(10))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfigWithDSByName(projectName, orgID,
					[]*matlas.ProjectTeam{
						{
							TeamID:    "5e0fa8c99ccf641c722fe683",
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamID:    "5e1dd7b4f2a30ba80a70cd3a",
							RoleNames: []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasProjectConfigWithDSByID(projectName, orgID string, teams []*matlas.ProjectTeam) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_project" "test" {
			project_id = "${mongodbatlas_project.test.id}"
		}
	`, testAccMongoDBAtlasPropjectConfig(projectName, orgID, teams))
}

func testAccMongoDBAtlasProjectConfigWithDSByName(projectName, orgID string, teams []*matlas.ProjectTeam) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_project" "test" {
			name = "${mongodbatlas_project.test.name}"
		}
	`, testAccMongoDBAtlasPropjectConfig(projectName, orgID, teams))
}
