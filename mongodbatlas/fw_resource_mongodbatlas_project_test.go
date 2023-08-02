package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// This test only helps to see what the plan is for debugging purposes
// func Test_DebugPlan(t *testing.T) {
// 	var (
// 		// resourceName = "mongodbatlas_project.test"
// 		projectName = acctest.RandomWithPrefix("test-acc-migration")
// 		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 	)
// 	t.Parallel()

// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: fmt.Sprintf(`resource "mongodbatlas_project" "main" {
// 					name   = "%s"
// 					org_id = "%s"
// 				  }`, projectName, orgID),
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

func TestAccRSProject_Teams(t *testing.T) {
	var (
		project      matlas.Project
		resourceName = "mongodbatlas_project.test"
		projectName  = acctest.RandomWithPrefix("test-acc-teams")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterCount = "0"
		teamsIds     = strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t); testCheckTeamsIds(t) },
		CheckDestroy: testAccCheckMongoDBAtlasProjectDestroy2,

		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
				Config: testAccMongoDBAtlasProjectConfig2(projectName, orgID,
					[]*matlas.ProjectTeam{
						{
							TeamID:    teamsIds[0],
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamID:    teamsIds[1],
							RoleNames: []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists2(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes2(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectDestroy2(s *terraform.State) error {
	conn := testMongoDBClient.(*MongoDBClient).Atlas

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

func testAccMongoDBAtlasProjectConfig2(projectName, orgID string, teams []*matlas.ProjectTeam) string {
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
			name  			 = "%s"
			org_id 			 = "%s"

			%s
		}
	`, projectName, orgID, ts)
}

func testAccCheckMongoDBAtlasProjectExists2(resourceName string, project *matlas.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testMongoDBClient.(*MongoDBClient).Atlas

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

func testAccCheckMongoDBAtlasProjectAttributes2(project *matlas.Project, projectName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if project.Name != projectName {
			return fmt.Errorf("bad project name: %s", project.Name)
		}

		return nil
	}
}
