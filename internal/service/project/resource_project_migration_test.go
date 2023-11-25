package project_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"go.mongodb.org/atlas-sdk/v20231115001/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccMigrationProjectRS_NoProps(t *testing.T) {
	var (
		resourceName          = "mongodbatlas_project.test"
		projectName           = acctest.RandomWithPrefix("test-acc-migration")
		orgID                 = os.Getenv("MONGODB_ATLAS_ORG_ID")
		lastVersionConstraint = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
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
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config: fmt.Sprintf(`resource "mongodbatlas_project" "test" {
					name   = "%s"
					org_id = "%s"
				  }`, projectName, orgID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationProjectRS_Teams(t *testing.T) {
	var teamsIds = strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
	if len(teamsIds) < 2 {
		t.Skip("`MONGODB_ATLAS_TEAMS_IDS` must have 2 team ids for this acceptance testing")
	}

	var (
		project         matlas.Project
		resourceName    = "mongodbatlas_project.test"
		projectName     = acctest.RandomWithPrefix("test-acc-teams")
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterCount    = "0"
		configWithTeams = acc.ConfigProject(projectName, orgID,
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
		lastVersionConstraint = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: configWithTeams,
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
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationProjectRS_WithFalseDefaultSettings(t *testing.T) {
	var (
		project               matlas.Project
		resourceName          = "mongodbatlas_project.test"
		projectName           = acctest.RandomWithPrefix("tf-acc-project")
		orgID                 = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID        = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
		configWithTeams       = acc.ConfigProjectWithFalseDefaultSettings(projectName, orgID, projectOwnerID)
		lastVersionConstraint = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicOwnerIDMigration(t) },
		CheckDestroy: acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: configWithTeams,
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
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationProjectRS_WithLimits(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project.test"
		projectName  = acctest.RandomWithPrefix("tf-acc-project")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
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
		lastVersionConstraint = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
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
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationProjectRSProjectIPAccesslist_SettingIPAddress(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_access_list.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)
	lastVersionConstraint := os.Getenv("MONGODB_ATLAS_LAST_VERSION")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: acc.CheckDestroyProjectIPAccessList,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: acc.ConfigProjectIPAccessListWithIPAddress(orgID, projectName, ipAddress, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   acc.ConfigProjectIPAccessListWithIPAddress(orgID, projectName, ipAddress, comment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationProjectRSProjectIPAccessList_SettingCIDRBlock(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_access_list.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	cidrBlock := fmt.Sprintf("179.154.226.%d/32", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for cidrBlock (%s)", cidrBlock)
	lastVersionConstraint := os.Getenv("MONGODB_ATLAS_LAST_VERSION")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: acc.CheckDestroyProjectIPAccessList,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: acc.ConfigProjectIPAccessListWithCIDRBlock(orgID, projectName, cidrBlock, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   acc.ConfigProjectIPAccessListWithCIDRBlock(orgID, projectName, cidrBlock, comment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationProjectRSProjectIPAccessList_Multiple_SettingMultiple(t *testing.T) {
	resourceName := "mongodbatlas_project_ip_access_list.test_1"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	const ipWhiteListCount = 20
	accessList := make([]map[string]string, 0)
	lastVersionConstraint := os.Getenv("MONGODB_ATLAS_LAST_VERSION")

	for i := 0; i < ipWhiteListCount; i++ {
		entry := make(map[string]string)
		entryName := ""
		ipAddr := ""

		if i%2 == 0 {
			entryName = "cidr_block"
			entry["cidr_block"] = fmt.Sprintf("%d.2.3.%d/32", i, acctest.RandIntRange(0, 255))
			ipAddr = entry["cidr_block"]
		} else {
			entryName = "ip_address"
			entry["ip_address"] = fmt.Sprintf("%d.2.3.%d", i, acctest.RandIntRange(0, 255))
			ipAddr = entry["ip_address"]
		}
		entry["comment"] = fmt.Sprintf("TestAcc for %s (%s)", entryName, ipAddr)

		accessList = append(accessList, entry)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: acc.CheckDestroyProjectIPAccessList,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: acc.ConfigProjectIPAccessListWithMultiple(projectName, orgID, accessList, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   acc.ConfigProjectIPAccessListWithMultiple(projectName, orgID, accessList, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}
