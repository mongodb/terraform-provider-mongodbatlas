package config

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
	matlasClient "go.mongodb.org/atlas/mongodbatlas"
	realmAuth "go.mongodb.org/realm/auth"
	"go.mongodb.org/realm/realm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/mongodb-forks/digest"
	adminpreview "github.com/mongodb/atlas-sdk-go/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/version"

	"golang.org/x/oauth2"
)

const (
	toolName              = "terraform-provider-mongodbatlas"
	terraformPlatformName = "Terraform"

	timeout               = 5 * time.Second
	keepAlive             = 30 * time.Second
	maxIdleConns          = 10
	maxIdleConnsPerHost   = 5
	idleConnTimeout       = 30 * time.Second
	expectContinueTimeout = 1 * time.Second
)

type AuthMethod int

const (
	Unknown AuthMethod = iota
	AccessToken
	ServiceAccount
	Digest
)

var baseTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   timeout,
		KeepAlive: keepAlive,
	}).DialContext,
	MaxIdleConns:          maxIdleConns,
	MaxIdleConnsPerHost:   maxIdleConnsPerHost,
	Proxy:                 http.ProxyFromEnvironment,
	IdleConnTimeout:       idleConnTimeout,
	ExpectContinueTimeout: expectContinueTimeout,
}

// networkLoggingBaseTransport should be used as a base for authentication transport so authentication requests can be logged.
func networkLoggingBaseTransport() http.RoundTripper {
	return NewTransportWithNetworkLogging(baseTransport, logging.IsDebugOrHigher())
}

// tfLoggingInterceptor should wrap the authentication transport to add Terraform logging.
func tfLoggingInterceptor(base http.RoundTripper) http.RoundTripper {
	// Don't change logging.NewTransport to NewSubsystemLoggingHTTPTransport until all resources are in TPF.
	return logging.NewTransport("Atlas", base)
}

// MongoDBClient contains the mongodbatlas clients and configurations
type MongoDBClient struct {
	Atlas           *matlasClient.Client
	AtlasV2         *admin.APIClient
	AtlasPreview    *adminpreview.APIClient
	AtlasV220240805 *admin20240805.APIClient // used in advanced_cluster to avoid adopting 2024-10-23 release with ISS autoscaling
	AtlasV220240530 *admin20240530.APIClient // used in advanced_cluster and cloud_backup_schedule for avoiding breaking changes (supporting deprecated replication_specs.id)
	AtlasV220241113 *admin20241113.APIClient // used in teams and atlas_users to avoiding breaking changes
	Config          *Config
}

// Config contains the configurations needed to use SDKs
type Config struct {
	AssumeRoleARN    string
	PublicKey        string
	PrivateKey       string
	BaseURL          string
	RealmBaseURL     string
	TerraformVersion string
	ClientID         string
	ClientSecret     string
	AccessToken      string
}

func NewClient(c *Credentials, terraformVersion string) (*MongoDBClient, error) {
	userAgent := userAgent(terraformVersion)
	client, err := getHTTPClient(c)
	if err != nil {
		return nil, err
	}

	// Initialize the old SDK
	optsAtlas := []matlasClient.ClientOpt{matlasClient.SetUserAgent(userAgent)}
	if c.BaseURL != "" {
		optsAtlas = append(optsAtlas, matlasClient.SetBaseURL(c.BaseURL))
	}
	atlasClient, err := matlasClient.New(client, optsAtlas...)
	if err != nil {
		return nil, err
	}

	// Initialize the new SDK for different versions
	sdkV2Client, err := newSDKV2Client(client, c.BaseURL, userAgent)
	if err != nil {
		return nil, err
	}
	sdkPreviewClient, err := newSDKPreviewClient(client, c.BaseURL, userAgent)
	if err != nil {
		return nil, err
	}
	sdkV220240530Client, err := newSDKV220240530Client(client, c.BaseURL, userAgent)
	if err != nil {
		return nil, err
	}
	sdkV220240805Client, err := newSDKV220240805Client(client, c.BaseURL, userAgent)
	if err != nil {
		return nil, err
	}
	sdkV220241113Client, err := newSDKV220241113Client(client, c.BaseURL, userAgent)
	if err != nil {
		return nil, err
	}

	clients := &MongoDBClient{
		Atlas:           atlasClient,
		AtlasV2:         sdkV2Client,
		AtlasPreview:    sdkPreviewClient,
		AtlasV220240530: sdkV220240530Client,
		AtlasV220240805: sdkV220240805Client,
		AtlasV220241113: sdkV220241113Client,
		// TODO: Config:          c,
	}
	return clients, nil
}

func getHTTPClient(c *Credentials) (*http.Client, error) {
	transport := networkLoggingBaseTransport()
	switch c.AuthMethod() {
	case AccessToken:
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: c.AccessToken,
			TokenType:   "Bearer", // Use a static bearer token with oauth2 transport.
		})
		transport = &oauth2.Transport{
			Source: tokenSource,
			Base:   networkLoggingBaseTransport(),
		}
	case ServiceAccount:
		tokenSource, err := getTokenSource(c.ClientID, c.ClientSecret, c.BaseURL, networkLoggingBaseTransport())
		if err != nil {
			return nil, err
		}
		transport = &oauth2.Transport{
			Source: tokenSource,
			Base:   networkLoggingBaseTransport(),
		}
	case Digest:
		transport = digest.NewTransportWithHTTPRoundTripper(c.PublicKey, c.PrivateKey, networkLoggingBaseTransport())
	case Unknown:
	}
	return &http.Client{Transport: tfLoggingInterceptor(transport)}, nil
}

func newSDKV2Client(client *http.Client, baseURL, userAgent string) (*admin.APIClient, error) {
	return admin.NewClient(
		admin.UseHTTPClient(client),
		admin.UseUserAgent(userAgent),
		admin.UseBaseURL(baseURL),
		admin.UseDebug(false),
	)
}

func newSDKPreviewClient(client *http.Client, baseURL, userAgent string) (*adminpreview.APIClient, error) {
	return adminpreview.NewClient(
		adminpreview.UseHTTPClient(client),
		adminpreview.UseUserAgent(userAgent),
		adminpreview.UseBaseURL(baseURL),
		adminpreview.UseDebug(false),
	)
}

func newSDKV220240530Client(client *http.Client, baseURL, userAgent string) (*admin20240530.APIClient, error) {
	return admin20240530.NewClient(
		admin20240530.UseHTTPClient(client),
		admin20240530.UseUserAgent(userAgent),
		admin20240530.UseBaseURL(baseURL),
		admin20240530.UseDebug(false),
	)
}

func newSDKV220240805Client(client *http.Client, baseURL, userAgent string) (*admin20240805.APIClient, error) {
	return admin20240805.NewClient(
		admin20240805.UseHTTPClient(client),
		admin20240805.UseUserAgent(userAgent),
		admin20240805.UseBaseURL(baseURL),
		admin20240805.UseDebug(false),
	)
}

func newSDKV220241113Client(client *http.Client, baseURL, userAgent string) (*admin20241113.APIClient, error) {
	return admin20241113.NewClient(
		admin20241113.UseHTTPClient(client),
		admin20241113.UseUserAgent(userAgent),
		admin20241113.UseBaseURL(baseURL),
		admin20241113.UseDebug(false),
	)
}

// TODO: lazy because it needs connection
func (c *MongoDBClient) GetRealmClient(ctx context.Context) (*realm.Client, error) {
	// Realm
	if c.Config.PublicKey == "" && c.Config.PrivateKey == "" {
		return nil, errors.New("please set `public_key` and `private_key` in order to use the realm client")
	}

	optsRealm := []realm.ClientOpt{
		realm.SetUserAgent(userAgent(c.Config.TerraformVersion)),
	}

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

	clientRealm := &http.Client{
		Transport: &realmAuth.Transport{
			Source: realmAuth.BasicTokenSource(token),
			Base:   logging.NewTransport("MongoDB Realm", baseTransport),
		},
	}

	// Initialize the MongoDB Realm API Client.
	realmClient, err := realm.New(clientRealm, optsRealm...)
	if err != nil {
		return nil, err
	}

	return realmClient, nil
}

type APICallParams struct {
	VersionHeader string
	RelativePath  string
	PathParams    map[string]string
	Method        string
}

func (c *MongoDBClient) UntypedAPICall(ctx context.Context, params *APICallParams, bodyReq []byte) (*http.Response, error) {
	localBasePath, _ := c.AtlasV2.GetConfig().ServerURLWithContext(ctx, "")
	localVarPath := localBasePath + params.RelativePath

	for key, value := range params.PathParams {
		localVarPath = strings.ReplaceAll(localVarPath, "{"+key+"}", url.PathEscape(value))
	}

	headerParams := make(map[string]string)
	headerParams["Content-Type"] = params.VersionHeader
	headerParams["Accept"] = params.VersionHeader

	var bodyPost any
	if bodyReq != nil { // if nil slice is sent with application/json content type SDK method returns an error
		bodyPost = bodyReq
	}
	untypedClient := c.AtlasV2.UntypedClient
	apiReq, err := untypedClient.PrepareRequest(ctx, localVarPath, params.Method, bodyPost, headerParams, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	apiResp, err := untypedClient.CallAPI(apiReq)
	if err != nil || apiResp == nil {
		return apiResp, err
	}

	// Returns a GenericOpenAPIError error if HTTP status code is not successful.
	if apiResp.StatusCode >= 300 {
		newErr := untypedClient.MakeApiError(apiResp, params.Method, localVarPath)
		return apiResp, newErr
	}

	return apiResp, err
}

func userAgent(terraformVersion string) string {
	metadata := []struct {
		Name  string
		Value string
	}{
		{toolName, version.ProviderVersion},
		{terraformPlatformName, terraformVersion},
	}
	var parts []string
	for _, info := range metadata {
		if info.Value == "" {
			continue
		}
		part := fmt.Sprintf("%s/%s", info.Name, info.Value)
		parts = append(parts, part)
	}
	return strings.Join(parts, " ")
}
