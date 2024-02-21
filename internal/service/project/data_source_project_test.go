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

const dataSourceName = "data.mongodbatlas_project.test"

func TestAccProjectDSProject_byID(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("test-acc")
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIDsWithMinCount(t, 2) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectDSByIDUsingRS(acc.ConfigProject(projectName, orgID,
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
				)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttr(dataSourceName, "teams.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ip_addresses.services.clusters.#"),
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
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIDsWithMinCount(t, 2) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectDSByNameUsingRS(acc.ConfigProject(projectName, orgID,
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
				)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttr(dataSourceName, "teams.#", "2"),
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
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIDsWithMinCount(t, 2) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectDSByNameUsingRS(acc.ConfigProject(projectName, orgID,
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
				)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "is_collect_database_specifics_statistics_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "is_data_explorer_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "is_extended_storage_sizes_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "is_performance_advisor_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "is_realtime_performance_panel_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "is_schema_advisor_enabled"),
					resource.TestCheckResourceAttr(dataSourceName, "teams.#", "2"),
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
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "limits.0.name"),
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
