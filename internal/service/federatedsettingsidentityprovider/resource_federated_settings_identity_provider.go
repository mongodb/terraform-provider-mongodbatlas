package federatedsettingsidentityprovider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const SAML = "SAML"
const OIDC = "OIDC"

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceMongoDBAtlasFederatedSettingsIdentityProviderRead,
		UpdateContext: resourceMongoDBAtlasFederatedSettingsIdentityProviderUpdate,
		DeleteContext: resourceMongoDBAtlasFederatedSettingsIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasFederatedSettingsIdentityProviderImportState,
		},
		Schema: map[string]*schema.Schema{
			"federation_settings_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_binding": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"response_signature_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"associated_domains": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sso_debug_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sso_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"okta_idp_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"idp_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"audience": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups_claim": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"requested_scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_claim": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if d.Get("protocol").(string) != OIDC {
		return diag.FromErr(errors.New("this resource must be imported"))
	}
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	createRequest := ExpandIdentityProviderOIDCCreate(d)
	federatedSettingsID := d.Get("federation_settings_id").(string)

	_, _, err := connV2.FederatedAuthenticationApi.CreateIdentityProvider(ctx, federatedSettingsID, createRequest).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating federation settings identity provider (%s): %s", federatedSettingsID, err))
	}
	return resourceMongoDBAtlasFederatedSettingsIdentityProviderRead(ctx, d, meta)
}

func resourceMongoDBAtlasFederatedSettingsIdentityProviderRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]

	// Since the migration of this resource to the latest version of the auto-generated SDK (v20231115) and the breaking changes of the API
	// the unique identifier used by the API & SDK of the supported identity providers is no longer "okta_idp_id", it is "idp_id". Nonetheless
	// "okta_idp_id" name was used to encode/decode the Terraform State Id. To ensure backwards compatibility, the format of this resource id remains the same but the key `okta_idp_id` will store either `okta_idp_id` or `idp_id` to identify the identity provider.
	// as few changes as possible, this name will remain.
	idpID := ids["okta_idp_id"]

	// latest version of v2 SDK
	federatedSettingsIdentityProvider, resp, err := connV2.FederatedAuthenticationApi.GetIdentityProvider(ctx, federationSettingsID, idpID).Execute()
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting federated settings identity provider: %s", err))
	}

	if federatedSettingsIdentityProvider.GetProtocol() == SAML {
		if err := d.Set("request_binding", federatedSettingsIdentityProvider.RequestBinding); err != nil {
			return diag.FromErr(fmt.Errorf("error setting request binding (%s): %s", d.Id(), err))
		}

		if err := d.Set("response_signature_algorithm", federatedSettingsIdentityProvider.ResponseSignatureAlgorithm); err != nil {
			return diag.FromErr(fmt.Errorf("error setting response signature algorithm (%s): %s", d.Id(), err))
		}

		if err := d.Set("sso_debug_enabled", federatedSettingsIdentityProvider.SsoDebugEnabled); err != nil {
			return diag.FromErr(fmt.Errorf("error setting sso debug enabled (%s): %s", d.Id(), err))
		}

		if err := d.Set("sso_url", federatedSettingsIdentityProvider.SsoUrl); err != nil {
			return diag.FromErr(fmt.Errorf("error setting sso url (%s): %s", d.Id(), err))
		}

		if err := d.Set("status", federatedSettingsIdentityProvider.Status); err != nil {
			return diag.FromErr(fmt.Errorf("error setting Status (%s): %s", d.Id(), err))
		}
	} else if federatedSettingsIdentityProvider.GetProtocol() == OIDC {
		if err := d.Set("audience", federatedSettingsIdentityProvider.Audience); err != nil {
			return diag.FromErr(fmt.Errorf("error setting audience claim list (%s): %s", d.Id(), err))
		}

		if err := d.Set("client_id", federatedSettingsIdentityProvider.ClientId); err != nil {
			return diag.FromErr(fmt.Errorf("error setting client id (%s): %s", d.Id(), err))
		}

		if err := d.Set("groups_claim", federatedSettingsIdentityProvider.GroupsClaim); err != nil {
			return diag.FromErr(fmt.Errorf("error setting groups claim (%s): %s", d.Id(), err))
		}

		if err := d.Set("requested_scopes", federatedSettingsIdentityProvider.RequestedScopes); err != nil {
			return diag.FromErr(fmt.Errorf("error setting requested scopes list (%s): %s", d.Id(), err))
		}

		if err := d.Set("user_claim", federatedSettingsIdentityProvider.UserClaim); err != nil {
			return diag.FromErr(fmt.Errorf("error setting user claim (%s): %s", d.Id(), err))
		}
	}

	if err := d.Set("federation_settings_id", federationSettingsID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting Identity Provider in Federation settings (%s): %s", d.Id(), err))
	}

	if err := d.Set("name", federatedSettingsIdentityProvider.DisplayName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name (%s): %s", d.Id(), err))
	}

	if err := d.Set("associated_domains", federatedSettingsIdentityProvider.AssociatedDomains); err != nil {
		return diag.FromErr(fmt.Errorf("error setting associated domains list (%s): %s", d.Id(), err))
	}

	if err := d.Set("okta_idp_id", federatedSettingsIdentityProvider.OktaIdpId); err != nil {
		return diag.FromErr(fmt.Errorf("error setting OktaIdpID (%s): %s", d.Id(), err))
	}

	if err := d.Set("issuer_uri", federatedSettingsIdentityProvider.IssuerUri); err != nil {
		return diag.FromErr(fmt.Errorf("error setting issuer uri (%s): %s", d.Id(), err))
	}

	if err := d.Set("idp_id", federatedSettingsIdentityProvider.Id); err != nil {
		return diag.FromErr(fmt.Errorf("error setting IdP Id (%s): %s", d.Id(), err))
	}

	if err := d.Set("protocol", federatedSettingsIdentityProvider.Protocol); err != nil {
		return diag.FromErr(fmt.Errorf("error setting protocol (%s): %s", d.Id(), err))
	}

	d.SetId(encodeStateID(federationSettingsID, federatedSettingsIdentityProvider.Id))

	return nil
}

func resourceMongoDBAtlasFederatedSettingsIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	oktaIdpID := ids["okta_idp_id"]

	existingIdentityProvider, _, err := connV2.FederatedAuthenticationApi.GetIdentityProvider(context.Background(), federationSettingsID, oktaIdpID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retreiving federation settings identity provider (%s): %s", federationSettingsID, err))
	}

	updateRequest := ExpandIdentityProviderUpdate(d, existingIdentityProvider)
	_, _, err = connV2.FederatedAuthenticationApi.UpdateIdentityProvider(ctx, federationSettingsID, oktaIdpID, updateRequest).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating federation settings identity provider (%s): %s", federationSettingsID, err))
	}

	return resourceMongoDBAtlasFederatedSettingsIdentityProviderRead(ctx, d, meta)
}

func resourceMongoDBAtlasFederatedSettingsIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourceMongoDBAtlasFederatedSettingsIdentityProviderImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	federationSettingsID, idpID, err := splitFederatedSettingsIdentityProviderImportID(d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(encodeStateID(*federationSettingsID, *idpID))

	return []*schema.ResourceData{d}, nil
}

func splitFederatedSettingsIdentityProviderImportID(id string) (federationSettingsID, oktaIdpID *string, err error) {
	var re = regexp.MustCompile(`(?s)^(.*)-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("import format error: to import a Federated SettingsIdentity Provider, use the format {federation_settings_id}-{okta_idp_id}")
		return
	}

	federationSettingsID = &parts[1]
	oktaIdpID = &parts[2]

	return
}

// Since the migration of this resource to the latest version of the auto-generated SDK (v20231115) and the breaking changes of the API
// the unique identifier used by the API & SDK of the supported identity providers is no longer "okta_idp_id", it is "idp_id". Nonetheless
// "okta_idp_id" name was used to encode/decode the Terraform State Id. To ensure backwards compatibility, the format of this resource id remains the same but the key `okta_idp_id` will store either `okta_idp_id` or `idp_id` to identify the identity provider.
// as few changes as possible, this name will remain.
func encodeStateID(federationSettingsID, idpID string) string {
	return conversion.EncodeStateID(map[string]string{
		"federation_settings_id": federationSettingsID,
		"okta_idp_id":            idpID,
	})
}
