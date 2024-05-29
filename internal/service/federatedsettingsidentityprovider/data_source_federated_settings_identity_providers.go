package federatedsettingsidentityprovider

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115013/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const WORKFORCE = "WORKFORCE"

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasFederatedSettingsIdentityProvidersRead,
		Schema: map[string]*schema.Schema{
			"federation_settings_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page_num": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: fmt.Sprintf(constant.DeprecationParamByVersion, "1.18.0"),
			},
			"items_per_page": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: fmt.Sprintf(constant.DeprecationParamByVersion, "1.18.0"),
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

	// once the SDK is upgraded above version "go.mongodb.org/atlas-sdk/v20231115012/mockadmin" we can use pagination parameters to iterate over all results (and adjust documentation)
	// pagination attributes are deprecated and can be removed as we move towards not exposing these pagination options to the user
	params := &admin.ListIdentityProvidersApiParams{
		FederationSettingsId: federationSettingsID.(string),
		Protocol:             &[]string{OIDC, SAML},
		IdpType:              &[]string{WORKFORCE},
	}

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
