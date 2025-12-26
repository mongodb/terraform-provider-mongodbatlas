package serviceaccountprojectassignment_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312011/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_service_account_project_assignment.test"
const dataSourceName = "data.mongodbatlas_service_account_project_assignment.test"
const dataSourcePluralName = "data.mongodbatlas_service_account_project_assignments.test"

func TestAccServiceAccountProjectAssignment_basic(t *testing.T) {
	var (
		orgID     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectID, []string{"GROUP_OWNER", "GROUP_READ_ONLY"}),
				Check:  checkBasic(projectID),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportStateVerifyIdentifierAttribute: "client_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccServiceAccountProjectAssignment_rolesOrdering(t *testing.T) {
	var (
		orgID     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID = acc.ProjectIDExecution(t)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectID, []string{"GROUP_CLUSTER_MANAGER", "GROUP_READ_ONLY"}),
			},
			{
				Config: configBasic(orgID, projectID, []string{"GROUP_READ_ONLY", "GROUP_CLUSTER_MANAGER"}),
			},
		},
	})
}

func configBasic(orgID, projectID string, roles []string) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	rolesHCL := fmt.Sprintf("[%s]", rolesStr)
	return fmt.Sprintf(`
		resource "mongodbatlas_service_account" "test" {
			org_id                     = %[1]q
			name                       = "tf-acc-test-sa"
			description                = "Terraform acceptance test Service Account"
			roles                      = ["ORG_OWNER"]
			secret_expires_after_hours = 8
		}

		resource "mongodbatlas_service_account_project_assignment" "test" {
			client_id  = mongodbatlas_service_account.test.client_id
			project_id = %[2]q
			roles      = %[3]s
		}

		data "mongodbatlas_service_account_project_assignment" "test" {
          client_id = mongodbatlas_service_account.test.client_id
		  project_id = %[2]q
		  depends_on = [mongodbatlas_service_account_project_assignment.test]
		}

		data "mongodbatlas_service_account_project_assignments" "test" {
		  org_id = %[1]q
		  client_id = mongodbatlas_service_account.test.client_id
		  depends_on = [mongodbatlas_service_account_project_assignment.test]
		}
	`, orgID, projectID, rolesHCL)
}

func checkBasic(projectID string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		acc.CheckRSAndDS(
			resourceName, admin.PtrString(dataSourceName), nil,
			[]string{"client_id", "project_id"},
			map[string]string{"roles.#": "2"},
			checkExists(resourceName),
		),
		resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.project_id", projectID),
	)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		projectID := rs.Primary.Attributes["project_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if projectID == "" || clientID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}

		_, _, err := acc.ConnV2().ServiceAccountsApi.GetGroupServiceAccount(context.Background(), projectID, clientID).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("project service account (%s/%s) does not exist", projectID, clientID)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_service_account_project_assignment" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if projectID == "" || clientID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}

		_, _, err := acc.ConnV2().ServiceAccountsApi.GetGroupServiceAccount(context.Background(), projectID, clientID).Execute()
		if err == nil {
			return fmt.Errorf("project service account (%s/%s) still exists", projectID, clientID)
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
		projectID := rs.Primary.Attributes["project_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if projectID == "" || clientID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", projectID, clientID), nil
	}
}
