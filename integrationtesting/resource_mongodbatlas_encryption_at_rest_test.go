package integrationtesting

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestTerraformResourceMongoDBAtlasEncryptionAtRestWithRole_basicAWS(t *testing.T) {
	t.Parallel()

	mongoSecrets := GetCredentialsFromEnv()
	awsSecrets := GetAWSCredentialsFromEnv()

	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/atlas-encryptionAtRest-roles",
		Vars: map[string]interface{}{
			"access_key":          awsSecrets.AccessKey,
			"secret_key":          awsSecrets.SecretKey,
			"customer_master_key": awsSecrets.CustomerMasterKey,
			"atlas_region":        awsSecrets.AwsRegion,
			"project_id":          mongoSecrets.ProjectID,
			"public_key":          mongoSecrets.PublicKey,
			"private_key":         mongoSecrets.PrivateKey,
			"base_url":            mongoSecrets.BaseURL,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the IP of the instance
	awsRoleARN := terraform.Output(t, terraformOptions, "aws_iam_role_arn")
	cpaRoleID := terraform.Output(t, terraformOptions, "cpa_role_id")

	fmt.Printf("awsRoleARN : %s", awsRoleARN)
	fmt.Printf("cpaRoleID : %s", cpaRoleID)

	terraformOptionsUpdated := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/atlas-encryptionAtRest-roles",
		Vars: map[string]interface{}{
			"access_key":          awsSecrets.AccessKey,
			"secret_key":          awsSecrets.SecretKey,
			"customer_master_key": awsSecrets.CustomerMasterKey,
			"atlas_region":        awsSecrets.AwsRegion,
			"project_id":          mongoSecrets.ProjectID,
			"public_key":          mongoSecrets.PublicKey,
			"private_key":         mongoSecrets.PrivateKey,
			"base_url":            mongoSecrets.BaseURL,
			"aws_iam_role_arn":    awsRoleARN,
		},
	})

	terraform.Apply(t, terraformOptionsUpdated)

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/atlas-encryptionAtRest-roles/second_step",
		Vars: map[string]interface{}{
			"customer_master_key": awsSecrets.CustomerMasterKey,
			"atlas_region":        awsSecrets.AwsRegion,
			"project_id":          mongoSecrets.ProjectID,
			"public_key":          mongoSecrets.PublicKey,
			"private_key":         mongoSecrets.PrivateKey,
			"base_url":            mongoSecrets.BaseURL,
			"cpa_role_id":         cpaRoleID,
		},
	})
	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptionsSecond)

	// Run `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptionsSecond)
}

func TestResourceEncryptionAtRestAws(t *testing.T) {
	t.Parallel()
	mongoSecrets := GetCredentialsFromEnv()
	awsSecrets := GetAWSCredentialsFromEnv()

	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		publicKey   = mongoSecrets.PublicKey
		privateKey  = mongoSecrets.PrivateKey

		awsAccess   = awsSecrets.AccessKey
		awsSecret   = awsSecrets.SecretKey
		awsCustomer = awsSecrets.CustomerMasterKey
	)

	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/encryption-at-rest/aws/v101",
		Vars: map[string]interface{}{
			"project_name":        projectName,
			"org_id":              orgID,
			"public_key":          publicKey,
			"private_key":         privateKey,
			"access_key":          awsAccess,
			"secret_key":          awsSecret,
			"customer_master_key": awsCustomer,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	if localPluginPath != "" {
		terraform.RunTerraformCommand(t, terraformOptions, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	} else {
		terraform.Init(t, terraformOptions)
	}

	terraform.Apply(t, terraformOptions)

	terraform.Plan(t, terraformOptions)
}

func TestResourceEncryptionAtRestAzure(t *testing.T) {
	t.Parallel()
	mongoSecrets := GetCredentialsFromEnv()

	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		publicKey   = mongoSecrets.PublicKey
		privateKey  = mongoSecrets.PrivateKey

		azureClientID      = os.Getenv("AZURE_CLIENT_ID")
		azureSuscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
		azureResourceName  = os.Getenv("AZURE_RESOURCE_GROUP_NAME")
		azureKeyVault      = os.Getenv("AZURE_KEY_VAULT_NAME")
		azureKeyIdentifier = os.Getenv("AZURE_KEY_IDENTIFIER")
		azureSecret        = os.Getenv("AZURE_SECRET")
		azureTenantID      = os.Getenv("AZURE_TENANT_ID")
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/encryption-at-rest/azure/v101",
		Vars: map[string]interface{}{
			"project_name":        projectName,
			"org_id":              orgID,
			"public_key":          publicKey,
			"private_key":         privateKey,
			"client_id":           azureClientID,
			"subscription_id":     azureSuscriptionID,
			"resource_group_name": azureResourceName,
			"key_vault_name":      azureKeyVault,
			"key_identifier":      azureKeyIdentifier,
			"client_secret":       azureSecret,
			"tenant_id":           azureTenantID,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.

	if localPluginPath != "" {
		terraform.RunTerraformCommand(t, terraformOptions, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	} else {
		terraform.Init(t, terraformOptions)
	}

	terraform.Apply(t, terraformOptions)

	terraform.Plan(t, terraformOptions)
}

func TestResourceEncryptionAtRestGCP(t *testing.T) {
	t.Parallel()
	mongoSecrets := GetCredentialsFromEnv()

	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		publicKey   = mongoSecrets.PublicKey
		privateKey  = mongoSecrets.PrivateKey

		gcpServiceAcc = os.Getenv("GCP_SERVICE_ACCOUNT_KEY")
		gcpVersion    = os.Getenv("GCP_KEY_VERSION_RESOURCE_ID")
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/encryption-at-rest/gcp/v101",
		Vars: map[string]interface{}{
			"project_name":                projectName,
			"org_id":                      orgID,
			"public_key":                  publicKey,
			"private_key":                 privateKey,
			"service_account_key":         gcpServiceAcc,
			"gcp_key_version_resource_id": gcpVersion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	if localPluginPath != "" {
		terraform.RunTerraformCommand(t, terraformOptions, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	} else {
		terraform.Init(t, terraformOptions)
	}

	terraform.Apply(t, terraformOptions)

	terraform.Plan(t, terraformOptions)
}
