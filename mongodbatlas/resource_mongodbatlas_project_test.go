package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasProject_basic(t *testing.T) {
	var project matlas.Project

	resourceName := "mongodbatlas_project.test"
	projectName := fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	clusterCount := "0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPropjectConfig(projectName, orgID,
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
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				Config: testAccMongoDBAtlasPropjectConfig(projectName, orgID,

					[]*matlas.ProjectTeam{
						{
							TeamID:    "5e1dd7b4f2a30ba80a70cd3a",
							RoleNames: []string{"GROUP_OWNER"},
						},
						{
							TeamID:    "5ddd59c29ccf64e86b7f68df",
							RoleNames: []string{"GROUP_DATA_ACCESS_READ_WRITE"},
						},
						{
							TeamID:    "5e0fa8c99ccf641c722fe683",
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
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
				Config: testAccMongoDBAtlasPropjectConfig(projectName, orgID,

					[]*matlas.ProjectTeam{
						{
							TeamID:    "5e1dd7b4f2a30ba80a70cd3a",
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_READ_ONLY"},
						},
						{
							TeamID:    "5ddd59c29ccf64e86b7f68df",
							RoleNames: []string{"GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"},
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
				Config: testAccMongoDBAtlasPropjectConfig(projectName, orgID, []*matlas.ProjectTeam{}),
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

func TestAccResourceMongoDBAtlasProject_withUpdatedRole(t *testing.T) {
	resourceName := "mongodbatlas_project.test"
	projectName := fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	teamID := "5e0fa8c99ccf641c722fe683"
	roleName := "GROUP_DATA_ACCESS_ADMIN"
	roleNameUpdated := "GROUP_READ_ONLY"
	clusterCount := "0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPropjectConfigWithUpdatedRole(projectName, orgID, teamID, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				Config: testAccMongoDBAtlasPropjectConfigWithUpdatedRole(projectName, orgID, teamID, roleNameUpdated),
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

	projectName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")

	resourceName := "mongodbatlas_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPropjectConfig(projectName, orgID,
					[]*matlas.ProjectTeam{},
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasProjectImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectExists(resourceName string, project *matlas.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

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

func testAccCheckMongoDBAtlasProjectDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

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
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccMongoDBAtlasPropjectConfig(projectName, orgID string, teams []*matlas.ProjectTeam) string {
	var ts string

	for _, t := range teams {
		ts += fmt.Sprintf(`
		teams {
			team_id = "%s"
			role_names = %s
		}
		`, t.TeamID, strings.ReplaceAll(fmt.Sprintf("%+q", t.RoleNames), " ", ","))
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"

			%s
		}
	`, projectName, orgID, ts)
}

func testAccMongoDBAtlasPropjectConfigWithUpdatedRole(projectName, orgID, teamID, roleName string) string {
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
