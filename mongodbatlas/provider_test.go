package mongodbatlas

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"mongodbatlas": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestGetPluginVersion(t *testing.T) {
	version := getPluginVersion()
	if version == "" {
		t.Fatal("version must not be empty")
	}
	t.Logf("\nVersion: %s", version)
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("MONGODB_ATLAS_PUBLIC_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PRIVATE_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PROJECT_ID") == "" ||
		os.Getenv("MONGODB_ATLAS_ORG_ID") == "" {
		t.Fatal("`MONGODB_ATLAS_PUBLIC_KEY`, `MONGODB_ATLAS_PRIVATE_KEY`, `MONGODB_ATLAS_PROJECT_ID` and `MONGODB_ATLAS_ORG_ID` must be set for acceptance testing")
	}
}

func checkPeeringEnvAWS(t *testing.T) {
	if os.Getenv("AWS_ACCOUNT_ID") == "" ||
		os.Getenv("AWS_VPC_ID") == "" ||
		os.Getenv("AWS_VPC_CIDR_BLOCK") == "" ||
		os.Getenv("AWS_REGION") == "" {
		t.Fatal("`AWS_ACCOUNT_ID`, `AWS_VPC_ID`, `AWS_VPC_CIDR_BLOCK` and `AWS_VPC_ID` must be set for  network peering acceptance testing")
	}
}

func checkPeeringEnvAzure(t *testing.T) {
	if os.Getenv("AZURE_DIRECTORY_ID") == "" ||
		os.Getenv("AZURE_SUBCRIPTION_ID") == "" ||
		os.Getenv("AZURE_VNET_NAME") == "" ||
		os.Getenv("AZURE_RESOURSE_GROUP_NAME") == "" {
		t.Fatal("`AZURE_DIRECTORY_ID`, `AZURE_SUBCRIPTION_ID`, `AZURE_VNET_NAME` and `AZURE_RESOURSE_GROUP_NAME` must be set for  network peering acceptance testing")
	}
}

func checkPeeringEnvGCP(t *testing.T) {
	if os.Getenv("GCP_PROJECT_ID") == "" {
		t.Fatal("`GCP_PROJECT_ID` must be set for network peering acceptance testing")
	}
}

func checkAwsEnv(t *testing.T) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" ||
		os.Getenv("AWS_SECRET_ACCESS_KEY") == "" ||
		os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID") == "" ||
		os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID_UPDATED") == "" ||
		os.Getenv("AWS_ACCESS_KEY_ID_UPDATED") == "" ||
		os.Getenv("AWS_SECRET_ACCESS_KEY_UPDATED") == "" {
		t.Fatal("`AWS_ACCESS_KEY_ID`, `AWS_VPC_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_CUSTOMER_MASTER_KEY_ID` must be set for acceptance testing")
	}
}
