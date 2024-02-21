package project_test

import (
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/atlas-sdk/v20231115007/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccProjectDSProjects_basic(t *testing.T) {
	var (
		projectName = fmt.Sprintf("test-datasource-project-%s", acctest.RandString(10))
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIDsWithMinCount(t, 2) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectsConfigWithDS(projectName, orgID,
					[]*admin.TeamRole{
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(0)),
							RoleNames: &[]string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(1)),
							RoleNames: &[]string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
					resource.TestCheckResourceAttr("mongodbatlas_project.test", "teams.#", "2"),
					// Test for Data source
					resource.TestCheckResourceAttrSet("data.mongodbatlas_projects.test", "total_count"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_projects.test", "results.#"),
				),
			},
		},
	})
}

func TestAccProjectDSProjects_withPagination(t *testing.T) {
	var (
		projectName = fmt.Sprintf("test-datasource-project-%s", acctest.RandString(10))
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIDsWithMinCount(t, 2) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectsConfigWithPagination(projectName, orgID,
					[]*admin.TeamRole{
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(0)),
							RoleNames: &[]string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(1)),
							RoleNames: &[]string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					},
					2, 5,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
					resource.TestCheckResourceAttr("mongodbatlas_project.test", "teams.#", "2"),
					// Test for Data source
					resource.TestCheckResourceAttrSet("data.mongodbatlas_projects.test", "total_count"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_projects.test", "results.#"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasProjectsConfigWithDS(projectName, orgID string, teams []*admin.TeamRole) string {
	config := fmt.Sprintf(`
		%s
		data "mongodbatlas_projects" "test" {}
	`, acc.ConfigProject(projectName, orgID, teams))
	return config
}

func testAccMongoDBAtlasProjectsConfigWithPagination(projectName, orgID string, teams []*admin.TeamRole, pageNum, itemPage int) string {
	return fmt.Sprintf(`
		%s
		data "mongodbatlas_projects" "test" {
			page_num = %d
			items_per_page = %d
		}
	`, acc.ConfigProject(projectName, orgID, teams), pageNum, itemPage)
}
