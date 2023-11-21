package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/mongodb-forks/digest"
	"github.com/mongodb/terraform-provider-mongodbatlas/version"
	"github.com/spf13/cast"
	atlasSDK "go.mongodb.org/atlas-sdk/v20231115002/admin"
	matlasClient "go.mongodb.org/atlas/mongodbatlas"
	realmAuth "go.mongodb.org/realm/auth"
	"go.mongodb.org/realm/realm"
)

const ToolName = "terraform-provider-mongodbatlas"

var userAgent = fmt.Sprintf("%s/%s", ToolName, version.ProviderVersion)

// Config contains the configurations needed to use SDKs
type Config struct {
	AssumeRole   *AssumeRole
	PublicKey    string
	PrivateKey   string
	BaseURL      string
	RealmBaseURL string
}

// MongoDBClient contains the mongodbatlas clients and configurations
type MongoDBClient struct {
	Atlas   *matlasClient.Client
	AtlasV2 *atlasSDK.APIClient
	Config  *Config
}

// NewClient func...
func (c *Config) NewClient(ctx context.Context) (any, error) {
	// setup a transport to handle digest
	transport := digest.NewTransport(cast.ToString(c.PublicKey), cast.ToString(c.PrivateKey))

	// initialize the client
	client, err := transport.Client()
	if err != nil {
		return nil, err
	}

	client.Transport = logging.NewTransport("MongoDB Atlas", transport)

	optsAtlas := []matlasClient.ClientOpt{matlasClient.SetUserAgent(userAgent)}
	if c.BaseURL != "" {
		optsAtlas = append(optsAtlas, matlasClient.SetBaseURL(c.BaseURL))
	}

	// Initialize the MongoDB Atlas API Client.
	atlasClient, err := matlasClient.New(client, optsAtlas...)
	if err != nil {
		return nil, err
	}

	sdkV2Client, err := c.newSDKV2Client(client)
	if err != nil {
		return nil, err
	}

	clients := &MongoDBClient{
		Atlas:   atlasClient,
		AtlasV2: sdkV2Client,
		Config:  c,
	}

	return clients, nil
}

func (c *Config) newSDKV2Client(client *http.Client) (*atlasSDK.APIClient, error) {
	opts := []atlasSDK.ClientModifier{
		atlasSDK.UseHTTPClient(client),
		atlasSDK.UseUserAgent(userAgent),
		atlasSDK.UseBaseURL(c.BaseURL),
		atlasSDK.UseDebug(false)}

	// Initialize the MongoDB Versioned Atlas Client.
	sdkv2, err := atlasSDK.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return sdkv2, nil
}

func (c *MongoDBClient) GetRealmClient(ctx context.Context) (*realm.Client, error) {
	// Realm
	if c.Config.PublicKey == "" && c.Config.PrivateKey == "" {
		return nil, errors.New("please set `public_key` and `private_key` in order to use the realm client")
	}

	optsRealm := []realm.ClientOpt{realm.SetUserAgent(userAgent)}
	authConfig := realmAuth.NewConfig(nil)
	if c.Config.BaseURL != "" && c.Config.RealmBaseURL != "" {
		adminURL := c.Config.RealmBaseURL + "api/admin/v3.0/"
		optsRealm = append(optsRealm, realm.SetBaseURL(adminURL))
		authConfig.AuthURL, _ = url.Parse(adminURL + "auth/providers/mongodb-cloud/login")
	}

	token, err := authConfig.NewTokenFromCredentials(ctx, c.Config.PublicKey, c.Config.PrivateKey)
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
