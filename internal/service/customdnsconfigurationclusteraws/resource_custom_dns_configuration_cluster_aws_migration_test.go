package customdnsconfigurationclusteraws_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationConfigRSCustomDNSConfigurationAWS_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		config      = configBasic(orgID, projectName, true)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			mig.TestStep(config),
		},
	})
}
