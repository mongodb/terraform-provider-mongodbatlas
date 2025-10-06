package config

import (
	"os"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

// Credentials has all the authentication fields, it also matches with fields that can be stored in AWS Secrets Manager.
type Credentials struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	PublicKey    string `json:"public_key"`
	PrivateKey   string `json:"private_key"`
	BaseURL      string `json:"base_url"`
	RealmBaseURL string `json:"realm_base_url"`
}

// GetCredentials follows the order of AWS Secrets Manager, provider vars and env vars.
func GetCredentials(providerVars, envVars *Vars, getAWSCredentials func(*AWSVars) (*Credentials, error)) (*Credentials, error) {
	if awsVars := CoalesceAWSVars(providerVars.GetAWS(), envVars.GetAWS()); awsVars != nil {
		awsCredentials, err := getAWSCredentials(awsVars)
		if err != nil {
			return nil, err
		}
		return awsCredentials, nil
	}
	if c := CoalesceCredentials(providerVars.GetCredentials(), envVars.GetCredentials()); c != nil {
		return c, nil
	}
	return &Credentials{}, nil
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
		return "Access Token will be used although API Keys is also set"
	}
	if c.HasServiceAccount() && c.HasDigest() {
		return "Service Account will be used although API Keys is also set"
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
	AccessToken        string
	ClientID           string
	ClientSecret       string
	PublicKey          string
	PrivateKey         string
	BaseURL            string
	RealmBaseURL       string
	AWSAssumeRoleARN   string
	AWSSecretName      string
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSSessionToken    string
	AWSEndpoint        string
}

func NewEnvVars() *Vars {
	return &Vars{
		AccessToken:        getEnv("MONGODB_ATLAS_ACCESS_TOKEN", "TF_VAR_ACCESS_TOKEN"),
		ClientID:           getEnv("MONGODB_ATLAS_CLIENT_ID", "TF_VAR_CLIENT_ID"),
		ClientSecret:       getEnv("MONGODB_ATLAS_CLIENT_SECRET", "TF_VAR_CLIENT_SECRET"),
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
		AccessToken:  e.AccessToken,
		ClientID:     e.ClientID,
		ClientSecret: e.ClientSecret,
		PublicKey:    e.PublicKey,
		PrivateKey:   e.PrivateKey,
		BaseURL:      e.BaseURL,
		RealmBaseURL: e.RealmBaseURL,
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
