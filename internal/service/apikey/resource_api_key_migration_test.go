package apikey_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigConfigAPIKey_basic(t *testing.T) {
	skipSupersededMigTest(t)
	var (
		resourceName = "mongodbatlas_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		description  = acc.RandomName()
		roleName     = "ORG_MEMBER"
		config       = configBasic(orgID, description, roleName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:            config,
				ExternalProviders: mig.ExternalProviders(),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

// skipSupersededMigTest skips the manual-era migration tests while mongodbatlas_api_key is served by
// the autogen implementation (spike pilot): superseded by the standalone TestMigAPIKey_* tests in
// internal/serviceapi/apikey, restored or removed at the full migration.
func skipSupersededMigTest(t *testing.T) {
	t.Helper()
	t.Skip("superseded by TestMigAPIKey_* in internal/serviceapi/apikey (autogen lift-and-shift pilot)")
}
