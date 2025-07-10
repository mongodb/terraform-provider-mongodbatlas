package clouduserorgassignment_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var resourceName = "mongodbatlas_cloud_user_org_assignment.test"

func TestAccCloudUserOrgAssignmentRS_basic(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	username := "test-cloud-user-org-assignment@example.com"
	roles := []string{"ORG_MEMBER"}
	rolesUpdated := []string{"ORG_MEMBER", "ORG_GROUP_CREATOR"}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudUserOrgAssignmentConfig(orgID, username, roles),
				Check:  cloudUserOrgAssignmentChecks(orgID, username, "PENDING", roles),
			},
			{
				Config: testAccCloudUserOrgAssignmentConfig(orgID, username, rolesUpdated),
				Check:  cloudUserOrgAssignmentChecks(orgID, username, "PENDING", rolesUpdated),
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "user_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attrs := s.RootModule().Resources[resourceName].Primary.Attributes
					orgID := attrs["org_id"]
					userID := attrs["user_id"]
					return orgID + "/" + userID, nil
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "user_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attrs := s.RootModule().Resources[resourceName].Primary.Attributes
					orgID := attrs["org_id"]
					username := attrs["username"]
					return orgID + "/" + username, nil
				},
			},
		},
	})
}

func testAccCloudUserOrgAssignmentConfig(orgID, username string, roles []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`

	return fmt.Sprintf(`
resource "mongodbatlas_cloud_user_org_assignment" "test" {
  org_id   = "%s"
  username = "%s"
  roles = {
    org_roles = [%s]
  }
}
`, orgID, username, rolesStr)
}

func cloudUserOrgAssignmentChecks(orgID, username, orgMembershipStatus string, roles []string) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{}
	attributes := map[string]string{
		"org_id":                orgID,
		"username":              username,
		"org_membership_status": orgMembershipStatus,
		"roles.org_roles.#":     strconv.Itoa(len(roles)),
	}
	checks = acc.AddAttrChecks(resourceName, checks, attributes)

	if orgMembershipStatus == "PENDING" {
		checks = acc.AddAttrSetChecks(resourceName, checks, "user_id", "invitation_created_at", "invitation_expires_at", "inviter_username")
	} else {
		checks = acc.AddAttrSetChecks(resourceName, checks, "user_id", "country", "created_at", "first_name", "last_auth", "last_name", "mobile_number")
	}

	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_user_org_assignment" {
			continue
		}
		orgID := rs.Primary.Attributes["org_id"]
		userID := rs.Primary.Attributes["user_id"]
		username := rs.Primary.Attributes["username"]
		conn := acc.ConnV2()

		if userID != "" {
			_, resp, err := conn.MongoDBCloudUsersApi.GetOrganizationUser(context.Background(), orgID, userID).Execute()
			if err == nil && resp != nil && resp.StatusCode != 404 {
				return fmt.Errorf("cloud user org assignment (%s) still exists", userID)
			}
		} else if username != "" {
			params := &admin.ListOrganizationUsersApiParams{
				OrgId:    orgID,
				Username: &username,
			}

			users, _, err := conn.MongoDBCloudUsersApi.ListOrganizationUsersWithParams(context.Background(), params).Execute()
			if err == nil && users != nil && len(*users.Results) > 0 {
				return fmt.Errorf("cloud user org assignment (%s) still exists", username)
			}
		}
	}
	return nil
}
