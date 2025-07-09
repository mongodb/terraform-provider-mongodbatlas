package clouduserorgassignment_test

import (
	"context"
	"fmt"
	"os"
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
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudUserOrgAssignmentConfig(orgID, username),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "roles.org_roles.0", "ORG_MEMBER"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attrs := s.RootModule().Resources[resourceName].Primary.Attributes
					orgID := attrs["org_id"]
					userID := attrs["user_id"]
					return orgID + "/" + userID, nil
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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

func TestAccCloudUserOrgAssignmentRS_importByUsername(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	username := os.Getenv("MONGODB_ATLAS_USERNAME")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudUserOrgAssignmentConfig(orgID, username),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "roles.org_roles.0", "ORG_MEMBER"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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

func testAccCloudUserOrgAssignmentConfig(orgID, username string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cloud_user_org_assignment" "test" {
  org_id   = "%s"
  username = "%s"
  roles = {
    org_roles = ["ORG_MEMBER"]
  }
}
`, orgID, username)
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
