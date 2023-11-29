package project_test

import (
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/atlas-sdk/v20231115001/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccProjectDSProject_byID(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIdsWithMinCount(t, 2) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectDSByIDUsingRS(acc.ConfigProject(projectName, orgID,
					[]*matlas.ProjectTeam{
						{
							TeamID:    acc.GetProjectTeamsIdsWithPos(0),
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamID:    acc.GetProjectTeamsIdsWithPos(1),
							RoleNames: []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					},
				)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
					resource.TestCheckResourceAttr("mongodbatlas_project.test", "teams.#", "2"),
				),
			},
		},
	})
}

func TestAccProjectDSProject_byName(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIdsWithMinCount(t, 2) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectDSByNameUsingRS(acc.ConfigProject(projectName, orgID,
					[]*matlas.ProjectTeam{
						{
							TeamID:    acc.GetProjectTeamsIdsWithPos(0),
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{

							TeamID:    acc.GetProjectTeamsIdsWithPos(1),
							RoleNames: []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					},
				)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
					resource.TestCheckResourceAttr("mongodbatlas_project.test", "teams.#", "2"),
				),
			},
		},
	})
}

func TestAccProjectDSProject_defaultFlags(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIdsWithMinCount(t, 2) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectDSByNameUsingRS(acc.ConfigProject(projectName, orgID,
					[]*matlas.ProjectTeam{
						{
							TeamID:    acc.GetProjectTeamsIdsWithPos(0),
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{

							TeamID:    acc.GetProjectTeamsIdsWithPos(1),
							RoleNames: []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					},
				)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "is_collect_database_specifics_statistics_enabled"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "is_data_explorer_enabled"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "is_extended_storage_sizes_enabled"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "is_performance_advisor_enabled"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "is_realtime_performance_panel_enabled"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "is_schema_advisor_enabled"),
					resource.TestCheckResourceAttr("mongodbatlas_project.test", "teams.#", "2"),
				),
			},
		},
	})
}

func TestAccProjectDSProject_limits(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectDSByNameUsingRS(acc.ConfigProjectWithLimits(projectName, orgID, []*admin.DataFederationLimit{})),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_project.test", "org_id"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_project.test", "limits.0.name"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasProjectDSByNameUsingRS(rs string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_project" "test" {
			name = "${mongodbatlas_project.test.name}"
		}
	`, rs)
}

func testAccMongoDBAtlasProjectDSByIDUsingRS(rs string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_project" "test" {
			project_id = "${mongodbatlas_project.test.id}"
		}
	`, rs)
}
