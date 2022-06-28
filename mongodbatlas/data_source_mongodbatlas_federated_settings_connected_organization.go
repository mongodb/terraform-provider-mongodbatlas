package mongodbatlas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasFederatedSettingsOrganizationConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasFederatedSettingsOrganizationConfigRead,
		Schema: map[string]*schema.Schema{
			"federation_settings_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_allow_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_restriction_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"identity_provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"post_auth_role_grants": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"role_mappings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_assignments": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"org_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"role": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"user_conflicts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"federation_settings_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceMongoDBAtlasFederatedSettingsOrganizationConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}

	orgID, orgIDOk := d.GetOk("org_id")

	if !orgIDOk {
		return diag.FromErr(errors.New("org_id must be configured"))
	}

	federatedSettingsConnectedOrganization, _, err := conn.FederatedSettings.GetConnectedOrg(ctx, federationSettingsID.(string), orgID.(string))
	if err != nil {
		return diag.Errorf("error getting federatedSettings connected organizations assigned (%s): %s", federationSettingsID, err)
	}

	if err := d.Set("domain_allow_list", federatedSettingsConnectedOrganization.DomainAllowList); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `domain_allow_list` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("domain_restriction_enabled", federatedSettingsConnectedOrganization.DomainRestrictionEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `domain_restriction_enabled` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("identity_provider_id", federatedSettingsConnectedOrganization.IdentityProviderID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `identity_provider_id` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("org_id", federatedSettingsConnectedOrganization.OrgID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `org_id` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("post_auth_role_grants", federatedSettingsConnectedOrganization.PostAuthRoleGrants); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `post_auth_role_grants` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("role_mappings", flattenRoleMappings(federatedSettingsConnectedOrganization.RoleMappings)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `role_mappings` for federatedSettings IdentityProviders: %s", err))
	}
	if federatedSettingsConnectedOrganization.UserConflicts == nil {
		if err := d.Set("user_conflicts", federatedSettingsConnectedOrganization.UserConflicts); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `user_conflicts` for federatedSettings IdentityProviders: %s", err))
		}
	} else {
		if err := d.Set("user_conflicts", flattenUserConflicts(*federatedSettingsConnectedOrganization.UserConflicts)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `user_conflicts` for federatedSettings IdentityProviders: %s", err))
		}
	}

	d.SetId(federatedSettingsConnectedOrganization.OrgID)

	return nil
}
