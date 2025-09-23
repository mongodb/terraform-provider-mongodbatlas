package acc

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/require"
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
				key_vault_name      = "%s"
				key_identifier      = "%s"
				secret  		    = "%s"
				tenant_id  		    = "%s"
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
		locals {
			project_id                 = %[1]q
			aws_ear_enabled			   = %[2]t
			aws_iam_role_name          = %[3]q
			aws_iam_role_policy_name   = %[4]q
			aws_customer_master_key_id = %[5]q
			aws_region				   = %[6]q
		}

		resource "aws_iam_role_policy" "test_policy" {
			name = local.aws_iam_role_policy_name
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
							"${local.aws_customer_master_key_id}"
						]
					}
			  ]
			})
		}
	  
		resource "aws_iam_role" "test_role" {
			name = local.aws_iam_role_name
	  
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

		resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
			project_id    = local.project_id
			provider_name = "AWS"
		}
	  
		resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
			project_id = local.project_id
			role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id
	  
			aws {
				iam_assumed_role_arn = aws_iam_role.test_role.arn
			}
		}

		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = local.project_id

			aws_kms_config {
				enabled                = local.aws_ear_enabled
				customer_master_key_id = local.aws_customer_master_key_id
				region                 = local.aws_region
				role_id                = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
				%[7]s
			}
			%[8]s
		}

		%[9]s
	`, projectID, awsKms.GetEnabled(), awsIAMRoleName, awsIAMRolePolicyName, awsKms.GetCustomerMasterKeyID(), awsKms.GetRegion(), requirePrivateNetworkingStr, enabledForSearchNodesStr, datasourceStr)
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

// EncryptionAtRestExecution creates an encryption at rest configuration for test execution.
func EncryptionAtRestExecution(tb testing.TB) string {
	tb.Helper()
	SkipInUnitTest(tb)
	require.True(tb, sharedInfo.init, "SetupSharedResources must called from TestMain test package")

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_EAR_PE_AWS_ID")

	sharedInfo.mu.Lock()
	defer sharedInfo.mu.Unlock()

	// lazy creation so it's only done if really needed
	if !sharedInfo.encryptionAtRestEnabled {
		tb.Logf("Creating execution encryption at rest configuration for project: %s\n", projectID)

		// Create encryption at rest configuration using environment variables
		awsKms := &admin.AWSKMSConfiguration{
			Enabled:                  conversion.Pointer(true),
			CustomerMasterKeyID:      conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
			Region:                   conversion.StringPtr(conversion.AWSRegionToMongoDBRegion(os.Getenv("AWS_REGION"))),
			RoleId:                   conversion.StringPtr(os.Getenv("AWS_EAR_ROLE_ID")),
			RequirePrivateNetworking: conversion.Pointer(true),
		}

		createEncryptionAtRest(tb, projectID, awsKms)
		sharedInfo.encryptionAtRestEnabled = true
	}

	return projectID
}

func createEncryptionAtRest(tb testing.TB, projectID string, aws *admin.AWSKMSConfiguration) {
	tb.Helper()

	encryptionAtRestReq := &admin.EncryptionAtRest{
		AwsKms: aws,
	}

	_, _, err := ConnV2().EncryptionAtRestUsingCustomerKeyManagementApi.UpdateEncryptionAtRest(tb.Context(), projectID, encryptionAtRestReq).Execute()
	require.NoError(tb, err, "Failed to create encryption at rest configuration for project: %s", projectID)
}

func deleteEncryptionAtRest(projectID string) {
	// Disable encryption at rest by setting all providers to disabled
	encryptionAtRestReq := &admin.EncryptionAtRest{
		AwsKms: &admin.AWSKMSConfiguration{
			Enabled: conversion.Pointer(false),
		},
		AzureKeyVault: &admin.AzureKeyVault{
			Enabled: conversion.Pointer(false),
		},
		GoogleCloudKms: &admin.GoogleCloudKMS{
			Enabled: conversion.Pointer(false),
		},
	}

	_, _, err := ConnV2().EncryptionAtRestUsingCustomerKeyManagementApi.UpdateEncryptionAtRest(context.Background(), projectID, encryptionAtRestReq).Execute()
	if err != nil {
		fmt.Printf("Failed to delete encryption at rest for project %s: %s\n", projectID, err)
	}
}
