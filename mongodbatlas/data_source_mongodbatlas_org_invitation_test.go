package mongodbatlas

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMongoDBAtlasOrgInvitation_basic(t *testing.T) {
	var (
		dataSourceName = "mongodbatlas_org_invitation.test"
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
				Config: testAccDataSourceMongoDBAtlasOrgInvitationConfig(orgID, name, initialRole),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "username"),
					resource.TestCheckResourceAttrSet(dataSourceName, "invitation_id"),
					resource.TestCheckResourceAttr(dataSourceName, "username", name),
					resource.TestCheckResourceAttr(dataSourceName, "roles.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceMongoDBAtlasOrgInvitationConfig(orgID, username string, roles []string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_org_invitation" "test" {
			org_id   = %[1]q
			username = %[2]q
			roles  	 = ["%[3]s"]
		}

		data "mongodbatlas_org_invitation" "test" {
			org_id        = mongodbatlas_org_invitation.test.org_id
			username      = mongodbatlas_org_invitation.test.username
			invitation_id = mongodbatlas_org_invitation.test.invitation_id
		}`, orgID, username,
		strings.Join(roles, `", "`),
	)
}
