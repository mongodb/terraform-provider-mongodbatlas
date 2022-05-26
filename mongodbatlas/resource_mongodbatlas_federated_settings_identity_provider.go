package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	federatedSettingsConnectedOrganization, resp, err := conn.FederatedSettings.GetIdentityProvider(context.Background(), federationSettingsID, oktaIdpID)
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting federated settings identity provider: %s", err))
	}

	if err := d.Set("sso_debug_enabled", federatedSettingsConnectedOrganization.SsoDebugEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("error setting sso debug enabled (%s): %s", d.Id(), err))
	}

	if err := d.Set("associated_domains", federatedSettingsConnectedOrganization.AssociatedDomains); err != nil {
		return diag.FromErr(fmt.Errorf("error setting associated domains list (%s): %s", d.Id(), err))
	}

	if err := d.Set("okta_idp_id", federatedSettingsConnectedOrganization.OktaIdpID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting OktaIdpID (%s): %s", d.Id(), err))
	}

	if err := d.Set("status", federatedSettingsConnectedOrganization.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting Status (%s): %s", d.Id(), err))
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

	federatedSettingsConnectedOrganizationUpdate, _, err := conn.FederatedSettings.GetIdentityProvider(context.Background(), federationSettingsID, oktaIdpID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retreiving federation settings identity provider (%s): %s", federationSettingsID, err))
	}

	if d.HasChange("sso_debug_enabled") {
		ssoDebugEnabled := d.Get("sso_debug_enabled").(bool)
		federatedSettingsConnectedOrganizationUpdate.SsoDebugEnabled = &ssoDebugEnabled
	}

	if d.HasChange("associated_domains") {
		associatedDomains := d.Get("associated_domains")
		federatedSettingsConnectedOrganizationUpdate.AssociatedDomains = cast.ToStringSlice(associatedDomains)
	}

	if d.HasChange("name") {
		identityName := d.Get("name").(string)
		federatedSettingsConnectedOrganizationUpdate.DisplayName = identityName
	}

	if d.HasChange("status") {
		status := d.Get("status").(string)
		federatedSettingsConnectedOrganizationUpdate.Status = status
	}

	federatedSettingsConnectedOrganizationUpdate.PemFileInfo = nil

	_, _, err = conn.FederatedSettings.UpdateIdentityProvider(ctx, federationSettingsID, oktaIdpID, federatedSettingsConnectedOrganizationUpdate)
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

	if err := d.Set("sso_debug_enabled", federatedSettingsIdentityProvider.SsoDebugEnabled); err != nil {
		return nil, fmt.Errorf("error setting sso debug enabled (%s): %s", d.Id(), err)
	}

	if err := d.Set("associated_domains", federatedSettingsIdentityProvider.AssociatedDomains); err != nil {
		return nil, fmt.Errorf("error setting associaed domains list (%s): %s", d.Id(), err)
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
		err = errors.New("import format error: to import a Federated SettingsIdentity Provider, use the format {federation_settings_id}-{org_id}-{okta_idp_id}")
		return
	}

	federationSettingsID = &parts[1]
	oktaIdpID = &parts[2]

	return
}
