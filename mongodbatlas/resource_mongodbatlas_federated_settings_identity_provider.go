package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"
)

func resourceMongoDBAtlasFederatedSettingsIdentityProvider() *schema.Resource {
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
		},
	}
}

func resourceMongoDBAtlasFederatedSettingsIdentityProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	if d.Id() == "" {
		d.SetId("")
		return nil
	}

	ids := decodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	oktaIdpID := ids["okta_idp_id"]

	federatedSettingsIdentityProvider, resp, err := conn.FederatedSettings.GetIdentityProvider(context.Background(), federationSettingsID, oktaIdpID)
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

	if err := d.Set("okta_idp_id", federatedSettingsIdentityProvider.OktaIdpID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting OktaIdpID (%s): %s", d.Id(), err))
	}

	if err := d.Set("status", federatedSettingsIdentityProvider.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting Status (%s): %s", d.Id(), err))
	}

	if err := d.Set("issuer_uri", federatedSettingsIdentityProvider.IssuerURI); err != nil {
		return diag.FromErr(fmt.Errorf("error setting issuer uri (%s): %s", d.Id(), err))
	}

	if err := d.Set("request_binding", federatedSettingsIdentityProvider.RequestBinding); err != nil {
		return diag.FromErr(fmt.Errorf("error setting request binding (%s): %s", d.Id(), err))
	}

	if err := d.Set("response_signature_algorithm", federatedSettingsIdentityProvider.ResponseSignatureAlgorithm); err != nil {
		return diag.FromErr(fmt.Errorf("error setting response signature algorithm (%s): %s", d.Id(), err))
	}

	if err := d.Set("sso_url", federatedSettingsIdentityProvider.SsoURL); err != nil {
		return diag.FromErr(fmt.Errorf("error setting sso url (%s): %s", d.Id(), err))
	}

	d.SetId(encodeStateID(map[string]string{
		"federation_settings_id": federationSettingsID,
		"okta_idp_id":            oktaIdpID,
	}))

	return nil
}

func resourceMongoDBAtlasFederatedSettingsIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	oktaIdpID := ids["okta_idp_id"]

	federatedSettingsIdentityProviderUpdate, _, err := conn.FederatedSettings.GetIdentityProvider(context.Background(), federationSettingsID, oktaIdpID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retreiving federation settings identity provider (%s): %s", federationSettingsID, err))
	}

	if d.HasChange("sso_debug_enabled") {
		ssoDebugEnabled := d.Get("sso_debug_enabled").(bool)
		federatedSettingsIdentityProviderUpdate.SsoDebugEnabled = &ssoDebugEnabled
	}

	if d.HasChange("associated_domains") {
		associatedDomains := d.Get("associated_domains")
		federatedSettingsIdentityProviderUpdate.AssociatedDomains = cast.ToStringSlice(associatedDomains)
	}

	if d.HasChange("name") {
		identityName := d.Get("name").(string)
		federatedSettingsIdentityProviderUpdate.DisplayName = identityName
	}

	if d.HasChange("status") {
		status := d.Get("status").(string)
		federatedSettingsIdentityProviderUpdate.Status = status
	}

	if d.HasChange("issuer_uri") {
		status := d.Get("issuer_uri").(string)
		federatedSettingsIdentityProviderUpdate.IssuerURI = status
	}

	if d.HasChange("request_binding") {
		status := d.Get("request_binding").(string)
		federatedSettingsIdentityProviderUpdate.RequestBinding = status
	}

	if d.HasChange("response_signature_algorithm") {
		status := d.Get("response_signature_algorithm").(string)
		federatedSettingsIdentityProviderUpdate.ResponseSignatureAlgorithm = status
	}

	if d.HasChange("sso_url") {
		status := d.Get("sso_url").(string)
		federatedSettingsIdentityProviderUpdate.SsoURL = status
	}

	federatedSettingsIdentityProviderUpdate.PemFileInfo = nil

	_, _, err = conn.FederatedSettings.UpdateIdentityProvider(ctx, federationSettingsID, oktaIdpID, federatedSettingsIdentityProviderUpdate)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating federation settings identity provider (%s): %s", federationSettingsID, err))
	}

	return resourceMongoDBAtlasFederatedSettingsIdentityProviderRead(ctx, d, meta)
}

func resourceMongoDBAtlasFederatedSettingsIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourceMongoDBAtlasFederatedSettingsIdentityProviderImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas
	federationSettingsID, oktaIdpID, err := splitFederatedSettingsIdentityProviderImportID(d.Id())
	if err != nil {
		return nil, err
	}

	federatedSettingsIdentityProvider, _, err := conn.FederatedSettings.GetIdentityProvider(context.Background(), *federationSettingsID, *oktaIdpID)
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

	if err := d.Set("issuer_uri", federatedSettingsIdentityProvider.IssuerURI); err != nil {
		return nil, fmt.Errorf("error setting issuer uri (%s): %s", d.Id(), err)
	}

	if err := d.Set("request_binding", federatedSettingsIdentityProvider.RequestBinding); err != nil {
		return nil, fmt.Errorf("error setting request binding (%s): %s", d.Id(), err)
	}

	if err := d.Set("response_signature_algorithm", federatedSettingsIdentityProvider.ResponseSignatureAlgorithm); err != nil {
		return nil, fmt.Errorf("error setting response signature algorithm (%s): %s", d.Id(), err)
	}

	if err := d.Set("sso_url", federatedSettingsIdentityProvider.SsoURL); err != nil {
		return nil, fmt.Errorf("error setting sso url (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
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
