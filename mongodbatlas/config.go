package mongodbatlas

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	digest "github.com/mongodb-forks/digest"
	matlasClient "go.mongodb.org/atlas/mongodbatlas"
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
func (c *Config) NewClient() interface{} {
	// setup a transport to handle digest
	transport := digest.NewTransport(c.PublicKey, c.PrivateKey)

	// initialize the client
	client, err := transport.Client()
	if err != nil {
		return err
	}

	client.Transport = logging.NewTransport("MongoDB Atlas", transport)

	optsAtlas := []matlasClient.ClientOpt{matlasClient.SetUserAgent("terraform-provider-mongodbatlas/" + ProviderVersion)}
	if c.BaseURL != "" {
		optsAtlas = append(optsAtlas, matlasClient.SetBaseURL(c.BaseURL))
	}

	optsRealm := []realm.ClientOpt{realm.SetUserAgent("terraform-provider-mongodbatlas/" + ProviderVersion)}
	if c.BaseURL != "" {
		optsRealm = append(optsRealm, realm.SetBaseURL(c.BaseURL))
	}

	// Initialize the MongoDB Atlas API Client.
	atlasClient, err := matlasClient.New(client, optsAtlas...)
	if err != nil {
		return err
	}

	// Initialize the MongoDB Realm API Client.
	realmClient, err := realm.New(client, optsRealm...)
	if err != nil {
		return err
	}

	clients := &MongoDBClient{
		Atlas: atlasClient,
		Realm: realmClient,
	}

	return clients
}
