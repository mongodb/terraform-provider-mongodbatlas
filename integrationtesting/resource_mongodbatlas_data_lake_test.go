package integrationtesting

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestTerraformResourceMongoDBAtlasDataLake_basicAWS(t *testing.T) {
	t.Parallel()

	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-project")
		dataLakeName = acctest.RandomWithPrefix("test-acc-data-lake")
		testS3Bucket = os.Getenv("AWS_S3_BUCKET")
		pluginPath   = os.Getenv("TERRATEST_PLUGIN_PATH")
	)
	mongoSecrets := GetCredentialsFromEnv()
	awsSecrets := GetAWSCredentialsFromEnv()
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/atlas-dataLake-roles",
		Vars: map[string]interface{}{
			"public_key":     mongoSecrets.PublicKey,
			"private_key":    mongoSecrets.PrivateKey,
			"access_key":     awsSecrets.AccessKey,
			"secret_key":     awsSecrets.SecretKey,
			"aws_region":     awsSecrets.AwsRegion,
			"project_name":   projectName,
			"org_id":         orgID,
			"data_lake_name": dataLakeName,
			"test_s3_bucket": testS3Bucket,
			"base_url":       mongoSecrets.BaseURL,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	if pluginPath != "" {
		terraform.RunTerraformCommand(t, terraformOptions, "init", fmt.Sprintf("--plugin-dir=%s", pluginPath))
		terraform.Apply(t, terraformOptions)
	} else {
		terraform.InitAndApply(t, terraformOptions)
	}

	terraform.Plan(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	roleID := terraform.Output(t, terraformOptions, "role_id")
	roleName := terraform.Output(t, terraformOptions, "role_name")
	policyName := terraform.Output(t, terraformOptions, "policy_name")
	lakeName := terraform.Output(t, terraformOptions, "data_lake_name")
	s3Bucket := terraform.Output(t, terraformOptions, "s3_bucket")

	tempTestFolder := CleanUpState(t, "examples/atlas-dataLake-roles/import")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
		Vars: map[string]interface{}{
			"public_key":  mongoSecrets.PublicKey,
			"private_key": mongoSecrets.PrivateKey,
			"access_key":  awsSecrets.AccessKey,
			"secret_key":  awsSecrets.SecretKey,
			"base_url":    mongoSecrets.BaseURL,
		},
	})

	if pluginPath != "" {
		terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", pluginPath))
	} else {
		terraform.Init(t, terraformOptionsSecond)
	}

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "aws_iam_role_policy.test_policy", fmt.Sprintf("%s:%s", roleName, policyName))
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "aws_iam_role.test_role", roleName)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_cloud_provider_access_setup.setup_only", fmt.Sprintf("%s-%s-%s", projectID, "AWS", roleID))
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_data_lake.test", fmt.Sprintf("%s--%s--%s", projectID, lakeName, s3Bucket))

	terraform.Plan(t, terraformOptionsSecond)
}
