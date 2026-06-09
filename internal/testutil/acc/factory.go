package acc

import (
	"fmt"
	"log"
	"maps"
	"os"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	adminpreview "github.com/mongodb/atlas-sdk-go/admin"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/provider"
	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"
	"go.mongodb.org/atlas-sdk/v20250312020/admin"
)

const (
	// Provider name for single configuration testing
	ProviderNameMongoDBAtlas = "mongodbatlas"
)

// TestAccProviderV6Factories is used in all tests for ProtoV6ProviderFactories.
var TestAccProviderV6Factories map[string]func() (tfprotov6.ProviderServer, error)

// sharedClient holds the acceptance-test client; access it through MongoDBClient().
var sharedClient *config.MongoDBClient

// clientInitErr records any error from initializing sharedClient, surfaced by MongoDBClient().
var clientInitErr error

// MongoDBClient returns the shared acceptance-test client, panicking with the init error
// (e.g. a failed service-account token fetch) instead of a nil pointer dereference.
func MongoDBClient() *config.MongoDBClient {
	if sharedClient == nil {
		panic(fmt.Sprintf("acceptance test Atlas client was not initialized: %v", clientInitErr))
	}
	return sharedClient
}

func Conn() *matlas.Client {
	return MongoDBClient().Atlas
}

func ConnV2() *admin.APIClient {
	return MongoDBClient().AtlasV2
}

func ConnPreview() *adminpreview.APIClient {
	return MongoDBClient().AtlasPreview
}

func ConnV220241113() *admin20241113.APIClient {
	return MongoDBClient().AtlasV220241113
}

func ConnV2UsingGov() *admin.APIClient {
	c := &config.Credentials{
		PublicKey:  os.Getenv("MONGODB_ATLAS_GOV_PUBLIC_KEY"),
		PrivateKey: os.Getenv("MONGODB_ATLAS_GOV_PRIVATE_KEY"),
		BaseURL:    os.Getenv("MONGODB_ATLAS_GOV_BASE_URL"),
	}
	client, err := config.NewClient(c, "")
	if err != nil {
		panic(fmt.Sprintf("acceptance test Atlas (gov) client could not be created: %v", err))
	}
	return client.AtlasV2
}

// ProtoV6FactoriesWithEcho returns provider factories that include both the
// mongodbatlas provider and the echo provider (for testing ephemeral values).
func ProtoV6FactoriesWithEcho() map[string]func() (tfprotov6.ProviderServer, error) {
	factories := make(map[string]func() (tfprotov6.ProviderServer, error))
	maps.Copy(factories, TestAccProviderV6Factories)
	factories["echo"] = echoprovider.NewProviderServer()
	return factories
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
		AccessToken:  os.Getenv("MONGODB_ATLAS_ACCESS_TOKEN"),
		BaseURL:      os.Getenv("MONGODB_ATLAS_BASE_URL"),
		RealmBaseURL: os.Getenv("MONGODB_REALM_BASE_URL"),
	}
	sharedClient, clientInitErr = config.NewClient(c, "")
	if clientInitErr != nil {
		log.Printf("[ERROR] failed to initialize Atlas client for acceptance tests: %v", clientInitErr)
	}
}
