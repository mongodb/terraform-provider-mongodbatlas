package acc

import (
	"context"
	"fmt"
	"strconv"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

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
			project_id = %[1]q

			azure_key_vault_config {
				enabled             = %[2]t
				client_id           = %[3]q
				azure_environment   = %[4]q
				subscription_id     = %[5]q
				resource_group_name = %[6]q
				key_vault_name      = %[7]q
				key_identifier      = %[8]q
				secret  		    = %[9]q
				tenant_id  		    = %[10]q
				%[11]s
			}
		}
	`, projectID, *azure.Enabled, azure.GetClientID(), azure.GetAzureEnvironment(), azure.GetSubscriptionID(), azure.GetResourceGroupName(),
		azure.GetKeyVaultName(), azure.GetKeyIdentifier(), azure.GetSecret(), azure.GetTenantID(), requirePrivateNetworkingAttr)

	if useDatasource {
		return fmt.Sprintf(`%s %s`, config, EARDatasourceConfig())
	}
	return config
}

func ConfigAwsKmsWithRole(projectID, awsIAMRoleName, awsIAMRolePolicyName string, awsKms *admin.AWSKMSConfiguration, useDatasource, useRequirePrivateNetworking, useEnabledForSearchNodes bool) string {
	requirePrivateNetworkingStr := ""
	if useRequirePrivateNetworking {
		requirePrivateNetworkingStr = fmt.Sprintf("require_private_networking = %t", awsKms.GetRequirePrivateNetworking())
	}
	enabledForSearchNodesStr := ""
	if useEnabledForSearchNodes {
		enabledForSearchNodesStr = fmt.Sprintf("enabled_for_search_nodes = %t", useEnabledForSearchNodes)
	}
	datasourceStr := ""
	if useDatasource {
		datasourceStr = EARDatasourceConfig()
	}
	config := fmt.Sprintf(`
		%[2]s

		resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
			project_id    = %[1]q
			provider_name = "AWS"
		}
	  
		resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
			project_id = %[1]q
			role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
	  
			aws {
				iam_assumed_role_arn = aws_iam_role.test_role.arn
			}
		}

		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = %[1]q

			aws_kms_config {
				enabled                = %[3]t
				customer_master_key_id = %[4]q
				region                 = %[5]q
				role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
				%[6]s
			}
			%[7]s
		}

		%[8]s
	`, projectID, configAwsRoleAndPolicy(awsIAMRoleName, awsIAMRolePolicyName, awsKms), awsKms.GetEnabled(), awsKms.GetCustomerMasterKeyID(), awsKms.GetRegion(), requirePrivateNetworkingStr, enabledForSearchNodesStr, datasourceStr)

	return config
}

func ConfigProjectWithAwsKmsPrivateNetworking(projectName, orgID, awsIAMRoleName, awsIAMRolePolicyName string, awsKms *admin.AWSKMSConfiguration, useDatasource, useRequirePrivateNetworking, useEnabledForSearchNodes bool) string {
	config := fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
		   name   = %[1]q
		   org_id = %[2]q
		}

		%[3]s

		resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
			project_id    = mongodbatlas_project.test.id
			provider_name = "AWS"
		}
	  
		resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
			project_id = mongodbatlas_project.test.id
			role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
	  
			aws {
				iam_assumed_role_arn = aws_iam_role.test_role.arn
			}
		}

		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = mongodbatlas_project.test.id

			aws_kms_config {
				enabled                    = %[4]t
				customer_master_key_id     = %[5]q
				region                     = %[6]q
				role_id                    = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
				require_private_networking = %[7]t, 
			}
		}
	`, projectName, orgID, configAwsRoleAndPolicy(awsIAMRoleName, awsIAMRolePolicyName, awsKms), awsKms.GetEnabled(), awsKms.GetCustomerMasterKeyID(), awsKms.GetRegion(), awsKms.GetRequirePrivateNetworking())
	return config
}

func configAwsRoleAndPolicy(awsIamRoleName, awsIAMRolePolicyName string, awsKms *admin.AWSKMSConfiguration) string {
	config := fmt.Sprintf(`
		resource "aws_iam_role" "test_role" {
			name = %[1]q
	  
			assume_role_policy = jsonencode({
				"Version" : "2012-10-17",
				"Statement" : [
					{
						"Effect" : "Allow",
						"Principal" : {
							"AWS" : "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_aws_account_arn}"
						},
						"Action" : "sts:AssumeRole",
						"Condition" : {
							"StringEquals" : {
								"sts:ExternalId" : "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config[0].atlas_assumed_role_external_id}"
							}
						}
					}
				]
			})
		}

		resource "aws_iam_role_policy" "test_policy" {
			name = %[2]q
			role = aws_iam_role.test_role.id
	  
			policy = jsonencode({
				"Version" : "2012-10-17",
				"Statement" : [
					{
						"Effect" : "Allow",
						"Action" : [
							"kms:Decrypt",
							"kms:Encrypt",
							"kms:DescribeKey"
						],
						"Resource" : [
 							%[3]q
						]
					}
			  ]
			})
		}
	`, awsIamRoleName, awsIAMRolePolicyName, awsKms.GetCustomerMasterKeyID())
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
