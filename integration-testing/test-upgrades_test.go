package integration_testing

import (
	"fmt"
	"os"
	"testing"

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
		TerraformDir: "../examples/test-upgrade/network-container/old",
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
		TerraformDir: "../examples/test-upgrade/network-container/updated",
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
		TerraformDir: "../examples/test-upgrade/database-user/old",
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

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/database-user/updated")

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

func TestUpgradeClusterDeprecationEBSVolume(t *testing.T) {
	t.Parallel()

	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		clusterName = acctest.RandomWithPrefix("test-acc")
		publicKey   = os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
		privateKey  = os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/cluster/old",
		Vars: map[string]interface{}{
			"project_name": projectName,
			"org_id":       orgID,
			"cluster_name": clusterName,
			"public_key":   publicKey,
			"private_key":  privateKey,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	clusterNameOutput := terraform.Output(t, terraformOptions, "cluster_name")

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/cluster/updated")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
		Vars: map[string]interface{}{
			"project_name": projectName,
			"org_id":       orgID,
			"cluster_name": clusterName,
			"public_key":   publicKey,
			"private_key":  privateKey,
		},
	})

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	//Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_cluster.test", fmt.Sprintf("%s-%s", projectID, clusterNameOutput))
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
