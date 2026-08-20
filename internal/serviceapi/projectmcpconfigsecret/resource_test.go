package projectmcpconfigsecret_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_project_mcp_config_secret.test"
const dataSourceName = "data.mongodbatlas_project_mcp_config_secret.test"
const dataSourcePluralName = "data.mongodbatlas_project_mcp_config_secrets.test"

// NOTE: the Atlas Go SDK (atlas-sdk-go admin.RemoteMCPConfigurationsApi) does not yet expose
// create/read/delete operations for MCP config secrets (only Get/List for the config itself).
// Because of this, checkExists below can only verify state, not the live API.
// Update this once the SDK adds secret support.

func TestAccProjectMcpConfigSecret_basic(t *testing.T) {
	var (
		projectID = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name      = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, name, 720),
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

func TestAccProjectMcpConfigSecret_maxSecretsLimit(t *testing.T) {
	var (
		projectID = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name      = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				// A config's ingress SA allows a maximum of 2 concurrent secrets.
				// Creating both here confirms rotation overlap is possible.
				Config: configTwoSecrets(projectID, name, 720),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName+"_1"),
					checkExists(resourceName+"_2"),
				),
			},
		},
	})
}

func configBasic(projectID, name string, secretExpiresAfterHours int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_mcp_config" "test" {
			project_id        = %[1]q
			mcp_config_name = %[2]q
			roles           = ["GROUP_READ_ONLY"]
		}

		resource "mongodbatlas_project_mcp_config_secret" "test" {
			project_id                   = %[1]q
			mcp_config_id              = mongodbatlas_project_mcp_config.test.mcp_config_id
			secret_expires_after_hours = %[3]d
		}

		data "mongodbatlas_project_mcp_config_secret" "test" {
			project_id      = %[1]q
			mcp_config_id = mongodbatlas_project_mcp_config.test.mcp_config_id
			secret_id     = mongodbatlas_project_mcp_config_secret.test.secret_id
		}

		data "mongodbatlas_project_mcp_config_secrets" "test" {
			project_id      = %[1]q
			mcp_config_id = mongodbatlas_project_mcp_config.test.mcp_config_id
			depends_on    = [mongodbatlas_project_mcp_config_secret.test]
		}
	`, projectID, name, secretExpiresAfterHours)
}

func configTwoSecrets(projectID, name string, secretExpiresAfterHours int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_mcp_config" "test" {
			project_id        = %[1]q
			mcp_config_name = %[2]q
			roles           = ["GROUP_READ_ONLY"]
		}

		resource "mongodbatlas_project_mcp_config_secret" "test_1" {
			project_id                   = %[1]q
			mcp_config_id              = mongodbatlas_project_mcp_config.test.mcp_config_id
			secret_expires_after_hours = %[3]d
		}

		resource "mongodbatlas_project_mcp_config_secret" "test_2" {
			project_id                   = %[1]q
			mcp_config_id              = mongodbatlas_project_mcp_config.test.mcp_config_id
			secret_expires_after_hours = %[3]d

			depends_on = [mongodbatlas_project_mcp_config_secret.test_1]
		}
	`, projectID, name, secretExpiresAfterHours)
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

// checkExists verifies state only. See NOTE at top of file re: SDK gap for secrets.
func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.Attributes["secret_id"] == "" {
			return fmt.Errorf("checkExists, secret_id not found for: %s", resourceName)
		}
		return nil
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		secretID := rs.Primary.Attributes["secret_id"]
		if projectID == "" || mcpConfigID == "" || secretID == "" {
			return "", fmt.Errorf("import, attributes not found for: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s", projectID, mcpConfigID, secretID), nil
	}
}
