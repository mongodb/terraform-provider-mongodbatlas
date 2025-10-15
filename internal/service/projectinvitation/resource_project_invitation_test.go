package projectinvitation_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccProjectRSProjectInvitation_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project_invitation.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		name         = acc.RandomEmail()
		initialRole  = []string{"GROUP_OWNER"}
		updateRoles  = []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_CLUSTER_MANAGER"}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, name, initialRole),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "invitation_id"),
					resource.TestCheckResourceAttr(resourceName, "username", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", initialRole[0]),
				),
			},
			{
				Config: configBasic(orgID, projectName, name, updateRoles),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "invitation_id"),
					resource.TestCheckResourceAttr(resourceName, "username", name),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", updateRoles[0]),
					resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", updateRoles[0]),
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

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]
		username := ids["username"]
		invitationID := ids["invitation_id"]
		if projectID == "" && username == "" && invitationID == "" {
			return fmt.Errorf("no ID is set")
		}
		_, _, err := acc.ConnV2().ProjectsApi.GetGroupInvite(context.Background(), projectID, invitationID).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("invitation (%s) does not exist", invitationID)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_invitations" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]
		invitationID := ids["invitation_id"]

		_, _, err := acc.ConnV2().ProjectsApi.GetGroupInvite(context.Background(), projectID, invitationID).Execute()
		if err == nil {
			return fmt.Errorf("invitation (%s) still exists", invitationID)
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["username"]), nil
	}
}

func configBasic(orgID, projectName, username string, roles []string) string {
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
