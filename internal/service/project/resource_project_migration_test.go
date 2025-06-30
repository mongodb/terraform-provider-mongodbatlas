package project_test

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigProject_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		config      = configBasic(orgID, projectName, "", false, nil, nil)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigProject_withTeams(t *testing.T) {
	var teamsIDs = strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
	if len(teamsIDs) < 2 {
		t.Skip("`MONGODB_ATLAS_TEAMS_IDS` must have 2 team ids for this acceptance testing")
	}

	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		clusterCount = "0"
		config       = configBasic(orgID, projectName, "", false,
			[]*admin.TeamRole{
				{
					TeamId:    &teamsIDs[0],
					RoleNames: &[]string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
				},
				{
					TeamId:    &teamsIDs[1],
					RoleNames: &[]string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
				},
			}, nil)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigProject_withFalseDefaultSettings(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
		projectName    = acc.RandomProjectName()
		config         = configWithFalseDefaultSettings(orgID, projectName, projectOwnerID)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasicOwnerID(t) },
		CheckDestroy: acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigProject_withLimits(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		config      = configWithLimits(orgID, projectName, []*admin.DataFederationLimit{
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "limits.0.name", "atlas.project.deployment.clusters"),
					resource.TestCheckResourceAttr(resourceName, "limits.0.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "limits.1.name", "atlas.project.deployment.nodesPerPrivateLinkRegion"),
					resource.TestCheckResourceAttr(resourceName, "limits.1.value", "2"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

// based on bug report: https://github.com/mongodb/terraform-provider-mongodbatlas/issues/2263
func TestMigGovProject_regionUsageRestrictionsDefault(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_GOV_ORG_ID")
		projectName = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckGovBasic(t) },
		CheckDestroy: acc.CheckDestroyProjectGov,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders("1.15.3"),
				Config:            configGovSimple(orgID, projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExistsGov(resourceName),
				),
			},
			{
				ExternalProviders: acc.ExternalProviders("1.16.0"),
				Config:            configGovSimple(orgID, projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExistsGov(resourceName),
				),
				ExpectError: regexp.MustCompile("Provider produced inconsistent result after apply"),
			},
			mig.TestStepCheckEmptyPlan(configGovSimple(orgID, projectName)),
		},
	})
}

func configGovSimple(orgID, projectName string) string {
	return acc.ConfigGovProvider() + fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id 			 = %[1]q
			name   			 = %[2]q
		}
	`, orgID, projectName)
}
