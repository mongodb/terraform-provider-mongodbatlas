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
		config       = configBasic(orgID, name, []string{username})
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
