package mongodbatlas

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccProjectDSProject_byID(t *testing.T) {
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); checkTeamsIds(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfigWithDSByID(projectName, orgID,
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
					[]*apiKey{
						{
							id:    apiKeysIds[0],
							roles: []string{"GROUP_READ_ONLY"},
						},
						{
							id:    apiKeysIds[1],
							roles: []string{"GROUP_OWNER"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccProjectDSProject_byName(t *testing.T) {
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); checkTeamsIds(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfigWithDSByName(projectName, orgID,
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
					[]*apiKey{
						{
							id:    apiKeysIds[0],
							roles: []string{"GROUP_READ_ONLY"},
						},
						{
							id:    apiKeysIds[1],
							roles: []string{"GROUP_OWNER"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccProjectDSProject_defaultFlags(t *testing.T) {
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); checkTeamsIds(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfigWithDSByName(projectName, orgID,
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
					[]*apiKey{
						{
							id:    apiKeysIds[0],
							roles: []string{"GROUP_READ_ONLY"},
						},
						{
							id:    apiKeysIds[1],
							roles: []string{"GROUP_OWNER"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "is_collect_database_specifics_statistics_enabled"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "is_data_explorer_enabled"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "is_performance_advisor_enabled"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "is_realtime_performance_panel_enabled"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "is_schema_advisor_enabled"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccMongoDBAtlasProjectConfigWithDSByID(projectName, orgID string, teams []*matlas.ProjectTeam, apiKeys []*apiKey) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_project" "test" {
			project_id = "${mongodbatlas_project.test.id}"
		}
	`, testAccMongoDBAtlasProjectConfig(projectName, orgID, teams, apiKeys))
}

func testAccMongoDBAtlasProjectConfigWithDSByName(projectName, orgID string, teams []*matlas.ProjectTeam, apiKeys []*apiKey) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_project" "test" {
			name = "${mongodbatlas_project.test.name}"
		}
	`, testAccMongoDBAtlasProjectConfig(projectName, orgID, teams, apiKeys))
}
