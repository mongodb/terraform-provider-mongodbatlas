package serviceaccountprojectassignment_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312013/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_service_account_project_assignment.test"
const dataSourceName = "data.mongodbatlas_service_account_project_assignment.test"
const dataSourcePluralName = "data.mongodbatlas_service_account_project_assignments.test"

func TestAccServiceAccountProjectAssignment_singleAssignment(t *testing.T) {
	var (
		orgID         = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectIDs    = acc.MultipleProjectIDsExecution(t, 1)
		resourceName0 = resourceName + "_0"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectIDs, []string{"GROUP_OWNER", "GROUP_READ_ONLY"}),
				Check:  checkBasic(projectIDs),
			},
			{
				ResourceName:                         resourceName0,
				ImportStateIdFunc:                    importStateIDFunc(resourceName0),
				ImportStateVerifyIdentifierAttribute: "project_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func TestAccServiceAccountProjectAssignment_multipleAssignments(t *testing.T) {
	var (
		orgID         = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectIDs    = acc.MultipleProjectIDsExecution(t, 2)
		resourceName0 = resourceName + "_0"
		resourceName1 = resourceName + "_1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectIDs, []string{"GROUP_OWNER", "GROUP_READ_ONLY"}),
				Check:  checkBasic(projectIDs),
			},
			{
				ResourceName:                         resourceName0,
				ImportStateIdFunc:                    importStateIDFunc(resourceName0),
				ImportStateVerifyIdentifierAttribute: "project_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				ResourceName:                         resourceName1,
				ImportStateIdFunc:                    importStateIDFunc(resourceName1),
				ImportStateVerifyIdentifierAttribute: "project_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func configBasic(orgID string, projectIDs, roles []string) string {
	rolesStr := fmt.Sprintf("[%s]", `"`+strings.Join(roles, `", "`)+`"`)

	assignmentsStr := ""
	resourceNames := []string{}
	for i, projectID := range projectIDs {
		assignmentsStr += fmt.Sprintf(`
			resource "mongodbatlas_service_account_project_assignment" "test_%[1]d" {
				client_id  = mongodbatlas_service_account.test.client_id
				project_id = %[2]q
				roles      = %[3]s
			}

			data "mongodbatlas_service_account_project_assignment" "test_%[1]d" {
			  client_id = mongodbatlas_service_account.test.client_id
			  project_id = %[2]q
			  depends_on = [mongodbatlas_service_account_project_assignment.test_%[1]d]
			}
		`, i, projectID, rolesStr)
		resourceNames = append(resourceNames, fmt.Sprintf("%s_%d", resourceName, i))
	}

	resourceNamesStr := fmt.Sprintf("[%s]", `"`+strings.Join(resourceNames, `", "`)+`"`)

	return fmt.Sprintf(`
		resource "mongodbatlas_service_account" "test" {
			org_id                     = %[1]q
			name                       = "tf-acc-test-sa"
			description                = "Terraform acceptance test Service Account"
			roles                      = ["ORG_OWNER"]
			secret_expires_after_hours = 8
		}

		%[2]s

		data "mongodbatlas_service_account_project_assignments" "test" {
		  org_id = %[1]q
		  client_id = mongodbatlas_service_account.test.client_id
		  depends_on = %[3]s
		}
	`, orgID, assignmentsStr, resourceNamesStr)
}

func checkBasic(projectIDs []string) resource.TestCheckFunc {
	attrsSet := []string{"client_id", "project_id"}
	attrsMap := map[string]string{"roles.#": "2"}
	checks := []resource.TestCheckFunc{}
	for i := range projectIDs {
		resourceName := fmt.Sprintf("%s_%d", resourceName, i)
		dataSourceName := fmt.Sprintf("%s_%d", dataSourceName, i)
		checks = append(checks, acc.CheckRSAndDS(
			resourceName, admin.PtrString(dataSourceName), nil,
			attrsSet, attrsMap,
			checkExists(resourceName),
		))
	}
	checks = append(checks, resource.TestCheckResourceAttr(dataSourcePluralName, "results.#", strconv.Itoa(len(projectIDs))))
	return resource.ComposeAggregateTestCheckFunc(checks...)
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
