package integration_testing

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

const localPluginPath = "TERRAFORM PLUGINS PATH"

func TestUpgradeNetworkContainerRegionsGCP(t *testing.T) {
	t.Parallel()

	var (
		randInt        = acctest.RandIntRange(0, 255)
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		atlasCIDRBlock = fmt.Sprintf("10.%d.0.0/18", randInt)
		providerName   = "GCP"
		publicKey      = os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
		privateKey     = os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v090/network-container/v082",
		Vars: map[string]interface{}{
			"project_id":       projectID,
			"org_id":           orgID,
			"atlas_cidr_block": atlasCIDRBlock,
			"provider_name":    providerName,
			"public_key":       publicKey,
			"private_key":      privateKey,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	networkContainerID := terraform.Output(t, terraformOptions, "network_container_id")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v090/network-container/v090",
		Vars: map[string]interface{}{
			"project_id":       projectID,
			"org_id":           orgID,
			"atlas_cidr_block": atlasCIDRBlock,
			"provider_name":    providerName,
			"public_key":       publicKey,
			"private_key":      privateKey,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptionsSecond)

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_network_container.test", fmt.Sprintf("%s-%s", projectID, networkContainerID))
	// Run `terraform apply`. Fail the test if there are any errors.
	terraform.Plan(t, terraformOptionsSecond)

}

func TestUpgradeDatabaseUserLDAPAuthType(t *testing.T) {
	t.Parallel()

	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		roleName    = "atlasAdmin"
		username    = "CN=ellen@example.com,OU=users,DC=example,DC=com"
		publicKey   = os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
		privateKey  = os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v090/database-user/v082",
		Vars: map[string]interface{}{
			"project_name": projectName,
			"org_id":       orgID,
			"username":     username,
			"role_name":    roleName,
			"public_key":   publicKey,
			"private_key":  privateKey,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	userName := terraform.Output(t, terraformOptions, "username")
	authDatabaseName := terraform.Output(t, terraformOptions, "auth_database_name")

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/v090/database-user/v090")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
		Vars: map[string]interface{}{
			"project_name": projectName,
			"org_id":       orgID,
			"username":     username,
			"role_name":    roleName,
			"public_key":   publicKey,
			"private_key":  privateKey,
		},
	})

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	//Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_database_user.test", fmt.Sprintf("%s-%s-%s", projectID, userName, authDatabaseName))
	// Run `terraform apply`. Fail the test if there are any errors.
	terraform.Plan(t, terraformOptionsSecond)

}

func testAccGetMongoDBAtlasMajorVersion() string {
	conn, _ := matlas.New(http.DefaultClient, matlas.SetBaseURL(matlas.CloudURL))
	majorVersion, _, _ := conn.DefaultMongoDBMajorVersion.Get(context.Background())

	return majorVersion
}

func TestUpgradeClusterDeprecationEBSVolume(t *testing.T) {
	t.Parallel()

	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = acctest.RandomWithPrefix("test-acc")
		publicKey    = os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
		privateKey   = os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
		majorVersion = testAccGetMongoDBAtlasMajorVersion()
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v090/cluster/v082",
		Vars: map[string]interface{}{
			"project_name":          projectName,
			"org_id":                orgID,
			"cluster_name":          clusterName,
			"public_key":            publicKey,
			"private_key":           privateKey,
			"mongodb_major_version": majorVersion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	clusterNameOutput := terraform.Output(t, terraformOptions, "cluster_name")

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/v090/cluster/v090")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
		Vars: map[string]interface{}{
			"project_name":          projectName,
			"org_id":                orgID,
			"cluster_name":          clusterName,
			"public_key":            publicKey,
			"private_key":           privateKey,
			"mongodb_major_version": majorVersion,
		},
	})

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	//Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_cluster.test", fmt.Sprintf("%s-%s", projectID, clusterNameOutput))
	// Run `terraform apply`. Fail the test if there are any errors.
	terraform.Plan(t, terraformOptionsSecond)

}

func TestUpgradePrivateEndpoint(t *testing.T) {
	t.Parallel()

	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		publicKey   = os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
		privateKey  = os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
		baseURL     = os.Getenv("MONGODB_ATLAS_BASE_URL")
		awsAccess   = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecret   = os.Getenv("AWS_SECRET_ACCESS_KEY")
		awsVPC      = os.Getenv("AWS_VPC_ID")
		awsSubnets  = os.Getenv("AWS_SUBNET_ID")
		awsSG       = os.Getenv("AWS_SECURITY_GROUP_ID")
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v090/private-endpoint/v080",
		Vars: map[string]interface{}{
			"project_name":   projectName,
			"org_id":         orgID,
			"public_key":     publicKey,
			"private_key":    privateKey,
			"base_url":       baseURL,
			"aws_access_key": awsAccess,
			"aws_secret_key": awsSecret,
			"aws_vpc_id":     awsVPC,
			"aws_subnet_ids": awsSubnets,
			"aws_sg_ids":     awsSG,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	terraform.Plan(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	vpcEndpoint := terraform.Output(t, terraformOptions, "vpc_endpoint_id")
	privateEndpoint := terraform.Output(t, terraformOptions, "private_endpoint_id")

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/v090/private-endpoint/v090")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
		Vars: map[string]interface{}{
			"project_name":   projectName,
			"org_id":         orgID,
			"public_key":     publicKey,
			"private_key":    privateKey,
			"base_url":       baseURL,
			"aws_access_key": awsAccess,
			"aws_secret_key": awsSecret,
			"aws_vpc_id":     awsVPC,
			"aws_subnet_ids": awsSubnets,
			"aws_sg_ids":     awsSG,
		},
	})

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init")
	//Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_private_endpoint.test", fmt.Sprintf("%s-%s-%s-%s", projectID, privateEndpoint, "AWS", "us-east-1"))
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "aws_vpc_endpoint.ptfe_service", vpcEndpoint)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_private_endpoint_interface_link.test", fmt.Sprintf("%s-%s-%s", projectID, privateEndpoint, vpcEndpoint))
	// Run `terraform apply`. Fail the test if there are any errors.
	terraform.Plan(t, terraformOptionsSecond)

}

// This func means that the terraform state will be always clean to avoid error about resource already used
func CleanUpState(t *testing.T, path string) string {
	// Root folder where terraform files should be (relative to the test folder)
	rootFolder := ".."
	// Relative path to terraform module being tested from the root folder
	terraformFolderRelativeToRoot := path
	// Copy the terraform folder to a temp folder
	return test_structure.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)
}
