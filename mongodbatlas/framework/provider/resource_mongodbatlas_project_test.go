package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var _ plancheck.PlanCheck = debugPlan{}

type debugPlan struct{}

func (e debugPlan) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	rd, err := json.Marshal(req.Plan)
	if err != nil {
		fmt.Println("error marshalling machine-readable plan output:", err)
	}
	fmt.Printf("req.Plan - %s\n", string(rd))
}

func DebugPlan() plancheck.PlanCheck {
	return debugPlan{}
}

// This test only helps to see what the plan is for debugging purposes
//
// func Test_DebugPlan(t *testing.T) {
// 	t.Parallel()

// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: `resource "mongodbatlas_project" "main" {
// 					name   = "repro-region-config-order-issue"
// 					org_id = "63bec56c014da65b8f73c05e"
// 				  }`,
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

func TestProjectResource_FrameworkMigration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.10.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: `resource "mongodbatlas_project" "test" {
					name   = "tf-test-project-migration3"
					org_id = "63bec56c014da65b8f73c05e"
				  }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mongodbatlas_project.test", "org_id", "63bec56c014da65b8f73c05e"),
				),
			},
			{
				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
				Config: `resource "mongodbatlas_project" "test" {
					name   = "tf-test-project-migration3"
					org_id = "63bec56c014da65b8f73c05e"
				  }`,
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

func testAccMongoDBAtlasProjectConfigWithProjectOwner(projectName, orgID, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = "%[1]s"
			org_id 			 = "%[2]s"
		    project_owner_id = "%[3]s"
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
