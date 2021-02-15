package integration_testing

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformResourceMongoDBAtlasDataLake_basicAWS(t *testing.T) {
	SkipTestExtCred(t)
	t.Parallel()

	var (
		accessKey    = os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey    = os.Getenv("AWS_SECRET_ACCESS_KEY")
		awsRegion    = os.Getenv("AWS_REGION")
		publicKey    = os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
		privateKey   = os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-project")
		dataLakeName = acctest.RandomWithPrefix("test-acc-data-lake")
		testS3Bucket = os.Getenv("AWS_S3_BUCKET")
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/atlas-dataLake-roles",
		Vars: map[string]interface{}{
			"access_key":   accessKey,
			"secret_key":   secretKey,
			"atlas_region": awsRegion,
			"public_key":   publicKey,
			"private_key":  privateKey,
			"project_name": projectName,
			"org_id":       orgID,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// To run a local plugin with --plugin-dir
	//terraform.RunTerraformCommand(t, terraformOptions, terraform.FormatArgs(terraformOptions, "init",
	//	fmt.Sprintf("--plugin-dir=%s", "PLUGIN PATH"))...)
	terraform.Apply(t, terraformOptions)
	// Run `terraform output` to get the IP of the instance
	awsRoleARN := terraform.Output(t, terraformOptions, "aws_iam_role_arn")
	cpaRoleID := terraform.Output(t, terraformOptions, "cpa_role_id")
	projectID := terraform.Output(t, terraformOptions, "project_id")

	terraformOptionsUpdated := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/atlas-dataLake-roles",
		Vars: map[string]interface{}{
			"access_key":       accessKey,
			"secret_key":       secretKey,
			"atlas_region":     awsRegion,
			"org_id":           orgID,
			"project_name":     projectName,
			"public_key":       publicKey,
			"private_key":      privateKey,
			"aws_iam_role_arn": awsRoleARN,
		},
	})

	terraform.Apply(t, terraformOptionsUpdated)

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/atlas-dataLake-roles/second_step",
		Vars: map[string]interface{}{
			"project_id":     projectID,
			"public_key":     publicKey,
			"private_key":    privateKey,
			"cpa_role_id":    cpaRoleID,
			"data_lake_name": dataLakeName,
			"test_s3_bucket": testS3Bucket,
		},
	})
	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptionsSecond)

	// Run `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)
	// To run a local plugin with --plugin-dir
	//terraform.RunTerraformCommand(t, terraformOptions, terraform.FormatArgs(terraformOptions, "init",
	//	fmt.Sprintf("--plugin-dir=%s", "PLUGIN PATH"))...)
	terraform.Apply(t, terraformOptionsSecond)
}
