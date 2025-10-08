package rolesorgid_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSOrgID_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_roles_org_id.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDS(),
				Check:  resource.TestCheckResourceAttr(dataSourceName, "org_id", orgID),
			},
		},
	})
}

func configDS() string {
	return `data "mongodbatlas_roles_org_id" "test" {}`
}
