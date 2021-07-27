package mongodbatlas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"

	digest "github.com/mongodb-forks/digest"
	matlasClient "go.mongodb.org/atlas/mongodbatlas"
	realmAuth "go.mongodb.org/realm/auth"
	"go.mongodb.org/realm/realm"
)

// Config struct ...
type Config struct {
	PublicKey  string
	PrivateKey string
	BaseURL    string
}

// MongoDBClient client
type MongoDBClient struct {
	Atlas *matlasClient.Client
	Realm *realm.Client
}

// NewClient func...
func (c *Config) NewClient(ctx context.Context) (interface{}, diag.Diagnostics) {
	// setup a transport to handle digest
	transport := digest.NewTransport(c.PublicKey, c.PrivateKey)

	// initialize the client
	client, err := transport.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.Transport = logging.NewTransport("MongoDB Atlas", transport)

	optsAtlas := []matlasClient.ClientOpt{matlasClient.SetUserAgent("terraform-provider-mongodbatlas/" + ProviderVersion)}
	if c.BaseURL != "" {
		optsAtlas = append(optsAtlas, matlasClient.SetBaseURL(c.BaseURL))
	}

	// Initialize the MongoDB Atlas API Client.
	atlasClient, err := matlasClient.New(client, optsAtlas...)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	// Realm
	optsRealm := []realm.ClientOpt{realm.SetUserAgent("terraform-provider-mongodbatlas/" + ProviderVersion)}
	if c.BaseURL != "" {
		optsRealm = append(optsRealm, realm.SetBaseURL(c.BaseURL))
	}
	authConfig := realmAuth.NewConfig(client)
	token, err := authConfig.NewTokenFromCredentials(ctx, c.PublicKey, c.PrivateKey)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	clientRealm := realmAuth.NewClient(realmAuth.BasicTokenSource(token))
	clientRealm.Transport = logging.NewTransport("MongoDB Realm", clientRealm.Transport)

	// Initialize the MongoDB Realm API Client.
	realmClient, err := realm.New(clientRealm, optsRealm...)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	clients := &MongoDBClient{
		Atlas: atlasClient,
		Realm: realmClient,
	}

	return clients, nil
}
