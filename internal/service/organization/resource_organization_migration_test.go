package organization_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

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
		config       = configBasic(orgOwnerID, name, description, roleName, false, nil)
	)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
