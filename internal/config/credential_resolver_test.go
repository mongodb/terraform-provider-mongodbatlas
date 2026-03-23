package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func TestCredentialResolver_Order(t *testing.T) {
	cr := &config.CredentialResolver{
		ProviderData: &config.EphemeralResourceData{
			ClientID:     "provider-id",
			ClientSecret: "provider-secret",
			BaseURL:      "https://provider.example.com/",
		},
	}

	t.Run("resource attributes first", func(t *testing.T) {
		id, secret, baseURL, diags := cr.ResolveServiceAccountCredentials("resource-id", "resource-secret")
		require.False(t, diags.HasError())
		require.Equal(t, "resource-id", id)
		require.Equal(t, "resource-secret", secret)
		require.Equal(t, "https://provider.example.com/", baseURL)
	})

	t.Run("provider fallback when no resource attributes", func(t *testing.T) {
		id, secret, baseURL, diags := cr.ResolveServiceAccountCredentials("", "")
		require.False(t, diags.HasError())
		require.Equal(t, "provider-id", id)
		require.Equal(t, "provider-secret", secret)
		require.Equal(t, "https://provider.example.com/", baseURL)
	})
}

func TestCredentialResolver_NonSAProviderAuth(t *testing.T) {
	t.Run("nil provider data", func(t *testing.T) {
		cr := &config.CredentialResolver{ProviderData: nil}
		clientID, clientSecret, baseURL, diags := cr.ResolveServiceAccountCredentials("", "")
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
		require.Contains(t, diags.Errors()[0].Detail(), "No client_id and client_secret were found")
	})

	t.Run("provider configured with PAK (no SA credentials)", func(t *testing.T) {
		cr := &config.CredentialResolver{
			ProviderData: &config.EphemeralResourceData{
				BaseURL: "https://cloud.mongodb.com/",
			},
		}
		clientID, clientSecret, baseURL, diags := cr.ResolveServiceAccountCredentials("", "")
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
		require.Contains(t, diags.Errors()[0].Detail(), "No client_id and client_secret were found")
	})
}

func TestCredentialResolver_DoesNotMixSources(t *testing.T) {
	t.Run("partial resource credentials do not fallback", func(t *testing.T) {
		cr := &config.CredentialResolver{
			ProviderData: &config.EphemeralResourceData{
				ClientID:     "provider-id",
				ClientSecret: "provider-secret",
				BaseURL:      "https://provider.example.com/",
			},
		}
		clientID, clientSecret, baseURL, diags := cr.ResolveServiceAccountCredentials("resource-id", "")
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
	})

	t.Run("partial provider credentials do not fallback", func(t *testing.T) {
		cr := &config.CredentialResolver{
			ProviderData: &config.EphemeralResourceData{
				ClientID: "provider-id",
				BaseURL:  "https://provider.example.com/",
			},
		}
		clientID, clientSecret, baseURL, diags := cr.ResolveServiceAccountCredentials("", "")
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
	})
}
