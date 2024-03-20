package federatedsettingsidentityprovider

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	admin20231001002 "go.mongodb.org/atlas-sdk/v20231115008/admin"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasFederatedSettingsIdentityProviderRead,
		Schema: map[string]*schema.Schema{
			"federation_settings_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identity_provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},

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
			"idp_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"audience_claim": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups_claim": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"requested_scopes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_claim": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceMongoDBAtlasFederatedSettingsIdentityProviderRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}

	idpID, idpIDOk := d.GetOk("identity_provider_id")

	if !idpIDOk {
		return diag.FromErr(errors.New("identity_provider_id must be configured"))
	}

	// to be removed in terraform-provider-1.16.0
	if len(idpID.(string)) == 20 {
		return append(oldSDKDSRead(ctx, federationSettingsID.(string), idpID.(string), d, meta), getGracePeriodWarning())
	}

	federatedSettingsIdentityProvider, _, err := connV2.FederatedAuthenticationApi.GetIdentityProvider(ctx, federationSettingsID.(string), idpID.(string)).Execute()
	if err != nil {
		return diag.Errorf("error getting federatedSettings IdentityProviders assigned (%s): %s", federationSettingsID, err)
	}

	if federatedSettingsIdentityProvider.GetProtocol() == SAML {
		if err := d.Set("acs_url", federatedSettingsIdentityProvider.AcsUrl); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `acs_url` for federatedSettings IdentityProviders: %s", err))
		}

		if err := d.Set("pem_file_info", FlattenPemFileInfo(*federatedSettingsIdentityProvider.PemFileInfo)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `pem_file_info` for federatedSettings IdentityProviders: %s", err))
		}

		if err := d.Set("request_binding", federatedSettingsIdentityProvider.RequestBinding); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `request_binding` for federatedSettings IdentityProviders: %s", err))
		}

		if err := d.Set("response_signature_algorithm", federatedSettingsIdentityProvider.ResponseSignatureAlgorithm); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `response_signature_algorithm` for federatedSettings IdentityProviders: %s", err))
		}

		if err := d.Set("sso_debug_enabled", federatedSettingsIdentityProvider.SsoDebugEnabled); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `sso_debug_enabled` for federatedSettings IdentityProviders: %s", err))
		}

		if err := d.Set("sso_url", federatedSettingsIdentityProvider.SsoUrl); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `sso_url` for federatedSettings IdentityProviders: %s", err))
		}

		if err := d.Set("status", federatedSettingsIdentityProvider.Status); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `status` for federatedSettings IdentityProviders: %s", err))
		}
	}

	if federatedSettingsIdentityProvider.GetProtocol() == OIDC {
		if err := d.Set("audience_claim", federatedSettingsIdentityProvider.AudienceClaim); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `audience_claim` for federatedSettings IdentityProviders: %s", err))
		}

		if err := d.Set("client_id", federatedSettingsIdentityProvider.ClientId); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `client_id` for federatedSettings IdentityProviders: %s", err))
		}

		if err := d.Set("groups_claim", federatedSettingsIdentityProvider.GroupsClaim); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `groups_claim` for federatedSettings IdentityProviders: %s", err))
		}

		if err := d.Set("requested_scopes", federatedSettingsIdentityProvider.RequestedScopes); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `associated_domains` for federatedSettings IdentityProviders: %s", err))
		}

		if err := d.Set("user_claim", federatedSettingsIdentityProvider.UserClaim); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `user_claim` for federatedSettings IdentityProviders: %s", err))
		}
	}

	if err := d.Set("associated_domains", federatedSettingsIdentityProvider.AssociatedDomains); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `associated_domains` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("associated_orgs", FlattenAssociatedOrgs(federatedSettingsIdentityProvider.GetAssociatedOrgs())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `associated_orgs` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("display_name", federatedSettingsIdentityProvider.DisplayName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `display_name` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("issuer_uri", federatedSettingsIdentityProvider.IssuerUri); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `issuer_uri` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("protocol", federatedSettingsIdentityProvider.Protocol); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `protocol` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("okta_idp_id", federatedSettingsIdentityProvider.OktaIdpId); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `okta_idp_id` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("idp_id", federatedSettingsIdentityProvider.Id); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `idp_id` for federatedSettings IdentityProviders: %s", err))
	}

	d.SetId(federatedSettingsIdentityProvider.Id)

	return nil
}

func oldSDKDSRead(ctx context.Context, federationSettingsID, idpID string, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn20231001002 := meta.(*config.MongoDBClient).Atlas20231001002
	federatedSettingsIdentityProvider, _, err := conn20231001002.FederatedAuthenticationApi.GetIdentityProvider(ctx, federationSettingsID, idpID).Execute()
	if err != nil {
		return diag.Errorf("error getting federatedSettings IdentityProviders assigned (%s): %s", federationSettingsID, err)
	}

	if err := d.Set("acs_url", federatedSettingsIdentityProvider.AcsUrl); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `acs_url` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("associated_domains", federatedSettingsIdentityProvider.AssociatedDomains); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `associated_domains` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("associated_orgs", oldSDKFlattenAssociatedOrgs(federatedSettingsIdentityProvider.AssociatedOrgs)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `associated_orgs` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("display_name", federatedSettingsIdentityProvider.DisplayName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `display_name` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("issuer_uri", federatedSettingsIdentityProvider.IssuerUri); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `issuer_uri` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("okta_idp_id", federatedSettingsIdentityProvider.OktaIdpId); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `idp_id` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("pem_file_info", oldSDKFlattenPemFileInfo(*federatedSettingsIdentityProvider.PemFileInfo)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `pem_file_info` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("request_binding", federatedSettingsIdentityProvider.RequestBinding); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `request_binding` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("response_signature_algorithm", federatedSettingsIdentityProvider.ResponseSignatureAlgorithm); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `response_signature_algorithm` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("sso_debug_enabled", federatedSettingsIdentityProvider.SsoDebugEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `sso_debug_enabled` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("sso_url", federatedSettingsIdentityProvider.SsoUrl); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `sso_url` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("status", federatedSettingsIdentityProvider.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `status` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("idp_id", federatedSettingsIdentityProvider.Id); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `idp_id` for federatedSettings IdentityProviders: %s", err))
	}

	d.SetId(federatedSettingsIdentityProvider.OktaIdpId)

	return nil
}

func oldSDKFlattenAssociatedOrgs(associatedOrgs []admin20231001002.ConnectedOrgConfig) []map[string]any {
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
				"identity_provider_id":       associatedOrgs[i].IdentityProviderId,
				"org_id":                     associatedOrgs[i].OrgId,
				"post_auth_role_grants":      associatedOrgs[i].PostAuthRoleGrants,
				"role_mappings":              oldSDKFlattenAuthFederationRoleMapping(associatedOrgs[i].RoleMappings),
				"user_conflicts":             nil,
			}
		} else {
			associatedOrgsMap[i] = map[string]any{
				"domain_allow_list":          associatedOrgs[i].DomainAllowList,
				"domain_restriction_enabled": associatedOrgs[i].DomainRestrictionEnabled,
				"identity_provider_id":       associatedOrgs[i].IdentityProviderId,
				"org_id":                     associatedOrgs[i].OrgId,
				"post_auth_role_grants":      associatedOrgs[i].PostAuthRoleGrants,
				"role_mappings":              oldSDKFlattenAuthFederationRoleMapping(associatedOrgs[i].RoleMappings),
				"user_conflicts":             oldSDKFlattenFederatedUser(associatedOrgs[i].UserConflicts),
			}
		}
	}

	return associatedOrgsMap
}

func oldSDKFlattenFederatedUser(federatedUsers []admin20231001002.FederatedUser) []map[string]any {
	var userConflictsMap []map[string]any

	if len(federatedUsers) == 0 {
		return nil
	}
	userConflictsMap = make([]map[string]any, len(federatedUsers))

	for i := range federatedUsers {
		userConflictsMap[i] = map[string]any{
			"email_address":          federatedUsers[i].EmailAddress,
			"federation_settings_id": federatedUsers[i].FederationSettingsId,
			"first_name":             federatedUsers[i].FirstName,
			"last_name":              federatedUsers[i].LastName,
			"user_id":                federatedUsers[i].UserId,
		}
	}

	return userConflictsMap
}

type oldSDKAuthFederationoleMappingsByGroupName []admin20231001002.AuthFederationRoleMapping

func (ra oldSDKAuthFederationoleMappingsByGroupName) Len() int      { return len(ra) }
func (ra oldSDKAuthFederationoleMappingsByGroupName) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }

func (ra oldSDKAuthFederationoleMappingsByGroupName) Less(i, j int) bool {
	return ra[i].ExternalGroupName < ra[j].ExternalGroupName
}

func oldSDKFlattenAuthFederationRoleMapping(roleMappings []admin20231001002.AuthFederationRoleMapping) []map[string]any {
	sort.Sort(oldSDKAuthFederationoleMappingsByGroupName(roleMappings))

	var roleMappingsMap []map[string]any

	if len(roleMappings) > 0 {
		roleMappingsMap = make([]map[string]any, len(roleMappings))

		for i := range roleMappings {
			roleMappingsMap[i] = map[string]any{
				"external_group_name": roleMappings[i].ExternalGroupName,
				"id":                  roleMappings[i].Id,
				"role_assignments":    oldSDKFlattenRoleAssignments(roleMappings[i].RoleAssignments),
			}
		}
	}

	return roleMappingsMap
}

type mRoleAssignmentOldV2 []admin20231001002.RoleAssignment

func (ra mRoleAssignmentOldV2) Len() int      { return len(ra) }
func (ra mRoleAssignmentOldV2) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }
func (ra mRoleAssignmentOldV2) Less(i, j int) bool {
	compareVal := strings.Compare(ra[i].GetOrgId(), ra[j].GetOrgId())

	if compareVal != 0 {
		return compareVal < 0
	}

	compareVal = strings.Compare(ra[i].GetGroupId(), ra[j].GetGroupId())

	if compareVal != 0 {
		return compareVal < 0
	}

	return *ra[i].Role < *ra[j].Role
}

func oldSDKFlattenRoleAssignments(roleAssignments []admin20231001002.RoleAssignment) []map[string]any {
	sort.Sort(mRoleAssignmentOldV2(roleAssignments))

	var roleAssignmentsMap []map[string]any

	if len(roleAssignments) > 0 {
		roleAssignmentsMap = make([]map[string]any, len(roleAssignments))

		for i := range roleAssignments {
			roleAssignmentsMap[i] = map[string]any{
				"group_id": roleAssignments[i].GroupId,
				"org_id":   roleAssignments[i].OrgId,
				"role":     roleAssignments[i].Role,
			}
		}
	}

	return roleAssignmentsMap
}

func oldSDKFlattenPemFileInfo(pemFileInfo admin20231001002.PemFileInfo) []map[string]any {
	var pemFileInfoMap []map[string]any

	if len(pemFileInfo.Certificates) > 0 {
		pemFileInfoMap = make([]map[string]any, 1)

		pemFileInfoMap[0] = map[string]any{
			"certificates": oldSDKFlattenFederatedSettingsCertificates(pemFileInfo.Certificates),
			"file_name":    pemFileInfo.FileName,
		}
	}

	return pemFileInfoMap
}

func oldSDKFlattenFederatedSettingsCertificates(certificates []admin20231001002.X509Certificate) []map[string]any {
	var certificatesMap []map[string]any

	if len(certificates) > 0 {
		certificatesMap = make([]map[string]any, len(certificates))

		for i := range certificates {
			certificatesMap[i] = map[string]any{
				"not_after":  conversion.TimePtrToStringPtr(certificates[i].NotAfter),
				"not_before": conversion.TimePtrToStringPtr(certificates[i].NotBefore),
			}
		}
	}

	return certificatesMap
}
