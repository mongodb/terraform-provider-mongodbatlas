package federatedsettingsorgconfig

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312004/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
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
						"data_access_identity_provider_ids": {
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
					},
				},
			},
		},
	}
}
func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("either federation_settings_id must be configured"))
	}

	params := &admin.ListConnectedOrgConfigsApiParams{
		FederationSettingsId: federationSettingsID.(string),
		PageNum:              conversion.Pointer(d.Get("page_num").(int)),
		ItemsPerPage:         conversion.Pointer(d.Get("items_per_page").(int)),
	}

	federatedSettingsConnectedOrganizations, _, err := conn.FederatedAuthenticationApi.ListConnectedOrgConfigsWithParams(ctx, params).Execute()
	if err != nil {
		return diag.Errorf("error getting federatedSettings connected organizations assigned (%s): %s", federationSettingsID, err)
	}

	if err := d.Set("results", flattenOrganizationConfigs(*federatedSettingsConnectedOrganizations)); err != nil {
		return diag.Errorf("error setting `result` for federatedSettings connected orgs: %s", err)
	}

	d.SetId(federationSettingsID.(string))

	return nil
}

func flattenOrganizationConfigs(federatedSettingsConnectedOrganizations admin.PaginatedConnectedOrgConfigs) []map[string]any {
	var federatedSettingsConnectedOrganizationsMap []map[string]any

	if (federatedSettingsConnectedOrganizations.GetTotalCount()) > 0 {
		federatedSettingsConnectedOrganizationsMap = make([]map[string]any, federatedSettingsConnectedOrganizations.GetTotalCount())

		for i := range federatedSettingsConnectedOrganizations.GetResults() {
			federatedSettingsConnectedOrganizationsMap[i] = map[string]any{
				"domain_allow_list":                 federatedSettingsConnectedOrganizations.GetResults()[i].GetDomainAllowList(),
				"domain_restriction_enabled":        federatedSettingsConnectedOrganizations.GetResults()[i].GetDomainRestrictionEnabled(),
				"identity_provider_id":              federatedSettingsConnectedOrganizations.GetResults()[i].GetIdentityProviderId(),
				"org_id":                            federatedSettingsConnectedOrganizations.GetResults()[i].GetOrgId(),
				"post_auth_role_grants":             federatedSettingsConnectedOrganizations.GetResults()[i].GetPostAuthRoleGrants(),
				"role_mappings":                     FlattenRoleMappings(federatedSettingsConnectedOrganizations.GetResults()[i].GetRoleMappings()),
				"data_access_identity_provider_ids": federatedSettingsConnectedOrganizations.GetResults()[i].GetDataAccessIdentityProviderIds(),
				"user_conflicts":                    FlattenUserConflicts(federatedSettingsConnectedOrganizations.GetResults()[i].GetUserConflicts()),
			}
		}
	}

	return federatedSettingsConnectedOrganizationsMap
}
