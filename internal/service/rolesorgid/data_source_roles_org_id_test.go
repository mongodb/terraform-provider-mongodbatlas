package rolesorgid_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSOrgID_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_roles_org_id.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name           = fmt.Sprintf("test-acc-%s@mongodb.com", acctest.RandString(10))
		initialRole    = []string{"ORG_OWNER"}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyOrgInvitation,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasOrgIDConfig(orgID, name, initialRole),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
				),
			},
		},
	})
}

func testAccDataSourceMongoDBAtlasOrgIDConfig(orgID, username string, roles []string) string {
	return (`
	data "mongodbatlas_roles_org_id" "test" {
	}
	
	output "org_id" {
	 value = data.mongodbatlas_roles_org_id.test.org_id
	}`)
}
