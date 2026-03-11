package serviceaccountjwt

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"golang.org/x/oauth2"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/atlas-sdk-go/auth"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	resourceTypeName = "service_account_jwt"
	closeDataKey     = "revoke_data"
)

var _ ephemeral.EphemeralResource = &es{}
var _ ephemeral.EphemeralResourceWithConfigure = &es{}
var _ ephemeral.EphemeralResourceWithClose = &es{}

type es struct {
	config.ESCommon
}

type closeData struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	BaseURL      string `json:"base_url"`
}

func New() ephemeral.EphemeralResource {
	return &es{
		ESCommon: config.ESCommon{
			ResourceName: resourceTypeName,
		},
	}
}

func (r *es) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = EphemeralResourceSchema(ctx)
}

func (r *es) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var model TFModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientID, clientSecret, baseURL, localDiags := r.resolveCredentials(&model)
	resp.Diagnostics.Append(localDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.generateToken(ctx, clientID, clientSecret, baseURL)
	if err != nil {
		resp.Diagnostics.AddError("Error generating Service Account JWT", err.Error())
		return
	}

	model.AccessToken = types.StringValue(token.AccessToken)
	model.TokenType = types.StringValue(token.Type())
	model.ExpiresIn = types.Int64Value(token.ExpiresIn)

	if model.RevokeOnClosure.ValueBool() {
		revokeData, err := json.Marshal(closeData{
			AccessToken:  token.AccessToken,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			BaseURL:      baseURL,
		})
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Failed to prepare Service Account token revoke payload",
				"The token will not be automatically revoked when the ephemeral resource is closed: "+err.Error())
		} else {
			resp.Diagnostics.Append(resp.Private.SetKey(ctx, closeDataKey, revokeData)...)
		}
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &model)...)
}

func (r *es) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	raw, diags := req.Private.GetKey(ctx, closeDataKey)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(raw) == 0 {
		log.Printf("[DEBUG] %s Close: no private state found (key=%q), skipping revocation", resourceTypeName, closeDataKey)
		return
	}
	var data closeData
	if err := json.Unmarshal(raw, &data); err != nil {
		resp.Diagnostics.AddWarning("Failed to read revoke payload", err.Error())
		return
	}
	log.Printf("[DEBUG] %s Close: revoking access token", resourceTypeName)
	conf := config.GetServiceAccountConfig(data.ClientID, data.ClientSecret, config.NormalizeBaseURL(data.BaseURL))
	if err := conf.RevokeToken(r.withUserAgentClient(ctx), &oauth2.Token{AccessToken: data.AccessToken}); err != nil {
		resp.Diagnostics.AddWarning("Failed to revoke Service Account token on close", err.Error())
	}
}

func (r *es) resolveCredentials(model *TFModel) (clientID, clientSecret, baseURL string, diags diag.Diagnostics) {
	erd := r.EphemeralResourceData

	// 1. Resource attributes (explicit client_id and client_secret on the ephemeral resource block).
	if model.ClientID.IsUnknown() || model.ClientSecret.IsUnknown() {
		diags.AddError("Unknown credentials",
			"client_id and client_secret must be known at apply time to generate a token.")
		return "", "", "", diags
	}
	id := strings.TrimSpace(model.ClientID.ValueString())
	secret := strings.TrimSpace(model.ClientSecret.ValueString())
	if id != "" && secret != "" {
		return id, secret, providerBaseURL(erd), diags
	} else if id != "" || secret != "" {
		diags.AddError("Invalid Service Account credentials",
			"When setting credentials on this ephemeral resource, both client_id and client_secret must be provided.")
		return "", "", "", diags
	}

	// 2. Provider credentials (the provider already coalesces HCL config and env vars).
	if erd != nil {
		id = strings.TrimSpace(erd.ClientID)
		secret = strings.TrimSpace(erd.ClientSecret)
		if id != "" && secret != "" {
			return id, secret, providerBaseURL(erd), diags
		} else if id != "" || secret != "" {
			diags.AddError("Invalid Service Account credentials",
				"To use this resource please ensure both client_id and client_secret are configured for the provider.")
			return "", "", "", diags
		}
	}

	// 3. No SA credentials found, provider is configured with non-SA auth (PAK or Access Token).
	diags.AddError(
		"Service Account credentials required",
		"This ephemeral resource requires Service Account credentials (client_id and client_secret). "+
			"The provider is currently configured with a different authentication method (Programmatic Access Key or Access Token). "+
			"Set client_id and client_secret on the ephemeral resource block or configure the provider with Service Account credentials.",
	)
	return "", "", "", diags
}

func providerBaseURL(providerData *config.EphemeralResourceData) string {
	if providerData != nil {
		return strings.TrimSpace(providerData.BaseURL)
	}
	return ""
}

func (r *es) generateToken(ctx context.Context, clientID, clientSecret, baseURL string) (*auth.Token, error) {
	conf := config.GetServiceAccountConfig(clientID, clientSecret, config.NormalizeBaseURL(baseURL))
	token, err := conf.Token(r.withUserAgentClient(ctx))
	if err != nil {
		return nil, err
	}
	return token, nil
}
