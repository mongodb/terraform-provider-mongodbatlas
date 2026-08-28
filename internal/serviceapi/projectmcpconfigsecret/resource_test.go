package projectmcpconfigsecret_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_project_mcp_config_secret.test"
const dataSourceName = "data.mongodbatlas_project_mcp_config_secret.test"
const dataSourcePluralName = "data.mongodbatlas_project_mcp_config_secrets.test"

func TestAccProjectMcpConfigSecret_basic(t *testing.T) {
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

func TestAccProjectMcpConfigSecret_rotate(t *testing.T) {
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
				// create the original secret.
				Config: configSecrets(projectID, name, 720, "test"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					func(s *terraform.State) error {
						return getSecretID(s, resourceName, &firstSecretID)
					},
				),
			},
			{
				// add a second secret.
				Config: configSecrets(projectID, name, 720, "test", "test_2"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					checkExists(resourceName+"_2"),
				),
			},
			{
				// rotate the first secret.
				// `taint` is deprecated in favor of -replace (https://developer.hashicorp.com/terraform/cli/commands/taint)
				// but testing plugin doesn't support -replace so using taint instead.
				Taint:  []string{resourceName},
				Config: configSecrets(projectID, name, 720, "test", "test_2"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					checkExists(resourceName+"_2"),
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

// builds a mongodbatlas_project_mcp_config_secret resource for each given address
// without data sources.
func configSecrets(projectID, name string, secretExpiresAfterHours int, addrs ...string) string {
	var secretsHCL strings.Builder
	for _, addr := range addrs {
		fmt.Fprintf(&secretsHCL, `
			resource "mongodbatlas_project_mcp_config_secret" "%[1]s" {
				project_id                 = %[2]q
				mcp_config_id              = mongodbatlas_project_mcp_config.test.mcp_config_id
				secret_expires_after_hours = %[3]d
			}
		`, addr, projectID, secretExpiresAfterHours)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_project_mcp_config" "test" {
			project_id      = %[1]q
			mcp_config_name = %[2]q
			roles           = ["GROUP_READ_ONLY"]
		}

		%[3]s
	`, projectID, name, secretsHCL.String())
}

// builds a single mongodbatlas_project_mcp_config_secret resource + its singular/plural data sources.
func configBasic(projectID, name string, secretExpiresAfterHours int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_mcp_config" "test" {
			project_id      = %[1]q
			mcp_config_name = %[2]q
			roles           = ["GROUP_READ_ONLY"]
		}

		resource "mongodbatlas_project_mcp_config_secret" "test" {
			project_id                 = %[1]q
			mcp_config_id              = mongodbatlas_project_mcp_config.test.mcp_config_id
			secret_expires_after_hours = %[3]d
		}

		data "mongodbatlas_project_mcp_config_secret" "test" {
			project_id    = %[1]q
			mcp_config_id = mongodbatlas_project_mcp_config.test.mcp_config_id
			secret_id     = mongodbatlas_project_mcp_config_secret.test.secret_id
		}

		data "mongodbatlas_project_mcp_config_secrets" "test" {
			project_id    = %[1]q
			mcp_config_id = mongodbatlas_project_mcp_config.test.mcp_config_id
			depends_on    = [mongodbatlas_project_mcp_config_secret.test]
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
		projectID := rs.Primary.Attributes["project_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		secretID := rs.Primary.Attributes["secret_id"]
		if projectID == "" || mcpConfigID == "" || secretID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		_, _, err := acc.ConnPreview().RemoteMCPConfigurationsAPI.GetGroupMcpSecret(context.Background(), projectID, mcpConfigID, secretID).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("mcp config secret (%s/%s/%s) does not exist: %w", projectID, mcpConfigID, secretID, err)
	}
}

func checkDestroy(s *terraform.State) error {
	var errs []error
	for name, rs := range s.RootModule().Resources {
		if !strings.HasPrefix(name, "mongodbatlas_project_mcp_config_secret.") {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		mcpConfigID := rs.Primary.Attributes["mcp_config_id"]
		secretID := rs.Primary.Attributes["secret_id"]
		if projectID == "" || mcpConfigID == "" || secretID == "" {
			errs = append(errs, fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName))
			continue
		}
		stillExists := func() bool {
			_, _, err := acc.ConnPreview().RemoteMCPConfigurationsAPI.GetGroupMcpSecret(context.Background(), projectID, mcpConfigID, secretID).Execute()
			return err == nil
		}
		if !acc.WaitUntilGoneOk(stillExists) {
			errs = append(errs, fmt.Errorf("mcp config secret (%s/%s/%s) still exists", projectID, mcpConfigID, secretID))
		}
	}
	if err := acc.CheckDestroyDeleteOrgMcpConfigs(s); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
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
