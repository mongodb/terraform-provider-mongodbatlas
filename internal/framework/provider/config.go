package provider

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	matlasClient "go.mongodb.org/atlas/mongodbatlas"
	realmAuth "go.mongodb.org/realm/auth"
	"go.mongodb.org/realm/realm"

	"github.com/mongodb-forks/digest"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/version"
)

type Config struct {
	PublicKey    string
	PrivateKey   string
	BaseURL      string
	RealmBaseURL string
	AssumeRole   *AssumeRole
}

type AssumeRole struct {
	RoleARN           string
	Duration          time.Duration
	ExternalID        string
	Policy            string
	PolicyARNs        []string
	SessionName       string
	SourceIdentity    string
	Tags              map[string]string
	TransitiveTagKeys []string
}

// MongoDBClient client
type MongoDBClient struct {
	Atlas  *matlasClient.Client
	Config *Config
}

const ToolName = "terraform-provider-mongodbatlas"

var UserAgent = fmt.Sprintf("%s/%s (%s;%s)", ToolName, version.ProviderVersion, runtime.GOOS, runtime.GOARCH)

// NewClient func...
func (config *Config) NewClient(ctx context.Context) (*MongoDBClient, error) {
	// setup a transport to handle digest
	transport := digest.NewTransport(cast.ToString(config.PublicKey), cast.ToString(config.PrivateKey))

	// initialize the client
	client, err := transport.Client()
	if err != nil {
		return nil, err
	}

	// client.Transport = logging.NewTransport("MongoDB Atlas", transport)

	optsAtlas := []matlasClient.ClientOpt{matlasClient.SetUserAgent(UserAgent)}
	if config.BaseURL != "" {
		optsAtlas = append(optsAtlas, matlasClient.SetBaseURL(config.BaseURL))
	}

	// Initialize the MongoDB Atlas API Client.
	atlasClient, err := matlasClient.New(client, optsAtlas...)
	if err != nil {
		return nil, err
	}

	clients := &MongoDBClient{
		Atlas:  atlasClient,
		Config: config,
	}

	return clients, nil
}

func (config *MongoDBClient) GetRealmClient(ctx context.Context) (*realm.Client, error) {
	// Realm
	if config.Config.PublicKey == "" && config.Config.PrivateKey == "" {
		return nil, errors.New("please set `public_key` and `private_key` in order to use the realm client")
	}

	optsRealm := []realm.ClientOpt{realm.SetUserAgent(UserAgent)}
	if config.Config.BaseURL != "" && config.Config.RealmBaseURL != "" {
		optsRealm = append(optsRealm, realm.SetBaseURL(config.Config.RealmBaseURL))
	}
	authConfig := realmAuth.NewConfig(nil)
	token, err := authConfig.NewTokenFromCredentials(ctx, config.Config.PublicKey, config.Config.PrivateKey)
	if err != nil {
		return nil, err
	}

	clientRealm := realmAuth.NewClient(realmAuth.BasicTokenSource(token))
	// clientRealm.Transport = logging.NewTransport("MongoDB Realm", clientRealm.Transport)

	// Initialize the MongoDB Realm API Client.
	realmClient, err := realm.New(clientRealm, optsRealm...)
	if err != nil {
		return nil, err
	}

	return realmClient, nil
}
