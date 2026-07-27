package config

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

// Credentials has all the authentication fields, it also matches with fields that can be stored in AWS Secrets Manager.
type Credentials struct {
	AccessToken       string `json:"access_token"`
	ClientID          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	PublicKey         string `json:"public_key"`
	PrivateKey        string `json:"private_key"`
	BaseURL           string `json:"base_url"`
	RealmBaseURL      string `json:"realm_base_url"`
	IsMongodbGovCloud bool   `json:"is_mongodbgov_cloud"`
}

// applyGovBaseURL sets BaseURL to the gov control plane when IsMongodbGovCloud
// is set, unless BaseURL is already a recognized dev/qa gov URL.
func (c *Credentials) applyGovBaseURL() {
	const govBaseURL = "https://cloud.mongodbgov.com"
	// additionalBaseURLs are gov control planes that must be preserved as-is
	// instead of being replaced with govBaseURL.
	additionalBaseURLs := []string{
		"https://cloud-dev.mongodbgov.com",
		"https://cloud-qa.mongodbgov.com",
	}
	if c.IsMongodbGovCloud && !slices.Contains(additionalBaseURLs, NormalizeBaseURL(c.BaseURL)) {
		c.BaseURL = govBaseURL
	}
}

// GetCredentials follows the order of AWS Secrets Manager, provider vars and env vars.
func GetCredentials(ctx context.Context, providerVars, envVars *Vars, getAWSCredentials func(context.Context, *AWSVars) (*Credentials, error)) (*Credentials, error) {
	var creds *Credentials
	if awsVars := CoalesceAWSVars(providerVars.GetAWS(), envVars.GetAWS()); awsVars != nil {
		awsCredentials, err := getAWSCredentials(ctx, awsVars)
		if err != nil {
			return nil, err
		}
		creds = awsCredentials
	} else if c := CoalesceCredentials(providerVars.GetCredentials(), envVars.GetCredentials()); c != nil {
		creds = c
	} else {
		creds = &Credentials{}
	}
	creds.applyGovBaseURL()
	return creds, nil
}

// AWSSecretsManagerIgnoredWarning returns a warning when base_url, realm_base_url or
// is_mongodbgov_cloud are set (as provider attributes or env vars) while AWS Secrets Manager
// credentials are in use, since those values are ignored on that path and must be defined in
// the secret payload instead. It returns an empty string when there is nothing to warn about.
func AWSSecretsManagerIgnoredWarning(providerVars, envVars *Vars) string {
	if CoalesceAWSVars(providerVars.GetAWS(), envVars.GetAWS()) == nil {
		return ""
	}
	var ignored []string
	if providerVars.BaseURL != "" || envVars.BaseURL != "" {
		ignored = append(ignored, "base_url")
	}
	if providerVars.RealmBaseURL != "" || envVars.RealmBaseURL != "" {
		ignored = append(ignored, "realm_base_url")
	}
	if providerVars.IsMongodbGovCloud || envVars.IsMongodbGovCloud {
		ignored = append(ignored, "is_mongodbgov_cloud")
	}
	if len(ignored) == 0 {
		return ""
	}
	return fmt.Sprintf("The following attributes are ignored when using AWS Secrets Manager credentials: %s. "+
		"Define these values in the secret payload instead.", strings.Join(ignored, ", "))
}

// AuthMethod follows the order of token, SA and PAK.
func (c *Credentials) AuthMethod() AuthMethod {
	switch {
	case c.HasAccessToken():
		return AccessToken
	case c.HasServiceAccount():
		return ServiceAccount
	case c.HasDigest():
		return Digest
	default:
		return Unknown
	}
}

func (c *Credentials) HasAccessToken() bool {
	return c.AccessToken != ""
}

func (c *Credentials) HasServiceAccount() bool {
	return c.ClientID != "" || c.ClientSecret != ""
}

func (c *Credentials) HasDigest() bool {
	return c.PublicKey != "" || c.PrivateKey != ""
}

func (c *Credentials) IsPresent() bool {
	return c.AuthMethod() != Unknown
}

func (c *Credentials) Warnings() string {
	if !c.IsPresent() {
		return "No credentials set"
	}
	// Prefer specific checks over generic code as there are few combinations and code is clearer.
	if c.HasAccessToken() && c.HasServiceAccount() && c.HasDigest() {
		return "Access Token will be used although Service Account and API Keys are also set"
	}
	if c.HasAccessToken() && c.HasServiceAccount() {
		return "Access Token will be used although Service Account is also set"
	}
	if c.HasAccessToken() && c.HasDigest() {
		return "Access Token will be used although API Key is also set"
	}
	if c.HasServiceAccount() && c.HasDigest() {
		return "Service Account will be used although API Key is also set"
	}
	return ""
}

const serviceAccountPrefix = "mdb_sa"

func (c *Credentials) Errors() string {
	switch c.AuthMethod() {
	case ServiceAccount:
		if c.ClientID == "" {
			return "Service Account will be used but Client ID is required"
		}
		if c.ClientSecret == "" {
			return "Service Account will be used but Client Secret is required"
		}
	case Digest:
		if strings.HasPrefix(c.PublicKey, serviceAccountPrefix) {
			return "Service Account credentials (starting with 'mdb_sa') were provided in public_key/private_key which are meant for Programmatic Access Keys. " +
				"Please use client_id and client_secret arguments for Service Account authentication"
		}
		if c.PublicKey == "" {
			return "API Key will be used but Public Key is required"
		}
		if c.PrivateKey == "" {
			return "API Key will be used but Private Key is required"
		}
	case Unknown, AccessToken:
	}
	return ""
}

type AWSVars struct {
	AssumeRoleARN   string
	SecretName      string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Endpoint        string
}

func (a *AWSVars) IsPresent() bool {
	return a.AssumeRoleARN != ""
}

type Vars struct {
	RealmBaseURL       string
	AWSRegion          string
	ClientSecret       string
	PublicKey          string
	PrivateKey         string
	BaseURL            string
	ClientID           string
	AWSAssumeRoleARN   string
	AccessToken        string
	AWSSecretName      string
	AWSEndpoint        string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSSessionToken    string
	IsMongodbGovCloud  bool
}

func NewEnvVars() *Vars {
	return &Vars{
		AccessToken:        getEnv("MONGODB_ATLAS_ACCESS_TOKEN"),
		ClientID:           getEnv("MONGODB_ATLAS_CLIENT_ID"),
		ClientSecret:       getEnv("MONGODB_ATLAS_CLIENT_SECRET"),
		PublicKey:          getEnv("MONGODB_ATLAS_PUBLIC_API_KEY", "MONGODB_ATLAS_PUBLIC_KEY", "MCLI_PUBLIC_API_KEY"),
		PrivateKey:         getEnv("MONGODB_ATLAS_PRIVATE_API_KEY", "MONGODB_ATLAS_PRIVATE_KEY", "MCLI_PRIVATE_API_KEY"),
		BaseURL:            getEnv("MONGODB_ATLAS_BASE_URL", "MCLI_OPS_MANAGER_URL"),
		RealmBaseURL:       getEnv("MONGODB_REALM_BASE_URL"),
		AWSAssumeRoleARN:   getEnv("ASSUME_ROLE_ARN", "TF_VAR_ASSUME_ROLE_ARN"),
		AWSSecretName:      getEnv("SECRET_NAME", "TF_VAR_SECRET_NAME"),
		AWSRegion:          getEnv("AWS_REGION", "TF_VAR_AWS_REGION"),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", "TF_VAR_AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", "TF_VAR_AWS_SECRET_ACCESS_KEY"),
		AWSSessionToken:    getEnv("AWS_SESSION_TOKEN", "TF_VAR_AWS_SESSION_TOKEN"),
		AWSEndpoint:        getEnv("STS_ENDPOINT", "TF_VAR_STS_ENDPOINT"),
	}
}

func (e *Vars) GetCredentials() *Credentials {
	return &Credentials{
		AccessToken:       e.AccessToken,
		ClientID:          e.ClientID,
		ClientSecret:      e.ClientSecret,
		PublicKey:         e.PublicKey,
		PrivateKey:        e.PrivateKey,
		BaseURL:           e.BaseURL,
		RealmBaseURL:      e.RealmBaseURL,
		IsMongodbGovCloud: e.IsMongodbGovCloud,
	}
}

// GetAWS returns variables in the format AWS expects, e.g. region in lowercase.
func (e *Vars) GetAWS() *AWSVars {
	return &AWSVars{
		AssumeRoleARN:   e.AWSAssumeRoleARN,
		SecretName:      e.AWSSecretName,
		Region:          conversion.MongoDBRegionToAWSRegion(e.AWSRegion),
		AccessKeyID:     e.AWSAccessKeyID,
		SecretAccessKey: e.AWSSecretAccessKey,
		SessionToken:    e.AWSSessionToken,
		Endpoint:        e.AWSEndpoint,
	}
}

func getEnv(key ...string) string {
	for _, k := range key {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func CoalesceAWSVars(awsVars ...*AWSVars) *AWSVars {
	for _, awsVar := range awsVars {
		if awsVar.IsPresent() {
			return awsVar
		}
	}
	return nil
}

func CoalesceCredentials(credentials ...*Credentials) *Credentials {
	for _, credential := range credentials {
		if credential.IsPresent() {
			return credential
		}
	}
	return nil
}
