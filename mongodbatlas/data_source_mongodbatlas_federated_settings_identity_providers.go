package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func DataSourceMongoDBAtlasFederatedSettingsIdentityProviders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasFederatedSettingsIdentityProvidersRead,
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
						"acs_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associated_domains": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"associated_orgs": {
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
						"audience_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issuer_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"okta_idp_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pem_file_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificates": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"not_after": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"not_before": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"file_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"request_binding": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_signature_algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sso_debug_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sso_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceMongoDBAtlasFederatedSettingsIdentityProvidersRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}

	federatedSettingsIdentityProviders, _, err := conn.FederatedSettings.ListIdentityProviders(ctx, federationSettingsID.(string), options)
	if err != nil {
		return diag.Errorf("error getting federatedSettings IdentityProviders assigned (%s): %s", federationSettingsID, err)
	}

	if err := d.Set("results", flattenFederatedSettingsIdentityProvider(federatedSettingsIdentityProviders)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for federatedSettings IdentityProviders: %s", err))
	}

	d.SetId(federationSettingsID.(string))

	return nil
}

func flattenFederatedSettingsIdentityProvider(federatedSettingsIdentityProvider []matlas.FederatedSettingsIdentityProvider) []map[string]any {
	var federatedSettingsIdentityProviderMap []map[string]any

	if len(federatedSettingsIdentityProvider) > 0 {
		federatedSettingsIdentityProviderMap = make([]map[string]any, len(federatedSettingsIdentityProvider))

		for i := range federatedSettingsIdentityProvider {
			federatedSettingsIdentityProviderMap[i] = map[string]any{
				"acs_url":                      federatedSettingsIdentityProvider[i].AcsURL,
				"associated_domains":           federatedSettingsIdentityProvider[i].AssociatedDomains,
				"associated_orgs":              flattenAssociatedOrgs(federatedSettingsIdentityProvider[i].AssociatedOrgs),
				"audience_uri":                 federatedSettingsIdentityProvider[i].AudienceURI,
				"display_name":                 federatedSettingsIdentityProvider[i].DisplayName,
				"issuer_uri":                   federatedSettingsIdentityProvider[i].IssuerURI,
				"okta_idp_id":                  federatedSettingsIdentityProvider[i].OktaIdpID,
				"pem_file_info":                flattenPemFileInfo(*federatedSettingsIdentityProvider[i].PemFileInfo),
				"request_binding":              federatedSettingsIdentityProvider[i].RequestBinding,
				"response_signature_algorithm": federatedSettingsIdentityProvider[i].ResponseSignatureAlgorithm,
				"sso_debug_enabled":            federatedSettingsIdentityProvider[i].SsoDebugEnabled,
				"sso_url":                      federatedSettingsIdentityProvider[i].SsoURL,
				"status":                       federatedSettingsIdentityProvider[i].Status,
			}
		}
	}

	return federatedSettingsIdentityProviderMap
}

func flattenAssociatedOrgs(associatedOrgs []*matlas.AssociatedOrgs) []map[string]any {
	var associatedOrgsMap []map[string]any

	if len(associatedOrgs) == 0 {
		return nil
	}
	associatedOrgsMap = make([]map[string]any, len(associatedOrgs))

	for i := range associatedOrgs {
		if associatedOrgs[i].UserConflicts == nil {
			associatedOrgsMap[i] = map[string]any{
				"domain_allow_list":          associatedOrgs[i].DomainAllowList,
				"domain_restriction_enabled": associatedOrgs[i].DomainRestrictionEnabled,
				"identity_provider_id":       associatedOrgs[i].IdentityProviderID,
				"org_id":                     associatedOrgs[i].OrgID,
				"post_auth_role_grants":      associatedOrgs[i].PostAuthRoleGrants,
				"role_mappings":              flattenRoleMappings(associatedOrgs[i].RoleMappings),
				"user_conflicts":             nil,
			}
		} else {
			associatedOrgsMap[i] = map[string]any{
				"domain_allow_list":          associatedOrgs[i].DomainAllowList,
				"domain_restriction_enabled": associatedOrgs[i].DomainRestrictionEnabled,
				"identity_provider_id":       associatedOrgs[i].IdentityProviderID,
				"org_id":                     associatedOrgs[i].OrgID,
				"post_auth_role_grants":      associatedOrgs[i].PostAuthRoleGrants,
				"role_mappings":              flattenRoleMappings(associatedOrgs[i].RoleMappings),
				"user_conflicts":             flattenUserConflicts(*associatedOrgs[i].UserConflicts),
			}
		}
	}

	return associatedOrgsMap
}

func flattenUserConflicts(userConflicts matlas.UserConflicts) []map[string]any {
	var userConflictsMap []map[string]any

	if len(userConflicts) == 0 {
		return nil
	}
	userConflictsMap = make([]map[string]any, len(userConflicts))

	for i := range userConflicts {
		userConflictsMap[i] = map[string]any{
			"email_address":          userConflicts[i].EmailAddress,
			"federation_settings_id": userConflicts[i].FederationSettingsID,
			"first_name":             userConflicts[i].FirstName,
			"last_name":              userConflicts[i].LastName,
			"user_id":                userConflicts[i].UserID,
		}
	}

	return userConflictsMap
}

func flattenPemFileInfo(pemFileInfo matlas.PemFileInfo) []map[string]any {
	var pemFileInfoMap []map[string]any

	if len(pemFileInfo.Certificates) > 0 {
		pemFileInfoMap = make([]map[string]any, 1)

		pemFileInfoMap[0] = map[string]any{
			"certificates": flattenFederatedSettingsCertificates(pemFileInfo.Certificates),
			"file_name":    pemFileInfo.FileName,
		}
	}

	return pemFileInfoMap
}

func flattenFederatedSettingsCertificates(certificates []*matlas.Certificates) []map[string]any {
	var certificatesMap []map[string]any

	if len(certificates) > 0 {
		certificatesMap = make([]map[string]any, len(certificates))

		for i := range certificates {
			certificatesMap[i] = map[string]any{
				"not_after":  certificates[i].NotAfter.String(),
				"not_before": certificates[i].NotBefore.String(),
			}
		}
	}

	return certificatesMap
}

type mRoleAssignment []*matlas.RoleAssignments

func (ra mRoleAssignment) Len() int      { return len(ra) }
func (ra mRoleAssignment) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }
func (ra mRoleAssignment) Less(i, j int) bool {
	compareVal := strings.Compare(ra[i].OrgID, ra[j].OrgID)

	if compareVal != 0 {
		return compareVal < 0
	}

	compareVal = strings.Compare(ra[i].GroupID, ra[j].GroupID)

	if compareVal != 0 {
		return compareVal < 0
	}

	return ra[i].Role < ra[j].Role
}

type roleMappingsByGroupName []*matlas.RoleMappings

func (ra roleMappingsByGroupName) Len() int      { return len(ra) }
func (ra roleMappingsByGroupName) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }

func (ra roleMappingsByGroupName) Less(i, j int) bool {
	return ra[i].ExternalGroupName < ra[j].ExternalGroupName
}

func flattenRoleMappings(roleMappings []*matlas.RoleMappings) []map[string]any {
	sort.Sort(roleMappingsByGroupName(roleMappings))

	var roleMappingsMap []map[string]any

	if len(roleMappings) > 0 {
		roleMappingsMap = make([]map[string]any, len(roleMappings))

		for i := range roleMappings {
			roleMappingsMap[i] = map[string]any{
				"external_group_name": roleMappings[i].ExternalGroupName,
				"id":                  roleMappings[i].ID,
				"role_assignments":    flattenRoleAssignments(roleMappings[i].RoleAssignments),
			}
		}
	}

	return roleMappingsMap
}

func flattenRoleAssignments(roleAssignments []*matlas.RoleAssignments) []map[string]any {
	sort.Sort(mRoleAssignment(roleAssignments))

	var roleAssignmentsMap []map[string]any

	if len(roleAssignments) > 0 {
		roleAssignmentsMap = make([]map[string]any, len(roleAssignments))

		for i := range roleAssignments {
			roleAssignmentsMap[i] = map[string]any{
				"group_id": roleAssignments[i].GroupID,
				"org_id":   roleAssignments[i].OrgID,
				"role":     roleAssignments[i].Role,
			}
		}
	}

	return roleAssignmentsMap
}
