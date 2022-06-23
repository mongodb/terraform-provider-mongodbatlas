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

func resourceMongoDBAtlasFederatedSettingsOrganizationConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasFederatedSettingsOrganizationConfigRead,
		ReadContext:   resourceMongoDBAtlasFederatedSettingsOrganizationConfigRead,
		UpdateContext: resourceMongoDBAtlasFederatedSettingsOrganizationConfigUpdate,
		DeleteContext: resourceMongoDBAtlasFederatedSettingsOrganizationConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasFederatedSettingsOrganizationConfigImportState,
		},
		Schema: map[string]*schema.Schema{
			"federation_settings_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identity_provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_allow_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"post_auth_role_grants": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_restriction_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceMongoDBAtlasFederatedSettingsOrganizationConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	if d.Id() == "" {
		d.SetId("")
		return nil
	}
	ids := decodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]

	federatedSettingsConnectedOrganization, resp, err := conn.FederatedSettings.GetConnectedOrg(context.Background(), federationSettingsID, orgID)
	if err != nil {
		// case 404
		// deleted in the backend case
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error getting federated settings organization config: %s", err))
	}

	if err := d.Set("domain_restriction_enabled", federatedSettingsConnectedOrganization.DomainRestrictionEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("error setting domain restriction enabled (%s): %s", d.Id(), err))
	}

	if err := d.Set("domain_allow_list", federatedSettingsConnectedOrganization.DomainAllowList); err != nil {
		return diag.FromErr(fmt.Errorf("error setting domain allow list (%s): %s", d.Id(), err))
	}

	if err := d.Set("post_auth_role_grants", federatedSettingsConnectedOrganization.PostAuthRoleGrants); err != nil {
		return diag.FromErr(fmt.Errorf("error setting post_auth_role_grants (%s): %s", d.Id(), err))
	}

	d.SetId(encodeStateID(map[string]string{
		"federation_settings_id": federationSettingsID,
		"org_id":                 orgID,
	}))

	return nil
}

func resourceMongoDBAtlasFederatedSettingsOrganizationConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]

	federatedSettingsConnectedOrganizationUpdate, _, err := conn.FederatedSettings.GetConnectedOrg(context.Background(), federationSettingsID, orgID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retreiving federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	if d.HasChange("domain_restriction_enabled") {
		domainRestrictionEnabled := d.Get("domain_restriction_enabled").(bool)
		federatedSettingsConnectedOrganizationUpdate.DomainRestrictionEnabled = &domainRestrictionEnabled
	}

	if d.HasChange("domain_allow_list") {
		domainAllowList := d.Get("domain_allow_list")
		federatedSettingsConnectedOrganizationUpdate.DomainAllowList = cast.ToStringSlice(domainAllowList)
	}

	if d.HasChange("identity_provider_id") {
		identityProviderID := d.Get("identity_provider_id").(string)
		federatedSettingsConnectedOrganizationUpdate.IdentityProviderID = identityProviderID
	}

	if d.HasChange("post_auth_role_grants") {
		postAuthRoleGrants := d.Get("post_auth_role_grants")
		federatedSettingsConnectedOrganizationUpdate.PostAuthRoleGrants = cast.ToStringSlice(postAuthRoleGrants)
	}

	_, _, err = conn.FederatedSettings.UpdateConnectedOrg(ctx, federationSettingsID, orgID, federatedSettingsConnectedOrganizationUpdate)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	return resourceMongoDBAtlasFederatedSettingsOrganizationConfigRead(ctx, d, meta)
}

func resourceMongoDBAtlasFederatedSettingsOrganizationConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]

	_, err := conn.FederatedSettings.DeleteConnectedOrg(ctx, federationSettingsID, orgID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	return nil
}

func resourceMongoDBAtlasFederatedSettingsOrganizationConfigImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas
	federationSettingsID, orgID, err := splitFederatedSettingsOrganizationConfigImportID(d.Id())
	if err != nil {
		return nil, err
	}

	federatedSettingsConnectedOrganization, _, err := conn.FederatedSettings.GetConnectedOrg(context.Background(), *federationSettingsID, *orgID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import Organization config (%s) in Federation settings (%s), error: %s", *orgID, *federationSettingsID, err)
	}

	if err := d.Set("federation_settings_id", *federationSettingsID); err != nil {
		return nil, fmt.Errorf("error setting Organization config in Federation settings (%s): %s", d.Id(), err)
	}

	if err := d.Set("domain_restriction_enabled", federatedSettingsConnectedOrganization.DomainRestrictionEnabled); err != nil {
		return nil, fmt.Errorf("error setting domain restriction enabled (%s): %s", d.Id(), err)
	}

	if err := d.Set("domain_allow_list", federatedSettingsConnectedOrganization.DomainAllowList); err != nil {
		return nil, fmt.Errorf("error setting domain allow list (%s): %s", d.Id(), err)
	}

	if err := d.Set("org_id", federatedSettingsConnectedOrganization.OrgID); err != nil {
		return nil, fmt.Errorf("error setting org id (%s): %s", d.Id(), err)
	}

	if err := d.Set("identity_provider_id", federatedSettingsConnectedOrganization.IdentityProviderID); err != nil {
		return nil, fmt.Errorf("error setting identity provider id (%s): %s", d.Id(), err)
	}

	d.SetId(encodeStateID(map[string]string{
		"federation_settings_id": *federationSettingsID,
		"org_id":                 *orgID,
	}))

	return []*schema.ResourceData{d}, nil
}

func splitFederatedSettingsOrganizationConfigImportID(id string) (federationSettingsID, orgID *string, err error) {
	var re = regexp.MustCompile(`(?s)^(.*)-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("import format error: to import a Federated Settings Orgnization Config, use the format {federation_settings_id}-{org_id}")
		return
	}

	federationSettingsID = &parts[1]
	orgID = &parts[2]

	return
}
