package projectinvitation_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccProjectDSProjectInvitation_basic(t *testing.T) {
	var (
		dataSourceName = "mongodbatlas_project_invitation.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
		name           = acc.RandomEmail()
		initialRole    = []string{"GROUP_OWNER"}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasProjectInvitationConfig(orgID, projectName, name, initialRole),
				Check: resource.ComposeAggregateTestCheckFunc(
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
