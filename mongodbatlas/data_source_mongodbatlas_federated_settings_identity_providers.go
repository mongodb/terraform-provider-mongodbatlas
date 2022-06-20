package mongodbatlas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasFederatedSettingsIdentityProviders() *schema.Resource {
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
func dataSourceMongoDBAtlasFederatedSettingsIdentityProvidersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("Federation_settings_id must be configured"))
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

func flattenFederatedSettingsIdentityProvider(federatedSettingsIdentityProvider []matlas.FederatedSettingsIdentityProvider) []map[string]interface{} {
	var federatedSettingsIdentityProviderMap []map[string]interface{}

	if len(federatedSettingsIdentityProvider) > 0 {
		federatedSettingsIdentityProviderMap = make([]map[string]interface{}, len(federatedSettingsIdentityProvider))

		for i := range federatedSettingsIdentityProvider {
			federatedSettingsIdentityProviderMap[i] = map[string]interface{}{
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

func flattenAssociatedOrgs(associatedOrgs []*matlas.AssociatedOrgs) []map[string]interface{} {
	var associatedOrgsMap []map[string]interface{}

	if len(associatedOrgs) == 0 {
		return nil
	}
	associatedOrgsMap = make([]map[string]interface{}, len(associatedOrgs))

	for i := range associatedOrgs {
		if associatedOrgs[i].UserConflicts == nil {
			associatedOrgsMap[i] = map[string]interface{}{
				"domain_allow_list":          associatedOrgs[i].DomainAllowList,
				"domain_restriction_enabled": associatedOrgs[i].DomainRestrictionEnabled,
				"identity_provider_id":       associatedOrgs[i].IdentityProviderID,
				"org_id":                     associatedOrgs[i].OrgID,
				"post_auth_role_grants":      associatedOrgs[i].PostAuthRoleGrants,
				"role_mappings":              flattenRoleMappings(associatedOrgs[i].RoleMappings),
				"user_conflicts":             nil,
			}
		} else {
			associatedOrgsMap[i] = map[string]interface{}{
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

func flattenUserConflicts(userConflicts matlas.UserConflicts) []map[string]interface{} {
	var userConflictsMap []map[string]interface{}

	if len(userConflicts) == 0 {
		return nil
	}
	userConflictsMap = make([]map[string]interface{}, len(userConflicts))

	for i := range userConflicts {
		userConflictsMap[i] = map[string]interface{}{
			"email_address":          userConflicts[i].EmailAddress,
			"federation_settings_id": userConflicts[i].FederationSettingsID,
			"first_name":             userConflicts[i].FirstName,
			"last_name":              userConflicts[i].LastName,
			"user_id":                userConflicts[i].UserID,
		}
	}

	return userConflictsMap
}

func flattenPemFileInfo(pemFileInfo matlas.PemFileInfo) []map[string]interface{} {
	var pemFileInfoMap []map[string]interface{}

	if len(pemFileInfo.Certificates) > 0 {
		pemFileInfoMap = make([]map[string]interface{}, 1)

		pemFileInfoMap[0] = map[string]interface{}{
			"certificates": flattenFederatedSettingsCertificates(pemFileInfo.Certificates),
			"file_name":    pemFileInfo.FileName,
		}
	}

	return pemFileInfoMap
}

func flattenFederatedSettingsCertificates(certificates []*matlas.Certificates) []map[string]interface{} {
	var certificatesMap []map[string]interface{}

	if len(certificates) > 0 {
		certificatesMap = make([]map[string]interface{}, len(certificates))

		for i := range certificates {
			certificatesMap[i] = map[string]interface{}{
				"not_after":  certificates[i].NotAfter.String(),
				"not_before": certificates[i].NotBefore.String(),
			}
		}
	}

	return certificatesMap
}

func flattenRoleMappings(roleMappings []*matlas.RoleMappings) []map[string]interface{} {
	var roleMappingsMap []map[string]interface{}

	if len(roleMappings) > 0 {
		roleMappingsMap = make([]map[string]interface{}, len(roleMappings))

		for i := range roleMappings {
			roleMappingsMap[i] = map[string]interface{}{
				"external_group_name": roleMappings[i].ExternalGroupName,
				"id":                  roleMappings[i].ID,
				"role_assignments":    flattenRoleAssignments(roleMappings[i].RoleAssignments),
			}
		}
	}

	return roleMappingsMap
}

func flattenRoleAssignments(roleAssignments []*matlas.RoleAssignments) []map[string]interface{} {
	var roleAssignmentsMap []map[string]interface{}

	if len(roleAssignments) > 0 {
		roleAssignmentsMap = make([]map[string]interface{}, len(roleAssignments))

		for i := range roleAssignments {
			roleAssignmentsMap[i] = map[string]interface{}{
				"group_id": roleAssignments[i].GroupID,
				"org_id":   roleAssignments[i].OrgID,
				"role":     roleAssignments[i].Role,
			}
		}
	}

	return roleAssignmentsMap
}
