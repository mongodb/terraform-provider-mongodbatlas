package clouduserteamassignment_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var resourceName = "mongodbatlas_cloud_user_team_assignment.test"

func TestAccCloudUserTeamAssignmentRS_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	userID := os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	teamName := acc.RandomName()

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: cloudUserTeamAssignmentConfig(orgID, userID, teamName),
				Check:  cloudUserTeamAssignmentAttributeChecks(orgID, userID),
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "user_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attrs := s.RootModule().Resources[resourceName].Primary.Attributes
					orgID := attrs["org_id"]
					teamID := attrs["team_id"]
					userID := attrs["user_id"]
					return orgID + "/" + teamID + "/" + userID, nil
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
					teamID := attrs["team_id"]
					username := attrs["username"]
					return orgID + "/" + teamID + "/" + username, nil
				},
			},
		},
	}
}

func cloudUserTeamAssignmentConfig(orgID, userID, teamName string) string {
	return fmt.Sprintf(` 
		resource "mongodbatlas_team" "test" {
			org_id     = %[1]q
			name       = %[3]q
		}
		resource "mongodbatlas_cloud_user_team_assignment" "test" {  
			org_id  = %[1]q  
			team_id = mongodbatlas_team.test.team_id
			user_id = %[2]q    
		}  
		`,
		orgID, userID, teamName)
}

func cloudUserTeamAssignmentAttributeChecks(orgID, userID string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
		resource.TestCheckResourceAttr(resourceName, "user_id", userID),

		resource.TestCheckResourceAttrSet(resourceName, "username"),
		resource.TestCheckResourceAttrWith(resourceName, "username", acc.IsUsername()),
		resource.TestCheckResourceAttrWith(resourceName, "created_at", acc.IsTimestamp()),

		resource.TestCheckResourceAttrWith(resourceName, "team_ids.#", acc.IntGreatThan(0)),
	)
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_user_team_assignment" {
			continue
		}
		orgID := rs.Primary.Attributes["org_id"]
		teamID := rs.Primary.Attributes["team_id"]
		userID := rs.Primary.Attributes["user_id"]
		username := rs.Primary.Attributes["username"]
		conn := acc.ConnV2()

		userListResp, _, err := conn.MongoDBCloudUsersApi.ListTeamUsers(context.Background(), orgID, teamID).Execute()
		if err != nil {
			continue
		}

		if userListResp != nil && userListResp.Results != nil {
			results := *userListResp.Results
			for i := range results {
				if userID != "" && results[i].GetId() == userID {
					return fmt.Errorf("cloud user team assignment for user (%s) in team (%s) still exists", userID, teamID)
				}
				if username != "" && results[i].GetUsername() == username {
					return fmt.Errorf("cloud user team assignment for user (%s) in team (%s) still exists", username, teamID)
				}
			}
		}
	}
	return nil
}
