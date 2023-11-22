package acc

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
)

const (
	// Provider name for single configuration testing
	ProviderNameMongoDBAtlas = "mongodbatlas"
)

var testAccProviderV6Factories map[string]func() (tfprotov6.ProviderServer, error)

// only being used in tests obtaining client: .Meta().(*MongoDBClient)
// this provider instance has to be passed into mux server factory for its configure method to be invoked
var testAccProviderSdkV2 *schema.Provider

// testMongoDBClient is used to configure client required for Framework-based acceptance tests
var testMongoDBClient any

func init() {
	testAccProviderSdkV2 = NewSdkV2Provider()

	testAccProviderV6Factories = map[string]func() (tfprotov6.ProviderServer, error){
		ProviderNameMongoDBAtlas: func() (tfprotov6.ProviderServer, error) {
			return provider.MuxedProviderFactory(testAccProviderSdkV2)(), nil
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
