package acc

import (
	"os"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	adminpreview "github.com/mongodb/atlas-sdk-go/admin"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
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

func ConnPreview() *adminpreview.APIClient {
	return MongoDBClient.AtlasPreview
}

func ConnV220241113() *admin20241113.APIClient {
	return MongoDBClient.AtlasV220241113
}

func ConnV2UsingGov() *admin.APIClient {
	c := &config.Credentials{
		PublicKey:  os.Getenv("MONGODB_ATLAS_GOV_PUBLIC_KEY"),
		PrivateKey: os.Getenv("MONGODB_ATLAS_GOV_PRIVATE_KEY"),
		BaseURL:    os.Getenv("MONGODB_ATLAS_GOV_BASE_URL"),
	}
	client, _ := config.NewClient(c, "")
	return client.AtlasV2
}

func init() {
	TestAccProviderV6Factories = map[string]func() (tfprotov6.ProviderServer, error){
		ProviderNameMongoDBAtlas: func() (tfprotov6.ProviderServer, error) {
			return provider.MuxProviderFactory()(), nil
		},
	}
	c := &config.Credentials{
		PublicKey:    os.Getenv("MONGODB_ATLAS_PUBLIC_KEY"),
		PrivateKey:   os.Getenv("MONGODB_ATLAS_PRIVATE_KEY"),
		ClientID:     os.Getenv("MONGODB_ATLAS_CLIENT_ID"),
		ClientSecret: os.Getenv("MONGODB_ATLAS_CLIENT_SECRET"),
		BaseURL:      os.Getenv("MONGODB_ATLAS_BASE_URL"),
		RealmBaseURL: os.Getenv("MONGODB_REALM_BASE_URL"),
	}
	MongoDBClient, _ = config.NewClient(c, "")
}
