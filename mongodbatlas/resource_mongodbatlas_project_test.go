package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasProject_basic(t *testing.T) {
	var (
		project      matlas.Project
		resourceName = "mongodbatlas_project.test"
		projectName  = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterCount = "0"
		teamsIds     = strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
		apiKeysIds   = strings.Split(os.Getenv("MONGODB_ATLAS_API_KEYS_IDS"), ",")
	)
	if len(teamsIds) < 3 {
		t.Skip("`MONGODB_ATLAS_TEAMS_IDS` must have 3 team ids for this acceptance testing")
	}
	if len(apiKeysIds) < 2 {
		t.Skip("`MONGODB_ATLAS_API_KEYS_IDS` must have 2 api key ids for this acceptance testing")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); checkTeamsIds(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
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
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,

					[]*matlas.ProjectTeam{
						{
							TeamID:    teamsIds[0],
							RoleNames: []string{"GROUP_OWNER"},
						},
						{
							TeamID:    teamsIds[1],
							RoleNames: []string{"GROUP_DATA_ACCESS_READ_WRITE"},
						},
						{
							TeamID:    teamsIds[2],
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
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
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,

					[]*matlas.ProjectTeam{
						{
							TeamID:    teamsIds[0],
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_READ_ONLY"},
						},
						{
							TeamID:    teamsIds[1],
							RoleNames: []string{"GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"},
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
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID, []*matlas.ProjectTeam{}, []*apiKey{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasProject_CreateWithProjectOwner(t *testing.T) {
	var (
		project        matlas.Project
		resourceName   = "mongodbatlas_project.test"
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
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

func TestAccResourceMongoDBAtlasProject_CreateWithFalseDefaultSettings(t *testing.T) {
	var (
		project        matlas.Project
		resourceName   = "mongodbatlas_project.test"
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
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

func TestAccResourceMongoDBAtlasProject_CreateWithFalseDefaultAdvSettings(t *testing.T) {
	var (
		project        matlas.Project
		resourceName   = "mongodbatlas_project.test"
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
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

func TestAccResourceMongoDBAtlasProject_withUpdatedRole(t *testing.T) {
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
		PreCheck:          func() { testAccPreCheck(t); checkTeamsIds(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
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

func TestAccResourceMongoDBAtlasProject_importBasic(t *testing.T) {
	var (
		projectName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		resourceName = "mongodbatlas_project.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,
					[]*matlas.ProjectTeam{},
					[]*apiKey{},
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

func testAccCheckMongoDBAtlasProjectExists(resourceName string, project *matlas.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

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

func TestAccResourceMongoDBAtlasProject_CreateWithAdvancedCluster(t *testing.T) {
	var (
		project             matlas.Project
		cluster             matlas.AdvancedCluster
		clusterResourceName = "mongodbatlas_advanced_cluster.test"
		resourceName        = "mongodbatlas_project.test"
		clusterName         = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		projectName         = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID               = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID      = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectConfigWithAdvancedCluster(projectName, orgID, projectOwnerID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasAdvancedClusterExists(clusterResourceName, &cluster),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectAttributes(project *matlas.Project, projectName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if project.Name != projectName {
			return fmt.Errorf("bad project name: %s", project.Name)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasProjectDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

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

func testAccMongoDBAtlasProjectConfig(projectName, orgID string, teams []*matlas.ProjectTeam, apiKeys []*apiKey) string {
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

func testAccMongoDBAtlasProjectConfigWithFalseDefaultAdvSettings(projectName, orgID, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = "%[1]s"
			org_id 			 = "%[2]s"
			project_owner_id = "%[3]s"
			with_default_alerts_settings = false
			is_collect_database_specifics_statistics_enabled = false
			is_data_explorer_enabled = false
			is_performance_advisor_enabled = false
			is_realtime_performance_panel_enabled = false
			is_schema_advisor_enabled = false
		}
	`, projectName, orgID, projectOwnerID)
}

func testAccMongoDBAtlasProjectConfigWithAdvancedCluster(projectName, orgID, projectOwnerID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name                         = %[1]q
			org_id                       = %[2]q
			project_owner_id             = %[3]q
			with_default_alerts_settings = false
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = mongodbatlas_project.test.id
			name         = %[4]q
			cluster_type = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
				}
			}
		}
	`, projectName, orgID, projectOwnerID, clusterName)
}
