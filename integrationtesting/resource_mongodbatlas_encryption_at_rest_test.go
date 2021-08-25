//go:build integration
// +build integration

package integrationtesting

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
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
