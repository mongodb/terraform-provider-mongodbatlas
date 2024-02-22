package auditing_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccGenericAuditingDS_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_auditing.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
		auditAuth      = true
		auditFilter    = "{ 'atype': 'authenticate', 'param': {   'user': 'auditAdmin',   'db': 'admin',   'mechanism': 'SCRAM-SHA-1' }}"
		enabled        = true
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, auditFilter, auditAuth, enabled),
				Check: resource.ComposeTestCheckFunc(
					checkExists("mongodbatlas_auditing.test"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "audit_filter", auditFilter),
					resource.TestCheckResourceAttr(dataSourceName, "audit_authorization_success", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "configuration_type", "FILTER_JSON"),
				),
			},
		},
	})
}
