package config_test

import (
	"errors"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentials_AuthMethod(t *testing.T) {
	testCases := map[string]struct {
		credentials config.Credentials
		want        config.AuthMethod
	}{
		"Empty credentials returns Unknown": {
			credentials: config.Credentials{},
			want:        config.Unknown,
		},
		"Access token takes priority": {
			credentials: config.Credentials{
				AccessToken:  "token",
				ClientID:     "id",
				ClientSecret: "secret",
				PublicKey:    "public",
				PrivateKey:   "private",
			},
			want: config.AccessToken,
		},
		"Service account when no access token": {
			credentials: config.Credentials{
				ClientID:     "id",
				ClientSecret: "secret",
				PublicKey:    "public",
				PrivateKey:   "private",
			},
			want: config.ServiceAccount,
		},
		"Service account with only ClientID": {
			credentials: config.Credentials{
				ClientID: "id",
			},
			want: config.ServiceAccount,
		},
		"Service account with only ClientSecret": {
			credentials: config.Credentials{
				ClientSecret: "secret",
			},
			want: config.ServiceAccount,
		},
		"Digest when only digest credentials": {
			credentials: config.Credentials{
				PublicKey:  "public",
				PrivateKey: "private",
			},
			want: config.Digest,
		},
		"Digest with only PublicKey": {
			credentials: config.Credentials{
				PublicKey: "public",
			},
			want: config.Digest,
		},
		"Digest with only PrivateKey": {
			credentials: config.Credentials{
				PrivateKey: "private",
			},
			want: config.Digest,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.credentials.AuthMethod()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCredentials_HasAccessToken(t *testing.T) {
	testCases := map[string]struct {
		credentials config.Credentials
		want        bool
	}{
		"Empty credentials": {
			credentials: config.Credentials{},
			want:        false,
		},
		"With access token": {
			credentials: config.Credentials{
				AccessToken: "token",
			},
			want: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.credentials.HasAccessToken()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCredentials_HasServiceAccount(t *testing.T) {
	testCases := map[string]struct {
		credentials config.Credentials
		want        bool
	}{
		"Empty credentials": {
			credentials: config.Credentials{},
			want:        false,
		},
		"With ClientID only": {
			credentials: config.Credentials{
				ClientID: "id",
			},
			want: true,
		},
		"With ClientSecret only": {
			credentials: config.Credentials{
				ClientSecret: "secret",
			},
			want: true,
		},
		"With both ClientID and ClientSecret": {
			credentials: config.Credentials{
				ClientID:     "id",
				ClientSecret: "secret",
			},
			want: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.credentials.HasServiceAccount()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCredentials_HasDigest(t *testing.T) {
	testCases := map[string]struct {
		credentials config.Credentials
		want        bool
	}{
		"Empty credentials": {
			credentials: config.Credentials{},
			want:        false,
		},
		"With PublicKey only": {
			credentials: config.Credentials{
				PublicKey: "public",
			},
			want: true,
		},
		"With PrivateKey only": {
			credentials: config.Credentials{
				PrivateKey: "private",
			},
			want: true,
		},
		"With both PublicKey and PrivateKey": {
			credentials: config.Credentials{
				PublicKey:  "public",
				PrivateKey: "private",
			},
			want: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.credentials.HasDigest()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCredentials_IsPresent(t *testing.T) {
	testCases := map[string]struct {
		credentials config.Credentials
		want        bool
	}{
		"Empty credentials": {
			credentials: config.Credentials{},
			want:        false,
		},
		"With access token": {
			credentials: config.Credentials{
				AccessToken: "token",
			},
			want: true,
		},
		"With service account": {
			credentials: config.Credentials{
				ClientID: "id",
			},
			want: true,
		},
		"With digest": {
			credentials: config.Credentials{
				PublicKey: "public",
			},
			want: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.credentials.IsPresent()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCredentials_Warnings(t *testing.T) {
	testCases := map[string]struct {
		credentials config.Credentials
		want        string
	}{
		"No credentials": {
			credentials: config.Credentials{},
			want:        "No credentials set",
		},
		"Only access token - no warning": {
			credentials: config.Credentials{
				AccessToken: "token",
			},
			want: "",
		},
		"Only service account - no warning": {
			credentials: config.Credentials{
				ClientID: "id",
			},
			want: "",
		},
		"Only digest - no warning": {
			credentials: config.Credentials{
				PublicKey: "public",
			},
			want: "",
		},
		"Access token and service account": {
			credentials: config.Credentials{
				AccessToken:  "token",
				ClientID:     "id",
				ClientSecret: "secret",
			},
			want: "Access Token will be used although Service Account is also set",
		},
		"Access token and digest": {
			credentials: config.Credentials{
				AccessToken: "token",
				PublicKey:   "public",
				PrivateKey:  "private",
			},
			want: "Access Token will be used although API Key is also set",
		},
		"Service account and digest": {
			credentials: config.Credentials{
				ClientID:   "id",
				PublicKey:  "public",
				PrivateKey: "private",
			},
			want: "Service Account will be used although API Key is also set",
		},
		"All three methods": {
			credentials: config.Credentials{
				AccessToken:  "token",
				ClientID:     "id",
				ClientSecret: "secret",
				PublicKey:    "public",
				PrivateKey:   "private",
			},
			want: "Access Token will be used although Service Account and API Keys are also set",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.credentials.Warnings()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCredentials_Errors(t *testing.T) {
	testCases := map[string]struct {
		credentials config.Credentials
		want        string
	}{
		"No credentials - no error": {
			credentials: config.Credentials{},
			want:        "",
		},
		"Valid access token - no error": {
			credentials: config.Credentials{
				AccessToken: "token",
			},
			want: "",
		},
		"Service account missing ClientID": {
			credentials: config.Credentials{
				ClientSecret: "secret",
			},
			want: "Service Account is being used but Client ID is required",
		},
		"Service account missing ClientSecret": {
			credentials: config.Credentials{
				ClientID: "id",
			},
			want: "Service Account is being used but Client Secret is required",
		},
		"Service account with both - no error": {
			credentials: config.Credentials{
				ClientID:     "id",
				ClientSecret: "secret",
			},
			want: "",
		},
		"Digest missing PublicKey": {
			credentials: config.Credentials{
				PrivateKey: "private",
			},
			want: "API Key is being used but Public Key is required",
		},
		"Digest missing PrivateKey": {
			credentials: config.Credentials{
				PublicKey: "public",
			},
			want: "API Key is being used but Private Key is required",
		},
		"Digest with both - no error": {
			credentials: config.Credentials{
				PublicKey:  "public",
				PrivateKey: "private",
			},
			want: "",
		},
		"Access token takes priority - no error even with incomplete service account": {
			credentials: config.Credentials{
				AccessToken: "token",
				ClientID:    "id",
				// Missing ClientSecret, but should not error since AccessToken takes priority
			},
			want: "",
		},
		"Access token takes priority - no error even with incomplete digest": {
			credentials: config.Credentials{
				AccessToken: "token",
				PublicKey:   "public",
				// Missing PrivateKey, but should not error since AccessToken takes priority
			},
			want: "",
		},
		"Service account takes priority over incomplete digest": {
			credentials: config.Credentials{
				ClientID:     "id",
				ClientSecret: "secret",
				PublicKey:    "public",
				// Missing PrivateKey, but should not error since ServiceAccount takes priority
			},
			want: "",
		},
		"Service account incomplete but takes priority over digest": {
			credentials: config.Credentials{
				ClientID: "id",
				// Missing ClientSecret
				PublicKey:  "public",
				PrivateKey: "private",
			},
			want: "Service Account is being used but Client Secret is required",
		},
		"All credentials present - no error": {
			credentials: config.Credentials{
				AccessToken:  "token",
				ClientID:     "id",
				ClientSecret: "secret",
				PublicKey:    "public",
				PrivateKey:   "private",
			},
			want: "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.credentials.Errors()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetCredentials(t *testing.T) {
	mockGetAWSCredentials := func(awsVars *config.AWSVars) (*config.Credentials, error) {
		if awsVars.AssumeRoleARN == "error" {
			return nil, errors.New("AWS error")
		}
		return &config.Credentials{
			AccessToken: "aws-token",
		}, nil
	}

	testCases := map[string]struct {
		providerVars *config.Vars
		envVars      *config.Vars
		want         *config.Credentials
		wantErr      bool
	}{
		"AWS credentials take priority": {
			providerVars: &config.Vars{
				AWSAssumeRoleARN: "arn",
				PublicKey:        "provider-public",
			},
			envVars: &config.Vars{
				PublicKey: "env-public",
			},
			want: &config.Credentials{
				AccessToken: "aws-token",
			},
			wantErr: false,
		},
		"AWS credentials error": {
			providerVars: &config.Vars{
				AWSAssumeRoleARN: "error",
			},
			envVars: &config.Vars{},
			want:    nil,
			wantErr: true,
		},
		"Provider vars take priority over env vars": {
			providerVars: &config.Vars{
				PublicKey: "provider-public",
			},
			envVars: &config.Vars{
				PublicKey: "env-public",
			},
			want: &config.Credentials{
				PublicKey: "provider-public",
			},
			wantErr: false,
		},
		"Env vars when no provider vars": {
			providerVars: &config.Vars{},
			envVars: &config.Vars{
				PublicKey: "env-public",
			},
			want: &config.Credentials{
				PublicKey: "env-public",
			},
			wantErr: false,
		},
		"Empty credentials when nothing provided": {
			providerVars: &config.Vars{},
			envVars:      &config.Vars{},
			want:         &config.Credentials{},
			wantErr:      false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := config.GetCredentials(tc.providerVars, tc.envVars, mockGetAWSCredentials)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestAWSVars_IsPresent(t *testing.T) {
	testCases := map[string]struct {
		awsVars *config.AWSVars
		want    bool
	}{
		"Empty AWS vars": {
			awsVars: &config.AWSVars{},
			want:    false,
		},
		"With AssumeRoleARN": {
			awsVars: &config.AWSVars{
				AssumeRoleARN: "arn",
			},
			want: true,
		},
		"With other fields but no AssumeRoleARN": {
			awsVars: &config.AWSVars{
				SecretName: "secret",
				Region:     "us-east-1",
			},
			want: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := tc.awsVars.IsPresent()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNewEnvVars(t *testing.T) {
	// Test the first env var for each attribute.
	t.Setenv("MONGODB_ATLAS_ACCESS_TOKEN", "env-token")
	t.Setenv("MONGODB_ATLAS_CLIENT_ID", "env-client-id")
	t.Setenv("MONGODB_ATLAS_CLIENT_SECRET", "env-client-secret")
	t.Setenv("MONGODB_ATLAS_PUBLIC_API_KEY", "env-public")
	t.Setenv("MONGODB_ATLAS_PRIVATE_API_KEY", "env-private")
	t.Setenv("MONGODB_ATLAS_BASE_URL", "url1")
	t.Setenv("MONGODB_REALM_BASE_URL", "url2")
	t.Setenv("ASSUME_ROLE_ARN", "arn")
	t.Setenv("SECRET_NAME", "env-secret")
	t.Setenv("AWS_REGION", "us-west-2")
	t.Setenv("AWS_ACCESS_KEY_ID", "env-access")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "env-secret-key")
	t.Setenv("AWS_SESSION_TOKEN", "env-token")
	t.Setenv("STS_ENDPOINT", "https://sts.amazonaws.com")

	vars := config.NewEnvVars()
	assert.Equal(t, "env-token", vars.AccessToken)
	assert.Equal(t, "env-client-id", vars.ClientID)
	assert.Equal(t, "env-client-secret", vars.ClientSecret)
	assert.Equal(t, "env-public", vars.PublicKey)
	assert.Equal(t, "env-private", vars.PrivateKey)
	assert.Equal(t, "url1", vars.BaseURL)
	assert.Equal(t, "url2", vars.RealmBaseURL)
	assert.Equal(t, "arn", vars.AWSAssumeRoleARN)
	assert.Equal(t, "env-secret", vars.AWSSecretName)
	assert.Equal(t, "us-west-2", vars.AWSRegion)
	assert.Equal(t, "env-access", vars.AWSAccessKeyID)
	assert.Equal(t, "env-secret-key", vars.AWSSecretAccessKey)
	assert.Equal(t, "env-token", vars.AWSSessionToken)
	assert.Equal(t, "https://sts.amazonaws.com", vars.AWSEndpoint)
}

func TestCoalesceAWSVars(t *testing.T) {
	awsVars1 := &config.AWSVars{AssumeRoleARN: "arn1"}
	awsVars2 := &config.AWSVars{AssumeRoleARN: "arn2"}
	awsVarsEmpty := &config.AWSVars{}

	testCases := map[string]struct {
		want    *config.AWSVars
		awsVars []*config.AWSVars
	}{
		"First present AWS vars": {
			awsVars: []*config.AWSVars{awsVars1, awsVars2},
			want:    awsVars1,
		},
		"Skip empty, return first present": {
			awsVars: []*config.AWSVars{awsVarsEmpty, awsVars2},
			want:    awsVars2,
		},
		"All empty returns nil": {
			awsVars: []*config.AWSVars{awsVarsEmpty, awsVarsEmpty},
			want:    nil,
		},
		"No vars returns nil": {
			awsVars: []*config.AWSVars{},
			want:    nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := config.CoalesceAWSVars(tc.awsVars...)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCoalesceCredentials(t *testing.T) {
	creds1 := &config.Credentials{PublicKey: "key1"}
	creds2 := &config.Credentials{PublicKey: "key2"}
	credsEmpty := &config.Credentials{}

	testCases := map[string]struct {
		want        *config.Credentials
		credentials []*config.Credentials
	}{
		"First present credentials": {
			credentials: []*config.Credentials{creds1, creds2},
			want:        creds1,
		},
		"Skip empty, return first present": {
			credentials: []*config.Credentials{credsEmpty, creds2},
			want:        creds2,
		},
		"All empty returns nil": {
			credentials: []*config.Credentials{credsEmpty, credsEmpty},
			want:        nil,
		},
		"No credentials returns nil": {
			credentials: []*config.Credentials{},
			want:        nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := config.CoalesceCredentials(tc.credentials...)
			assert.Equal(t, tc.want, got)
		})
	}
}
