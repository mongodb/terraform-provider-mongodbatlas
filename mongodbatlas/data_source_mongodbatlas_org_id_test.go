package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigDSOrgID_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_roles_org_id.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name           = fmt.Sprintf("test-acc-%s@mongodb.com", acctest.RandString(10))
		initialRole    = []string{"ORG_OWNER"}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasOrgInvitationDestroy,
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
