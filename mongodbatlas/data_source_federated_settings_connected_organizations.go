package mongodbatlas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/federatedsettingsidentityprovider"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func DataSourceFederatedSettingsOrganizationConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasFederatedSettingsOrganizationConfigsRead,
		Schema: map[string]*schema.Schema{
			"federation_settings_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"items_per_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"org_id": {
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
				},
			},
		},
	}
}
func dataSourceMongoDBAtlasFederatedSettingsOrganizationConfigsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("either federation_settings_id must be configured"))
	}

	federatedSettingsConnectedOrganizations, _, err := conn.FederatedSettings.ListConnectedOrgs(ctx, federationSettingsID.(string), options)
	if err != nil {
		return diag.Errorf("error getting federatedSettings connected organizations assigned (%s): %s", federationSettingsID, err)
	}

	if err := d.Set("results", flattenFederatedSettingsOrganizationConfigs(*federatedSettingsConnectedOrganizations)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for federatedSettings IdentityProviders: %s", err))
	}

	d.SetId(federationSettingsID.(string))

	return nil
}

func flattenFederatedSettingsOrganizationConfigs(federatedSettingsConnectedOrganizations matlas.FederatedSettingsConnectedOrganizations) []map[string]any {
	var federatedSettingsConnectedOrganizationsMap []map[string]any

	if (federatedSettingsConnectedOrganizations.TotalCount) > 0 {
		federatedSettingsConnectedOrganizationsMap = make([]map[string]any, federatedSettingsConnectedOrganizations.TotalCount)

		for i := range federatedSettingsConnectedOrganizations.Results {
			if federatedSettingsConnectedOrganizations.Results[i].UserConflicts == nil {
				federatedSettingsConnectedOrganizationsMap[i] = map[string]any{
					"domain_allow_list":          federatedSettingsConnectedOrganizations.Results[i].DomainAllowList,
					"domain_restriction_enabled": federatedSettingsConnectedOrganizations.Results[i].DomainRestrictionEnabled,
					"identity_provider_id":       federatedSettingsConnectedOrganizations.Results[i].IdentityProviderID,
					"org_id":                     federatedSettingsConnectedOrganizations.Results[i].OrgID,
					"post_auth_role_grants":      federatedSettingsConnectedOrganizations.Results[i].PostAuthRoleGrants,
					"role_mappings":              federatedsettingsidentityprovider.FlattenRoleMappings(federatedSettingsConnectedOrganizations.Results[i].RoleMappings),
					"user_conflicts":             nil,
				}
			} else {
				federatedSettingsConnectedOrganizationsMap[i] = map[string]any{
					"domain_allow_list":          federatedSettingsConnectedOrganizations.Results[i].DomainAllowList,
					"domain_restriction_enabled": federatedSettingsConnectedOrganizations.Results[i].DomainRestrictionEnabled,
					"identity_provider_id":       federatedSettingsConnectedOrganizations.Results[i].IdentityProviderID,
					"org_id":                     federatedSettingsConnectedOrganizations.Results[i].OrgID,
					"post_auth_role_grants":      federatedSettingsConnectedOrganizations.Results[i].PostAuthRoleGrants,
					"role_mappings":              federatedsettingsidentityprovider.FlattenRoleMappings(federatedSettingsConnectedOrganizations.Results[i].RoleMappings),
					"user_conflicts":             federatedsettingsidentityprovider.FlattenUserConflicts(*federatedSettingsConnectedOrganizations.Results[i].UserConflicts),
				}
			}
		}
	}

	return federatedSettingsConnectedOrganizationsMap
}
