// +build integration

package integration_testing

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

const (
	defaultTerratestFilesCPA = "../examples/atlas-cloud-provider-access/aws/"
)

func TestTerraformResourceMongoDBAtlasCloudProviderAccess_basicAWS(t *testing.T) {
	t.Parallel()

	mongoSecrets := GetCredentialsFromEnv()
	awsSecrets := GetAWSCredentialsFromEnv()

	testFiles := os.Getenv("TERRATEST_CLOUD_PROVIDER_ACCESS_AWS")
	if testFiles == "" {
		testFiles = defaultTerratestFilesCPA
	}

	terraformOptions := &terraform.Options{
		TerraformDir: testFiles,
		Vars: map[string]interface{}{
			"project_id":                 mongoSecrets.ProjectID,
			"cloud_provider_access_name": "AWS",
			"public_key":                 mongoSecrets.PublicKey,
			"private_key":                mongoSecrets.PrivateKey,
			"base_url":                   mongoSecrets.BaseURL,
			"access_key":                 awsSecrets.AccessKey,
			"secret_key":                 awsSecrets.SecretKey,
			"aws_region":                 awsSecrets.AwsRegion,
		},
	}

	terraformTest := terraform.WithDefaultRetryableErrors(t, terraformOptions)

	defer terraform.Destroy(t, terraformTest)
	terraform.InitAndApply(t, terraformTest)
}
