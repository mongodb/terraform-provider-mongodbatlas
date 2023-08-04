package mongodbatlas

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var _ plancheck.PlanCheck = debugPlan{}

type debugPlan struct{}

func (e debugPlan) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	rd, err := json.Marshal(req.Plan)
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("error marshaling machine-readable plan output: %s", err))
	}
	tflog.Info(ctx, fmt.Sprintf("req.Plan - %s\n", string(rd)))
}

func DebugPlan() plancheck.PlanCheck {
	return debugPlan{}
}

func TestAccRSProject_Migration_NoProps(t *testing.T) {
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
						VersionConstraint: "1.11.0",
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
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config: fmt.Sprintf(`resource "mongodbatlas_project" "test" {
					name   = "%s"
					org_id = "%s"
				  }`, projectName, orgID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccRSProject_Migration_Teams(t *testing.T) {
	var (
		project         matlas.Project
		resourceName    = "mongodbatlas_project.test"
		projectName     = acctest.RandomWithPrefix("test-acc-teams")
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterCount    = "0"
		teamsIds        = strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
		configWithTeams = testAccMongoDBAtlasProjectConfig(projectName, orgID,
			[]*matlas.ProjectTeam{
				{
					TeamID:    teamsIds[0],
					RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
				},
				{
					TeamID:    teamsIds[1],
					RoleNames: []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
				},
			})
	)
	if len(teamsIds) < 2 {
		t.Skip("`MONGODB_ATLAS_TEAMS_IDS` must have 2 team ids for this acceptance testing")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.11.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: configWithTeams,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   configWithTeams,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccRSProject_Migration_WithFalseDefaultSettings(t *testing.T) {
	var (
		project         matlas.Project
		resourceName    = "mongodbatlas_project.test"
		projectName     = acctest.RandomWithPrefix("tf-acc-project")
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID  = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
		configWithTeams = testAccMongoDBAtlasProjectConfigWithFalseDefaultSettings(projectName, orgID, projectOwnerID)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasicOwnerID(t) },
		CheckDestroy: testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.11.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: configWithTeams,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   configWithTeams,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccRSProject_Migration_WithLimits(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project.test"
		projectName  = acctest.RandomWithPrefix("tf-acc-project")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		config       = testAccMongoDBAtlasProjectConfigWithLimits(projectName, orgID, []*projectLimit{
			{
				name:  "atlas.project.deployment.clusters",
				value: 1,
			},
			{
				name:  "atlas.project.deployment.nodesPerPrivateLinkRegion",
				value: 2,
			},
		})
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.11.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "limits.0.name", "atlas.project.deployment.clusters"),
					resource.TestCheckResourceAttr(resourceName, "limits.0.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "limits.1.name", "atlas.project.deployment.nodesPerPrivateLinkRegion"),
					resource.TestCheckResourceAttr(resourceName, "limits.1.value", "2"),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProviderV6Factories,
				Config:                   config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}
