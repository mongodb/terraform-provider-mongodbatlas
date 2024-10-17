package encryptionatrest_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"go.mongodb.org/atlas-sdk/v20240805005/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigEncryptionAtRest_basicAWS(t *testing.T) {
	acc.SkipTestForCI(t) // needs AWS configuration

	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		awsKms = admin.AWSKMSConfiguration{
			Enabled:             conversion.Pointer(true),
			CustomerMasterKeyID: conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
			Region:              conversion.StringPtr(conversion.AWSRegionToMongoDBRegion(os.Getenv("AWS_REGION"))),
			RoleId:              conversion.StringPtr(os.Getenv("AWS_ROLE_ID")),
		}
		useDatasource = mig.IsProviderVersionAtLeast("1.19.0") // data source introduced in this version
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheck(t); acc.PreCheckAwsEnv(t) },
		CheckDestroy: acc.EARDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configAwsKms(projectID, &awsKms, useDatasource),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.region", awsKms.GetRegion()),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.role_id", awsKms.GetRoleId()),
				),
			},
			mig.TestStepCheckEmptyPlan(configAwsKms(projectID, &awsKms, useDatasource)),
		},
	})
}

func TestMigEncryptionAtRest_withRole_basicAWS(t *testing.T) {
	acc.SkipTestForCI(t) // needs AWS configuration

	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		awsIAMRoleName       = acc.RandomIAMRole()
		awsIAMRolePolicyName = fmt.Sprintf("%s-policy", awsIAMRoleName)
		awsKeyName           = acc.RandomName()

		awsKms = admin.AWSKMSConfiguration{
			Enabled:             conversion.Pointer(true),
			Region:              conversion.StringPtr(conversion.AWSRegionToMongoDBRegion(os.Getenv("AWS_REGION"))),
			CustomerMasterKeyID: conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
		}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheck(t); acc.PreCheckAwsEnv(t) },
		CheckDestroy: acc.EARDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProvidersWithAWS(),
				Config:            testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(projectID, awsIAMRoleName, awsIAMRolePolicyName, awsKeyName, &awsKms),
			},
			{
				ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(projectID, awsIAMRoleName, awsIAMRolePolicyName, awsKeyName, &awsKms),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.region", awsKms.GetRegion()),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(), // special case using AWS resources
					},
				},
			},
		},
	})
}

func TestMigEncryptionAtRest_basicAzure(t *testing.T) {
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = acc.ProjectIDExecution(t)

		azureKeyVault = admin.AzureKeyVault{
			Enabled:           conversion.Pointer(true),
			ClientID:          conversion.StringPtr(os.Getenv("AZURE_CLIENT_ID")),
			AzureEnvironment:  conversion.StringPtr("AZURE"),
			SubscriptionID:    conversion.StringPtr(os.Getenv("AZURE_SUBSCRIPTION_ID")),
			ResourceGroupName: conversion.StringPtr(os.Getenv("AZURE_RESOURCE_GROUP_NAME")),
			KeyVaultName:      conversion.StringPtr(os.Getenv("AZURE_KEY_VAULT_NAME")),
			KeyIdentifier:     conversion.StringPtr(os.Getenv("AZURE_KEY_IDENTIFIER")),
			Secret:            conversion.StringPtr(os.Getenv("AZURE_APP_SECRET")),
			TenantID:          conversion.StringPtr(os.Getenv("AZURE_TENANT_ID")),
		}

		attrMap = map[string]string{
			"enabled":             strconv.FormatBool(azureKeyVault.GetEnabled()),
			"azure_environment":   azureKeyVault.GetAzureEnvironment(),
			"resource_group_name": azureKeyVault.GetResourceGroupName(),
			"key_vault_name":      azureKeyVault.GetKeyVaultName(),
			"client_id":           azureKeyVault.GetClientID(),
			"key_identifier":      azureKeyVault.GetKeyIdentifier(),
			"subscription_id":     azureKeyVault.GetSubscriptionID(),
			"tenant_id":           azureKeyVault.GetTenantID(),
		}

		useDatasource = mig.IsProviderVersionAtLeast("1.19.0") // data source introduced in this version
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t); acc.PreCheckEncryptionAtRestEnvAzure(t) },
		CheckDestroy: acc.EARDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            acc.ConfigEARAzureKeyVault(projectID, &azureKeyVault, false, useDatasource),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					acc.EARCheckResourceAttr(resourceName, "azure_key_vault_config.0", attrMap),
				),
			},
			mig.TestStepCheckEmptyPlan(acc.ConfigEARAzureKeyVault(projectID, &azureKeyVault, false, useDatasource)),
		},
	})
}

func TestMigEncryptionAtRest_basicGCP(t *testing.T) {
	acc.SkipTestForCI(t) // needs GCP configuration

	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		googleCloudKms = admin.GoogleCloudKMS{
			Enabled:              conversion.Pointer(true),
			ServiceAccountKey:    conversion.StringPtr(os.Getenv("GCP_SERVICE_ACCOUNT_KEY")),
			KeyVersionResourceID: conversion.StringPtr(os.Getenv("GCP_KEY_VERSION_RESOURCE_ID")),
		}
		useDatasource = mig.IsProviderVersionAtLeast("1.19.0") // data source introduced in this version
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheck(t); acc.PreCheckGPCEnv(t) },
		CheckDestroy: acc.EARDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configGoogleCloudKms(projectID, &googleCloudKms, useDatasource),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms_config.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "google_cloud_kms_config.0.key_version_resource_id"),
				),
			},
			mig.TestStepCheckEmptyPlan(configGoogleCloudKms(projectID, &googleCloudKms, useDatasource)),
		},
	})
}

func TestMigEncryptionAtRest_basicAWS_from_v1_11_0(t *testing.T) {
	acc.SkipTestForCI(t) // needs AWS configuration

	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		awsKms = admin.AWSKMSConfiguration{
			Enabled:             conversion.Pointer(true),
			AccessKeyID:         conversion.StringPtr(os.Getenv("AWS_ACCESS_KEY_ID")),
			SecretAccessKey:     conversion.StringPtr(os.Getenv("AWS_SECRET_ACCESS_KEY")),
			CustomerMasterKeyID: conversion.StringPtr(os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID")),
			Region:              conversion.StringPtr(conversion.AWSRegionToMongoDBRegion(os.Getenv("AWS_REGION"))),
			RoleId:              conversion.StringPtr(os.Getenv("AWS_ROLE_ID")),
		}
		useDatasource = mig.IsProviderVersionAtLeast("1.19.0") // data source introduced in this version
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheck(t); acc.PreCheckAwsEnv(t) },
		CheckDestroy: acc.EARDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProvidersWithAWS("1.11.0"),
				Config:            configAwsKms(projectID, &awsKms, useDatasource),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckEARExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.region", awsKms.GetRegion()),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.role_id", awsKms.GetRoleId()),
				),
			},
			mig.TestStepCheckEmptyPlan(configAwsKms(projectID, &awsKms, useDatasource)),
		},
	})
}
