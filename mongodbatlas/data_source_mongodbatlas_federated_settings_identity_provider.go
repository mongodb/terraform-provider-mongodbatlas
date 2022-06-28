package mongodbatlas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasFederatedSettingsIdentityProvider() *schema.Resource {
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
		},
	}
}
func dataSourceMongoDBAtlasFederatedSettingsIdentityProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}

	idpID, idpIDOk := d.GetOk("identity_provider_id")

	if !idpIDOk {
		return diag.FromErr(errors.New("identity_provider_id must be configured"))
	}

	federatedSettingsIdentityProvider, _, err := conn.FederatedSettings.GetIdentityProvider(ctx, federationSettingsID.(string), idpID.(string))
	if err != nil {
		return diag.Errorf("error getting federatedSettings IdentityProviders assigned (%s): %s", federationSettingsID, err)
	}

	if err := d.Set("acs_url", federatedSettingsIdentityProvider.AcsURL); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `acs_url` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("associated_domains", federatedSettingsIdentityProvider.AssociatedDomains); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `associated_domains` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("associated_orgs", flattenAssociatedOrgs(federatedSettingsIdentityProvider.AssociatedOrgs)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `associated_orgs` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("display_name", federatedSettingsIdentityProvider.DisplayName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `display_name` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("issuer_uri", federatedSettingsIdentityProvider.IssuerURI); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `issuer_uri` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("okta_idp_id", federatedSettingsIdentityProvider.OktaIdpID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `idp_id` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("pem_file_info", flattenPemFileInfo(*federatedSettingsIdentityProvider.PemFileInfo)); err != nil {
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

	if err := d.Set("sso_url", federatedSettingsIdentityProvider.SsoURL); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `sso_url` for federatedSettings IdentityProviders: %s", err))
	}

	if err := d.Set("status", federatedSettingsIdentityProvider.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `status` for federatedSettings IdentityProviders: %s", err))
	}

	d.SetId(federatedSettingsIdentityProvider.OktaIdpID)

	return nil
}
