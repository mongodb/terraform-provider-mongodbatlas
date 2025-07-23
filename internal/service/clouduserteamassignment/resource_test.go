package clouduserteamassignment_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var resourceName = "mongodbatlas_cloud_user_team_assignment.test"
var dataSourceName1 = "data.mongodbatlas_cloud_user_team_assignment.test1"
var dataSourceName2 = "data.mongodbatlas_cloud_user_team_assignment.test2"

func TestAccCloudUserTeamAssignment_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func TestAccCloudUserTeamAssignmentDS_error(t *testing.T) {
	resource.ParallelTest(t, *errorTestCase(t))
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

func errorTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	teamName := acc.RandomName()

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configError(orgID, teamName),
				ExpectError: regexp.MustCompile("either username or user_id must be provided"),
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

func configError(orgID, teamName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_team" "test" {
			org_id     = %[1]q
			name       = %[2]q
		}


		data "mongodbatlas_cloud_user_team_assignment" "test" {
		org_id  = %[1]q
		team_id = mongodbatlas_team.test.team_id
		}
		`, orgID, teamName)
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
