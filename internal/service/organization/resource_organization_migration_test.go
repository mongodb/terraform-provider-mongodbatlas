package organization_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigConfigRSOrganization_Basic(t *testing.T) {
	acc.SkipTestForCI(t) // affects the org

	var (
		resourceName = "mongodbatlas_organization.test"
		orgOwnerID   = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name         = acc.RandomName()
		description  = "test Key for Acceptance tests"
		roleName     = "ORG_OWNER"
	)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, name, description, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, name, description, roleName),
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
