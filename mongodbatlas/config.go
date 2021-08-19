package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/mongodb/terraform-provider-mongodbatlas/version"
	"github.com/spf13/cast"
	realmAuth "go.mongodb.org/realm/auth"

	"github.com/mongodb-forks/digest"
	matlasClient "go.mongodb.org/atlas/mongodbatlas"
	"go.mongodb.org/realm/realm"
)

// MongoDBClient client
type MongoDBClient struct {
	PublicKey  *string
	PrivateKey *string
	BaseURL    *string
	Atlas      *matlasClient.Client
}

// NewClient func...
func (c *MongoDBClient) NewClient(ctx context.Context) (interface{}, diag.Diagnostics) {
	// setup a transport to handle digest
	transport := digest.NewTransport(cast.ToString(c.PublicKey), cast.ToString(c.PrivateKey))

	// initialize the client
	client, err := transport.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.Transport = logging.NewTransport("MongoDB Atlas", transport)

	optsAtlas := []matlasClient.ClientOpt{matlasClient.SetUserAgent("terraform-provider-mongodbatlas/" + version.ProviderVersion)}
	if cast.ToString(c.BaseURL) == "" {
		optsAtlas = append(optsAtlas, matlasClient.SetBaseURL(cast.ToString(c.BaseURL)))
	}

	// Initialize the MongoDB Atlas API Client.
	atlasClient, err := matlasClient.New(client, optsAtlas...)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	clients := &MongoDBClient{
		Atlas:      atlasClient,
		PublicKey:  c.PublicKey,
		PrivateKey: c.PrivateKey,
		BaseURL:    c.BaseURL,
	}

	return clients, nil
}

func (c *MongoDBClient) GetRealmClient(ctx context.Context) (*realm.Client, error) {
	// Realm
	if cast.ToString(c.PublicKey) == "" && cast.ToString(c.PrivateKey) == "" {
		return nil, fmt.Errorf("please set `public_key` and `private_key` in order to use the realm client")
	}

	optsRealm := []realm.ClientOpt{realm.SetUserAgent("terraform-provider-mongodbatlas/" + version.ProviderVersion)}
	if cast.ToString(c.BaseURL) == "" {
		optsRealm = append(optsRealm, realm.SetBaseURL(cast.ToString(c.BaseURL)))
	}
	authConfig := realmAuth.NewConfig(nil)
	token, err := authConfig.NewTokenFromCredentials(ctx, cast.ToString(c.PublicKey), cast.ToString(c.PrivateKey))
	if err != nil {
		return nil, err
	}

	clientRealm := realmAuth.NewClient(realmAuth.BasicTokenSource(token))
	clientRealm.Transport = logging.NewTransport("MongoDB Realm", clientRealm.Transport)

	// Initialize the MongoDB Realm API Client.
	realmClient, err := realm.New(clientRealm, optsRealm...)
	if err != nil {
		return nil, err
	}

	return realmClient, nil
}
