package project_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigProjectRS_withNoProps(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config: fmt.Sprintf(`resource "mongodbatlas_project" "test" {
					name   = "%s"
					org_id = "%s"
				  }`, projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config: fmt.Sprintf(`resource "mongodbatlas_project" "test" {
					name   = "%s"
					org_id = "%s"
				  }`, projectName, orgID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestMigProjectRS_withTeams(t *testing.T) {
	var teamsIDs = strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
	if len(teamsIDs) < 2 {
		t.Skip("`MONGODB_ATLAS_TEAMS_IDS` must have 2 team ids for this acceptance testing")
	}

	var (
		project         admin.Group
		resourceName    = "mongodbatlas_project.test"
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName     = acc.RandomProjectName()
		clusterCount    = "0"
		configWithTeams = acc.ConfigProject(projectName, orgID,
			[]*admin.TeamRole{
				{
					TeamId:    &teamsIDs[0],
					RoleNames: &[]string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
				},
				{
					TeamId:    &teamsIDs[1],
					RoleNames: &[]string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
				},
			})
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configWithTeams,
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &project),
					acc.CheckProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configWithTeams,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestMigProjectRS_withFalseDefaultSettings(t *testing.T) {
	var (
		project         admin.Group
		resourceName    = "mongodbatlas_project.test"
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID  = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
		projectName     = acc.RandomProjectName()
		configWithTeams = acc.ConfigProjectWithFalseDefaultSettings(projectName, orgID, projectOwnerID)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasicOwnerID(t) },
		CheckDestroy: acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configWithTeams,
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &project),
					acc.CheckProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configWithTeams,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestMigProjectRS_withLimits(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		config       = acc.ConfigProjectWithLimits(projectName, orgID, []*admin.DataFederationLimit{
			{
				Name:  "atlas.project.deployment.clusters",
				Value: 1,
			},
			{
				Name:  "atlas.project.deployment.nodesPerPrivateLinkRegion",
				Value: 2,
			},
		})
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
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
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
