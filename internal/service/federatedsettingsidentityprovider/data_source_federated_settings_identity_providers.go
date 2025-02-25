package federatedsettingsidentityprovider

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
		Schema: map[string]*schema.Schema{
			"federation_settings_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"idp_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"protocols": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
						"audience": {
							Type:     schema.TypeString,
							Computed: true,
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"authorization_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idp_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	federationSettingsID, federationSettingsIDOk := d.GetOk("federation_settings_id")

	if !federationSettingsIDOk {
		return diag.FromErr(errors.New("federation_settings_id must be configured"))
	}
	idpTypes := conversion.ExpandStringList(d.Get("idp_types").([]any))
	protocols := conversion.ExpandStringList(d.Get("protocols").([]any))

	params := &admin.ListIdentityProvidersApiParams{
		FederationSettingsId: federationSettingsID.(string),
		Protocol:             &protocols,
		IdpType:              &idpTypes,
	}

	// iterating all results to be implemented as part of CLOUDP-227485
	providers, _, err := connV2.FederatedAuthenticationApi.ListIdentityProvidersWithParams(ctx, params).Execute()
	if err != nil {
		return diag.Errorf("error getting federatedSettings Identity Providers assigned (%s): %s", federationSettingsID, err)
	}

	if err := d.Set("results", FlattenFederatedSettingsIdentityProvider(providers.GetResults())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for federatedSettings IdentityProviders: %s", err))
	}

	d.SetId(federationSettingsID.(string))
	return nil
}
