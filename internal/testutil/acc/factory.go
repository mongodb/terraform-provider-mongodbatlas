package acc

import (
	"context"
	"os"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

const (
	// Provider name for single configuration testing
	ProviderNameMongoDBAtlas = "mongodbatlas"
)

var TestAccProviderV6Factories map[string]func() (tfprotov6.ProviderServer, error)

// TestAccProviderSdkV2 is only being used in tests obtaining client: .Meta().(*MongoDBClient)
// this provider instance has to be passed into mux server factory for its configure method to be invoked
var TestAccProviderSdkV2 *schema.Provider

// testMongoDBClient is used to configure client required for Framework-based acceptance tests
var testMongoDBClient any

func Conn() *matlas.Client {
	return testMongoDBClient.(*config.MongoDBClient).Atlas
}

func ConnV2() *admin.APIClient {
	return testMongoDBClient.(*config.MongoDBClient).AtlasV2
}

func init() {
	TestAccProviderSdkV2 = provider.NewSdkV2Provider()

	TestAccProviderV6Factories = map[string]func() (tfprotov6.ProviderServer, error){
		ProviderNameMongoDBAtlas: func() (tfprotov6.ProviderServer, error) {
			return provider.MuxedProviderFactoryFn(TestAccProviderSdkV2)(), nil
		},
	}

	cfg := config.Config{
		PublicKey:    os.Getenv("MONGODB_ATLAS_PUBLIC_KEY"),
		PrivateKey:   os.Getenv("MONGODB_ATLAS_PRIVATE_KEY"),
		BaseURL:      os.Getenv("MONGODB_ATLAS_BASE_URL"),
		RealmBaseURL: os.Getenv("MONGODB_REALM_BASE_URL"),
	}
	testMongoDBClient, _ = cfg.NewClient(context.Background())
}
