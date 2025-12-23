package serviceaccountsecret_test

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

const resourceName = "mongodbatlas_service_account_secret.test"

func TestAccServiceAccountSecret_basic(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name),
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

func TestAccServiceAccountSecret_rotateWithTaint(t *testing.T) {
	var (
		orgID         = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name          = acc.RandomName()
		firstSecretID string
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, name),
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
				Config: configBasic(orgID, name),
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

func TestAccServiceAccountSecret_dataSourceErrors(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configDataSourceInvalidClientID(orgID, name),
				ExpectError: regexp.MustCompile("The requested resource does not exist"),
			},
			{
				Config:      configDataSourceInvalidSecretID(orgID, name),
				ExpectError: regexp.MustCompile("The requested resource does not exist"),
			},
		},
	})
}

func configBasic(orgID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_service_account" "test" {
			org_id                     = %[1]q
			name                       = %[2]q
			description                = "Acceptance Test SA"
			roles                      = ["ORG_OWNER"]
			secret_expires_after_hours = 12
		}
		resource "mongodbatlas_service_account_secret" "test" {
			org_id                     = %[1]q
			client_id 				   = mongodbatlas_service_account.test.client_id
			secret_expires_after_hours = 12
		}

		data "mongodbatlas_service_account_secret" "test" {
			org_id    = mongodbatlas_service_account_secret.test.org_id
			client_id = mongodbatlas_service_account_secret.test.client_id
			secret_id = mongodbatlas_service_account_secret.test.secret_id
		}
	`, orgID, name)
}

func configDataSourceInvalidClientID(orgID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_service_account" "test" {
			org_id                     = %[1]q
			name                       = %[2]q
			description                = "Acceptance Test SA for data source errors"
			roles                      = ["ORG_OWNER"]
			secret_expires_after_hours = 12
		}

		data "mongodbatlas_service_account_secret" "test" {
			org_id    = %[1]q
			client_id = "000000000000000000000000"
			secret_id = mongodbatlas_service_account.test.secrets.0.secret_id
		}
	`, orgID, name)
}

func configDataSourceInvalidSecretID(orgID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_service_account" "test" {
			org_id                     = %[1]q
			name                       = %[2]q
			description                = "Acceptance Test SA for data source errors"
			roles                      = ["ORG_OWNER"]
			secret_expires_after_hours = 12
		}

		data "mongodbatlas_service_account_secret" "test" {
			org_id    = %[1]q
			client_id = mongodbatlas_service_account.test.client_id
			secret_id = "000000000000000000000000"
		}
	`, orgID, name)
}

func checkBasic() resource.TestCheckFunc {
	dataSourceName := "data.mongodbatlas_service_account_secret.test"
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
		orgID := rs.Primary.Attributes["org_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if id == "" || orgID == "" || clientID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		orgServiceAccount, _, err := acc.ConnV2().ServiceAccountsApi.GetOrgServiceAccount(context.Background(), orgID, clientID).Execute()
		if err != nil {
			return fmt.Errorf("failed to get org service account: %w", err)
		}
		if orgServiceAccount.Secrets != nil {
			for _, secret := range *orgServiceAccount.Secrets {
				if secret.Id == id {
					return nil
				}
			}
		}
		return fmt.Errorf("service account secret (%s/%s/%s) does not exist", id, orgID, clientID)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_service_account_secret" {
			continue
		}
		id := rs.Primary.Attributes["secret_id"]
		orgID := rs.Primary.Attributes["org_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if id == "" || orgID == "" || clientID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}

		orgServiceAccount, _, err := acc.ConnV2().ServiceAccountsApi.GetOrgServiceAccount(context.Background(), orgID, clientID).Execute()
		if err == nil && orgServiceAccount.Secrets != nil {
			for _, secret := range *orgServiceAccount.Secrets {
				if secret.Id == id {
					return fmt.Errorf("org service account secret (%s/%s/%s) still exists", id, orgID, clientID)
				}
			}
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
		id := rs.Primary.Attributes["secret_id"]
		orgID := rs.Primary.Attributes["org_id"]
		clientID := rs.Primary.Attributes["client_id"]
		if id == "" || orgID == "" || clientID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s", orgID, clientID, id), nil
	}
}
