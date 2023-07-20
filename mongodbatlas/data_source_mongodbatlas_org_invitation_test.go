package mongodbatlas

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigDSOrgInvitation_basic(t *testing.T) {
	var (
		dataSourceName = "mongodbatlas_org_invitation.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name           = fmt.Sprintf("test-acc-%s@mongodb.com", acctest.RandString(10))
		initialRole    = []string{"ORG_OWNER"}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
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
