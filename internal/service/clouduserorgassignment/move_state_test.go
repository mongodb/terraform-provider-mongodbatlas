package clouduserorgassignment_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccCloudUserOrgAssignmentRS_moveFromOrgInvitation(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	username := "test-move-from-org-invitation@example.com"
	roles := []string{"ORG_MEMBER", "ORG_GROUP_CREATOR"}
	teamsIDs := []string{acc.GetProjectTeamsIDsWithPos(0), acc.GetProjectTeamsIDsWithPos(1)}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configOrgInvitationFirst(orgID, username, roles, teamsIDs),
			},
			{
				Config: configMoveFromOrgInvitationSecond(orgID, username, roles),
				Check: resource.ComposeTestCheckFunc(
					cloudUserOrgAssignmentChecks(orgID, username, "PENDING", roles),
					resource.TestCheckResourceAttr("mongodbatlas_cloud_user_org_assignment.test", "team_ids.#", "2"),
				),
			},
		},
	})
}

func configOrgInvitationFirst(orgID, username string, roles, teamsIDs []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	teamsIDsStr := `"` + strings.Join(teamsIDs, `", "`) + `"`

	return fmt.Sprintf(`
resource "mongodbatlas_org_invitation" "old" {
  org_id   = "%s"
  username = "%s"
  roles    = [%s]
  teams_ids = [%s]
}
`, orgID, username, rolesStr, teamsIDsStr)
}

func configMoveFromOrgInvitationSecond(orgID, username string, roles []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`

	return fmt.Sprintf(`
moved {
  from = mongodbatlas_org_invitation.old
  to   = mongodbatlas_cloud_user_org_assignment.test
}

resource "mongodbatlas_cloud_user_org_assignment" "test" {
  org_id   = "%s"
  username = "%s"
  roles = {
    org_roles = [%s]
  }
}
`, orgID, username, rolesStr)
}
