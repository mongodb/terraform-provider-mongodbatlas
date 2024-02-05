package rolesorgid_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSOrgID_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_roles_org_id.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDS(),
				Check:  resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
			},
		},
	})
}

func configDS() string {
	return `data "mongodbatlas_roles_org_id" "test" {}`
}
