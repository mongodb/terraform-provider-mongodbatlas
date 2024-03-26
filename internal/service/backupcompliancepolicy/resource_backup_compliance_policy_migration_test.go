package backupcompliancepolicy_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigBackupCompliancePolicy_basic(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName() // No ProjectIDExecution to avoid conflicts with backup comliance policy
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
		config         = configBasic(projectName, orgID, projectOwnerID)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { mig.PreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             resource.ComposeTestCheckFunc(checks()...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
