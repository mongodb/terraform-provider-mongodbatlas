package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mwielbut/pointy"
)

const (
	initialConfigEncryptionRestRoleAWS = `
provider "aws" {
	region     = lower(replace("%[1]s", "_", "-"))
	access_key = "%[2]s"
	secret_key = "%[3]s"
}

%[7]s

resource "mongodbatlas_cloud_provider_access" "test" {
	project_id = "%[4]s"
	provider_name = "AWS"
	%[8]s
		
}

resource "aws_iam_role_policy" "test_policy" {
  name = "%[5]s"
  role = aws_iam_role.test_role.id

  policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
		"Action": "*",
		"Resource": "*"
      }
    ]
  }
  EOF
}

resource "aws_iam_role" "test_role" {
 name = "%[6]s"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "${mongodbatlas_cloud_provider_access.test.atlas_aws_account_arn}"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "${mongodbatlas_cloud_provider_access.test.atlas_assumed_role_external_id}"
        }
      }
    }
  ]
}
EOF

}

%[9]s

`
	configEncryptionRest = `
resource "mongodbatlas_encryption_at_rest" "test" {
	project_id = "%s"

	aws_kms_config {
		enabled                = %t
		customer_master_key_id = "%s"
		region                 = "%s"
		role_id = mongodbatlas_cloud_provider_access.test.role_id
	}
}`
	dataAWSARNConfig = `
data "aws_iam_role" "test" {
  name = "%s"
}

`
)

func TestAccAdvRSEncryptionAtRest_basicAWS(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		awsKms = matlas.AwsKms{
			Enabled:             pointy.Bool(true),
			AccessKeyID:         os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey:     os.Getenv("AWS_SECRET_ACCESS_KEY"),
			CustomerMasterKeyID: os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID"),
			Region:              os.Getenv("AWS_REGION"),
			RoleID:              os.Getenv("AWS_ROLE_ID"),
		}

		awsKmsUpdated = matlas.AwsKms{
			Enabled:             pointy.Bool(true),
			AccessKeyID:         os.Getenv("AWS_ACCESS_KEY_ID_UPDATED"),
			SecretAccessKey:     os.Getenv("AWS_SECRET_ACCESS_KEY_UPDATED"),
			CustomerMasterKeyID: os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID_UPDATED"),
			Region:              os.Getenv("AWS_REGION_UPDATED"),
			RoleID:              os.Getenv("AWS_ROLE_ID"),
		}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); testCheckAwsEnv(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(projectID, &awsKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.region", awsKms.Region),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.role_id", awsKms.RoleID),
				),
			},
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKms(projectID, &awsKmsUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.region", awsKmsUpdated.Region),
					resource.TestCheckResourceAttr(resourceName, "aws_kms_config.0.role_id", awsKmsUpdated.RoleID),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasEncryptionAtRestImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"google_cloud_kms_config", "azure_key_vault_config"},
			},
		},
	})
}

func TestAccAdvRSEncryptionAtRest_basicAzure(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")

		azureKeyVault = matlas.AzureKeyVault{
			Enabled:           pointy.Bool(true),
			ClientID:          os.Getenv("AZURE_CLIENT_ID"),
			AzureEnvironment:  "AZURE",
			SubscriptionID:    os.Getenv("AZURE_SUBSCRIPTION_ID"),
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
			SubscriptionID:    os.Getenv("AZURE_SUBSCRIPTION_ID"),
			ResourceGroupName: os.Getenv("AZURE_RESOURCE_GROUP_NAME_UPDATED"),
			KeyVaultName:      os.Getenv("AZURE_KEY_VAULT_NAME_UPDATED"),
			KeyIdentifier:     os.Getenv("AZURE_KEY_IDENTIFIER_UPDATED"),
			Secret:            os.Getenv("AZURE_SECRET_UPDATED"),
			TenantID:          os.Getenv("AZURE_TENANT_ID"),
		}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); testCheckEncryptionAtRestEnvAzure(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAzureKeyVault(projectID, &azureKeyVault),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.azure_environment", azureKeyVault.AzureEnvironment),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.resource_group_name", azureKeyVault.ResourceGroupName),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.key_vault_name", azureKeyVault.KeyVaultName),
				),
			},
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAzureKeyVault(projectID, &azureKeyVaultUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.azure_environment", azureKeyVaultUpdated.AzureEnvironment),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.resource_group_name", azureKeyVaultUpdated.ResourceGroupName),
					resource.TestCheckResourceAttr(resourceName, "azure_key_vault_config.0.key_vault_name", azureKeyVaultUpdated.KeyVaultName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEncryptionAtRestImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// "azure_key_vault_config.0.secret" is a sensitive value not returned by the API
				ImportStateVerifyIgnore: []string{"google_cloud_kms_config", "aws_kms_config", "azure_key_vault_config.0.secret"},
			},
		},
	})
}

func TestAccAdvRSEncryptionAtRest_basicGCP(t *testing.T) {
	SkipTestExtCred(t)
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); testAccPreCheckGPCEnv(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigGoogleCloudKms(projectID, &googleCloudKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms_config.0.enabled", "true"),
				),
			},
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigGoogleCloudKms(projectID, &googleCloudKmsUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_kms_config.0.enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasEncryptionAtRestImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// "google_cloud_kms_config.0.service_account_key" is a sensitive value not returned by the API
				ImportStateVerifyIgnore: []string{"aws_kms_config", "azure_key_vault_config", "google_cloud_kms_config.0.service_account_key"},
			},
		},
	})
}

func TestAccAdvRSEncryptionAtRestWithRole_basicAWS(t *testing.T) {
	SkipTest(t) // For now it will skipped because of aws errors reasons, already made another test using terratest.
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_encryption_at_rest.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		accessKeyID  = os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey    = os.Getenv("AWS_SECRET_ACCESS_KEY")
		policyName   = acctest.RandomWithPrefix("test-aws-policy")
		roleName     = acctest.RandomWithPrefix("test-aws-role")

		awsKms = matlas.AwsKms{
			Enabled:             pointy.Bool(true),
			CustomerMasterKeyID: os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID"),
			Region:              os.Getenv("AWS_REGION"),
		}
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t); testCheckAwsEnv(t) },
		ExternalProviders: map[string]resource.ExternalProvider{
			"aws": {
				VersionConstraint: "5.1.0",
				Source:            "hashicorp/aws",
			},
		},
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasEncryptionAtRestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(awsKms.Region, accessKeyID, secretKey, projectID, policyName, roleName, false, &awsKms),
			},
			{
				Config: testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(awsKms.Region, accessKeyID, secretKey, projectID, policyName, roleName, true, &awsKms),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasEncryptionAtRestImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"google_cloud_kms_config", "azure_key_vault_config"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasEncryptionAtRestExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testMongoDBClient.(*MongoDBClient).Atlas

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
	conn := testMongoDBClient.(*MongoDBClient).Atlas

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

		  aws_kms_config {
				enabled                = %t
				customer_master_key_id = "%s"
				region                 = "%s"
				role_id              = "%s"
			}
		}
	`, projectID, *aws.Enabled, aws.CustomerMasterKeyID, aws.Region, aws.RoleID)
}

func testAccMongoDBAtlasEncryptionAtRestConfigAzureKeyVault(projectID string, azure *matlas.AzureKeyVault) string {
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
			}
		}
	`, projectID, *azure.Enabled, azure.ClientID, azure.AzureEnvironment, azure.SubscriptionID, azure.ResourceGroupName,
		azure.KeyVaultName, azure.KeyIdentifier, azure.Secret, azure.TenantID)
}

func testAccMongoDBAtlasEncryptionAtRestConfigGoogleCloudKms(projectID string, google *matlas.GoogleCloudKms) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_encryption_at_rest" "test" {
			project_id = "%s"

		  google_cloud_kms_config {
				enabled                 = %t
				service_account_key     = "%s"
				key_version_resource_id = "%s"
			}
		}
	`, projectID, *google.Enabled, google.ServiceAccountKey, google.KeyVersionResourceID)
}

func testAccMongoDBAtlasEncryptionAtRestConfigAwsKmsWithRole(region, awsAccesKey, awsSecretKey, projectID, policyName, awsRoleName string, isUpdate bool, aws *matlas.AwsKms) string {
	config := fmt.Sprintf(initialConfigEncryptionRestRoleAWS, region, awsAccesKey, awsSecretKey, projectID, policyName, awsRoleName, "", "", "")
	if isUpdate {
		configEncrypt := fmt.Sprintf(configEncryptionRest, projectID, *aws.Enabled, aws.CustomerMasterKeyID, aws.Region)
		dataAWSARN := fmt.Sprintf(dataAWSARNConfig, awsRoleName)
		dataARN := `iam_assumed_role_arn = data.aws_iam_role.test.arn`
		config = fmt.Sprintf(initialConfigEncryptionRestRoleAWS, region, awsAccesKey, awsSecretKey, projectID, policyName, awsRoleName, dataAWSARN, dataARN, configEncrypt)
	}
	return config
}

func testAccCheckMongoDBAtlasEncryptionAtRestImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
