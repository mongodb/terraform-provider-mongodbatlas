package projectserviceaccount_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312012/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_project_service_account.test"
const dataSourceName = "data.mongodbatlas_project_service_account.test"
const dataSourcePluralName = "data.mongodbatlas_project_service_accounts.test"

func TestAccProjectServiceAccount_basic(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t)
		name1        = acc.RandomName()
		name2        = fmt.Sprintf("%s-updated", name1)
		description1 = "Acceptance Test SA"
		description2 = "Updated Description"
		roles1       = []string{"GROUP_OWNER"}
		roles2       = []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, name1, description1, roles1, 24),
				Check:  checkBasic(true, roles1),
			},
			{
				Config: configBasic(projectID, name2, description2, roles2, 24),
				Check:  checkBasic(false, roles2),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportStateVerifyIdentifierAttribute: "client_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"secret_expires_after_hours"},
			},
		},
	})
}

func TestAccProjectServiceAccount_createOnlyAttributes(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		name        = acc.RandomName()
		description = "description"
		roles       = []string{"GROUP_READ_ONLY"}
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, name, description, roles, 24),
				Check:  checkExists(resourceName),
			},
			{
				Config:      configBasic(projectID, name, description, roles, 48),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("secret_expires_after_hours cannot be updated"),
			},
			{
				Config:      configBasic("updated-project-id", name, description, roles, 24),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("project_id cannot be updated"),
			},
		},
	})
}

func configBasic(projectID, name, description string, roles []string, secretExpiresAfterHours int) string {
	rolesStr := `"` + strings.Join(roles, `", "`) + `"`
	rolesHCL := fmt.Sprintf("[%s]", rolesStr)
	return fmt.Sprintf(`
		resource "mongodbatlas_project_service_account" "test" {
			project_id                 = %[1]q
			name                       = %[2]q
			description                = %[3]q
			roles                      = %[4]s
			secret_expires_after_hours = %[5]d
		}

		data "mongodbatlas_project_service_account" "test" {
			project_id = %[1]q
			client_id = mongodbatlas_project_service_account.test.client_id
		}

		data "mongodbatlas_project_service_accounts" "test" {
			project_id = %[1]q
			depends_on = [mongodbatlas_project_service_account.test]
		}
	`, projectID, name, description, rolesHCL, secretExpiresAfterHours)
}

func checkBasic(isCreate bool, roles []string) resource.TestCheckFunc {
	commonAttrsSet := []string{"client_id", "created_at", "secrets.0.secret_id", "secrets.0.created_at", "secrets.0.expires_at"}
	commonAttrsMap := map[string]string{"secrets.#": "1", "roles.#": strconv.Itoa(len(roles))}

	checks := acc.CheckRSAndDS(resourceName, admin.PtrString(dataSourceName), nil, commonAttrsSet, commonAttrsMap, checkExists(resourceName))

	additionalChecks := []resource.TestCheckFunc{}
	if isCreate {
		additionalChecks = acc.AddAttrSetChecks(resourceName, additionalChecks, "secrets.0.secret")
	} else {
		additionalChecks = acc.AddAttrSetChecks(resourceName, additionalChecks, "secrets.0.masked_secret_value")
	}

	additionalChecks = acc.AddAttrSetChecks(dataSourceName, additionalChecks, "secrets.0.masked_secret_value")
	additionalChecks = acc.AddAttrSetChecksPrefix(dataSourcePluralName, additionalChecks, []string{"secrets.0.masked_secret_value"}, "results.0")

	return resource.ComposeAggregateTestCheckFunc(checks, resource.ComposeAggregateTestCheckFunc(additionalChecks...))
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
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	for name, rs := range s.RootModule().Resources {
		if name != resourceName {
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

		// Delete the service account (project_service_account DELETE only removes the project assignment)
		_, err = acc.ConnV2().ServiceAccountsApi.DeleteOrgServiceAccount(context.Background(), clientID, orgID).Execute()
		if err != nil {
			return fmt.Errorf("failed to cleanup service account (%s/%s): %w", orgID, clientID, err)
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
