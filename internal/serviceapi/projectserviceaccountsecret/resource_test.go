package projectserviceaccountsecret_test

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

const resourceName = "mongodbatlas_project_service_account_secret.test"

func TestAccProjectServiceAccountSecret_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		name      = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, name),
				Check:  checkBasic(),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportStateVerifyIdentifierAttribute: "secret_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"secret", "masked_secret_value"}, // create returns secret only, import returns masked secret only
			},
		},
	})
}

func TestAccProjectServiceAccountSecret_rotateWithTaint(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		name          = acc.RandomName()
		firstSecretID string
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					checkBasic(),
					func(s *terraform.State) error {
						return getSecretID(s, &firstSecretID)
					},
				),
			},
			{
				// The `taint` command is deprecated in favor of the `-replace` flag: https://developer.hashicorp.com/terraform/cli/commands/taint.
				// The testing plugin does not facilitate testing with replace, but it does enable tainting so using taint here.
				Taint:  []string{resourceName},
				Config: configBasic(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					checkBasic(),
					func(s *terraform.State) error {
						var secondSecretID string
						if err := getSecretID(s, &secondSecretID); err != nil {
							return err
						}
						if secondSecretID == firstSecretID {
							return fmt.Errorf("expected secret %s to be replaced but it still exists", firstSecretID)
						}
						return nil
					},
				),
			},
		},
	})
}

func getSecretID(s *terraform.State, secretID *string) error {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return fmt.Errorf("not found: %s", resourceName)
	}
	id := rs.Primary.Attributes["secret_id"]
	if id == "" {
		return fmt.Errorf("secret_id is empty")
	}
	*secretID = id
	return nil
}

func TestAccProjectServiceAccountSecret_dataSourceErrors(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		name      = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configDataSourceInvalidClientID(projectID, name),
				ExpectError: regexp.MustCompile("The requested resource does not exist"),
			},
			{
				Config:      configDataSourceInvalidSecretID(projectID, name),
				ExpectError: regexp.MustCompile("The requested resource does not exist"),
			},
		},
	})
}

func configBasic(projectID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_service_account" "test" {
			project_id                 = %[1]q
			name                       = %[2]q
			description                = "Acceptance Test Project SA"
			roles                      = ["GROUP_OWNER"]
			secret_expires_after_hours = 12
		}

		resource "mongodbatlas_project_service_account_secret" "test" {
			project_id                 = %[1]q
			client_id 				   = mongodbatlas_project_service_account.test.client_id
			secret_expires_after_hours = 12
		}

		data "mongodbatlas_project_service_account_secret" "test" {
			project_id = mongodbatlas_project_service_account_secret.test.project_id
			client_id  = mongodbatlas_project_service_account_secret.test.client_id
			secret_id  = mongodbatlas_project_service_account_secret.test.secret_id
		}
	`, projectID, name)
}

func configDataSourceInvalidClientID(projectID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_service_account" "test" {
			project_id                 = %[1]q
			name                       = %[2]q
			description                = "Acceptance Test SA for data source errors"
			roles                      = ["GROUP_OWNER"]
			secret_expires_after_hours = 12
		}

		data "mongodbatlas_project_service_account_secret" "test" {
			project_id = %[1]q
			client_id  = "000000000000000000000000"
			secret_id  = mongodbatlas_project_service_account.test.secrets.0.secret_id
		}
	`, projectID, name)
}

func configDataSourceInvalidSecretID(projectID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_service_account" "test" {
			project_id                 = %[1]q
			name                       = %[2]q
			description                = "Acceptance Test SA for data source errors"
			roles                      = ["GROUP_OWNER"]
			secret_expires_after_hours = 12
		}

		data "mongodbatlas_project_service_account_secret" "test" {
			project_id = %[1]q
			client_id  = mongodbatlas_project_service_account.test.client_id
			secret_id  = "000000000000000000000000"
		}
	`, projectID, name)
}

func checkBasic() resource.TestCheckFunc {
	dataSourceName := "data.mongodbatlas_project_service_account_secret.test"
	attrsSet := []string{"secret_id", "created_at", "expires_at"}
	extraChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceName, "secret"), // resource has secret, data source does not
		resource.TestCheckNoResourceAttr(resourceName, "masked_secret_value"),
		checkExists(resourceName),
	}
	return acc.CheckRSAndDS(resourceName, &dataSourceName, nil, attrsSet, nil, extraChecks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		id := rs.Primary.Attributes["secret_id"]
		projectID := rs.Primary.Attributes["project_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if id == "" || projectID == "" || clientID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		projectServiceAccount, _, err := acc.ConnV2().ServiceAccountsApi.GetGroupServiceAccount(context.Background(), projectID, clientID).Execute()
		if err != nil {
			return fmt.Errorf("failed to get project service account: %w", err)
		}
		if projectServiceAccount.Secrets != nil {
			for _, secret := range *projectServiceAccount.Secrets {
				if secret.Id == id {
					return nil
				}
			}
		}
		return fmt.Errorf("project service account secret (%s/%s/%s) does not exist", id, projectID, clientID)
	}
}

func checkDestroy(s *terraform.State) error {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	for name, rs := range s.RootModule().Resources {
		if name != resourceName {
			continue
		}
		id := rs.Primary.Attributes["secret_id"]
		projectID := rs.Primary.Attributes["project_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if id == "" || projectID == "" || clientID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}

		projectServiceAccount, _, err := acc.ConnV2().ServiceAccountsApi.GetGroupServiceAccount(context.Background(), projectID, clientID).Execute()
		if err == nil && projectServiceAccount.Secrets != nil {
			for _, secret := range *projectServiceAccount.Secrets {
				if secret.Id == id {
					return fmt.Errorf("project service account secret (%s/%s/%s) still exists", id, projectID, clientID)
				}
			}
		}

		// Delete the service account (project_service_account DELETE only removes the project assignment)
		_, _ = acc.ConnV2().ServiceAccountsApi.DeleteOrgServiceAccount(context.Background(), clientID, orgID).Execute()
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		id := rs.Primary.Attributes["secret_id"]
		projectID := rs.Primary.Attributes["project_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if id == "" || projectID == "" || clientID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s", projectID, clientID, id), nil
	}
}
