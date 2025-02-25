package config

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	adminpreview "github.com/mongodb/atlas-sdk-go/admin"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
	matlasClient "go.mongodb.org/atlas/mongodbatlas"
	realmAuth "go.mongodb.org/realm/auth"
	"go.mongodb.org/realm/realm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/mongodb-forks/digest"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/version"
)

const (
	toolName              = "terraform-provider-mongodbatlas"
	terraformPlatformName = "Terraform"
)

// MongoDBClient contains the mongodbatlas clients and configurations
type MongoDBClient struct {
	Atlas           *matlasClient.Client
	AtlasV2         *admin.APIClient
	AtlasPreview    *adminpreview.APIClient
	AtlasV220240805 *admin20240805.APIClient // used in advanced_cluster to avoid adopting 2024-10-23 release with ISS autoscaling
	AtlasV220240530 *admin20240530.APIClient // used in advanced_cluster and cloud_backup_schedule for avoiding breaking changes (supporting deprecated replication_specs.id)
	AtlasV220241113 *admin20241113.APIClient // used in project to avoiding breaking changes
	Config          *Config
}

// Config contains the configurations needed to use SDKs
type Config struct {
	AssumeRole       *AssumeRole
	ProxyPort        *int
	PublicKey        string
	PrivateKey       string
	BaseURL          string
	RealmBaseURL     string
	TerraformVersion string
}

type AssumeRole struct {
	Tags              map[string]string
	RoleARN           string
	ExternalID        string
	Policy            string
	SessionName       string
	SourceIdentity    string
	PolicyARNs        []string
	TransitiveTagKeys []string
	Duration          time.Duration
}

type SecretData struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type PlatformVersion struct {
	Name    string
	Version string
}

// NewClient func...
func (c *Config) NewClient(ctx context.Context) (any, error) {
	// setup a transport to handle digest
	transport := digest.NewTransport(cast.ToString(c.PublicKey), cast.ToString(c.PrivateKey))

	// proxy is only used for testing purposes to connect with hoverfly for capturing/replaying requests
	if c.ProxyPort != nil {
		proxyURL, _ := url.Parse(fmt.Sprintf("http://localhost:%d", *c.ProxyPort))
		transport.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	// initialize the client
	client, err := transport.Client()
	if err != nil {
		return nil, err
	}

	client.Transport = logging.NewTransport("MongoDB Atlas", transport)

	optsAtlas := []matlasClient.ClientOpt{matlasClient.SetUserAgent(userAgent(c))}
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

	sdkPreviewClient, err := c.newSDKPreviewClient(client)
	if err != nil {
		return nil, err
	}

	sdkV220240530Client, err := c.newSDKV220240530Client(client)
	if err != nil {
		return nil, err
	}

	sdkV220240805Client, err := c.newSDKV220240805Client(client)
	if err != nil {
		return nil, err
	}

	clients := &MongoDBClient{
		Atlas:           atlasClient,
		AtlasV2:         sdkV2Client,
		AtlasPreview:    sdkPreviewClient,
		AtlasV220240530: sdkV220240530Client,
		AtlasV220240805: sdkV220240805Client,
		Config:          c,
	}
	return clients, nil
}

func (c *Config) newSDKV2Client(client *http.Client) (*admin.APIClient, error) {
	opts := []admin.ClientModifier{
		admin.UseHTTPClient(client),
		admin.UseUserAgent(userAgent(c)),
		admin.UseBaseURL(c.BaseURL),
		admin.UseDebug(false)}

	sdk, err := admin.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func (c *Config) newSDKPreviewClient(client *http.Client) (*adminpreview.APIClient, error) {
	opts := []adminpreview.ClientModifier{
		adminpreview.UseHTTPClient(client),
		adminpreview.UseUserAgent(userAgent(c)),
		adminpreview.UseBaseURL(c.BaseURL),
		adminpreview.UseDebug(false)}

	sdk, err := adminpreview.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func (c *Config) newSDKV220240530Client(client *http.Client) (*admin20240530.APIClient, error) {
	opts := []admin20240530.ClientModifier{
		admin20240530.UseHTTPClient(client),
		admin20240530.UseUserAgent(userAgent(c)),
		admin20240530.UseBaseURL(c.BaseURL),
		admin20240530.UseDebug(false)}

	sdk, err := admin20240530.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func (c *Config) newSDKV220240805Client(client *http.Client) (*admin20240805.APIClient, error) {
	opts := []admin20240805.ClientModifier{
		admin20240805.UseHTTPClient(client),
		admin20240805.UseUserAgent(userAgent(c)),
		admin20240805.UseBaseURL(c.BaseURL),
		admin20240805.UseDebug(false)}

	sdk, err := admin20240805.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func (c *MongoDBClient) GetRealmClient(ctx context.Context) (*realm.Client, error) {
	// Realm
	if c.Config.PublicKey == "" && c.Config.PrivateKey == "" {
		return nil, errors.New("please set `public_key` and `private_key` in order to use the realm client")
	}

	optsRealm := []realm.ClientOpt{realm.SetUserAgent(userAgent(c.Config))}

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

func userAgent(c *Config) string {
	platformVersions := []PlatformVersion{
		{toolName, version.ProviderVersion},
		{terraformPlatformName, c.TerraformVersion},
	}

	var parts []string
	for _, info := range platformVersions {
		part := fmt.Sprintf("%s/%s", info.Name, info.Version)
		parts = append(parts, part)
	}

	return strings.Join(parts, " ")
}
