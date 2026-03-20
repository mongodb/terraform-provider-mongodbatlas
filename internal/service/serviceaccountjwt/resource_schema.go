package serviceaccountjwt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func EphemeralResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Client ID of the Service Account.",
			},
			"client_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The service account client secret.",
			},
			"revoke_on_closure": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "When true, the access token is revoked when the Terraform operation ends. Defaults to `false`.",
			},
			"access_token": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The service account access token for authenticating API requests.",
			},
			"token_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The mechanism for token authorization, always `Bearer`.",
			},
			"expires_in": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The duration the access token is valid, in seconds.",
			},
		},
	}
}

type TFModel struct {
	ClientID        types.String `tfsdk:"client_id"`
	ClientSecret    types.String `tfsdk:"client_secret"`
	AccessToken     types.String `tfsdk:"access_token"`
	TokenType       types.String `tfsdk:"token_type"`
	ExpiresIn       types.Int64  `tfsdk:"expires_in"`
	RevokeOnClosure types.Bool   `tfsdk:"revoke_on_closure"`
}
