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

func TestAccProjectDSProject_byID(t *testing.T) {
	projectName := acctest.RandomWithPrefix("test-acc")
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	teamsIds := strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
	if len(teamsIds) < 2 {
		t.Skip("`MONGODB_ATLAS_TEAMS_IDS` must have 2 team ids for this acceptance testing")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t); testCheckTeamsIds(t) },
<<<<<<< HEAD:mongodbatlas/data_source_mongodbatlas_project_test.go
		ProtoV6ProviderFactories: testAccProviderV6Factories,
=======
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
>>>>>>> f32c9d60 (data sources):mongodbatlas/fw_data_source_mongodbatlas_project_test.go
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectDSByIDUsingRS(testAccMongoDBAtlasProjectConfig(projectName, orgID,
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
				)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
				),
			},
		},
	})
}

func TestAccProjectDSProject_byName(t *testing.T) {
	projectName := acctest.RandomWithPrefix("test-acc")
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	teamsIds := strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
	if len(teamsIds) < 2 {
		t.Skip("`MONGODB_ATLAS_TEAMS_IDS` must have 2 team ids for this acceptance testing")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t); testCheckTeamsIds(t) },
<<<<<<< HEAD:mongodbatlas/data_source_mongodbatlas_project_test.go
		ProtoV6ProviderFactories: testAccProviderV6Factories,
=======
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
>>>>>>> f32c9d60 (data sources):mongodbatlas/fw_data_source_mongodbatlas_project_test.go
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectDSByNameUsingRS(testAccMongoDBAtlasProjectConfig(projectName, orgID,
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
				)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_project.test", "org_id"),
				),
			},
		},
	})
}

func TestAccProjectDSProject_defaultFlags(t *testing.T) {
	projectName := acctest.RandomWithPrefix("test-acc")
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	teamsIds := strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
	if len(teamsIds) < 2 {
		t.Skip("`MONGODB_ATLAS_TEAMS_IDS` must have 2 team ids for this acceptance testing")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t); testCheckTeamsIds(t) },
<<<<<<< HEAD:mongodbatlas/data_source_mongodbatlas_project_test.go
		ProtoV6ProviderFactories: testAccProviderV6Factories,
=======
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
>>>>>>> f32c9d60 (data sources):mongodbatlas/fw_data_source_mongodbatlas_project_test.go
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectDSByNameUsingRS(testAccMongoDBAtlasProjectConfig(projectName, orgID,
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
				),
			},
		},
	})
}

func TestAccProjectDSProject_limits(t *testing.T) {
	projectName := acctest.RandomWithPrefix("test-acc")
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
<<<<<<< HEAD:mongodbatlas/data_source_mongodbatlas_project_test.go
		ProtoV6ProviderFactories: testAccProviderV6Factories,
=======
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
>>>>>>> f32c9d60 (data sources):mongodbatlas/fw_data_source_mongodbatlas_project_test.go
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectDSByNameUsingRS(testAccMongoDBAtlasProjectConfigWithLimits(projectName, orgID, []*projectLimit{})),
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
