package acceptancetests

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	sdkv2schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	// Provider name for single configuration testing
	ProviderNameMongoDBAtlas = "mongodbatlas"
)

var TestAccProviderV6Factories map[string]func() (tfprotov6.ProviderServer, error)

// only being used in tests obtaining client: .Meta().(*MongoDBClient)
// this provider instance has to be passed into mux server factory for its configure method to be invoked
var testAccProviderSdkV2 *sdkv2schema.Provider

// testMongoDBClient is used to configure client required for Framework-based acceptance tests
var testMongoDBClient interface{}

func init() {
	testAccProviderSdkV2 = mongodbatlas.NewSdkV2Provider()
	TestAccProviderV6Factories = map[string]func() (tfprotov6.ProviderServer, error){
		ProviderNameMongoDBAtlas: func() (tfprotov6.ProviderServer, error) {
			return MuxedProviderFactoryWithProvider(testAccProviderSdkV2)(), nil
		},
	}

	config := mongodbatlas.Config{
		PublicKey:    os.Getenv("MONGODB_ATLAS_PUBLIC_KEY"),
		PrivateKey:   os.Getenv("MONGODB_ATLAS_PRIVATE_KEY"),
		BaseURL:      os.Getenv("MONGODB_ATLAS_BASE_URL"),
		RealmBaseURL: os.Getenv("MONGODB_REALM_BASE_URL"),
	}
	testMongoDBClient, _ = config.NewClient(context.Background())
}

func MuxedProviderFactory() func() tfprotov6.ProviderServer {
	return MuxedProviderFactoryWithProvider(mongodbatlas.NewSdkV2Provider())
}

// MuxedProviderFactoryWithProvider creates mux provider using existing sdk v2 provider passed as parameter and creating new instance of framework provider.
// Used in testing where existing sdk v2 provider has to be used.
func MuxedProviderFactoryWithProvider(sdkV2Provider *sdkv2schema.Provider) func() tfprotov6.ProviderServer {
	fwProvider := mongodbatlas.NewFrameworkProvider()

	ctx := context.Background()
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(ctx, sdkV2Provider.GRPCProvider)
	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},
		providerserver.NewProtocol6(fwProvider),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}
	return muxServer.ProviderServer
}

func TestSdkV2Provider(t *testing.T) {
	if err := mongodbatlas.NewSdkV2Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_PUBLIC_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PRIVATE_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PROJECT_ID") == "" ||
		os.Getenv("MONGODB_ATLAS_ORG_ID") == "" {
		tb.Fatal("`MONGODB_ATLAS_PUBLIC_KEY`, `MONGODB_ATLAS_PRIVATE_KEY`, `MONGODB_ATLAS_PROJECT_ID` and `MONGODB_ATLAS_ORG_ID` must be set for acceptance testing")
	}
}

func testAccPreCheckBasic(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_PUBLIC_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PRIVATE_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_ORG_ID") == "" {
		tb.Fatal("`MONGODB_ATLAS_PUBLIC_KEY`, `MONGODB_ATLAS_PRIVATE_KEY`, and `MONGODB_ATLAS_ORG_ID` must be set for acceptance testing")
	}
}

func testAccPreCheckCloudProviderAccessAzure(tb testing.TB) {
	testAccPreCheckBasic(tb)
	if os.Getenv("AZURE_ATLAS_APP_ID") == "" ||
		os.Getenv("AZURE_SERVICE_PRINCIPAL_ID") == "" ||
		os.Getenv("AZURE_TENANT_ID") == "" {
		tb.Fatal("`AZURE_ATLAS_APP_ID`, `AZURE_SERVICE_PRINCIPAL_ID`, and `AZURE_TENANT_ID` must be set for acceptance testing")
	}
}

func testAccPreCheckBasicOwnerID(tb testing.TB) {
	testAccPreCheckBasic(tb)
	if os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID") == "" {
		tb.Fatal("`MONGODB_ATLAS_PROJECT_OWNER_ID` must be set ")
	}
}

func testAccPreCheckGov(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_PUBLIC_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PRIVATE_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PROJECT_ID_GOV") == "" ||
		os.Getenv("MONGODB_ATLAS_ORG_ID_GOV") == "" {
		tb.Skip()
	}
}

func testAccPreCheckGPCEnv(tb testing.TB) {
	if os.Getenv("GCP_SERVICE_ACCOUNT_KEY") == "" || os.Getenv("GCP_KEY_VERSION_RESOURCE_ID") == "" {
		tb.Fatal("`GCP_SERVICE_ACCOUNT_KEY` and `GCP_KEY_VERSION_RESOURCE_ID` must be set for acceptance testing")
	}
}

func testCheckPeeringEnvAWS(tb testing.TB) {
	if os.Getenv("AWS_ACCOUNT_ID") == "" ||
		os.Getenv("AWS_VPC_ID") == "" ||
		os.Getenv("AWS_VPC_CIDR_BLOCK") == "" ||
		os.Getenv("AWS_REGION") == "" {
		tb.Fatal("`AWS_ACCOUNT_ID`, `AWS_VPC_ID`, `AWS_VPC_CIDR_BLOCK` and `AWS_VPC_ID` must be set for  network peering acceptance testing")
	}
}

func testCheckPeeringEnvAzure(tb testing.TB) {
	if os.Getenv("AZURE_DIRECTORY_ID") == "" ||
		os.Getenv("AZURE_SUBSCRIPTION_ID") == "" ||
		os.Getenv("AZURE_VNET_NAME") == "" ||
		os.Getenv("AZURE_RESOURCE_GROUP_NAME") == "" {
		tb.Fatal("`AZURE_DIRECTORY_ID`, `AZURE_SUBSCRIPTION_ID`, `AZURE_VNET_NAME` and `AZURE_RESOURCE_GROUP_NAME` must be set for  network peering acceptance testing")
	}
}

func testCheckEncryptionAtRestEnvAzure(tb testing.TB) {
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
		tb.Fatal(`'AZURE_CLIENT_ID','AZURE_CLIENT_ID_UPDATED', 'AZURE_SUBSCRIPTION_ID',
		'AZURE_RESOURCE_GROUP_NAME','AZURE_RESOURCE_GROUP_NAME_UPDATED', 'AZURE_SECRET',
		'AZURE_SECRET_UPDATED', 'AZURE_KEY_VAULT_NAME', 'AZURE_KEY_IDENTIFIER', 'AZURE_KEY_VAULT_NAME_UPDATED',
		'AZURE_KEY_IDENTIFIER_UPDATED', and 'AZURE_TENANT_ID' must be set for Encryption At Rest acceptance testing`)
	}
}

func testCheckPeeringEnvGCP(tb testing.TB) {
	if os.Getenv("GCP_PROJECT_ID") == "" ||
		os.Getenv("GCP_CLUSTER_REGION_NAME") == "" ||
		os.Getenv("GCP_REGION_NAME") == "" ||
		os.Getenv("GOOGLE_CLOUD_KEYFILE_JSON") == "" {
		tb.Fatal("`GCP_PROJECT_ID`,`GOOGLE_CLOUD_KEYFILE_JSON`, `GCP_CLUSTER_REGION_NAME`, `and GCP_REGION_NAME` must be set for network peering acceptance testing")
	}
}

func testCheckAwsEnv(tb testing.TB) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" ||
		os.Getenv("AWS_SECRET_ACCESS_KEY") == "" ||
		os.Getenv("AWS_CUSTOMER_MASTER_KEY_ID") == "" {
		tb.Fatal("`AWS_ACCESS_KEY_ID`, `AWS_VPC_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_CUSTOMER_MASTER_KEY_ID` must be set for acceptance testing")
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

func SkipTestExtCred(tb testing.TB) {
	if strings.EqualFold(os.Getenv("SKIP_TEST_EXTERNAL_CREDENTIALS"), "true") {
		tb.Skip()
	}
}

func testCheckDataLakePipelineRun(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_DATA_LAKE_PIPELINE_RUN_ID") == "" {
		tb.Skip("`MONGODB_ATLAS_DATA_LAKE_PIPELINE_RUN_ID` must be set for Projects acceptance testing")
	}
	testCheckDataLakePipelineRuns(tb)
}

func testCheckDataLakePipelineRuns(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_DATA_LAKE_PIPELINE_NAME") == "" {
		tb.Skip("`MONGODB_ATLAS_DATA_LAKE_PIPELINE_NAME` must be set for Projects acceptance testing")
	}
}

func testCheckTeamsIds(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_TEAMS_IDS") == "" {
		tb.Skip("`MONGODB_ATLAS_TEAMS_IDS` must be set for Projects acceptance testing")
	}
}

func SkipTest(tb testing.TB) {
	if strings.EqualFold(os.Getenv("SKIP_TEST"), "true") {
		tb.Skip()
	}
}

// SkipTestForCI is added to tests that cannot run as part of a CI
func SkipTestForCI(tb testing.TB) {
	if strings.EqualFold(os.Getenv("CI"), "true") {
		tb.Skip()
	}
}

func testCheckLDAP(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME") == "" ||
		os.Getenv("MONGODB_ATLAS_LDAP_USERNAME") == "" ||
		os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD") == "" ||
		os.Getenv("MONGODB_ATLAS_LDAP_PORT") == "" {
		tb.Fatal("`MONGODB_ATLAS_LDAP_HOSTNAME`, `MONGODB_ATLAS_LDAP_USERNAME`, `MONGODB_ATLAS_LDAP_PASSWORD` and `MONGODB_ATLAS_LDAP_PORT` must be set for ldap configuration/verify acceptance testing")
	}
}

func testCheckFederatedSettings(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_FEDERATED_PROJECT_ID") == "" ||
		os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID") == "" ||
		os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID") == "" {
		tb.Fatal("`MONGODB_ATLAS_FEDERATED_PROJECT_ID`, `MONGODB_ATLAS_FEDERATED_ORG_ID` and `MONGODB_ATLAS_FEDERATION_SETTINGS_ID` must be set for federated settings/verify acceptance testing")
	}
}

func testCheckPrivateEndpointServiceDataFederationOnlineArchiveRun(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_PRIVATE_ENDPOINT_ID") == "" {
		tb.Skip("`MONGODB_ATLAS_PRIVATE_ENDPOINT_ID` must be set for Private Endpoint Service Data Federation and Online Archive acceptance testing")
	}
}
func decodeStateID(stateID string) map[string]string {
	decode := func(d string) string {
		decodedString, err := base64.StdEncoding.DecodeString(d)
		if err != nil {
			log.Printf("[WARN] error decoding state ID: %s", err)
		}

		return string(decodedString)
	}
	decodedValues := make(map[string]string)
	encodedValues := strings.Split(stateID, "-")

	for _, value := range encodedValues {
		keyValue := strings.Split(value, ":")
		if len(keyValue) > 1 {
			decodedValues[decode(keyValue[0])] = decode(keyValue[1])
		}
	}

	return decodedValues
}

func encodeStateID(values map[string]string) string {
	encode := func(e string) string { return base64.StdEncoding.EncodeToString([]byte(e)) }
	encodedValues := make([]string, 0)

	// sort to make sure the same encoding is returned in case of same input
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		encodedValues = append(encodedValues, fmt.Sprintf("%s:%s", encode(key), encode(values[key])))
	}

	return strings.Join(encodedValues, "-")
}

func removeLabel(list []matlas.Label, item matlas.Label) []matlas.Label {
	var pos int

	for _, v := range list {
		if reflect.DeepEqual(v, item) {
			list = append(list[:pos], list[pos+1:]...)

			if pos > 0 {
				pos--
			}

			continue
		}
		pos++
	}

	return list
}
