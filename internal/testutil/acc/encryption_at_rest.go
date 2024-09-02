package acc

import (
	"fmt"

	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

func ConfigEARAzureKeyVault(projectID string, azure *admin.AzureKeyVault, useRequirePrivateNetworking bool) string {
	var requirePrivateNetworkingAttr string
	if useRequirePrivateNetworking {
		requirePrivateNetworkingAttr = fmt.Sprintf("require_private_networking = %t", azure.GetRequirePrivateNetworking())
	}

	return fmt.Sprintf(`
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

		%s
	`, projectID, *azure.Enabled, azure.GetClientID(), azure.GetAzureEnvironment(), azure.GetSubscriptionID(), azure.GetResourceGroupName(),
		azure.GetKeyVaultName(), azure.GetKeyIdentifier(), azure.GetSecret(), azure.GetTenantID(), requirePrivateNetworkingAttr, TestAccDatasourceConfig())
}

func TestAccDatasourceConfig() string {
	return `data "mongodbatlas_encryption_at_rest" "test" {
			project_id = mongodbatlas_encryption_at_rest.test.project_id
		}`
}
