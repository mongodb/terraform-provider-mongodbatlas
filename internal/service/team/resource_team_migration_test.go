package team_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigConfigTeams_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_team.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		username     = os.Getenv("MONGODB_ATLAS_USERNAME")
		name         = acc.RandomName()
		usernames    = []string{username}
		config       = configBasic(orgID, name, &usernames)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckAtlasUsername(t) },
		CheckDestroy: acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigConfigTeams_usernamesDeprecation(t *testing.T) {
	var (
		resourceName = "mongodbatlas_team.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		username     = os.Getenv("MONGODB_ATLAS_USERNAME")
		name         = acc.RandomName()
		usernames    = []string{username}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckAtlasUsername(t) },
		CheckDestroy: acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configBasic(orgID, name, &usernames),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "usernames.*", username),
				),
			},
			{
				Config: configBasic(orgID, name, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					// usernames should still be present in state (computed) but not in config
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "usernames.*", username),
				),
			},
		},
	})
}
