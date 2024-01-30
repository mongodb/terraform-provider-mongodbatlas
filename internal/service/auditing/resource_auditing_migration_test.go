package auditing_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccMigrationGenericAuditing_basic(t *testing.T) {
	var (
		auditing     matlas.Auditing
		resourceName = "mongodbatlas_auditing.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		auditAuth    = true
		auditFilter  = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
		enabled      = true
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:            configBasic(orgID, projectName, auditFilter, auditAuth, enabled),
				ExternalProviders: mig.ExternalProviders(),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &auditing),
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
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configBasic(orgID, projectName, auditFilter, auditAuth, enabled),
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
