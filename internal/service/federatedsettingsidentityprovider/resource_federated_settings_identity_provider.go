package federatedsettingsidentityprovider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	admin20231001002 "go.mongodb.org/atlas-sdk/v20231001002/admin"
	"go.mongodb.org/atlas-sdk/v20231115005/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/spf13/cast"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasFederatedSettingsIdentityProviderRead,
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
				Required: true,
			},
			"response_signature_algorithm": {
				Type:     schema.TypeString,
				Required: true,
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
				Required: true,
			},
			"sso_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"okta_idp_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"idp_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasFederatedSettingsIdentityProviderRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	if d.Id() == "" {
		d.SetId("")
		return nil
	}

	ids := conversion.DecodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]

	// Since the migration of this resource to the latest version of the auto-generated SDK (v20231115) and the breaking changes of the API
	// the unique identifier used by the API & SDK of the supported identity providers is no longer "okta_idp_id", it is "idp_id". Nonetheless
	// "okta_idp_id" name was used to encode/decode the Terraform State Id. To ensure backwards compatibility, the format of this resource id remains the same but the key `okta_idp_id` will store either `okta_idp_id` or `idp_id` to identify the identity provider.
	// as few changes as possible, this name will remain.
	idpID := ids["okta_idp_id"]

	// to be removed in terraform-provider-1.16.0
	if len(idpID) == 20 {
		// use old version of v2 SDK
		return append(oldSDKRead(federationSettingsID, idpID, d, meta), getGracePeriodWarning())
	}
	// latest version of v2 SDK
	federatedSettingsIdentityProvider, resp, err := connV2.FederatedAuthenticationApi.GetIdentityProvider(context.Background(), federationSettingsID, idpID).Execute()
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting federated settings identity provider: %s", err))
	}

	if err := d.Set("sso_debug_enabled", federatedSettingsIdentityProvider.SsoDebugEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("error setting sso debug enabled (%s): %s", d.Id(), err))
	}

	if err := d.Set("associated_domains", federatedSettingsIdentityProvider.AssociatedDomains); err != nil {
		return diag.FromErr(fmt.Errorf("error setting associated domains list (%s): %s", d.Id(), err))
	}

	if err := d.Set("okta_idp_id", federatedSettingsIdentityProvider.OktaIdpId); err != nil {
		return diag.FromErr(fmt.Errorf("error setting OktaIdpID (%s): %s", d.Id(), err))
	}

	if err := d.Set("status", federatedSettingsIdentityProvider.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting Status (%s): %s", d.Id(), err))
	}

	if err := d.Set("issuer_uri", federatedSettingsIdentityProvider.IssuerUri); err != nil {
		return diag.FromErr(fmt.Errorf("error setting issuer uri (%s): %s", d.Id(), err))
	}

	if err := d.Set("request_binding", federatedSettingsIdentityProvider.RequestBinding); err != nil {
		return diag.FromErr(fmt.Errorf("error setting request binding (%s): %s", d.Id(), err))
	}

	if err := d.Set("response_signature_algorithm", federatedSettingsIdentityProvider.ResponseSignatureAlgorithm); err != nil {
		return diag.FromErr(fmt.Errorf("error setting response signature algorithm (%s): %s", d.Id(), err))
	}

	if err := d.Set("sso_url", federatedSettingsIdentityProvider.SsoUrl); err != nil {
		return diag.FromErr(fmt.Errorf("error setting sso url (%s): %s", d.Id(), err))
	}

	if err := d.Set("idp_id", federatedSettingsIdentityProvider.Id); err != nil {
		return diag.FromErr(fmt.Errorf("error setting IdP Id (%s): %s", d.Id(), err))
	}

	d.SetId(encodeStateID(federationSettingsID, federatedSettingsIdentityProvider.Id))

	return nil
}

func oldSDKRead(federationSettingsID, oktaIdpID string, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn20231001002 := meta.(*config.MongoDBClient).Atlas20231001002

	federatedSettingsIdentityProvider, resp, err := conn20231001002.FederatedAuthenticationApi.GetIdentityProvider(context.Background(), federationSettingsID, oktaIdpID).Execute()
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting federated settings identity provider: %s", err))
	}

	if err := d.Set("sso_debug_enabled", federatedSettingsIdentityProvider.SsoDebugEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("error setting sso debug enabled (%s): %s", d.Id(), err))
	}

	if err := d.Set("associated_domains", federatedSettingsIdentityProvider.AssociatedDomains); err != nil {
		return diag.FromErr(fmt.Errorf("error setting associated domains list (%s): %s", d.Id(), err))
	}

	if err := d.Set("okta_idp_id", federatedSettingsIdentityProvider.OktaIdpId); err != nil {
		return diag.FromErr(fmt.Errorf("error setting OktaIdpID (%s): %s", d.Id(), err))
	}

	if err := d.Set("status", federatedSettingsIdentityProvider.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting Status (%s): %s", d.Id(), err))
	}

	if err := d.Set("issuer_uri", federatedSettingsIdentityProvider.IssuerUri); err != nil {
		return diag.FromErr(fmt.Errorf("error setting issuer uri (%s): %s", d.Id(), err))
	}

	if err := d.Set("request_binding", federatedSettingsIdentityProvider.RequestBinding); err != nil {
		return diag.FromErr(fmt.Errorf("error setting request binding (%s): %s", d.Id(), err))
	}

	if err := d.Set("response_signature_algorithm", federatedSettingsIdentityProvider.ResponseSignatureAlgorithm); err != nil {
		return diag.FromErr(fmt.Errorf("error setting response signature algorithm (%s): %s", d.Id(), err))
	}

	if err := d.Set("sso_url", federatedSettingsIdentityProvider.SsoUrl); err != nil {
		return diag.FromErr(fmt.Errorf("error setting sso url (%s): %s", d.Id(), err))
	}

	if err := d.Set("idp_id", federatedSettingsIdentityProvider.Id); err != nil {
		return diag.FromErr(fmt.Errorf("error setting IdP Id (%s): %s", d.Id(), err))
	}

	d.SetId(encodeStateID(federationSettingsID, oktaIdpID))

	return nil
}

func resourceMongoDBAtlasFederatedSettingsIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	oktaIdpID := ids["okta_idp_id"]

	// to be removed in terraform-provider-1.16.0
	if len(oktaIdpID) == 20 {
		return append(oldSDKUpdate(ctx, federationSettingsID, oktaIdpID, d, meta), getGracePeriodWarning())
	}

	updateRequest := new(admin.FederationIdentityProviderUpdate)
	_, _, err := connV2.FederatedAuthenticationApi.GetIdentityProvider(context.Background(), federationSettingsID, oktaIdpID).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error retreiving federation settings identity provider (%s): %s", federationSettingsID, err))
	}

	if d.HasChange("sso_debug_enabled") {
		ssoDebugEnabled := d.Get("sso_debug_enabled").(bool)
		updateRequest.SsoDebugEnabled = &ssoDebugEnabled
	}

	if d.HasChange("associated_domains") {
		associatedDomains := d.Get("associated_domains")
		associatedDomainsSlice := cast.ToStringSlice(associatedDomains)
		if associatedDomainsSlice == nil {
			associatedDomainsSlice = []string{}
		}
		updateRequest.AssociatedDomains = &associatedDomainsSlice
	}

	if d.HasChange("name") {
		identityName := d.Get("name").(string)
		updateRequest.DisplayName = &identityName
	}

	if d.HasChange("status") {
		status := d.Get("status").(string)
		updateRequest.Status = &status
	}

	if d.HasChange("issuer_uri") {
		status := d.Get("issuer_uri").(string)
		updateRequest.IssuerUri = &status
	}

	if d.HasChange("request_binding") {
		status := d.Get("request_binding").(string)
		updateRequest.RequestBinding = &status
	}

	if d.HasChange("response_signature_algorithm") {
		status := d.Get("response_signature_algorithm").(string)
		updateRequest.ResponseSignatureAlgorithm = &status
	}

	if d.HasChange("sso_url") {
		status := d.Get("sso_url").(string)
		updateRequest.SsoUrl = &status
	}

	updateRequest.PemFileInfo = nil

	_, _, err = connV2.FederatedAuthenticationApi.UpdateIdentityProvider(ctx, federationSettingsID, oktaIdpID, updateRequest).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating federation settings identity provider (%s): %s", federationSettingsID, err))
	}

	return resourceMongoDBAtlasFederatedSettingsIdentityProviderRead(ctx, d, meta)
}

func oldSDKUpdate(ctx context.Context, federationSettingsID, oktaIdpID string, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn20231001002 := meta.(*config.MongoDBClient).Atlas20231001002
	var updateRequest *admin20231001002.SamlIdentityProviderUpdate
	_, _, err := conn20231001002.FederatedAuthenticationApi.GetIdentityProvider(context.Background(), federationSettingsID, oktaIdpID).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error retreiving federation settings identity provider (%s): %s", federationSettingsID, err))
	}

	if d.HasChange("sso_debug_enabled") {
		ssoDebugEnabled := d.Get("sso_debug_enabled").(bool)
		updateRequest.SsoDebugEnabled = ssoDebugEnabled
	}

	if d.HasChange("associated_domains") {
		associatedDomains := d.Get("associated_domains")
		updateRequest.AssociatedDomains = cast.ToStringSlice(associatedDomains)
	}

	if d.HasChange("name") {
		identityName := d.Get("name").(string)
		updateRequest.DisplayName = &identityName
	}

	if d.HasChange("status") {
		status := d.Get("status").(string)
		updateRequest.Status = &status
	}

	if d.HasChange("issuer_uri") {
		status := d.Get("issuer_uri").(string)
		updateRequest.IssuerUri = &status
	}

	if d.HasChange("request_binding") {
		status := d.Get("request_binding").(string)
		updateRequest.RequestBinding = &status
	}

	if d.HasChange("response_signature_algorithm") {
		status := d.Get("response_signature_algorithm").(string)
		updateRequest.ResponseSignatureAlgorithm = &status
	}

	if d.HasChange("sso_url") {
		status := d.Get("sso_url").(string)
		updateRequest.SsoUrl = &status
	}

	updateRequest.PemFileInfo = nil

	_, _, err = conn20231001002.FederatedAuthenticationApi.UpdateIdentityProvider(ctx, federationSettingsID, oktaIdpID, updateRequest).Execute()
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
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	federationSettingsID, idpID, err := splitFederatedSettingsIdentityProviderImportID(d.Id())
	if err != nil {
		return nil, err
	}

	// to be removed in terraform-provider-1.16.0
	if len(*idpID) == 20 {
		return oldSDKImport(federationSettingsID, idpID, d, meta)
	}

	federatedSettingsIdentityProvider, _, err := connV2.FederatedAuthenticationApi.GetIdentityProvider(context.Background(), *federationSettingsID, *idpID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import Organization config (%s) in Federation settings (%s), error: %s", *idpID, *federationSettingsID, err)
	}

	if err := d.Set("federation_settings_id", *federationSettingsID); err != nil {
		return nil, fmt.Errorf("error setting Identity Provider in Federation settings (%s): %s", d.Id(), err)
	}

	if err := d.Set("name", federatedSettingsIdentityProvider.DisplayName); err != nil {
		return nil, fmt.Errorf("error setting name (%s): %s", d.Id(), err)
	}

	if err := d.Set("sso_debug_enabled", federatedSettingsIdentityProvider.SsoDebugEnabled); err != nil {
		return nil, fmt.Errorf("error setting sso debug enabled (%s): %s", d.Id(), err)
	}

	if err := d.Set("associated_domains", federatedSettingsIdentityProvider.AssociatedDomains); err != nil {
		return nil, fmt.Errorf("error setting associaed domains list (%s): %s", d.Id(), err)
	}

	if err := d.Set("issuer_uri", federatedSettingsIdentityProvider.IssuerUri); err != nil {
		return nil, fmt.Errorf("error setting issuer uri (%s): %s", d.Id(), err)
	}

	if err := d.Set("request_binding", federatedSettingsIdentityProvider.RequestBinding); err != nil {
		return nil, fmt.Errorf("error setting request binding (%s): %s", d.Id(), err)
	}

	if err := d.Set("response_signature_algorithm", federatedSettingsIdentityProvider.ResponseSignatureAlgorithm); err != nil {
		return nil, fmt.Errorf("error setting response signature algorithm (%s): %s", d.Id(), err)
	}

	if err := d.Set("sso_url", federatedSettingsIdentityProvider.SsoUrl); err != nil {
		return nil, fmt.Errorf("error setting sso url (%s): %s", d.Id(), err)
	}

	if err := d.Set("idp_id", federatedSettingsIdentityProvider.Id); err != nil {
		return nil, fmt.Errorf("error setting IdP Id (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(*federationSettingsID, *idpID))

	return []*schema.ResourceData{d}, nil
}

func oldSDKImport(federationSettingsID, oktaIdpID *string, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn20231001002 := meta.(*config.MongoDBClient).Atlas20231001002
	federatedSettingsIdentityProvider, _, err := conn20231001002.FederatedAuthenticationApi.GetIdentityProvider(context.Background(), *federationSettingsID, *oktaIdpID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import Organization config (%s) in Federation settings (%s), error: %s", *oktaIdpID, *federationSettingsID, err)
	}

	if err := d.Set("federation_settings_id", *federationSettingsID); err != nil {
		return nil, fmt.Errorf("error setting Identity Provider in Federation settings (%s): %s", d.Id(), err)
	}

	if err := d.Set("name", federatedSettingsIdentityProvider.DisplayName); err != nil {
		return nil, fmt.Errorf("error setting name (%s): %s", d.Id(), err)
	}

	if err := d.Set("sso_debug_enabled", federatedSettingsIdentityProvider.SsoDebugEnabled); err != nil {
		return nil, fmt.Errorf("error setting sso debug enabled (%s): %s", d.Id(), err)
	}

	if err := d.Set("associated_domains", federatedSettingsIdentityProvider.AssociatedDomains); err != nil {
		return nil, fmt.Errorf("error setting associaed domains list (%s): %s", d.Id(), err)
	}

	if err := d.Set("issuer_uri", federatedSettingsIdentityProvider.IssuerUri); err != nil {
		return nil, fmt.Errorf("error setting issuer uri (%s): %s", d.Id(), err)
	}

	if err := d.Set("request_binding", federatedSettingsIdentityProvider.RequestBinding); err != nil {
		return nil, fmt.Errorf("error setting request binding (%s): %s", d.Id(), err)
	}

	if err := d.Set("response_signature_algorithm", federatedSettingsIdentityProvider.ResponseSignatureAlgorithm); err != nil {
		return nil, fmt.Errorf("error setting response signature algorithm (%s): %s", d.Id(), err)
	}

	if err := d.Set("sso_url", federatedSettingsIdentityProvider.SsoUrl); err != nil {
		return nil, fmt.Errorf("error setting sso url (%s): %s", d.Id(), err)
	}

	if err := d.Set("idp_id", federatedSettingsIdentityProvider.Id); err != nil {
		return nil, fmt.Errorf("error setting IdP Id (%s): %s", d.Id(), err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"federation_settings_id": *federationSettingsID,
		"okta_idp_id":            *oktaIdpID,
	}))

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

func getGracePeriodWarning() diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Warning: deprecated identity provider id",
		Detail: "Identity provider id format defined in resource will be deprecated. Please import the resource with the new format.\n" +
			" Follow instructions here: https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/1.15.0-upgrade-guide",
	}
}
