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

var resourceNamePending = "mongodbatlas_cloud_user_project_assignment.test_pending"
var resourceNameActive = "mongodbatlas_cloud_user_project_assignment.test_active"
var DSNameUsername = "data.mongodbatlas_cloud_user_project_assignment.test_username"
var DSNameUserID = "data.mongodbatlas_cloud_user_project_assignment.test_user_id"

func TestAccCloudUserProjectAssignment_basic(t *testing.T) {
	resource.Test(t, *basicTestCase(t))
}

func TestAccCloudUserProjectAssignmentDS_error(t *testing.T) {
	resource.ParallelTest(t, *errorTestCase(t))
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	// Use MONGODB_ATLAS_USERNAME_2 to avoid USER_ALREADY_IN_GROUP.
	// The default MONGODB_ATLAS_USERNAME (Org Owner) is auto-assigned if no ProjectOwner is set.
	activeUsername := os.Getenv("MONGODB_ATLAS_USERNAME_2")
	pendingUsername := acc.RandomEmail()
	projectID := acc.ProjectIDExecution(t)
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	roles := []string{"GROUP_OWNER", "GROUP_CLUSTER_MANAGER"}
	updatedRoles := []string{"GROUP_OWNER", "GROUP_SEARCH_INDEX_EDITOR", "GROUP_READ_ONLY"}

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckAtlasUsernames(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, pendingUsername, activeUsername, projectID, roles),
				Check:  checks(pendingUsername, activeUsername, projectID, roles),
			},
			{
				Config: configBasic(orgID, pendingUsername, activeUsername, projectID, updatedRoles),
				Check:  checks(pendingUsername, activeUsername, projectID, updatedRoles),
			},
			{
				ResourceName:                         resourceNamePending,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "user_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attrs := s.RootModule().Resources[resourceNamePending].Primary.Attributes
					projectID := attrs["project_id"]
					userID := attrs["user_id"]
					return projectID + "/" + userID, nil
				},
			},
			{
				ResourceName:                         resourceNamePending,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "user_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attrs := s.RootModule().Resources[resourceNamePending].Primary.Attributes
					projectID := attrs["project_id"]
					username := attrs["username"]
					return projectID + "/" + username, nil
				},
			},
			{
				ResourceName:                         resourceNameActive,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "user_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attrs := s.RootModule().Resources[resourceNameActive].Primary.Attributes
					projectID := attrs["project_id"]
					userID := attrs["user_id"]
					return projectID + "/" + userID, nil
				},
			},
			{
				ResourceName:                         resourceNameActive,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "user_id",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					attrs := s.RootModule().Resources[resourceNameActive].Primary.Attributes
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
	projectID := acc.ProjectIDExecution(t)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configError(orgID, projectID),
				ExpectError: regexp.MustCompile("either username or user_id must be provided"),
			},
		},
	}
}

func configError(orgID, projectID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_cloud_user_project_assignment" "test" {
			project_id = %[1]q
		}
	`, projectID, orgID)
}

func configBasic(orgID, pendingUsername, activeUsername, projectID string, roles []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[1]q
			org_id = %[2]q
		}

		resource "mongodbatlas_cloud_user_project_assignment" "test_pending" {
			username = %[3]q
			project_id = %[1]q
			roles = [%[5]s]
		}

		resource "mongodbatlas_cloud_user_project_assignment" "test_active" {
			username = %[4]q
			project_id = %[1]q
			roles = [%[5]s]
		}

		data "mongodbatlas_cloud_user_project_assignment" "test_username" {
			project_id = %[1]q
			username = mongodbatlas_cloud_user_project_assignment.test_pending.username
		}

		data "mongodbatlas_cloud_user_project_assignment" "test_user_id" {
			project_id = %[1]q
			user_id = mongodbatlas_cloud_user_project_assignment.test_pending.user_id
		}`,
		projectID, orgID, pendingUsername, activeUsername, rolesStr)
}

func checks(pendingUsername, activeUsername, projectID string, roles []string) resource.TestCheckFunc {
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceNamePending, "username", pendingUsername),
		resource.TestCheckResourceAttr(resourceNamePending, "project_id", projectID),
		resource.TestCheckResourceAttr(resourceNamePending, "roles.#", fmt.Sprintf("%d", len(roles))),
		resource.TestCheckResourceAttr(resourceNameActive, "username", activeUsername),
		resource.TestCheckResourceAttr(resourceNameActive, "project_id", projectID),
		resource.TestCheckResourceAttr(resourceNameActive, "roles.#", fmt.Sprintf("%d", len(roles))),
		resource.TestCheckResourceAttr(DSNameUserID, "username", pendingUsername),
		resource.TestCheckResourceAttrPair(DSNameUserID, "username", resourceNamePending, "username"),
		resource.TestCheckResourceAttrPair(DSNameUsername, "user_id", DSNameUserID, "user_id"),
		resource.TestCheckResourceAttrPair(DSNameUsername, "project_id", DSNameUserID, "project_id"),
		resource.TestCheckResourceAttrPair(DSNameUsername, "roles.#", DSNameUserID, "roles.#"),
	}

	for _, role := range roles {
		checkFuncs = append(checkFuncs,
			resource.TestCheckTypeSetElemAttr(resourceNamePending, "roles.*", role),
			resource.TestCheckTypeSetElemAttr(resourceNameActive, "roles.*", role),
		)
	}

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func checkDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "mongodbatlas_cloud_user_project_assignment" {
			continue
		}

		userID := r.Primary.Attributes["user_id"]
		projectID := r.Primary.Attributes["project_id"]

		_, _, err := acc.ConnV2().MongoDBCloudUsersApi.GetGroupUser(context.Background(), projectID, userID).Execute()
		if err == nil {
			return fmt.Errorf("cloud user project assignment for user (%s) in project (%s) still exists", userID, projectID)
		}
	}
	return nil
}
