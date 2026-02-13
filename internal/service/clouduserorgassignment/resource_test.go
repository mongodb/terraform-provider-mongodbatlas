package clouduserorgassignment_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var resourceName = "mongodbatlas_cloud_user_org_assignment.test"

func TestAccCloudUserOrgAssignmentRS_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func TestAccCloudUserOrgAssignmentDS_basic(t *testing.T) {
	resource.ParallelTest(t, *dataSourceTestCase(t))
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	username := acc.RandomEmail()
	roles := []string{"ORG_MEMBER"}
	rolesUpdated := []string{"ORG_MEMBER", "ORG_GROUP_CREATOR"}

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudUserOrgAssignmentConfig(orgID, username, roles),
				Check:  cloudUserOrgAssignmentChecks(resourceName, orgID, username, "PENDING", roles),
			},
			{
				Config: testAccCloudUserOrgAssignmentConfig(orgID, username, rolesUpdated),
				Check:  cloudUserOrgAssignmentChecks(resourceName, orgID, username, "PENDING", rolesUpdated),
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
	}
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

func dataSourceTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	username := acc.RandomEmail()
	roles := []string{"ORG_MEMBER"}

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudUserOrgAssignmentWithDataSourceConfig(orgID, username, roles),
				Check: resource.ComposeTestCheckFunc(
					cloudUserOrgAssignmentChecks("data.mongodbatlas_cloud_user_org_assignment.by_username", orgID, username, "PENDING", roles),
					cloudUserOrgAssignmentChecks("data.mongodbatlas_cloud_user_org_assignment.by_user_id", orgID, username, "PENDING", roles),
				),
			},
		},
	}
}

func testAccCloudUserOrgAssignmentWithDataSourceConfig(orgID, username string, roles []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`

	return fmt.Sprintf(`
resource "mongodbatlas_cloud_user_org_assignment" "test" {
  org_id   = %[1]q
  username = %[2]q
  roles = {
    org_roles = [%[3]s]
  }
}

# Query by username - user_id should be computed
data "mongodbatlas_cloud_user_org_assignment" "by_username" {
  org_id   = %[1]q
  username = mongodbatlas_cloud_user_org_assignment.test.username
}

# Query by user_id - username should be computed
data "mongodbatlas_cloud_user_org_assignment" "by_user_id" {
  org_id   = %[1]q
  user_id  = mongodbatlas_cloud_user_org_assignment.test.user_id
}

# The following resources verify that user_id is properly marked as Computed
# When querying by username, user_id must be marked as "Computed: true" in the schema
# so that Terraform knows its value will be available for use in other resources.
# Without "Computed: true", Terraform will fail at plan time because it cannot validate
# that user_id will be available as an input to the team assignment resource below, where user_id is required.
#   - Without Computed: This test FAILS at plan time
#   - With Computed: This test SUCCEEDS
resource "mongodbatlas_team" "test" {
  org_id = %[1]q
  name   = "test-team"
}

resource "mongodbatlas_cloud_user_team_assignment" "test_team_assignment" {
  org_id  = %[1]q
  team_id = mongodbatlas_team.test.team_id
  user_id = data.mongodbatlas_cloud_user_org_assignment.by_username.user_id
}
`, orgID, username, rolesStr)
}

func cloudUserOrgAssignmentChecks(resourceName, orgID, username, orgMembershipStatus string, roles []string) resource.TestCheckFunc {
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
			_, resp, err := conn.MongoDBCloudUsersApi.GetOrgUser(context.Background(), orgID, userID).Execute()
			if err == nil && resp != nil && resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("cloud user org assignment (%s) still exists", userID)
			}
		} else if username != "" {
			params := &admin.ListOrgUsersApiParams{
				OrgId:    orgID,
				Username: &username,
			}

			users, _, err := conn.MongoDBCloudUsersApi.ListOrgUsersWithParams(context.Background(), params).Execute()
			if err == nil && users != nil && len(*users.Results) > 0 {
				return fmt.Errorf("cloud user org assignment (%s) still exists", username)
			}
		}
	}
	return nil
}
