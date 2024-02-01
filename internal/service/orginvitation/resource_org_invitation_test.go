package orginvitation_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSOrgInvitation_basic(t *testing.T) {
	var (
		invitation   matlas.Invitation
		resourceName = "mongodbatlas_org_invitation.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name         = fmt.Sprintf("test-acc-%s@mongodb.com", acctest.RandString(10))
		initialRole  = []string{"ORG_OWNER"}
		updateRoles  = []string{"ORG_GROUP_CREATOR", "ORG_BILLING_ADMIN"}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyOrgInvitation,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name, initialRole),
				Check: resource.ComposeTestCheckFunc(
					checkExists(t, resourceName, &invitation),
					resource.TestCheckResourceAttrSet(resourceName, "invitation_id"),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "username", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", initialRole[0]),
				),
			},
			{
				Config: configBasic(orgID, name, updateRoles),
				Check: resource.ComposeTestCheckFunc(
					checkExists(t, resourceName, &invitation),
					resource.TestCheckResourceAttrSet(resourceName, "invitation_id"),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "username", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", updateRoles[0]),
					resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", updateRoles[1]),
				),
			},
		},
	})
}

func TestAccConfigRSOrgInvitation_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_org_invitation.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name         = fmt.Sprintf("test-acc-%s@mongodb.com", acctest.RandString(10))
		initialRole  = []string{"ORG_OWNER"}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyOrgInvitation,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name, initialRole),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "invitation_id"),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "username", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", initialRole[0]),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkExists(t *testing.T, resourceName string, invitation *matlas.Invitation) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		orgID := ids["org_id"]
		username := ids["username"]
		invitationID := ids["invitation_id"]
		if orgID == "" && username == "" && invitationID == "" {
			return fmt.Errorf("no ID is set")
		}
		t.Logf("orgID: %s", orgID)
		t.Logf("username: %s", username)
		t.Logf("invitationID: %s", invitationID)
		invitationResp, _, err := acc.Conn().Organizations.Invitation(context.Background(), orgID, invitationID)
		if err == nil {
			*invitation = *invitationResp
			return nil
		}
		return fmt.Errorf("invitation(%s) does not exist", invitationID)
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["org_id"], rs.Primary.Attributes["username"]), nil
	}
}

func configBasic(orgID, username string, roles []string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_org_invitation" "test" {
			org_id   = %[1]q
			username = %[2]q
			roles  	 = ["%[3]s"]
		}`, orgID, username,
		strings.Join(roles, `", "`),
	)
}
