package teams_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationTeams_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_teams.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name         = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		username     = os.Getenv("MONGODB_ATLAS_USERNAME")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckAtlasUsername(t) },
		CheckDestroy: acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configBasic(orgID, name, []string{username}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "usernames.#", "1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configBasic(orgID, name, []string{username}),
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
