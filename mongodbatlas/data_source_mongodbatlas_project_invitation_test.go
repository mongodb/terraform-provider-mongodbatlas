package mongodbatlas

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectDSProjectInvitation_basic(t *testing.T) {
	var (
		dataSourceName = "mongodbatlas_project_invitation.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc")
		name           = fmt.Sprintf("test-acc-%s@mongodb.com", acctest.RandString(10))
		initialRole    = []string{"GROUP_OWNER"}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectInvitationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasProjectInvitationConfig(orgID, projectName, name, initialRole),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "username"),
					resource.TestCheckResourceAttrSet(dataSourceName, "invitation_id"),
					resource.TestCheckResourceAttr(dataSourceName, "username", name),
					resource.TestCheckResourceAttr(dataSourceName, "roles.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceMongoDBAtlasProjectInvitationConfig(orgID, projectName, username string, roles []string) string {
	return fmt.Sprintf(`

		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_project_invitation" "test" {
			project_id = mongodbatlas_project.test.id
			username   = %[3]q
			roles  	 = ["%[4]s"]
		}

		data "mongodbatlas_project_invitation" "test" {
			project_id    = mongodbatlas_project_invitation.test.project_id
			username 	  = mongodbatlas_project_invitation.test.username
			invitation_id = mongodbatlas_project_invitation.test.invitation_id
		}`, orgID, projectName, username,
		strings.Join(roles, `", "`),
	)
}
