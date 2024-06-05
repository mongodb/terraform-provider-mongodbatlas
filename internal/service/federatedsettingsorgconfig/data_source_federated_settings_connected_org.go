package federatedsettingsorgconfig

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrgRead,
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
				Elem:     userConflictsElemSchema(),
			},
			"data_access_identity_provider_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func userConflictsElemSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func dataSourceOrgRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}

	orgID, orgIDOk := d.GetOk("org_id")

	if !orgIDOk {
		return diag.FromErr(errors.New("org_id must be configured"))
	}

	federatedSettingsConnectedOrganization, _, err := conn.FederatedAuthenticationApi.GetConnectedOrgConfig(ctx, federationSettingsID.(string), orgID.(string)).Execute()

	if err != nil {
		return diag.Errorf("error getting federatedSettings connected organizations assigned (%s): %s", federationSettingsID, err)
	}

	if err := d.Set("domain_allow_list", federatedSettingsConnectedOrganization.GetDomainAllowList()); err != nil {
		return diag.Errorf("error setting `domain_allow_list` for federatedSettings Org(%s): %s", orgID, err)
	}

	if err := d.Set("domain_restriction_enabled", federatedSettingsConnectedOrganization.GetDomainRestrictionEnabled()); err != nil {
		return diag.Errorf("error setting `domain_restriction_enabled` for federatedSettings Org(%s): %s", orgID, err)
	}

	if err := d.Set("identity_provider_id", federatedSettingsConnectedOrganization.GetIdentityProviderId()); err != nil {
		return diag.Errorf("error setting `identity_provider_id` for federatedSettings Org(%s): %s", orgID, err)
	}

	if err := d.Set("org_id", federatedSettingsConnectedOrganization.GetOrgId()); err != nil {
		return diag.Errorf("error setting `org_id` for federatedSettings Org(%s): %s", orgID, err)
	}

	if err := d.Set("post_auth_role_grants", federatedSettingsConnectedOrganization.GetPostAuthRoleGrants()); err != nil {
		return diag.Errorf("error setting `post_auth_role_grants` for federatedSettings Org(%s): %s", orgID, err)
	}

	if err := d.Set("role_mappings", FlattenRoleMappings(federatedSettingsConnectedOrganization.GetRoleMappings())); err != nil {
		return diag.Errorf("error setting `role_mappings` for federatedSettings Org(%s): %s", orgID, err)
	}
	if federatedSettingsConnectedOrganization.UserConflicts == nil {
		if err := d.Set("user_conflicts", federatedSettingsConnectedOrganization.GetUserConflicts()); err != nil {
			return diag.Errorf("error setting `user_conflicts` for federatedSettings Org(%s): %s", orgID, err)
		}
	} else {
		if err := d.Set("user_conflicts", FlattenUserConflicts(federatedSettingsConnectedOrganization.GetUserConflicts())); err != nil {
			return diag.Errorf("error setting `user_conflicts` for federatedSettings Org(%s): %s", orgID, err)
		}
	}
	if err := d.Set("data_access_identity_provider_ids", federatedSettingsConnectedOrganization.GetDataAccessIdentityProviderIds()); err != nil {
		return diag.Errorf("error setting `data_access_identity_provider_ids` for federatedSettings Org(%s): %s", orgID, err)
	}

	d.SetId(federatedSettingsConnectedOrganization.GetOrgId())

	return nil
}
