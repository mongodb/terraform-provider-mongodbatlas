package clouduserteamassignment_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

var resourceName = "mongodbatlas_cloud_user_team_assignment.test"
var dataSourceName1 = "data.mongodbatlas_cloud_user_team_assignment.test1"
var dataSourceName2 = "data.mongodbatlas_cloud_user_team_assignment.test2"

func TestAccCloudUserTeamAssignment_basic(t *testing.T) {
	resource.Test(t, *basicTestCase(t))
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
				Config: configBasic(orgID, userID, teamName),
				Check:  checks(orgID, userID),
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
		},
	}
}

func configBasic(orgID, userID, teamName string) string {
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
		data "mongodbatlas_cloud_user_team_assignment" "test1" {
			org_id   = %[1]q
			team_id  = mongodbatlas_team.test.team_id
			user_id  = mongodbatlas_cloud_user_team_assignment.test.user_id
		}

		data "mongodbatlas_cloud_user_team_assignment" "test2" {
			org_id   = %[1]q
			team_id  = mongodbatlas_team.test.team_id
			username = mongodbatlas_cloud_user_team_assignment.test.username
		}
		`,
		orgID, userID, teamName)
}

func checks(orgID, userID string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
		resource.TestCheckResourceAttr(resourceName, "user_id", userID),
		resource.TestCheckResourceAttrSet(resourceName, "username"),
		resource.TestCheckResourceAttrWith(resourceName, "username", acc.IsUsername()),
		resource.TestCheckResourceAttrWith(resourceName, "created_at", acc.IsTimestamp()),
		resource.TestCheckResourceAttrWith(resourceName, "team_ids.#", acc.IntGreatThan(0)),

		resource.TestCheckResourceAttr(dataSourceName1, "user_id", userID),
		resource.TestCheckResourceAttrWith(dataSourceName1, "username", acc.IsUsername()),
		resource.TestCheckResourceAttr(dataSourceName1, "org_id", orgID),

		resource.TestCheckResourceAttr(dataSourceName2, "user_id", userID),
		resource.TestCheckResourceAttrWith(dataSourceName2, "username", acc.IsUsername()),
		resource.TestCheckResourceAttr(dataSourceName2, "org_id", orgID),
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
		conn := acc.ConnV2()
		ctx := context.Background()

		if userID != "" {
			userIDParams := &admin.ListTeamUsersApiParams{
				UserId: &userID,
				OrgId:  orgID,
				TeamId: teamID,
			}
			userListUserID, _, err := conn.MongoDBCloudUsersApi.ListTeamUsersWithParams(ctx, userIDParams).Execute()
			if userListUserID.HasResults() {
				return fmt.Errorf("cloud user team assignment for user (%s) in team (%s) still exists %s", userID, teamID, err)
			}
		}
	}
	return nil
}
