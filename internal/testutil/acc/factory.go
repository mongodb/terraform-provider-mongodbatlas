package acc

import (
	"context"
	"os"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

const (
	// Provider name for single configuration testing
	ProviderNameMongoDBAtlas = "mongodbatlas"
)

// TestAccProviderV6Factories is used in all tests for ProtoV6ProviderFactories.
var TestAccProviderV6Factories map[string]func() (tfprotov6.ProviderServer, error)

// MongoDBClient is used to configure client required for Framework-based acceptance tests.
var MongoDBClient *config.MongoDBClient

func Conn() *matlas.Client {
	return MongoDBClient.Atlas
}

func ConnV2() *admin.APIClient {
	return MongoDBClient.AtlasV2
}

func init() {
	TestAccProviderV6Factories = map[string]func() (tfprotov6.ProviderServer, error){
		ProviderNameMongoDBAtlas: func() (tfprotov6.ProviderServer, error) {
			return provider.MuxedProviderFactory()(), nil
		},
	}
	cfg := config.Config{
		PublicKey:    os.Getenv("MONGODB_ATLAS_PUBLIC_KEY"),
		PrivateKey:   os.Getenv("MONGODB_ATLAS_PRIVATE_KEY"),
		BaseURL:      os.Getenv("MONGODB_ATLAS_BASE_URL"),
		RealmBaseURL: os.Getenv("MONGODB_REALM_BASE_URL"),
	}
	client, _ := cfg.NewClient(context.Background())
	MongoDBClient = client.(*config.MongoDBClient)
}
