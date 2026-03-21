package serviceaccountjwt

import (
	"context"
	"encoding/json"
	"log"

	"golang.org/x/oauth2"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/atlas-sdk-go/auth"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	ResourceTypeName = "service_account_jwt"
	closeDataKey     = "revoke_data"
)

type TokenGenerator interface {
	GenerateToken(ctx context.Context, clientID, clientSecret, baseURL string) (*auth.Token, error)
}

var _ ephemeral.EphemeralResource = &ER{}
var _ ephemeral.EphemeralResourceWithConfigure = &ER{}
var _ ephemeral.EphemeralResourceWithClose = &ER{}

type ER struct {
	TokenGen TokenGenerator
	config.ESCommon
}

type closeData struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	BaseURL      string `json:"base_url"`
}

func New() ephemeral.EphemeralResource {
	r := &ER{
		ESCommon: config.ESCommon{
			ResourceName: ResourceTypeName,
		},
	}
	r.TokenGen = r
	return r
}

func (r *ER) Schema(ctx context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = EphemeralResourceSchema(ctx)
}

func (r *ER) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var model TFModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if model.ClientID.IsUnknown() || model.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddError("Unknown credentials",
			"client_id and client_secret must be known at apply time to generate a token.")
		return
	}

	resolver := &config.CredentialResolver{ProviderData: r.EphemeralResourceData}
	clientID, clientSecret, baseURL, localDiags := resolver.ResolveServiceAccountCredentials(
		model.ClientID.ValueString(),
		model.ClientSecret.ValueString(),
	)
	resp.Diagnostics.Append(localDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	token, err := r.TokenGen.GenerateToken(ctx, clientID, clientSecret, baseURL)
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

func (r *ER) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	revokeData, diags := req.Private.GetKey(ctx, closeDataKey)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(revokeData) == 0 {
		return
	}
	var data closeData
	if err := json.Unmarshal(revokeData, &data); err != nil {
		resp.Diagnostics.AddWarning("Failed to read revoke payload",
			"Could not deserialize the token revocation data from private state. The token will not be revoked.")
		return
	}
	log.Printf("[DEBUG] %s Close: revoking access token", ResourceTypeName)
	conf := config.GetServiceAccountConfig(data.ClientID, data.ClientSecret, config.NormalizeBaseURL(data.BaseURL))
	if err := conf.RevokeToken(r.oauthClientCtx(ctx), &oauth2.Token{AccessToken: data.AccessToken}); err != nil {
		resp.Diagnostics.AddWarning("Failed to revoke Service Account token on close", err.Error())
	}
}

// oauthClientCtx injects an HTTP client into the context via
// auth.HTTPClient so the atlas-sdk-go token exchange picks it up.
func (r *ER) oauthClientCtx(ctx context.Context) context.Context {
	client := config.NewOAuthHTTPClient(r.TerraformVersion())
	return context.WithValue(ctx, auth.HTTPClient, client)
}

func (r *ER) GenerateToken(ctx context.Context, clientID, clientSecret, baseURL string) (*auth.Token, error) {
	conf := config.GetServiceAccountConfig(clientID, clientSecret, config.NormalizeBaseURL(baseURL))
	token, err := conf.Token(r.oauthClientCtx(ctx))
	if err != nil {
		return nil, err
	}
	return token, nil
}
