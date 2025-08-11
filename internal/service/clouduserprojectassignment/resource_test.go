package clouduserprojectassignment_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var resourceName = "mongodbatlas_cloud_user_project_assignment.test"

func TestAccCloudUserProjectAssignmentRS_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	username := acc.RandomEmail()
	projectName := acc.RandomName()
	roles := []string{"GROUP_OWNER", "GROUP_CLUSTER_MANAGER"}
	updatedRoles := []string{"GROUP_OWNER", "GROUP_SEARCH_INDEX_EDITOR", "GROUP_READ_ONLY"}

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, username, projectName, roles),
				Check:  checks(username, projectName, roles),
			},
			{
				Config: configBasic(orgID, username, projectName, updatedRoles),
				Check:  checks(username, projectName, updatedRoles),
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "user_id",
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
				ImportStateVerifyIdentifierAttribute: "user_id",
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
		}`,
		projectName, orgID, username, rolesStr)
}

func checks(username, projectName string, roles []string) resource.TestCheckFunc {
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "username", username),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "roles.#", fmt.Sprintf("%d", len(roles))),
	}

	for _, role := range roles {
		checkFuncs = append(checkFuncs, resource.TestCheckTypeSetElemAttr(resourceName, "roles.*", role))
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

		_, _, err := acc.ConnV2().MongoDBCloudUsersApi.GetProjectUser(context.Background(), projectID, userID).Execute()
		if err == nil {
			return fmt.Errorf("cloud user project assignment for user (%s) in project (%s) still exists", userID, projectID)
		}
	}
	return nil
}
