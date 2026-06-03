package organization3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"go.mongodb.org/atlas-sdk/v20250312020/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func (r *organization3RS) atlasV2(ctx context.Context, state *TFModel) *admin.APIClient {
	if state == nil {
		tflog.Debug(ctx, "organization3: state is nil, using provider-configured API credentials")
		return r.Client.AtlasV2
	}
	if state.ClientID.IsNull() || state.ClientSecret.IsNull() {
		tflog.Debug(ctx, "organization3: client_id or client_secret missing in state, using provider-configured API credentials",
			map[string]any{
				"client_id_null":     state.ClientID.IsNull(),
				"client_secret_null": state.ClientSecret.IsNull(),
				"org_id":             state.OrgID.ValueString(),
			})
		return r.Client.AtlasV2
	}
	clientID := state.ClientID.ValueString()
	secret := state.ClientSecret.ValueString()
	orgID := state.OrgID.ValueString()
	if clientID == "" || secret == "" {
		tflog.Debug(ctx, "organization3: client_id or client_secret empty in state, using provider-configured API credentials",
			map[string]any{"org_id": orgID})
		return r.Client.AtlasV2
	}
	if saClient := newSAClient(ctx, orgID, clientID, secret, r.Client); saClient != nil {
		tflog.Debug(ctx, "organization3: using service account credentials from state",
			map[string]any{"org_id": orgID, "client_id": clientID})
		return saClient
	}
	tflog.Debug(ctx, "organization3: service account credentials from state failed validation, using provider-configured API credentials",
		map[string]any{"org_id": orgID, "client_id": clientID})
	return r.Client.AtlasV2
}

func newSAClient(ctx context.Context, orgID, clientID, secretValue string, currentClient *config.MongoDBClient) *admin.APIClient {
	c := &config.Credentials{
		ClientID:     clientID,
		ClientSecret: secretValue,
		BaseURL:      currentClient.BaseURL,
	}
	newClient, err := config.NewClient(c, currentClient.TerraformVersion)
	if err != nil {
		return nil
	}
	if orgID != "" {
		if _, _, err := newClient.AtlasV2.OrganizationsApi.GetOrg(ctx, orgID).Execute(); err != nil {
			return nil
		}
	}
	return newClient.AtlasV2
}

func providerAtlasV2(client *config.MongoDBClient) *admin.APIClient {
	return client.AtlasV2
}
