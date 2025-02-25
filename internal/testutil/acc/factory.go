package acc

import (
	"context"
	"os"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	adminpreview "github.com/mongodb/atlas-sdk-go/admin"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

const (
	// Provider name for single configuration testing
	ProviderNameMongoDBAtlas = "mongodbatlas"
)

// TestAccProviderV6Factories is used in all tests for ProtoV6ProviderFactories.
var TestAccProviderV6Factories map[string]func() (tfprotov6.ProviderServer, error)

func TestAccProviderV6FactoriesWithProxy(proxyPort *int) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		ProviderNameMongoDBAtlas: func() (tfprotov6.ProviderServer, error) {
			return provider.MuxProviderFactoryForTesting(proxyPort)(), nil
		},
	}
}

// MongoDBClient is used to configure client required for Framework-based acceptance tests.
var MongoDBClient *config.MongoDBClient

func Conn() *matlas.Client {
	return MongoDBClient.Atlas
}

func ConnV2() *admin.APIClient {
	return MongoDBClient.AtlasV2
}

func ConnPreview() *adminpreview.APIClient {
	return MongoDBClient.AtlasPreview
}

func ConnV2UsingProxy(proxyPort *int) *admin.APIClient {
	cfg := config.Config{
		PublicKey:    os.Getenv("MONGODB_ATLAS_PUBLIC_KEY"),
		PrivateKey:   os.Getenv("MONGODB_ATLAS_PRIVATE_KEY"),
		BaseURL:      os.Getenv("MONGODB_ATLAS_BASE_URL"),
		RealmBaseURL: os.Getenv("MONGODB_REALM_BASE_URL"),
		ProxyPort:    proxyPort,
	}
	client, _ := cfg.NewClient(context.Background())
	return client.(*config.MongoDBClient).AtlasV2
}

func ConnV2UsingGov() *admin.APIClient {
	cfg := config.Config{
		PublicKey:  os.Getenv("MONGODB_ATLAS_GOV_PUBLIC_KEY"),
		PrivateKey: os.Getenv("MONGODB_ATLAS_GOV_PRIVATE_KEY"),
		BaseURL:    os.Getenv("MONGODB_ATLAS_GOV_BASE_URL"),
	}
	client, _ := cfg.NewClient(context.Background())
	return client.(*config.MongoDBClient).AtlasV2
}

func init() {
	TestAccProviderV6Factories = map[string]func() (tfprotov6.ProviderServer, error){
		ProviderNameMongoDBAtlas: func() (tfprotov6.ProviderServer, error) {
			return provider.MuxProviderFactory()(), nil
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
