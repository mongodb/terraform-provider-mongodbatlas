package mongodbatlas

import (
	"fmt"
	"os"
	"strings"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectDSProjects_basic(t *testing.T) {
	projectName := fmt.Sprintf("test-datasource-project-%s", acctest.RandString(10))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	teamsIds := strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
	if len(teamsIds) < 2 {
		t.Skip("`MONGODB_ATLAS_TEAMS_IDS` must have 2 team ids for this acceptance testing")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t); testCheckTeamsIds(t) },
<<<<<<< HEAD:mongodbatlas/data_source_mongodbatlas_projects_test.go
		ProtoV6ProviderFactories: testAccProviderV6Factories,
=======
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
>>>>>>> f32c9d60 (data sources):mongodbatlas/fw_data_source_mongodbatlas_projects_test.go
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectsConfigWithDS(projectName, orgID,
					[]*matlas.ProjectTeam{
						{
							TeamID:    teamsIds[0],
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamID:    teamsIds[1],
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

func TestAccProjectDSProjects_withPagination(t *testing.T) {
	projectName := fmt.Sprintf("test-datasource-project-%s", acctest.RandString(10))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	teamsIds := strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
	apiKeysIds := strings.Split(os.Getenv("MONGODB_ATLAS_API_KEYS_IDS"), ",")
	if len(teamsIds) < 2 {
		t.Skip("`MONGODB_ATLAS_TEAMS_IDS` must have 2 team ids for this acceptance testing")
	}
	if len(apiKeysIds) < 2 {
		t.Skip("`MONGODB_ATLAS_API_KEYS_IDS` must have 2 api key ids for this acceptance testing")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t); testCheckTeamsIds(t) },
<<<<<<< HEAD:mongodbatlas/data_source_mongodbatlas_projects_test.go
		ProtoV6ProviderFactories: testAccProviderV6Factories,
=======
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
>>>>>>> f32c9d60 (data sources):mongodbatlas/fw_data_source_mongodbatlas_projects_test.go
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectsConfigWithPagination(projectName, orgID,
					[]*matlas.ProjectTeam{
						{
							TeamID:    teamsIds[0],
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamID:    teamsIds[1],
							RoleNames: []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					},
					2, 5,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasProjectsConfigWithDS(projectName, orgID string, teams []*matlas.ProjectTeam) string {
	config := fmt.Sprintf(`
		%s
		data "mongodbatlas_projects" "test" {}
	`, testAccMongoDBAtlasProjectConfig(projectName, orgID, teams))
	return config
}

func testAccMongoDBAtlasProjectsConfigWithPagination(projectName, orgID string, teams []*matlas.ProjectTeam, pageNum, itemPage int) string {
	return fmt.Sprintf(`
		%s
		data "mongodbatlas_projects" "test" {
			page_num = %d
			items_per_page = %d
		}
	`, testAccMongoDBAtlasProjectConfig(projectName, orgID, teams), pageNum, itemPage)
}
