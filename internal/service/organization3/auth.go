package organization3

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312020/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func (r *organization3RS) atlasV2(ctx context.Context, state *TFModel) *admin.APIClient {
	if state != nil && !state.ClientID.IsNull() && !state.ClientSecret.IsNull() {
		clientID := state.ClientID.ValueString()
		secret := state.ClientSecret.ValueString()
		orgID := state.OrgID.ValueString()
		if clientID != "" && secret != "" {
			if saClient := newSAClient(ctx, orgID, clientID, secret, r.Client); saClient != nil {
				return saClient
			}
		}
	}
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
