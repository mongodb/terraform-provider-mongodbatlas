package mcpconfigsecret_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_mcp_config_secret.test"
const dataSourceName = "data.mongodbatlas_mcp_config_secret.test"
const dataSourcePluralName = "data.mongodbatlas_mcp_config_secrets.test"

func TestAccMcpConfigSecret_basic(t *testing.T) {
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
				Config: configBasic(orgID, name, 720, false),
				Check:  checkBasic(true),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportStateVerifyIdentifierAttribute: "secret_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"secret_expires_after_hours", "secret", "masked_secret_value"},
			},
		},
	})
}

func TestAccMcpConfigSecret_rotateWithTaint(t *testing.T) {
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
				Config: configBasic(orgID, name, 720, false),
				Check: resource.ComposeTestCheckFunc(
					checkBasic(true),
					func(s *terraform.State) error {
						return getSecretID(s, resourceName, &firstSecretID)
					},
				),
			},
			{
				// The `taint` command is deprecated in favor of the `-replace` flag: https://developer.hashicorp.com/terraform/cli/commands/taint.
				// The testing plugin does not facilitate testing with replace, but it does enable tainting so using taint here.
				Taint:  []string{resourceName},
				Config: configBasic(orgID, name, 720, false),
				Check: resource.ComposeTestCheckFunc(
					checkBasic(true),
					func(s *terraform.State) error {
						var secondSecretID string
						if err := getSecretID(s, resourceName, &secondSecretID); err != nil {
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

func TestAccMcpConfigSecret_maxSecretsLimit(t *testing.T) {
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
				// A config's ingress SA allows a maximum of 2 concurrent secrets.
				// Creating both here confirms rotation overlap is possible.
				Config: configBasic(orgID, name, 720, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName+"_1"),
					checkExists(resourceName+"_2"),
				),
			},
		},
	})
}

func configBasic(orgID, name string, secretExpiresAfterHours int, twoSecrets bool) string {
	if twoSecrets {
		return fmt.Sprintf(`
			resource "mongodbatlas_mcp_config" "test" {
				org_id          = %[1]q
				mcp_config_name = %[2]q
				roles           = ["ORG_READ_ONLY"]
			}

			resource "mongodbatlas_mcp_config_secret" "test_1" {
				org_id                     = %[1]q
				mcp_config_id              = mongodbatlas_mcp_config.test.mcp_config_id
				secret_expires_after_hours = %[3]d
			}

			resource "mongodbatlas_mcp_config_secret" "test_2" {
				org_id                     = %[1]q
				mcp_config_id              = mongodbatlas_mcp_config.test.mcp_config_id
				secret_expires_after_hours = %[3]d

				depends_on = [mongodbatlas_mcp_config_secret.test_1]
			}
		`, orgID, name, secretExpiresAfterHours)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_mcp_config" "test" {
			org_id          = %[1]q
			mcp_config_name = %[2]q
			roles           = ["ORG_READ_ONLY"]
		}

		resource "mongodbatlas_mcp_config_secret" "test" {
			org_id                     = %[1]q
			mcp_config_id              = mongodbatlas_mcp_config.test.mcp_config_id
			secret_expires_after_hours = %[3]d
		}

		data "mongodbatlas_mcp_config_secret" "test" {
			org_id        = %[1]q
			mcp_config_id = mongodbatlas_mcp_config.test.mcp_config_id
			secret_id     = mongodbatlas_mcp_config_secret.test.secret_id
		}

		data "mongodbatlas_mcp_config_secrets" "test" {
			org_id        = %[1]q
			mcp_config_id = mongodbatlas_mcp_config.test.mcp_config_id
			depends_on    = [mongodbatlas_mcp_config_secret.test]
		}
	`, orgID, name, secretExpiresAfterHours)
}

func checkBasic(isCreate bool) resource.TestCheckFunc {
	commonAttrsSet := []string{"secret_id", "created_at", "expires_at"}
	checks := acc.CheckRSAndDS(resourceName, new(dataSourceName), nil, commonAttrsSet, nil, checkExists(resourceName))

	additionalChecks := []resource.TestCheckFunc{}
	if isCreate {
		additionalChecks = acc.AddAttrSetChecks(resourceName, additionalChecks, "secret")
	}
	additionalChecks = acc.AddAttrSetChecks(dataSourceName, additionalChecks, "masked_secret_value")
	additionalChecks = acc.AddAttrSetChecksPrefix(dataSourcePluralName, additionalChecks, []string{"masked_secret_value"}, "results.0")

	return resource.ComposeAggregateTestCheckFunc(checks, resource.ComposeAggregateTestCheckFunc(additionalChecks...))
}

func getSecretID(s *terraform.State, resourceName string, secretID *string) error {
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

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["org_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		secretID := rs.Primary.Attributes["secret_id"]
		if orgID == "" || mcpConfigID == "" || secretID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnPreview().RemoteMCPConfigurationsAPI.GetOrgMcpSecret(context.Background(), orgID, mcpConfigID, secretID).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("mcp config secret (%s/%s/%s) does not exist: %w", orgID, mcpConfigID, secretID, err)
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_mcp_config_secret" {
			continue
		}
		orgID := rs.Primary.Attributes["org_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		secretID := rs.Primary.Attributes["secret_id"]
		if orgID == "" || mcpConfigID == "" || secretID == "" {
			return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnPreview().RemoteMCPConfigurationsAPI.GetOrgMcpSecret(context.Background(), orgID, mcpConfigID, secretID).Execute()
		if err == nil {
			return fmt.Errorf("mcp config secret (%s/%s/%s) still exists", orgID, mcpConfigID, secretID)
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
		orgID := rs.Primary.Attributes["org_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		secretID := rs.Primary.Attributes["secret_id"]
		if orgID == "" || mcpConfigID == "" || secretID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s", orgID, mcpConfigID, secretID), nil
	}
}
