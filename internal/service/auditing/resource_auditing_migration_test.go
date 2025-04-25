package auditing_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigGenericAuditing_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.34.0") // Version where JSON comparison in audit_filter field in mongodbatlas_auditing was fixed
	var (
		projectID   = acc.ProjectIDExecution(t)
		auditFilter = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
		config      = configBasic(projectID, auditFilter, true, true)
	)

	// Serial so it doesn't conflict with TestAccGenericAuditing_basic
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:            config,
				ExternalProviders: mig.ExternalProviders(),
				Check:             resource.ComposeAggregateTestCheckFunc(checks(auditFilter, true, true)...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
