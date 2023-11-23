package todoacc

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
)

const (
	// Provider name for single configuration testing
	ProviderNameMongoDBAtlas = "mongodbatlas"
)

var TestAccProviderV6Factories map[string]func() (tfprotov6.ProviderServer, error)

// TestAccProviderSdkV2 is only being used in tests obtaining client: .Meta().(*MongoDBClient)
// this provider instance has to be passed into mux server factory for its configure method to be invoked
var TestAccProviderSdkV2 *schema.Provider

// TestMongoDBClient is used to configure client required for Framework-based acceptance tests
var TestMongoDBClient any

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
	TestMongoDBClient, _ = cfg.NewClient(context.Background())
}
