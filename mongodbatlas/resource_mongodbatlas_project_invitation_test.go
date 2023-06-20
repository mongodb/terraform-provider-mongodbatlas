package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccProjectRSProjectInvitation_basic(t *testing.T) {
	var (
		invitation   matlas.Invitation
		resourceName = "mongodbatlas_project_invitation.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = fmt.Sprintf("test-acc-%s@mongodb.com", acctest.RandString(10))
		initialRole  = []string{"GROUP_OWNER"}
		updateRoles  = []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_CLUSTER_MANAGER"}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectInvitationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectInvitationConfig(orgID, projectName, name, initialRole),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectInvitationExists(t, resourceName, &invitation),
					testAccCheckMongoDBAtlasProjectInvitationUsernameAttribute(&invitation, name),
					testAccCheckMongoDBAtlasProjectInvitationRoleAttribute(&invitation, initialRole),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "invitation_id"),
					resource.TestCheckResourceAttrSet(resourceName, "roles.#"),
					resource.TestCheckResourceAttr(resourceName, "username", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectInvitationConfig(orgID, projectName, name, updateRoles),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectInvitationExists(t, resourceName, &invitation),
					testAccCheckMongoDBAtlasProjectInvitationUsernameAttribute(&invitation, name),
					testAccCheckMongoDBAtlasProjectInvitationRoleAttribute(&invitation, updateRoles),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "invitation_id"),
					resource.TestCheckResourceAttrSet(resourceName, "roles.#"),
					resource.TestCheckResourceAttr(resourceName, "username", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
				),
			},
		},
	})
}

func TestAccProjectRSProjectInvitation_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project_invitation.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = fmt.Sprintf("test-acc-%s@mongodb.com", acctest.RandString(10))
		initialRole  = []string{"GROUP_OWNER"}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectInvitationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectInvitationConfig(orgID, projectName, name, initialRole),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "roles.#"),
					resource.TestCheckResourceAttrSet(resourceName, "invitation_id"),
					resource.TestCheckResourceAttr(resourceName, "username", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasProjectInvitationStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectInvitationExists(t *testing.T, resourceName string, invitation *matlas.Invitation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		projectID := ids["project_id"]
		username := ids["username"]
		invitationID := ids["invitation_id"]

		if projectID == "" && username == "" && invitationID == "" {
			return fmt.Errorf("no ID is set")
		}

		t.Logf("projectID: %s", projectID)
		t.Logf("username: %s", username)
		t.Logf("invitationID: %s", invitationID)

		invitationResp, _, err := conn.Projects.Invitation(context.Background(), projectID, invitationID)
		if err == nil {
			*invitation = *invitationResp
			return nil
		}

		return fmt.Errorf("invitation(%s) does not exist", invitationID)
	}
}

func testAccCheckMongoDBAtlasProjectInvitationUsernameAttribute(invitation *matlas.Invitation, username string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if invitation.Username != username {
			return fmt.Errorf("bad name: %s", invitation.Username)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasProjectInvitationRoleAttribute(invitation *matlas.Invitation, roles []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(roles) > 0 {
			for _, role := range roles {
				for _, currentRole := range invitation.Roles {
					if currentRole == role {
						return nil
					}
				}
			}
		}

		return fmt.Errorf("bad role: %s", invitation.Roles)
	}
}

func testAccCheckMongoDBAtlasProjectInvitationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_invitations" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		projectID := ids["project_id"]
		invitationID := ids["invitation_id"]

		// Try to find the invitation
		_, _, err := conn.Projects.Invitation(context.Background(), projectID, invitationID)
		if err == nil {
			return fmt.Errorf("invitation (%s) still exists", invitationID)
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasProjectInvitationStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["username"]), nil
	}
}

func testAccMongoDBAtlasProjectInvitationConfig(orgID, projectName, username string, roles []string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_project_invitation" "test" {
			project_id = mongodbatlas_project.test.id
			username   = %[3]q
			roles  	 = ["%[4]s"]
		}`, orgID, projectName, username,
		strings.Join(roles, `", "`),
	)
}
