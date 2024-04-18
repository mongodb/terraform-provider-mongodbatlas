package auditing_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigGenericAuditing_basic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		auditFilter = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
		config      = configBasic(projectID, auditFilter, true, true)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:            config,
				ExternalProviders: mig.ExternalProviders(),
				Check:             resource.ComposeTestCheckFunc(checks(auditFilter, true, true)...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
