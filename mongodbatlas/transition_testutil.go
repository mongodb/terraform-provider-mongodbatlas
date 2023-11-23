package mongodbatlas

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func testAccPreCheckBasic(tb testing.TB) {
	acc.PreCheckBasic(tb)
}

func testCheckAwsEnv(tb testing.TB) {
	acc.PreCheckAwsEnv(tb)
}

func testAccPreCheck(tb testing.TB) {
	acc.PreCheck(tb)
}

func SkipTest(tb testing.TB) {
	acc.SkipTest(tb)
}

func SkipTestForCI(tb testing.TB) {
	acc.SkipTestForCI(tb)
}

func SkipTestExtCred(tb testing.TB) {
	acc.SkipTestExtCred(tb)
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

func testAccPreCheckSearchIndex(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_PUBLIC_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PRIVATE_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_ORG_ID") == "" ||
		os.Getenv("MONGODB_ATLAS_PROJECT_ID") == "" {
		tb.Fatal("`MONGODB_ATLAS_PUBLIC_KEY`, `MONGODB_ATLAS_PRIVATE_KEY`,  `MONGODB_ATLAS_ORG_ID`, and `MONGODB_ATLAS_PROJECT_ID` must be set for acceptance testing")
	}
}

// TODO INITIALIZE OR LINK TO INTERNAL ************
// TODO INITIALIZE OR LINK TO INTERNAL ************
var testAccProviderV6Factories map[string]func() (tfprotov6.ProviderServer, error)
var testAccProviderSdkV2 *schema.Provider
var testMongoDBClient any
