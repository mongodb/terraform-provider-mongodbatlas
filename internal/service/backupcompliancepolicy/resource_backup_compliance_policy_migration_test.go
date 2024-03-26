package backupcompliancepolicy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigBackupCompliancePolicy_basic(t *testing.T) {
	var (
		projectID = mig.ProjectIDGlobal(t)
		config    = configBasic(projectID)
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
