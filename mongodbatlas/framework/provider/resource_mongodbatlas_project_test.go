package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"mongodbatlas": providerserver.NewProtocol6WithError(New()()),
	}
}

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

func Test_DebugPlan(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		// ExternalProviders: map[string]r.ExternalProvider{
		// 	"random": {
		// 		Source: "registry.terraform.io/hashicorp/random",
		// 	},
		// },
		// ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "mongodbatlas_project" "main" {
					name   = "repro-region-config-order-issue"
					org_id = "63bec56c014da65b8f73c05e"
				  }`,
				PlanOnly: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						DebugPlan(),
					},
				},
			},
		},
	})
}

func TestResource_UpgradeFromVersion(t *testing.T) {
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
				// ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ProtoV6ProviderFactories: protoV6ProviderFactories(),
				Config: `resource "mongodbatlas_project" "test" {
					name   = "tf-test-project-migration3"
					org_id = "63bec56c014da65b8f73c05e"
				  }`,
				PlanOnly: true,
			},
		},
	})
}

// func TestAccProjectRSProject_CreateWithProjectOwner(t *testing.T) {
// 	var (
// 		project        matlas.Project
// 		resourceName   = "mongodbatlas_project.test"
// 		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
// 		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
// 	)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckBasicOwnerID(t) },
// 		ProviderFactories: TestAccProtoV6ProviderFactories,
// 		// CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMongoDBAtlasProjectConfigWithProjectOwner(projectName, orgID, projectOwnerID),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
// 					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
// 					resource.TestCheckResourceAttr(resourceName, "name", projectName),
// 					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccProjectRSProject_CreateWithFalseDefaultSettings(t *testing.T) {
// 	var (
// 		project        matlas.Project
// 		resourceName   = "mongodbatlas_project.test"
// 		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
// 		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
// 	)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckBasicOwnerID(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMongoDBAtlasProjectConfigWithFalseDefaultSettings(projectName, orgID, projectOwnerID),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
// 					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
// 					resource.TestCheckResourceAttr(resourceName, "name", projectName),
// 					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccProjectRSProject_CreateWithFalseDefaultAdvSettings(t *testing.T) {
// 	var (
// 		project        matlas.Project
// 		resourceName   = "mongodbatlas_project.test"
// 		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
// 		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
// 	)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckBasicOwnerID(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMongoDBAtlasProjectConfigWithFalseDefaultAdvSettings(projectName, orgID, projectOwnerID),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
// 					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
// 					resource.TestCheckResourceAttr(resourceName, "name", projectName),
// 					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccProjectRSProject_withUpdatedRole(t *testing.T) {
// 	var (
// 		resourceName    = "mongodbatlas_project.test"
// 		projectName     = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
// 		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 		roleName        = "GROUP_DATA_ACCESS_ADMIN"
// 		roleNameUpdated = "GROUP_READ_ONLY"
// 		clusterCount    = "0"
// 		teamsIds        = strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
// 	)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckBasic(t); testCheckTeamsIds(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMongoDBAtlasProjectConfigWithUpdatedRole(projectName, orgID, teamsIds[0], roleName),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceName, "name", projectName),
// 					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
// 					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
// 				),
// 			},
// 			{
// 				Config: testAccMongoDBAtlasProjectConfigWithUpdatedRole(projectName, orgID, teamsIds[0], roleNameUpdated),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr(resourceName, "name", projectName),
// 					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
// 					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccProjectRSProject_importBasic(t *testing.T) {
// 	var (
// 		projectName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
// 		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
// 		resourceName = "mongodbatlas_project.test"
// 	)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheckBasic(t) },
// 		ProviderFactories: testAccProviderFactories,
// 		CheckDestroy:      testAccCheckMongoDBAtlasProjectDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,
// 					[]*matlas.ProjectTeam{},
// 					[]*apiKey{},
// 				),
// 			},
// 			{
// 				ResourceName:            resourceName,
// 				ImportStateIdFunc:       testAccCheckMongoDBAtlasProjectImportStateIDFunc(resourceName),
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"with_default_alerts_settings"},
// 			},
// 		},
// 	})
// }
