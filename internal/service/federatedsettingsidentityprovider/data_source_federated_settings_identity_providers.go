package federatedsettingsidentityprovider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20231115005/admin"
)

func PluralDataSource() *schema.Resource {
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
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idp_id": {
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
				},
			},
		},
	}
}
func dataSourceMongoDBAtlasFederatedSettingsIdentityProvidersRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}

	oidcParams := &admin.ListIdentityProvidersApiParams{
		FederationSettingsId: federationSettingsID.(string),
		Protocol:             conversion.StringPtr(OIDC),
	}
	samlParams := &admin.ListIdentityProvidersApiParams{
		FederationSettingsId: federationSettingsID.(string),
		Protocol:             conversion.StringPtr(SAML),
	}

	samlFederatedSettingsIdentityProviders, _, samlErr := connV2.FederatedAuthenticationApi.ListIdentityProvidersWithParams(ctx, samlParams).Execute()
	if samlErr != nil {
		return diag.Errorf("error getting federatedSettings IdentityProviders assigned (%s): %s", federationSettingsID, samlErr)
	}
	oidcFederatedSettingsIdentityProviders, _, oidcErr := connV2.FederatedAuthenticationApi.ListIdentityProvidersWithParams(ctx, oidcParams).Execute()
	if oidcErr != nil {
		return diag.Errorf("error getting federatedSettings IdentityProviders assigned (%s): %s", federationSettingsID, oidcErr)
	}
	allFederatedSettingsIdentityProviders := append(samlFederatedSettingsIdentityProviders.GetResults(), oidcFederatedSettingsIdentityProviders.GetResults()...)

	if err := d.Set("results", FlattenFederatedSettingsIdentityProvider(allFederatedSettingsIdentityProviders)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for federatedSettings IdentityProviders: %s", err))
	}

	d.SetId(federationSettingsID.(string))

	return nil
}
