package config

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const ErrPartialCreds = "Both client_id and client_secret must be provided together." //nolint:gosec // not a credential

type CredentialResolver struct {
	ProviderData *EphemeralResourceData
}

// Resolve returns the client_id, client_secret, and base URL to use for token generation.
// It checks specified arguments first (clientID/clientSecret), then falls back
// to provider credentials from ProviderData. Both client_id and client_secret must come from
// the same source.
func (cr *CredentialResolver) ResolveServiceAccountCredentials(argClientID, argClientSecret string) (clientID, clientSecret, baseURL string, diags diag.Diagnostics) {
	id := strings.TrimSpace(argClientID)
	secret := strings.TrimSpace(argClientSecret)

	// Resource attributes (explicit client_id and client_secret).
	if id != "" && secret != "" {
		return id, secret, cr.providerBaseURL(), diags
	} else if id != "" || secret != "" {
		diags.AddError("Invalid Service Account credentials",
			ErrPartialCreds)
		return "", "", "", diags
	}

	// Provider credentials (the provider already coalesces HCL config and env vars).
	if cr.ProviderData != nil {
		id = strings.TrimSpace(cr.ProviderData.ClientID)
		secret = strings.TrimSpace(cr.ProviderData.ClientSecret)
		if id != "" && secret != "" {
			return id, secret, cr.providerBaseURL(), diags
		} else if id != "" || secret != "" {
			diags.AddError("Invalid Service Account credentials",
				"Both client_id and client_secret must be configured for the provider.")
			return "", "", "", diags
		}
	}

	diags.AddError("No Service Account credentials found",
		"No client_id and client_secret were found in the resource attributes or provider configuration.")
	return "", "", "", diags
}

func (cr *CredentialResolver) providerBaseURL() string {
	if cr.ProviderData != nil {
		return strings.TrimSpace(cr.ProviderData.BaseURL)
	}
	return ""
}
