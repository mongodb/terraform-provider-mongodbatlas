package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProvider provider.Provider
var testProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
var testMongoDBClient *MongoDBClient

func init() {
	testAccProvider = New()()

	testProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"mongodbatlas": providerserver.NewProtocol6WithError(New()()),
	}
	config := Config{
		PublicKey:    os.Getenv("MONGODB_ATLAS_PUBLIC_KEY"),
		PrivateKey:   os.Getenv("MONGODB_ATLAS_PRIVATE_KEY"),
		BaseURL:      os.Getenv("MONGODB_ATLAS_BASE_URL"),
		RealmBaseURL: os.Getenv("MONGODB_REALM_BASE_URL"),
	}

	testMongoDBClient, _ = config.NewClient(context.Background())
}

func testAccPreCheckBasicOwnerID(tb testing.TB) {
	testAccPreCheckBasic(tb)
	if os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID") == "" {
		tb.Fatal("`MONGODB_ATLAS_PROJECT_OWNER_ID` must be set ")
	}
}

func testAccPreCheckBasic(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_PUBLIC_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_PRIVATE_KEY") == "" ||
		os.Getenv("MONGODB_ATLAS_ORG_ID") == "" {
		tb.Fatal("`MONGODB_ATLAS_PUBLIC_KEY`, `MONGODB_ATLAS_PRIVATE_KEY`, and `MONGODB_ATLAS_ORG_ID` must be set for acceptance testing")
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

func testCheckTeamsIds(tb testing.TB) {
	if os.Getenv("MONGODB_ATLAS_TEAMS_IDS") == "" {
		tb.Skip("`MONGODB_ATLAS_TEAMS_IDS` must be set for Projects acceptance testing")
	}
}
