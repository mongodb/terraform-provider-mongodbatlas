package acc

import (
	"context"
	"fmt"
	"strconv"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

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

func ConfigAwsKms(projectID string, aws *admin.AWSKMSConfiguration, useDatasource, useRequirePrivateNetworking, useEnabledForSearchNodes bool) string {
	requirePrivateNetworkingStr := ""
	if useRequirePrivateNetworking {
		requirePrivateNetworkingStr = fmt.Sprintf("require_private_networking = %t", aws.GetRequirePrivateNetworking())
	}
	enabledForSearchNodes := ""
	if useEnabledForSearchNodes {
		enabledForSearchNodes = fmt.Sprintf("enabled_for_search_nodes = %t", useEnabledForSearchNodes)
	}
	config := fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = %[1]q

		  	aws_kms_config {
				enabled                = %[2]t
				customer_master_key_id = %[3]q
				region                 = %[4]q
				role_id              = %[5]q
				%[6]s
			}
			%[7]s
		}
	`, projectID, aws.GetEnabled(), aws.GetCustomerMasterKeyID(), aws.GetRegion(), aws.GetRoleId(), requirePrivateNetworkingStr, enabledForSearchNodes)

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

func ConvertToAwsKmsEARAttrMap(awsKms *admin.AWSKMSConfiguration) map[string]string {
	return map[string]string{
		"enabled":                    strconv.FormatBool(awsKms.GetEnabled()),
		"region":                     awsKms.GetRegion(),
		"role_id":                    awsKms.GetRoleId(),
		"customer_master_key_id":     awsKms.GetCustomerMasterKeyID(),
		"valid":                      "true",
		"require_private_networking": strconv.FormatBool(awsKms.GetRequirePrivateNetworking()),
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
