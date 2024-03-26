package auditing_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigGenericAuditing_basic(t *testing.T) {
	var (
		projectID   = mig.ProjectIDGlobal(t)
		auditAuth   = true
		auditFilter = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
		enabled     = true
		config      = configBasic(projectID, auditFilter, auditAuth, enabled)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:            config,
				ExternalProviders: mig.ExternalProviders(),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_filter"),
					resource.TestCheckResourceAttrSet(resourceName, "audit_authorization_success"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "audit_filter", auditFilter),
					resource.TestCheckResourceAttr(resourceName, "audit_authorization_success", "true"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configuration_type", "FILTER_JSON"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
