package mongodbatlas

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google"
	"github.com/terraform-providers/terraform-provider-aws/aws"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	// Provider name for single configuration testing
	ProviderNameMongoDBAtlas = "mongodbatlas"
	ProviderNameAws          = "aws"
	ProviderNameGoogle       = "google"
)

var testAccProviders map[string]*schema.Provider

var testAccProviderFactories map[string]func() (*schema.Provider, error)

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		ProviderNameMongoDBAtlas: testAccProvider,
		ProviderNameAws:          aws.Provider(),
		ProviderNameGoogle:       google.Provider(),
	}

	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		ProviderNameMongoDBAtlas: func() (*schema.Provider, error) { return testAccProvider, nil },
		ProviderNameAws:          func() (*schema.Provider, error) { return aws.Provider(), nil },
		ProviderNameGoogle:       func() (*schema.Provider, error) { return google.Provider(), nil },
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("MONGODB_ATLAS_PUBLIC_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PRIVATE_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PROJECT_ID") == "" ||
		os.Getenv("MONGODB_ATLAS_ORG_ID") == "" {
		t.Fatal("`MONGODB_ATLAS_PUBLIC_KEY`, `MONGODB_ATLAS_PRIVATE_KEY`, `MONGODB_ATLAS_PROJECT_ID` and `MONGODB_ATLAS_ORG_ID` must be set for acceptance testing")
	}
}

func testAccPreCheckGov(t *testing.T) {
	if os.Getenv("MONGODB_ATLAS_PUBLIC_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PRIVATE_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PROJECT_ID_GOV") == "" ||
		os.Getenv("MONGODB_ATLAS_ORG_ID_GOV") == "" {
		t.Skip()
	}
}

func testAccPreCheckBetaFeatures(t *testing.T) {
	enableFeatures, _ := strconv.ParseBool(os.Getenv("MONGODB_ATLAS_ENABLE_BETA"))
	if !enableFeatures {
		t.Fatal("`MONGODB_ATLAS_ENABLE_BETA` must be set to `true` in order to run beta features")
	}
}

func testAccPreCheckGPCEnv(t *testing.T) {
	if os.Getenv("GCP_SERVICE_ACCOUNT_KEY") == "" || os.Getenv("GCP_KEY_VERSION_RESOURCE_ID") == "" {
		t.Fatal("`GCP_SERVICE_ACCOUNT_KEY` and `GCP_KEY_VERSION_RESOURCE_ID` must be set for acceptance testing")
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
		os.Getenv("AZURE_SUBSCRIPTION_ID") == "" ||
		os.Getenv("AZURE_VNET_NAME") == "" ||
		os.Getenv("AZURE_RESOURCE_GROUP_NAME") == "" {
		t.Fatal("`AZURE_DIRECTORY_ID`, `AZURE_SUBSCRIPTION_ID`, `AZURE_VNET_NAME` and `AZURE_RESOURCE_GROUP_NAME` must be set for  network peering acceptance testing")
	}
}

func checkEncryptionAtRestEnvAzure(t *testing.T) {
	if os.Getenv("AZURE_CLIENT_ID") == "" ||
		os.Getenv("AZURE_CLIENT_ID_UPDATED") == "" ||
		os.Getenv("AZURE_SUBSCRIPTION_ID") == "" ||
		os.Getenv("AZURE_RESOURCE_GROUP_NAME") == "" ||
		os.Getenv("AZURE_RESOURCE_GROUP_NAME_UPDATED") == "" ||
		os.Getenv("AZURE_SECRET") == "" ||
		os.Getenv("AZURE_KEY_VAULT_NAME") == "" ||
		os.Getenv("AZURE_KEY_VAULT_NAME_UPDATED") == "" ||
		os.Getenv("AZURE_KEY_IDENTIFIER") == "" ||
		os.Getenv("AZURE_KEY_IDENTIFIER_UPDATED") == "" ||
		os.Getenv("AZURE_TENANT_ID") == "" {
		t.Fatal(`'AZURE_CLIENT_ID','AZURE_CLIENT_ID_UPDATED', 'AZURE_SUBSCRIPTION_ID',
		'AZURE_RESOURCE_GROUP_NAME','AZURE_RESOURCE_GROUP_NAME_UPDATED', 'AZURE_SECRET',
		'AZURE_SECRET_UPDATED', 'AZURE_KEY_VAULT_NAME', 'AZURE_KEY_IDENTIFIER', 'AZURE_KEY_VAULT_NAME_UPDATED',
		'AZURE_KEY_IDENTIFIER_UPDATED', and 'AZURE_TENANT_ID' must be set for Encryption At Rest acceptance testing`)
	}
}

func checkPeeringEnvGCP(t *testing.T) {
	if os.Getenv("GCP_PROJECT_ID") == "" ||
		os.Getenv("GCP_CLUSTER_REGION_NAME") == "" ||
		os.Getenv("GCP_REGION_NAME") == "" ||
		os.Getenv("GOOGLE_CLOUD_KEYFILE_JSON") == "" {
		t.Fatal("`GCP_PROJECT_ID`,`GOOGLE_CLOUD_KEYFILE_JSON`, `GCP_CLUSTER_REGION_NAME`, `and GCP_REGION_NAME` must be set for network peering acceptance testing")
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

func SkipTestExtCred(t *testing.T) {
	if strings.EqualFold(os.Getenv("SKIP_TEST_EXTERNAL_CREDENTIALS"), "true") {
		t.Skip()
	}
}

func checkTeamsIds(t *testing.T) {
	if os.Getenv("MONGODB_ATLAS_TEAMS_IDS") == "" {
		t.Skip("`MONGODB_ATLAS_TEAMS_IDS` must be set for Projects acceptance testing")
	}
}

func SkipTest(t *testing.T) {
	if strings.EqualFold(os.Getenv("SKIP_TEST"), "true") {
		t.Skip()
	}
}

func checkLDAP(t *testing.T) {
	if os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME") == "" ||
		os.Getenv("MONGODB_ATLAS_LDAP_USERNAME") == "" ||
		os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD") == "" ||
		os.Getenv("MONGODB_ATLAS_LDAP_PORT") == "" {
		t.Fatal("`MONGODB_ATLAS_LDAP_HOSTNAME`, `MONGODB_ATLAS_LDAP_USERNAME`, `MONGODB_ATLAS_LDAP_PASSWORD` and `MONGODB_ATLAS_LDAP_PORT` must be set for ldap configuration/verify acceptance testing")
	}
}

func checkFederatedSettings(t *testing.T) {
	if os.Getenv("MONGODB_ATLAS_FEDERATED_PROJECT_ID") == "" ||
		os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID") == "" ||
		os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID") == "" {
		t.Fatal("`MONGODB_ATLAS_FEDERATED_PROJECT_ID`, `MONGODB_ATLAS_FEDERATED_ORG_ID` and `MONGODB_ATLAS_FEDERATION_SETTINGS_ID` must be set for federated settings/verify acceptance testing")
	}
}
