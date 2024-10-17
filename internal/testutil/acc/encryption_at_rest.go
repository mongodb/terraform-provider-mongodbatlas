package acc

import (
	"context"
	"fmt"
	"strconv"

	"go.mongodb.org/atlas-sdk/v20240805005/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func ConfigEARAzureKeyVault(projectID string, azure *admin.AzureKeyVault, useRequirePrivateNetworking, useDatasource bool) string {
	var requirePrivateNetworkingAttr string
	if useRequirePrivateNetworking {
		requirePrivateNetworkingAttr = fmt.Sprintf("require_private_networking = %t", azure.GetRequirePrivateNetworking())
	}

	config := fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "%s"

		  azure_key_vault_config {
				enabled             = %t
				client_id           = "%s"
				azure_environment   = "%s"
				subscription_id     = "%s"
				resource_group_name = "%s"
				key_vault_name  	  = "%s"
				key_identifier  	  = "%s"
				secret  						= "%s"
				tenant_id  					= "%s"
				%s
			}
		}
	`, projectID, *azure.Enabled, azure.GetClientID(), azure.GetAzureEnvironment(), azure.GetSubscriptionID(), azure.GetResourceGroupName(),
		azure.GetKeyVaultName(), azure.GetKeyIdentifier(), azure.GetSecret(), azure.GetTenantID(), requirePrivateNetworkingAttr)

	if useDatasource {
		return fmt.Sprintf(`%s %s`, config, EARDatasourceConfig())
	}
	return config
}

func EARDatasourceConfig() string {
	return `data "mongodbatlas_encryption_at_rest" "test" {
			project_id = mongodbatlas_encryption_at_rest.test.project_id
		}`
}

func CheckEARExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		if _, _, err := ConnV2().EncryptionAtRestUsingCustomerKeyManagementApi.GetEncryptionAtRest(context.Background(), rs.Primary.ID).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("encryptionAtRest (%s) does not exist", rs.Primary.ID)
	}
}

func ConvertToAzureKeyVaultEARAttrMap(az *admin.AzureKeyVault) map[string]string {
	return map[string]string{
		"enabled":                    strconv.FormatBool(az.GetEnabled()),
		"azure_environment":          az.GetAzureEnvironment(),
		"resource_group_name":        az.GetResourceGroupName(),
		"key_vault_name":             az.GetKeyVaultName(),
		"client_id":                  az.GetClientID(),
		"key_identifier":             az.GetKeyIdentifier(),
		"subscription_id":            az.GetSubscriptionID(),
		"tenant_id":                  az.GetTenantID(),
		"require_private_networking": strconv.FormatBool(az.GetRequirePrivateNetworking()),
	}
}

func EARCheckResourceAttr(resourceName, prefix string, attrsMap map[string]string) resource.TestCheckFunc {
	checks := AddAttrChecksPrefix(resourceName, []resource.TestCheckFunc{}, attrsMap, prefix)

	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func EARDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_encryption_at_rest" {
			continue
		}
		res, _, err := ConnV2().EncryptionAtRestUsingCustomerKeyManagementApi.GetEncryptionAtRest(context.Background(), rs.Primary.ID).Execute()
		if err != nil ||
			(res.AwsKms.GetEnabled() ||
				res.AzureKeyVault.GetEnabled() ||
				res.GoogleCloudKms.GetEnabled()) {
			return fmt.Errorf("encryptionAtRest (%s) still exists: err: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func EARImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
