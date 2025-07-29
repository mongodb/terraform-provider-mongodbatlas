package clouduserprojectassignment_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var resourceName = "mongodbatlas_cloud_user_project_assignment.test"
var dataSourceName1 = "data.mongodbatlas_cloud_user_project_assignment.testUsername"
var dataSourceName2 = "data.mongodbatlas_cloud_user_project_assignment.testUserID"

func TestAccCloudUserProjectAssignment_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func TestAccCloudUserProjectAssignmentDS_error(t *testing.T) {
	resource.ParallelTest(t, *errorTestCase(t))
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	username := acc.RandomEmail()
	projectName := acc.RandomName()
	roles := []string{"GROUP_OWNER", "GROUP_CLUSTER_MANAGER"}
	// updatedRoles := []string{"GROUP_OWNER", "GROUP_SEARCH_INDEX_EDITOR", "GROUP_READ_ONLY"}

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, username, projectName, roles),
				Check:  checks(username, roles),
			},
			/*{
				Config: configBasic(orgID, username, projectName, updatedRoles),
				Check:  checks(username, updatedRoles),
			},*/
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "project_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attrs := s.RootModule().Resources[resourceName].Primary.Attributes
					projectID := attrs["project_id"]
					userID := attrs["user_id"]
					return projectID + "/" + userID, nil
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "project_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attrs := s.RootModule().Resources[resourceName].Primary.Attributes
					projectID := attrs["project_id"]
					username := attrs["username"]
					return projectID + "/" + username, nil
				},
			},
		},
	}
}

func errorTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acc.RandomName()

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configError(orgID, projectName),
				ExpectError: regexp.MustCompile("either username or user_id must be provided"),
			},
		},
	}
}

func configBasic(orgID, username, projectName string, roles []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[1]q
			org_id = %[2]q
		}

		resource "mongodbatlas_cloud_user_project_assignment" "test" {
			username = %[3]q
			project_id = mongodbatlas_project.test.id
			roles = [%[4]s]
		}
			
		data "mongodbatlas_cloud_user_project_assignment" "testUsername" {
			project_id = mongodbatlas_project.test.id
			username = mongodbatlas_cloud_user_project_assignment.test.username
		}
			
		data "mongodbatlas_cloud_user_project_assignment" "testUserID" {
			project_id = mongodbatlas_project.test.id
			user_id = mongodbatlas_cloud_user_project_assignment.test.user_id
		}`,
		projectName, orgID, username, rolesStr)
}

func configError(orgID, projectName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[1]q
			org_id = %[2]q
		}
		data "mongodbatlas_cloud_user_project_assignment" "test" {
			project_id = mongodbatlas_project.test.id
		}
	`, projectName, orgID)
}

func checks(username string, roles []string) resource.TestCheckFunc {
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "username", username),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "roles.#", fmt.Sprintf("%d", len(roles))),
	}

	for _, role := range roles {
		checkFuncs = append(checkFuncs, resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", role))
	}
	dataCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dataSourceName1, "username", username),
		resource.TestCheckResourceAttr(dataSourceName2, "username", username),

		// resource.TestCheckResourceAttrPair(dataSourceName1, "user_id", dataSourceName2, "user_id"),
		resource.TestCheckResourceAttrPair(dataSourceName1, "project_id", dataSourceName2, "project_id"),
		resource.TestCheckResourceAttrPair(dataSourceName1, "roles.#", dataSourceName2, "roles.#"),
	}

	checkFuncs = append(checkFuncs, dataCheckFuncs...)
	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func checkDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "mongodbatlas_cloud_user_project_assignment" {
			continue
		}

		userID := r.Primary.Attributes["user_id"]
		username := r.Primary.Attributes["username"]
		projectID := r.Primary.Attributes["project_id"]
		conn := acc.ConnV2()

		userListResp, _, err := conn.MongoDBCloudUsersApi.ListProjectUsers(context.Background(), projectID).Execute()
		if err != nil {
			continue
		}

		if userListResp != nil {
			results := userListResp.GetResults()
			for i := range results {
				if userID != "" && results[i].GetId() == userID {
					return fmt.Errorf("cloud user project assignment for user (%s) in project (%s) still exists", r.Primary.Attributes["username"], projectID)
				}
				if username != "" && results[i].GetUsername() == username {
					return fmt.Errorf("cloud user project assignment for user (%s) in project (%s) still exists", r.Primary.Attributes["username"], projectID)
				}
			}
		}
	}
	return nil
}
