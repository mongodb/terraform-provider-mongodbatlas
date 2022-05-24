package mongodbatlas

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasFederatedSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasFederatedSettingsRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"org_id"},
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

func dataSourceMongoDBAtlasFederatedSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	orgID, orgIDOk := d.GetOk("org_id")
	name, nameOk := d.GetOk("name")

	if !orgIDOk && !nameOk {
		return diag.FromErr(errors.New("either org_id or name must be configured"))
	}

	var (
		err  error
		org  *matlas.Organization
		orgs *matlas.Organizations
	)

	if orgIDOk {
		org, _, err = conn.Organizations.Get(ctx, orgID.(string))
	} else {
		orgs, _, err = conn.Organizations.List(ctx, nil)
		if err != nil {
			return diag.Errorf("Organizations.List returned error: %v", err)
		}
		for _, o := range orgs.Results {
			if o.Name == name.(string) {
				org, _, err = conn.Organizations.Get(ctx, o.ID)
			}
		}
	}

	if err != nil {
		return diag.Errorf("Error reading Organization %s %s", orgID, err)
	}

	federationSettings, _, err := conn.FederatedSettings.Get(ctx, org.ID)
	if err != nil {
		return diag.Errorf("error getting Federated settings (%s): %s", orgID, err)
	}

	if err := d.Set("org_id", org.ID); err != nil {
		return diag.Errorf("error getting Federated settings (%s): %s %s", `org_id`, org.ID, err)
	}

	if err := d.Set("federated_domains", federationSettings.FederatedDomains); err != nil {
		return diag.Errorf("error getting Federated settings (%s): %s %s", `federated_domains`, federationSettings.FederatedDomains, err)
	}

	if err := d.Set("identity_provider_status", federationSettings.IdentityProviderStatus); err != nil {
		return diag.Errorf("error getting Federated settings (%s): %s %s", `identityProviderStatus`, federationSettings.IdentityProviderStatus, err)
	}

	if err := d.Set("identity_provider_id", federationSettings.IdentityProviderID); err != nil {
		return diag.Errorf("error getting Federated settings (%s): %s %s", `IdentityProviderID`, federationSettings.IdentityProviderID, err)
	}

	if err := d.Set("has_role_mappings", federationSettings.HasRoleMappings); err != nil {
		return diag.Errorf("error getting Federated settings (%s): flag  %s ", `HasRoleMappings`, err)
	}

	d.SetId(federationSettings.ID)

	return nil
}
