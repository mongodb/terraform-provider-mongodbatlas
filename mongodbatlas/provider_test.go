package mongodbatlas

import (
	"os"
	"testing"

	"github.com/go-test/deep"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
	"github.com/terraform-providers/terraform-provider-aws/aws"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"mongodbatlas": testAccProvider,
		"aws":          aws.Provider(),
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

func checkEncryptionAtRestEnvAzure(t *testing.T) {
	if os.Getenv("AZURE_CLIENT_ID") == "" ||
		os.Getenv("AZURE_CLIENT_ID_UPDATED") == "" ||
		os.Getenv("AZURE_SUBCRIPTION_ID") == "" ||
		os.Getenv("AZURE_RESOURSE_GROUP_NAME") == "" ||
		os.Getenv("AZURE_RESOURSE_GROUP_NAME_UPDATED") == "" ||
		os.Getenv("AZURE_SECRET") == "" ||
		os.Getenv("AZURE_KEY_VAULT_NAME") == "" ||
		os.Getenv("AZURE_KEY_VAULT_NAME_UPDATED") == "" ||
		os.Getenv("AZURE_KEY_IDENTIFIER") == "" ||
		os.Getenv("AZURE_KEY_IDENTIFIER_UPDATED") == "" ||
		os.Getenv("AZURE_TENANT_ID") == "" {
		t.Fatal(`'AZURE_CLIENT_ID','AZURE_CLIENT_ID_UPDATED', 'AZURE_SUBCRIPTION_ID',
		'AZURE_RESOURSE_GROUP_NAME','AZURE_RESOURSE_GROUP_NAME_UPDATED', 'AZURE_SECRET',
		'AZURE_SECRET_UPDATED', 'AZURE_KEY_VAULT_NAME', 'AZURE_KEY_IDENTIFIER', 'AZURE_KEY_VAULT_NAME_UPDATED',
		'AZURE_KEY_IDENTIFIER_UPDATED', and 'AZURE_TENANT_ID' must be set for Encryption At Rest acceptance testing`)
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
		os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID") == "" {
		t.Fatal("`AWS_ACCESS_KEY_ID`, `AWS_VPC_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_CUSTOMER_MASTER_KEY_ID` must be set for acceptance testing")
	}
}

func TestEncodeDecodeID(t *testing.T) {
	expected := map[string]string{
		"project_id":   "5cf5a45a9ccf6400e60981b6",
		"cluster_name": "test-acc-q4y272zo9y",
		"snapshot_id":  "5e42e646553855a5aee40138",
	}

	got := decodeStateID(encodeStateID(expected))

	if diff := deep.Equal(expected, got); diff != nil {
		t.Fatalf("Bad testEncodeDecodeID return \n got = %#v\nwant = %#v \ndiff = %#v", got, expected, diff)
	}
}

func TestDecodeID(t *testing.T) {
	expected := "Y2x1c3Rlcl9uYW1l:dGVzdC1hY2MtcTR5Mjcyem85eQ==-c25hcHNob3RfaWQ=:NWU0MmU2NDY1NTM4NTVhNWFlZTQwMTM4-cHJvamVjdF9pZA==:NWNmNWE0NWE5Y2NmNjQwMGU2MDk4MWI2"
	expected2 := "c25hcHNob3RfaWQ=:NWU0MmU2NDY1NTM4NTVhNWFlZTQwMTM4-cHJvamVjdF9pZA==:NWNmNWE0NWE5Y2NmNjQwMGU2MDk4MWI2-Y2x1c3Rlcl9uYW1l:dGVzdC1hY2MtcTR5Mjcyem85eQ=="

	got := decodeStateID(expected)
	got2 := decodeStateID(expected2)

	if diff := deep.Equal(got, got2); diff != nil {
		t.Fatalf("Bad TestDecodeID return \n got = %#v\nwant = %#v \ndiff = %#v", got, got2, diff)
	}
}

func TestRemoveLabel(t *testing.T) {
	toRemove := matlas.Label{Key: "To Remove", Value: "To remove value"}

	expected := []matlas.Label{
		{Key: "Name", Value: "Test"},
		{Key: "Version", Value: "1.0"},
		{Key: "Type", Value: "testing"},
	}

	labels := []matlas.Label{
		{Key: "Name", Value: "Test"},
		{Key: "Version", Value: "1.0"},
		{Key: "To Remove", Value: "To remove value"},
		{Key: "Type", Value: "testing"},
	}

	got := removeLabel(labels, toRemove)

	if diff := deep.Equal(expected, got); diff != nil {
		t.Fatalf("Bad removeLabel return \n got = %#v\nwant = %#v \ndiff = %#v", got, expected, diff)
	}
}
