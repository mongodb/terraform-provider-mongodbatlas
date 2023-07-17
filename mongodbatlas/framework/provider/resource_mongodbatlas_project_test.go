package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type testAPIKey struct {
	id    string
	roles []string
}

var _ plancheck.PlanCheck = debugPlan{}

type debugPlan struct{}

func (e debugPlan) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	rd, err := json.Marshal(req.Plan)
	if err != nil {
		fmt.Println("error marshaling machine-readable plan output:", err)
	}
	fmt.Printf("req.Plan - %s\n", string(rd))
}

func DebugPlan() plancheck.PlanCheck {
	return debugPlan{}
}

// This test only helps to see what the plan is for debugging purposes
//
// func Test_DebugPlan(t *testing.T) {
// var (
// 	resourceName = "mongodbatlas_project.test"
// 	projectName  = acctest.RandomWithPrefix("test-acc-migration")
// 	orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
// )
// 	t.Parallel()

// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: fmt.Sprintf(`resource "mongodbatlas_project" "main" {
// 					name   = "repro-region-config-order-issue"
// 					org_id = "%s"
// 				  }`,orgID),
// 				PlanOnly: true,
// 				ConfigPlanChecks: resource.ConfigPlanChecks{
// 					PostApplyPreRefresh: []plancheck.PlanCheck{
// 						DebugPlan(),
// 					},
// 				},
// 			},
// 		},
// 	})
// }

func TestProjectResource_Basic_Migration(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project.test"
		projectName  = acctest.RandomWithPrefix("test-acc-migration")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.10.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: fmt.Sprintf(`resource "mongodbatlas_project" "test" {
					name   = "%s"
					org_id = "%s"
				  }`, projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
			{
				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
				Config: fmt.Sprintf(`resource "mongodbatlas_project" "test" {
					name   = "%s"
					org_id = "%s"
				  }`, projectName, orgID),
				PlanOnly: true,
			},
		},
	})
}

func TestProjectResource_Migration(t *testing.T) {
	var (
		project      matlas.Project
		resourceName = "mongodbatlas_project.test"
		projectName  = acctest.RandomWithPrefix("test-acc-migration")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterCount = "0"
		teamsIds     = strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
		apiKeysIds   = strings.Split(os.Getenv("MONGODB_ATLAS_API_KEYS_IDS"), ",")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t); testCheckTeamsIds(t) },
		CheckDestroy: testAccCheckMongoDBAtlasProjectDestroy,

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.10.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,
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
					[]*testAPIKey{
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
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,
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
					[]*testAPIKey{
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
				PlanOnly: true,
			},
		},
	})
}

func TestAccProjectRSProject_CreateWithProjectOwner(t *testing.T) {
	var (
		project        matlas.Project
		resourceName   = "mongodbatlas_project.test"
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfigWithProjectOwner(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
		},
	})
}

func TestAccProjectRSGovProject_CreateWithProjectOwner(t *testing.T) {
	var (
		project        matlas.Project
		resourceName   = "mongodbatlas_project.test"
		projectName    = fmt.Sprintf("testacc-project-gov-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID_GOV")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID_GOV")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckGov(t) },
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasGovProjectConfigWithProjectOwner(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
		},
	})
}
func TestAccProjectRSProject_CreateWithFalseDefaultSettings(t *testing.T) {
	var (
		project        matlas.Project
		resourceName   = "mongodbatlas_project.test"
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfigWithFalseDefaultSettings(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
		},
	})
}

func TestAccProjectRSProject_CreateWithFalseDefaultAdvSettings(t *testing.T) {
	var (
		project        matlas.Project
		resourceName   = "mongodbatlas_project.test"
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfigWithFalseDefaultAdvSettings(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
		},
	})
}

func TestAccProjectRSProject_withUpdatedRole(t *testing.T) {
	var (
		resourceName    = "mongodbatlas_project.test"
		projectName     = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		roleName        = "GROUP_DATA_ACCESS_ADMIN"
		roleNameUpdated = "GROUP_READ_ONLY"
		clusterCount    = "0"
		teamsIds        = strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t); testCheckTeamsIds(t) },
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfigWithUpdatedRole(projectName, orgID, teamsIds[0], roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectConfigWithUpdatedRole(projectName, orgID, teamsIds[0], roleNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
		},
	})
}

func TestAccProjectRSProject_importBasic(t *testing.T) {
	var (
		projectName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		resourceName = "mongodbatlas_project.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,
					[]*matlas.ProjectTeam{},
					[]*testAPIKey{},
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasProjectImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"with_default_alerts_settings"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectDestroy(s *terraform.State) error {
	conn := testMongoDBClient.Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project" {
			continue
		}

		projectRes, _, _ := conn.Projects.GetOneProjectByName(context.Background(), rs.Primary.ID)
		if projectRes != nil {
			return fmt.Errorf("project (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasProjectImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

func testAccMongoDBAtlasProjectConfig(projectName, orgID string, teams []*matlas.ProjectTeam, apiKeys []*testAPIKey) string {
	var ts string

	for _, t := range teams {
		ts += fmt.Sprintf(`
		teams {
			team_id = "%s"
			role_names = %s
		}
		`, t.TeamID, strings.ReplaceAll(fmt.Sprintf("%+q", t.RoleNames), " ", ","))
	}

	for _, apiKey := range apiKeys {
		ts += fmt.Sprintf(`
		api_keys {
			api_key_id = "%s"
			role_names = %s
		}
		`, apiKey.id, strings.ReplaceAll(fmt.Sprintf("%+q", apiKey.roles), " ", ","))
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name  			 = "%s"
			org_id 			 = "%s"

			%s
		}
	`, projectName, orgID, ts)
}

func testAccMongoDBAtlasProjectConfigWithUpdatedRole(projectName, orgID, teamID, roleName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"

			teams {
				team_id = "%s"
				role_names = ["%s"]
			}
		}
	`, projectName, orgID, teamID, roleName)
}

func testAccMongoDBAtlasProjectConfigWithProjectOwner(projectName, orgID, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = "%[1]s"
			org_id 			 = "%[2]s"
		    project_owner_id = "%[3]s"
		}
	`, projectName, orgID, projectOwnerID)
}

func testAccMongoDBAtlasGovProjectConfigWithProjectOwner(projectName, orgID, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = "%[1]s"
			org_id 			 = "%[2]s"
		    project_owner_id = "%[3]s"
			region_usage_restrictions = "GOV_REGIONS_ONLY"
		}
	`, projectName, orgID, projectOwnerID)
}

func testAccMongoDBAtlasProjectConfigWithFalseDefaultSettings(projectName, orgID, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = "%[1]s"
			org_id 			 = "%[2]s"
			project_owner_id = "%[3]s"
			with_default_alerts_settings = false
		}
	`, projectName, orgID, projectOwnerID)
}

func testAccCheckMongoDBAtlasProjectExists(resourceName string, project *matlas.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testMongoDBClient.Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] projectID: %s", rs.Primary.ID)

		if projectResp, _, err := conn.Projects.GetOneProjectByName(context.Background(), rs.Primary.Attributes["name"]); err == nil {
			*project = *projectResp
			return nil
		}

		return fmt.Errorf("project (%s) does not exist", rs.Primary.ID)
	}
}

func testAccCheckMongoDBAtlasProjectAttributes(project *matlas.Project, projectName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if project.Name != projectName {
			return fmt.Errorf("bad project name: %s", project.Name)
		}

		return nil
	}
}

func testAccMongoDBAtlasProjectConfigWithFalseDefaultAdvSettings(projectName, orgID, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = "%[1]s"
			org_id 			 = "%[2]s"
			project_owner_id = "%[3]s"
			with_default_alerts_settings = false
			is_collect_database_specifics_statistics_enabled = false
			is_data_explorer_enabled = false
			is_extended_storage_sizes_enabled = false
			is_performance_advisor_enabled = false
			is_realtime_performance_panel_enabled = false
			is_schema_advisor_enabled = false
		}
	`, projectName, orgID, projectOwnerID)
}
