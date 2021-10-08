package integrationtesting

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

var (
	localPluginPath = os.Getenv("TERRATEST_PLUGIN_PATH")
)

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
	// Remove states
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
	// Remove states
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
	// Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_private_endpoint.test", fmt.Sprintf("%s-%s-%s-%s", projectID, privateEndpoint, "AWS", "us-east-1"))
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "aws_vpc_endpoint.ptfe_service", vpcEndpoint)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_private_endpoint_interface_link.test", fmt.Sprintf("%s-%s-%s", projectID, privateEndpoint, vpcEndpoint))
	// Run `terraform apply`. Fail the test if there are any errors.
	terraform.Plan(t, terraformOptionsSecond)
}

func TestUpgradeProjectIPWhitelistDeprecation(t *testing.T) {
	t.Parallel()

	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		publicKey   = os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
		privateKey  = os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
		ipAddress   = fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
		comment     = fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v100/ip-whitelist-accestList/v091",
		Vars: map[string]interface{}{
			"project_name": projectName,
			"org_id":       orgID,
			"ip_address":   ipAddress,
			"public_key":   publicKey,
			"private_key":  privateKey,
			"comment":      comment,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	entry := terraform.Output(t, terraformOptions, "entry")

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/v100/ip-whitelist-accestList/v100")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
		Vars: map[string]interface{}{
			"project_name": projectName,
			"org_id":       orgID,
			"ip_address":   ipAddress,
			"public_key":   publicKey,
			"private_key":  privateKey,
			"comment":      comment,
		},
	})

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	// Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project_ip_access_list.test", fmt.Sprintf("%s-%s", projectID, entry))
	// Run `terraform apply`. Fail the test if there are any errors.
	terraform.Plan(t, terraformOptionsSecond)
}

func TestUpgradeDesignIDState(t *testing.T) {
	t.Parallel()
	mongoSecrets := GetCredentialsFromEnv()
	awsSecrets := GetAWSCredentialsFromEnv()

	var (
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName     = acctest.RandomWithPrefix("test-acc")
		clusterName     = acctest.RandomWithPrefix("test-acc")
		description     = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays = "1"
		publicKey       = mongoSecrets.PublicKey
		privateKey      = mongoSecrets.PrivateKey
		awsAccess       = awsSecrets.AccessKey
		awsSecret       = awsSecrets.SecretKey
		awsVPC          = os.Getenv("AWS_VPC_ID")
		awsSubnets      = os.Getenv("AWS_SUBNET_ID")
		awsSG           = os.Getenv("AWS_SECURITY_GROUP_ID")
		vpcCIDRBlock    = os.Getenv("AWS_VPC_CIDR_BLOCK")
		awsAccountID    = os.Getenv("AWS_ACCOUNT_ID")
		awsRegion       = awsSecrets.AwsRegion
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptionsProject := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v100/design-id-reference/project",
		Vars: map[string]interface{}{
			"project_name": projectName,
			"org_id":       orgID,
			"public_key":   publicKey,
			"private_key":  privateKey,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptionsProject)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.RunTerraformCommand(t, terraformOptionsProject, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	terraform.Apply(t, terraformOptionsProject)

	terraform.Plan(t, terraformOptionsProject)

	// Alert Configuration

	terraformOptionsAlertConfiguration := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v100/design-id-reference/alert-configuration",
		Vars: map[string]interface{}{
			"project_name": projectName,
		},
	})
	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptionsAlertConfiguration)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.RunTerraformCommand(t, terraformOptionsAlertConfiguration, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	terraform.Apply(t, terraformOptionsAlertConfiguration)

	terraform.Plan(t, terraformOptionsAlertConfiguration)

	// Network container/peering

	terraformOptionsNetwork := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v100/design-id-reference/network",
		Vars: map[string]interface{}{
			"project_name":           projectName,
			"region_name":            awsRegion,
			"route_table_cidr_block": vpcCIDRBlock,
			"vpc_id":                 awsVPC,
			"aws_account_id":         awsAccountID,
		},
	})
	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptionsNetwork)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.RunTerraformCommand(t, terraformOptionsNetwork, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	terraform.Apply(t, terraformOptionsNetwork)

	terraform.Plan(t, terraformOptionsNetwork)

	// PrivateLink

	terraformOptionsPrivateLink := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v100/design-id-reference/privatelink",
		Vars: map[string]interface{}{
			"project_name":   projectName,
			"aws_access_key": awsAccess,
			"aws_secret_key": awsSecret,
			"aws_vpc_id":     awsVPC,
			"aws_subnet_ids": awsSubnets,
			"aws_sg_ids":     awsSG,
		},
	})
	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptionsPrivateLink)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.RunTerraformCommand(t, terraformOptionsPrivateLink, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	terraform.Apply(t, terraformOptionsPrivateLink)

	terraform.Plan(t, terraformOptionsPrivateLink)

	// Snapshot Restore

	terraformOptionsSnapshotRestore := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v100/design-id-reference/snapshot-restore",
		Vars: map[string]interface{}{
			"project_name":      projectName,
			"cluster_name":      clusterName,
			"description":       description,
			"retention_in_days": retentionInDays,
		},
	})
	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptionsSnapshotRestore)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.RunTerraformCommand(t, terraformOptionsSnapshotRestore, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	terraform.Apply(t, terraformOptionsSnapshotRestore)

	terraform.Plan(t, terraformOptionsSnapshotRestore)
}

func TestUpgradePrivateLinkEndpointDeprecation(t *testing.T) {
	t.Parallel()

	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		publicKey   = os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
		privateKey  = os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
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
		TerraformDir: "../examples/test-upgrade/v100/privatelink-endpoint/v091",
		Vars: map[string]interface{}{
			"project_name":   projectName,
			"org_id":         orgID,
			"public_key":     publicKey,
			"private_key":    privateKey,
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

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/v100/privatelink-endpoint/v100")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
		Vars: map[string]interface{}{
			"project_name":   projectName,
			"org_id":         orgID,
			"public_key":     publicKey,
			"private_key":    privateKey,
			"aws_access_key": awsAccess,
			"aws_secret_key": awsSecret,
			"aws_vpc_id":     awsVPC,
			"aws_subnet_ids": awsSubnets,
			"aws_sg_ids":     awsSG,
		},
	})

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	// Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_privatelink_endpoint.test", fmt.Sprintf("%s-%s-%s-%s", projectID, privateEndpoint, "AWS", "us-east-1"))
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "aws_vpc_endpoint.ptfe_service", vpcEndpoint)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_privatelink_endpoint_service.test", fmt.Sprintf("%s--%s--%s--%s", projectID, privateEndpoint, vpcEndpoint, "AWS"))
	// Run `terraform apply`. Fail the test if there are any errors.
	terraform.Plan(t, terraformOptionsSecond)
}

func TestUpgradeCloudBackupPolicies(t *testing.T) {
	t.Parallel()

	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		publicKey   = os.Getenv("MONGODB_ATLAS_PUBLIC_KEY")
		privateKey  = os.Getenv("MONGODB_ATLAS_PRIVATE_KEY")
		clusterName = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v100/cloud-backup-policies/v091",
		Vars: map[string]interface{}{
			"project_name": projectName,
			"org_id":       orgID,
			"public_key":   publicKey,
			"private_key":  privateKey,
			"cluster_name": clusterName,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	cluster := terraform.Output(t, terraformOptions, "cluster_name")

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/v100/cloud-backup-policies/v100")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
		Vars: map[string]interface{}{
			"project_name": projectName,
			"org_id":       orgID,
			"public_key":   publicKey,
			"private_key":  privateKey,
			"cluster_name": clusterName,
		},
	})

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	// Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.project_test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_cluster.cluster_test", fmt.Sprintf("%s-%s", projectID, cluster))
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_cloud_backup_schedule.test", fmt.Sprintf("%s-%s", projectID, cluster))
	// Run `terraform apply`. Fail the test if there are any errors.

	terraform.Plan(t, terraformOptionsSecond)
}

func TestUpgradeEncryptionAtRestAws(t *testing.T) {
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
		TerraformDir: "../examples/test-upgrade/encryption-at-rest/aws/v091",
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
	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	roleName := terraform.Output(t, terraformOptions, "role_name")
	policyName := terraform.Output(t, terraformOptions, "role_policy_name")

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/encryption-at-rest/aws/v101")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
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

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	// Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_encryption_at_rest.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "aws_iam_role.test_role", roleName)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "aws_iam_role_policy.test_policy", fmt.Sprintf("%s:%s", roleName, policyName))
	// Run `terraform apply`. Fail the test if there are any errors.

	terraform.Plan(t, terraformOptionsSecond)
}

func TestUpgradeEncryptionAtRestAzure(t *testing.T) {
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
		TerraformDir: "../examples/test-upgrade/encryption-at-rest/azure/v091",
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
	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/encryption-at-rest/azure/v101")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
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

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	// Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_encryption_at_rest.test", projectID)
	// Run `terraform apply`. Fail the test if there are any errors.

	terraform.Plan(t, terraformOptionsSecond)
}

func TestUpgradeEncryptionAtRestGCP(t *testing.T) {
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
		TerraformDir: "../examples/test-upgrade/encryption-at-rest/gcp/v091",
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
	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/encryption-at-rest/gcp/v101")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
		Vars: map[string]interface{}{
			"project_name":                projectName,
			"org_id":                      orgID,
			"public_key":                  publicKey,
			"private_key":                 privateKey,
			"service_account_key":         gcpServiceAcc,
			"gcp_key_version_resource_id": gcpVersion,
		},
	})

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	// Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_encryption_at_rest.test", projectID)
	// Run `terraform apply`. Fail the test if there are any errors.

	terraform.Plan(t, terraformOptionsSecond)
}

func TestUpgradeCloudBackupSnapshot(t *testing.T) {
	t.Parallel()

	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
		clusterName = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/test-upgrade/v110/cloud-backup-snapshot/v102",
		Vars: map[string]interface{}{
			"project_name": projectName,
			"org_id":       orgID,
			"cluster_name": clusterName,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	projectID := terraform.Output(t, terraformOptions, "project_id")
	cluster := terraform.Output(t, terraformOptions, "cluster_name")
	snapshotID := terraform.Output(t, terraformOptions, "snapshot_id")
	restoreJobID := terraform.Output(t, terraformOptions, "snapshot_restore_job_id")

	tempTestFolder := CleanUpState(t, "examples/test-upgrade/v110/cloud-backup-snapshot/v110")

	terraformOptionsSecond := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tempTestFolder,
		Vars: map[string]interface{}{
			"project_name": projectName,
			"org_id":       orgID,
			"cluster_name": clusterName,
		},
	})

	terraform.RunTerraformCommand(t, terraformOptionsSecond, "init", fmt.Sprintf("--plugin-dir=%s", localPluginPath))
	// Remove states
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_project.project_test", projectID)
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_cluster.cluster_test", fmt.Sprintf("%s-%s", projectID, cluster))
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_cloud_backup_snapshot.test", fmt.Sprintf("%s-%s-%s", projectID, cluster, snapshotID))
	terraform.RunTerraformCommand(t, terraformOptionsSecond, "import", "mongodbatlas_cloud_backup_snapshot_restore_job.test", fmt.Sprintf("%s-%s-%s", projectID, cluster, restoreJobID))
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
