package federatedsettingsorgconfig

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/spf13/cast"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateNotAllowed,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourcDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
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
				Optional: true,
			},
			"domain_allow_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"data_access_identity_provider_ids": {
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
			"user_conflicts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     userConflictsElemSchema(),
			},
		},
	}
}

func resourceCreateNotAllowed(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	return diag.FromErr(errors.New("this resource must be imported"))
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]

	federatedSettingsConnectedOrganization, resp, err := conn.FederatedAuthenticationApi.GetConnectedOrgConfig(context.Background(), federationSettingsID, orgID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting federated settings organization config: %s", err)
	}

	if err := d.Set("domain_restriction_enabled", federatedSettingsConnectedOrganization.DomainRestrictionEnabled); err != nil {
		return diag.Errorf("error setting domain restriction enabled (%s): %s", orgID, err)
	}

	if err := d.Set("domain_allow_list", federatedSettingsConnectedOrganization.DomainAllowList); err != nil {
		return diag.Errorf("error setting domain allow list (%s): %s", orgID, err)
	}
	if err := d.Set("data_access_identity_provider_ids", federatedSettingsConnectedOrganization.GetDataAccessIdentityProviderIds()); err != nil {
		return diag.Errorf("error setting data_access_identity_provider_ids (%s): %s", orgID, err)
	}

	if err := d.Set("post_auth_role_grants", federatedSettingsConnectedOrganization.PostAuthRoleGrants); err != nil {
		return diag.Errorf("error setting post_auth_role_grants (%s): %s", orgID, err)
	}
	if err := d.Set("user_conflicts", FlattenUserConflicts(federatedSettingsConnectedOrganization.GetUserConflicts())); err != nil {
		return diag.Errorf("error setting `user_conflicts` (%s): %s", orgID, err)
	}
	if err := d.Set("identity_provider_id", federatedSettingsConnectedOrganization.GetIdentityProviderId()); err != nil {
		return diag.Errorf("error setting identity provider id (%s): %s", orgID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"federation_settings_id": federationSettingsID,
		"org_id":                 orgID,
	}))

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]

	federatedSettingsConnectedOrganizationUpdate, _, err := conn.FederatedAuthenticationApi.GetConnectedOrgConfig(ctx, federationSettingsID, orgID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retreiving federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	if d.HasChange("domain_restriction_enabled") {
		domainRestrictionEnabled := d.Get("domain_restriction_enabled").(bool)
		federatedSettingsConnectedOrganizationUpdate.SetDomainRestrictionEnabled(domainRestrictionEnabled)
	}

	if d.HasChange("domain_allow_list") {
		domainAllowList := d.Get("domain_allow_list")
		federatedSettingsConnectedOrganizationUpdate.SetDomainAllowList(cast.ToStringSlice(domainAllowList))
	}
	if d.HasChange("data_access_identity_provider_ids") {
		dataAccessIdentityProviderIDs := d.Get("data_access_identity_provider_ids")
		federatedSettingsConnectedOrganizationUpdate.SetDataAccessIdentityProviderIds(cast.ToStringSlice(dataAccessIdentityProviderIDs))
	}

	if d.HasChange("identity_provider_id") {
		identityProviderID := d.Get("identity_provider_id").(string)
		// if identityProviderId is not part of the PATCH payload, it will be detached, "" will raise VALIDATION_ERROR
		if identityProviderID == "" {
			federatedSettingsConnectedOrganizationUpdate.IdentityProviderId = nil
		} else {
			federatedSettingsConnectedOrganizationUpdate.SetIdentityProviderId(identityProviderID)
		}
	}

	if d.HasChange("post_auth_role_grants") {
		postAuthRoleGrants := d.Get("post_auth_role_grants")
		federatedSettingsConnectedOrganizationUpdate.SetPostAuthRoleGrants(cast.ToStringSlice(postAuthRoleGrants))
	}
	// role mappings are managed by the `mongodbatlas_federated_settings_org_role_mapping` resource, no updates when it is excluded in the payload
	// keeping existing value [] will raise VALIDATION_ERROR if identity_provider_id is not set
	federatedSettingsConnectedOrganizationUpdate.RoleMappings = nil

	_, _, err = conn.FederatedAuthenticationApi.UpdateConnectedOrgConfig(ctx, federationSettingsID, orgID, federatedSettingsConnectedOrganizationUpdate).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating federation settings connected organization (%s): %s", federationSettingsID, err))
	}

	return resourceRead(ctx, d, meta)
}

func resourcDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	federationSettingsID := ids["federation_settings_id"]
	orgID := ids["org_id"]

	_, _, err := conn.FederatedAuthenticationApi.RemoveConnectedOrgConfig(ctx, federationSettingsID, orgID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting federation settings connected organization (%s): %s", federationSettingsID, err))
	}
	d.SetId("")
	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).AtlasV2
	federationSettingsID, orgID, err := splitImportID(d.Id())
	if err != nil {
		return nil, err
	}

	_, _, err = conn.FederatedAuthenticationApi.GetConnectedOrgConfig(context.Background(), *federationSettingsID, *orgID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import Organization config (%s) in Federation settings (%s), error: %s", *orgID, *federationSettingsID, err)
	}

	if err := d.Set("federation_settings_id", *federationSettingsID); err != nil {
		return nil, fmt.Errorf("error setting Organization config in Federation settings (%s): %s", d.Id(), err)
	}
	if err := d.Set("org_id", *orgID); err != nil {
		return nil, fmt.Errorf("error setting org id (%s): %s", d.Id(), err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"federation_settings_id": *federationSettingsID,
		"org_id":                 *orgID,
	}))

	return []*schema.ResourceData{d}, nil
}

func splitImportID(id string) (federationSettingsID, orgID *string, err error) {
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
