//nolint:testpackage // White-box tests for internal helper behavior.
package serviceaccountjwt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/require"
)

func TestResolveCredentials_Order(t *testing.T) {
	r := &es{
		ESCommon: config.ESCommon{
			EphemeralResourceData: &config.EphemeralResourceData{
				ClientID:         "provider-id",
				ClientSecret:     "provider-secret",
				BaseURL:          "https://provider.example.com/",
				TerraformVersion: "1.10.0",
			},
		},
	}

	t.Run("resource attributes first", func(t *testing.T) {
		model := TFModel{
			ClientID:     types.StringValue("resource-id"),
			ClientSecret: types.StringValue("resource-secret"),
		}
		id, secret, baseURL, diags := r.resolveCredentials(&model)
		require.False(t, diags.HasError())
		require.Equal(t, "resource-id", id)
		require.Equal(t, "resource-secret", secret)
		require.Equal(t, "https://provider.example.com/", baseURL)
	})

	t.Run("provider fallback when no resource attributes", func(t *testing.T) {
		model := TFModel{}
		id, secret, baseURL, diags := r.resolveCredentials(&model)
		require.False(t, diags.HasError())
		require.Equal(t, "provider-id", id)
		require.Equal(t, "provider-secret", secret)
		require.Equal(t, "https://provider.example.com/", baseURL)
	})
}

func TestResolveCredentials_NonSAProviderAuth(t *testing.T) {
	t.Run("nil provider data", func(t *testing.T) {
		r := &es{}
		model := TFModel{}
		clientID, clientSecret, baseURL, diags := r.resolveCredentials(&model)
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
		require.Contains(t, diags.Errors()[0].Detail(), "Service Account credentials")
	})

	t.Run("provider configured with PAK (no SA credentials)", func(t *testing.T) {
		r := &es{
			ESCommon: config.ESCommon{
				EphemeralResourceData: &config.EphemeralResourceData{
					BaseURL:          "https://cloud.mongodb.com/",
					TerraformVersion: "1.11.0",
				},
			},
		}
		model := TFModel{}
		clientID, clientSecret, baseURL, diags := r.resolveCredentials(&model)
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
		require.Contains(t, diags.Errors()[0].Detail(), "different authentication method")
	})
}

func TestResolveCredentials_DoesNotMixSources(t *testing.T) {
	t.Run("partial resource credentials do not fallback", func(t *testing.T) {
		r := &es{
			ESCommon: config.ESCommon{
				EphemeralResourceData: &config.EphemeralResourceData{
					ClientID:     "provider-id",
					ClientSecret: "provider-secret",
					BaseURL:      "https://provider.example.com/",
				},
			},
		}
		model := TFModel{ClientID: types.StringValue("resource-id")}
		clientID, clientSecret, baseURL, diags := r.resolveCredentials(&model)
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
	})

	t.Run("partial provider credentials do not fallback", func(t *testing.T) {
		r := &es{
			ESCommon: config.ESCommon{
				EphemeralResourceData: &config.EphemeralResourceData{
					ClientID: "provider-id",
					BaseURL:  "https://provider.example.com/",
				},
			},
		}
		model := TFModel{}
		clientID, clientSecret, baseURL, diags := r.resolveCredentials(&model)
		require.Empty(t, clientID)
		require.Empty(t, clientSecret)
		require.Empty(t, baseURL)
		require.True(t, diags.HasError())
	})
}
