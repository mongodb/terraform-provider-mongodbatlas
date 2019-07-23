package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
)

func TestAccResourceMongoDBAtlasEncryptionAtRest_basicAWS(t *testing.T) {
	var encryptionAtRest = matlas.EncryptionAtRest{}

	resourceName := "mongodbatlas_encryption_at_rest.test"
	projectID := "5d0f1f73cf09a29120e173cf"

	awsKms := matlas.AwsKms{
		Enabled:             pointy.Bool(true),
		AccessKeyID:         os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey:     os.Getenv("AWS_SECRET_ACCESS_KEY"),
		CustomerMasterKeyID: os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID"),
		Region:              "US_EAST_2",
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkAwsEnv(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(&awsKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName, &encryptionAtRest),
					testAccCheckMongoDBAtlasEncryptionAtRestAttributes(&encryptionAtRest, pointy.Bool(true)),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.enabled", cast.ToString(awsKms.Enabled)),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.access_key_id", awsKms.AccessKeyID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.secret_access_key", awsKms.SecretAccessKey),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.customer_master_key_id", awsKms.CustomerMasterKeyID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms.region", awsKms.Region),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasEncryptionAtRest_basicAzure(t *testing.T) {
	t.Skip()
	var encryptionAtRest = matlas.EncryptionAtRest{}

	resourceName := "mongodbatlas_encryption_at_rest.test"
	projectID := "5d0f1f73cf09a29120e173cf"

	if os.Getenv("AZURE_KEY_IDENTIFIER") == "" || os.Getenv("AZURE_SECRET") == "" {
		t.Fatal("`AZURE_KEY_IDENTIFIER` and `AZURE_SECRET` must be set for acceptance testing")
	}

	azureKeyVault := matlas.AzureKeyVault{
		Enabled:           pointy.Bool(true),
		ClientID:          "g54f9e2-89e3-40fd-8188-EXAMPLEID",
		AzureEnvironment:  "AZURE",
		SubscriptionID:    "0ec944e3-g725-44f9-a147-EXAMPLEID",
		ResourceGroupName: "ExampleRGName",
		KeyVaultName:      "EXAMPLEKeyVault",
		KeyIdentifier:     os.Getenv("AZURE_KEY_IDENTIFIER"),
		Secret:            os.Getenv("AZURE_SECRET"),
		TenantID:          "e8e4b6ba-ff32-4c88-a9af-EXAMPLEID",
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAzureKeyVault(&azureKeyVault),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName, &encryptionAtRest),
					testAccCheckMongoDBAtlasEncryptionAtRestAttributes(&encryptionAtRest, pointy.Bool(true)),
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
		},
	})
}

func TestAccResourceMongoDBAtlasEncryptionAtRest_basicGCP(t *testing.T) {
	t.Skip()

	var encryptionAtRest = matlas.EncryptionAtRest{}

	resourceName := "mongodbatlas_encryption_at_rest.test"
	projectID := "5d0f1f73cf09a29120e173cf"

	if os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY") == "" || os.Getenv("GOOGLE_KEY_VERSION_RESOURCE_ID") == "" {
		t.Fatal("`GOOGLE_SERVICE_ACCOUNT_KEY` and `GOOGLE_KEY_VERSION_RESOURCE_ID` must be set for acceptance testing")
	}

	googleCloudKms := matlas.GoogleCloudKms{
		Enabled:              pointy.Bool(true),
		ServiceAccountKey:    os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY"),
		KeyVersionResourceID: os.Getenv("GOOGLE_KEY_VERSION_RESOURCE_ID"),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigGoogleCloudKms(&googleCloudKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName, &encryptionAtRest),
					testAccCheckMongoDBAtlasEncryptionAtRestAttributes(&encryptionAtRest, pointy.Bool(true)),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms.enabled", cast.ToString(googleCloudKms.Enabled)),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms.service_account_key", googleCloudKms.ServiceAccountKey),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms.key_version_resource_id", googleCloudKms.KeyVersionResourceID),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName string, encryptionAtRest *matlas.EncryptionAtRest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] encryptionAtRest ID: %s", rs.Primary.ID)

		if encryptionRes, _, err := conn.EncryptionsAtRest.Get(context.Background(), rs.Primary.ID); err == nil {
			*encryptionAtRest = *encryptionRes
			return nil
		}
		return fmt.Errorf("encryptionAtRest (%s) does not exist", rs.Primary.ID)
	}
}

func testAccCheckMongoDBAtlasEncryptionAtRestAttributes(encryptionAtRest *matlas.EncryptionAtRest, enabled *bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *encryptionAtRest.AwsKms.Enabled != *enabled {
			return fmt.Errorf("bad encryptionAtRest enabled: %s", cast.ToString(encryptionAtRest.AwsKms.Enabled))
		}
		return nil
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

func testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(aws *matlas.AwsKms) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "5d0f1f73cf09a29120e173cf"

		  aws_kms = {
				enabled                = %s
				access_key_id          = "%s"
				secret_access_key      = "%s"
				customer_master_key_id = "%s"
				region                 = "%s"
			}
		}
	`, cast.ToString(*aws.Enabled), aws.AccessKeyID, aws.SecretAccessKey, aws.CustomerMasterKeyID, aws.Region)
}

func testAccMongoDBAtlasEncryptionAtRestConfigAzureKeyVault(azure *matlas.AzureKeyVault) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "5d0f1f73cf09a29120e173cf"

		  azure_key_vault = {
				enabled             = %s
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
	`, cast.ToString(*azure.Enabled), azure.ClientID, azure.AzureEnvironment, azure.SubscriptionID, azure.ResourceGroupName,
		azure.KeyVaultName, azure.KeyIdentifier, azure.Secret, azure.TenantID)
}

func testAccMongoDBAtlasEncryptionAtRestConfigGoogleCloudKms(google *matlas.GoogleCloudKms) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "5d0f1f73cf09a29120e173cf"

		  google_cloud_kms = {
				enabled                 = %s
				service_account_key     = "%s"
				key_version_resource_id = "%s"
			}
		}
	`, cast.ToString(*google.Enabled), google.ServiceAccountKey, google.KeyVersionResourceID)
}
