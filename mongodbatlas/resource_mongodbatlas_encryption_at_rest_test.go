package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasEncryptionAtRest_basicAWS(t *testing.T) {
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		awsKms = matlas.AwsKms{
			Enabled:             pointy.Bool(true),
			AccessKeyID:         os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey:     os.Getenv("AWS_SECRET_ACCESS_KEY"),
			CustomerMasterKeyID: os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID"),
			Region:              os.Getenv("AWS_REGION"),
		}

		awsKmsUpdated = matlas.AwsKms{
			Enabled:             pointy.Bool(true),
			AccessKeyID:         os.Getenv("AWS_ACCESS_KEY_ID_UPDATED"),
			SecretAccessKey:     os.Getenv("AWS_SECRET_ACCESS_KEY_UPDATED"),
			CustomerMasterKeyID: os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID_UPDATED"),
			Region:              os.Getenv("AWS_REGION_UPDATED"),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkAwsEnv(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(projectID, &awsKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.enabled", cast.ToString(awsKms.Enabled)),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.access_key_id", awsKms.AccessKeyID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.secret_access_key", awsKms.SecretAccessKey),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.customer_master_key_id", awsKms.CustomerMasterKeyID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.region", awsKms.Region),
				),
			},
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(projectID, &awsKmsUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.enabled", cast.ToString(awsKmsUpdated.Enabled)),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.access_key_id", awsKmsUpdated.AccessKeyID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.secret_access_key", awsKmsUpdated.SecretAccessKey),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.customer_master_key_id", awsKmsUpdated.CustomerMasterKeyID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.region", awsKmsUpdated.Region),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasEncryptionAtRest_basicAzure(t *testing.T) {
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		azureKeyVault = matlas.AzureKeyVault{
			Enabled:           pointy.Bool(true),
			ClientID:          os.Getenv("AZURE_CLIENT_ID"),
			AzureEnvironment:  "AZURE",
			SubscriptionID:    os.Getenv("AZURE_SUBCRIPTION_ID"),
			ResourceGroupName: os.Getenv("AZURE_RESOURCE_GROUP_NAME"),
			KeyVaultName:      os.Getenv("AZURE_KEY_VAULT_NAME"),
			KeyIdentifier:     os.Getenv("AZURE_KEY_IDENTIFIER"),
			Secret:            os.Getenv("AZURE_SECRET"),
			TenantID:          os.Getenv("AZURE_TENANT_ID"),
		}

		azureKeyVaultUpdated = matlas.AzureKeyVault{
			Enabled:           pointy.Bool(true),
			ClientID:          os.Getenv("AZURE_CLIENT_ID_UPDATED"),
			AzureEnvironment:  "AZURE",
			SubscriptionID:    os.Getenv("AZURE_SUBCRIPTION_ID"),
			ResourceGroupName: os.Getenv("AZURE_RESOURCE_GROUP_NAME_UPDATED"),
			KeyVaultName:      os.Getenv("AZURE_KEY_VAULT_NAME_UPDATED"),
			KeyIdentifier:     os.Getenv("AZURE_KEY_IDENTIFIER_UPDATED"),
			Secret:            os.Getenv("AZURE_SECRET_UPDATED"),
			TenantID:          os.Getenv("AZURE_TENANT_ID"),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkEncryptionAtRestEnvAzure(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAzureKeyVault(projectID, &azureKeyVault),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.enabled", cast.ToString(azureKeyVault.Enabled)),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.client_id", azureKeyVault.ClientID),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.azure_environment", azureKeyVault.AzureEnvironment),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.subscription_id", azureKeyVault.SubscriptionID),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.resource_group_name", azureKeyVault.ResourceGroupName),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.key_vault_name", azureKeyVault.KeyVaultName),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.key_identifier", azureKeyVault.KeyIdentifier),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.secret", azureKeyVault.Secret),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.tenant_id", azureKeyVault.TenantID),
				),
			},
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAzureKeyVault(projectID, &azureKeyVaultUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.enabled", cast.ToString(azureKeyVaultUpdated.Enabled)),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.client_id", azureKeyVaultUpdated.ClientID),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.azure_environment", azureKeyVaultUpdated.AzureEnvironment),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.subscription_id", azureKeyVaultUpdated.SubscriptionID),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.resource_group_name", azureKeyVaultUpdated.ResourceGroupName),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.key_vault_name", azureKeyVaultUpdated.KeyVaultName),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.key_identifier", azureKeyVaultUpdated.KeyIdentifier),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.secret", azureKeyVaultUpdated.Secret),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault.tenant_id", azureKeyVaultUpdated.TenantID),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasEncryptionAtRest_basicGCP(t *testing.T) {
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		googleCloudKms = matlas.GoogleCloudKms{
			Enabled:              pointy.Bool(true),
			ServiceAccountKey:    os.Getenv("GCP_SERVICE_ACCOUNT_KEY"),
			KeyVersionResourceID: os.Getenv("GCP_KEY_VERSION_RESOURCE_ID"),
		}

		googleCloudKmsUpdated = matlas.GoogleCloudKms{
			Enabled:              pointy.Bool(true),
			ServiceAccountKey:    os.Getenv("GCP_SERVICE_ACCOUNT_KEY_UPDATED"),
			KeyVersionResourceID: os.Getenv("GCP_KEY_VERSION_RESOURCE_ID_UPDATED"),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckGPCEnv(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigGoogleCloudKms(projectID, &googleCloudKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms.enabled", cast.ToString(googleCloudKms.Enabled)),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms.service_account_key", googleCloudKms.ServiceAccountKey),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms.key_version_resource_id", googleCloudKms.KeyVersionResourceID),
				),
			},
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigGoogleCloudKms(projectID, &googleCloudKmsUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms.enabled", cast.ToString(googleCloudKmsUpdated.Enabled)),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms.service_account_key", googleCloudKmsUpdated.ServiceAccountKey),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms.key_version_resource_id", googleCloudKmsUpdated.KeyVersionResourceID),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		if _, _, err := conn.EncryptionsAtRest.Get(context.Background(), rs.Primary.ID); err == nil {
			return nil
		}

		return fmt.Errorf("encryptionAtRest (%s) does not exist", rs.Primary.ID)
	}
}

func testAccCheckMongoDBAtlasEncryptionAtRestDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_encryption_at_rest" {
			continue
		}

		res, _, err := conn.EncryptionsAtRest.Get(context.Background(), rs.Primary.ID)

		if err != nil ||
			(*res.AwsKms.Enabled != false &&
				*res.AzureKeyVault.Enabled != false &&
				*res.GoogleCloudKms.Enabled != false) {
			return fmt.Errorf("encryptionAtRest (%s) still exists: err: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(projectID string, aws *matlas.AwsKms) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "%s"

		  aws_kms = {
				enabled                = %t
				access_key_id          = "%s"
				secret_access_key      = "%s"
				customer_master_key_id = "%s"
				region                 = "%s"
			}
		}
	`, projectID, *aws.Enabled, aws.AccessKeyID, aws.SecretAccessKey, aws.CustomerMasterKeyID, aws.Region)
}

func testAccMongoDBAtlasEncryptionAtRestConfigAzureKeyVault(projectID string, azure *matlas.AzureKeyVault) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "%s"

		  azure_key_vault = {
				enabled             = %t
				client_id           = "%s"
				azure_environment   = "%s"
				subscription_id     = "%s"
				resource_group_name = "%s"
				key_vault_name  	  = "%s"
				key_identifier  	  = "%s"
				secret  						= "%s"
				tenant_id  					= "%s"
			}
		}
	`, projectID, *azure.Enabled, azure.ClientID, azure.AzureEnvironment, azure.SubscriptionID, azure.ResourceGroupName,
		azure.KeyVaultName, azure.KeyIdentifier, azure.Secret, azure.TenantID)
}

func testAccMongoDBAtlasEncryptionAtRestConfigGoogleCloudKms(projectID string, google *matlas.GoogleCloudKms) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "%s"

		  google_cloud_kms = {
				enabled                 = %t
				service_account_key     = "%s"
				key_version_resource_id = "%s"
			}
		}
	`, projectID, *google.Enabled, google.ServiceAccountKey, google.KeyVersionResourceID)
}
