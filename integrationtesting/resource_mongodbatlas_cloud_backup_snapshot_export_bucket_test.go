package integrationtesting

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestResourceCloudBackupSnapshotExportBucket(t *testing.T) {
	t.Parallel()
	mongoSecrets := GetCredentialsFromEnv()
	awsSecrets := GetAWSCredentialsFromEnv()

	var (
		publicKey  = mongoSecrets.PublicKey
		privateKey = mongoSecrets.PrivateKey
		projectID  = mongoSecrets.ProjectID
		awsAccess  = awsSecrets.AccessKey
		awsSecret  = awsSecrets.SecretKey
	)

	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/atlas-provider-snapshot-export-bucket",
		Vars: map[string]interface{}{
			"project_id":  projectID,
			"public_key":  publicKey,
			"private_key": privateKey,
			"access_key":  awsAccess,
			"secret_key":  awsSecret,
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
