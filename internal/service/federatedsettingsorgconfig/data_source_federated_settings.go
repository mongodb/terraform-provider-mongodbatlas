package federatedsettingsorgconfig

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

func DataSourceSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"federated_domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"has_role_mappings": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_provider_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).AtlasV2

	orgID, orgIDOk := d.GetOk("org_id")

	if !orgIDOk {
		return diag.FromErr(errors.New("org_id must be configured"))
	}

	var (
		err error
		org *admin.AtlasOrganization
	)

	if orgIDOk {
		org, _, err = conn.OrganizationsApi.GetOrg(ctx, orgID.(string)).Execute()
	}

	if err != nil {
		return diag.Errorf("Error reading Organization %s %s", orgID, err)
	}

	federationSettings, _, err := conn.FederatedAuthenticationApi.GetFederationSettings(ctx, org.GetId()).Execute()
	if err != nil {
		return diag.Errorf("error getting Federated settings (%s): %s", orgID, err)
	}

	if err := d.Set("org_id", org.GetId()); err != nil {
		return diag.Errorf("error getting Federated settings (%s): %s %s", `org_id`, org.GetId(), err)
	}

	if err := d.Set("federated_domains", federationSettings.GetFederatedDomains()); err != nil {
		return diag.Errorf("error getting Federated settings (%s): %s %s", `federated_domains`, federationSettings.GetFederatedDomains(), err)
	}

	if err := d.Set("identity_provider_status", federationSettings.GetIdentityProviderStatus()); err != nil {
		return diag.Errorf("error getting Federated settings (%s): %s %s", `identityProviderStatus`, federationSettings.GetIdentityProviderStatus(), err)
	}

	if err := d.Set("identity_provider_id", federationSettings.GetIdentityProviderId()); err != nil {
		return diag.Errorf("error getting Federated settings (%s): %s %s", `IdentityProviderID`, federationSettings.GetIdentityProviderId(), err)
	}

	if err := d.Set("has_role_mappings", federationSettings.GetHasRoleMappings()); err != nil {
		return diag.Errorf("error getting Federated settings (%s): flag  %s ", `HasRoleMappings`, err)
	}

	d.SetId(federationSettings.GetId())

	return nil
}
