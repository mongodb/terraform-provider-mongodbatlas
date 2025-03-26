package organization

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer:      nil, // import is not supported. See CLOUDP-215155
		Schema: map[string]*schema.Schema{
			"org_owner_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"role_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"federation_settings_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_access_list_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"multi_factor_auth_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"restrict_employee_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"gen_ai_features_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			// skip_default_alerts_settings defaults to `true` to prevent Atlas from automatically creating organization-level alerts not explicitly managed through Terraform.
			// Note that this deviates from the API default of `false` for this attribute.
			"skip_default_alerts_settings": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := ValidateAPIKeyIsOrgOwner(conversion.ExpandStringList(d.Get("role_names").(*schema.Set).List())); err != nil {
		return diag.FromErr(err)
	}

	conn := meta.(*config.MongoDBClient).AtlasV2
	organization, resp, err := conn.OrganizationsApi.CreateOrganization(ctx, newCreateOrganizationRequest(d)).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) && !strings.Contains(err.Error(), "USER_NOT_FOUND") {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error creating Organization: %s", err))
	}

	orgID := organization.Organization.GetId()

	// update settings using new keys for this created organization because
	// the provider/requesting API keys are not applicable for performing updates/delete for this new organization
	cfg := config.Config{
		PublicKey:        *organization.ApiKey.PublicKey,
		PrivateKey:       *organization.ApiKey.PrivateKey,
		BaseURL:          meta.(*config.MongoDBClient).Config.BaseURL,
		TerraformVersion: meta.(*config.MongoDBClient).Config.TerraformVersion,
	}

	clients, _ := cfg.NewClient(ctx)
	conn = clients.(*config.MongoDBClient).AtlasV2

	_, _, errUpdate := conn.OrganizationsApi.UpdateOrganizationSettings(ctx, orgID, newOrganizationSettings(d)).Execute()
	if errUpdate != nil {
		if _, _, err := conn.OrganizationsApi.DeleteOrganization(ctx, orgID).Execute(); err != nil {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("an error occurred when updating Organization settings: %s.\n Unable to delete organization, there may be dangling resources: %s", errUpdate.Error(), err.Error()))
		}
		d.SetId("")
		return diag.FromErr(fmt.Errorf("an error occurred when updating Organization settings: %s", err))
	}

	if err := d.Set("private_key", organization.ApiKey.GetPrivateKey()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `private_key`: %s", err))
	}

	if err := d.Set("public_key", organization.ApiKey.GetPublicKey()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	if err := d.Set("org_id", organization.Organization.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `org_id`: %s", err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id": organization.Organization.GetId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	cfg := config.Config{
		PublicKey:        d.Get("public_key").(string),
		PrivateKey:       d.Get("private_key").(string),
		BaseURL:          meta.(*config.MongoDBClient).Config.BaseURL,
		TerraformVersion: meta.(*config.MongoDBClient).Config.TerraformVersion,
	}

	clients, _ := cfg.NewClient(ctx)
	conn := clients.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]

	organization, resp, err := conn.OrganizationsApi.GetOrganization(ctx, orgID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			log.Printf("warning Organization deleted will recreate: %s \n", err.Error())
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading organization information: %s", err))
	}

	if err := d.Set("name", organization.Name); err != nil {
		return diag.Errorf("error setting `name` for organization (%s): %s", *organization.Id, err)
	}
	if err := d.Set("skip_default_alerts_settings", organization.SkipDefaultAlertsSettings); err != nil {
		return diag.Errorf("error setting `skip_default_alerts_settings` for organization (%s): %s", orgID, err)
	}

	settings, _, err := conn.OrganizationsApi.GetOrganizationSettings(ctx, orgID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading organization settings: %s", err))
	}

	if err := d.Set("api_access_list_required", settings.ApiAccessListRequired); err != nil {
		return diag.Errorf("error setting `api_access_list_required` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("multi_factor_auth_required", settings.MultiFactorAuthRequired); err != nil {
		return diag.Errorf("error setting `multi_factor_auth_required` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("restrict_employee_access", settings.RestrictEmployeeAccess); err != nil {
		return diag.Errorf("error setting `restrict_employee_access` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("gen_ai_features_enabled", settings.GenAIFeaturesEnabled); err != nil {
		return diag.Errorf("error setting `gen_ai_features_enabled` for organization (%s): %s", orgID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id": organization.GetId(),
	}))
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	cfg := config.Config{
		PublicKey:        d.Get("public_key").(string),
		PrivateKey:       d.Get("private_key").(string),
		BaseURL:          meta.(*config.MongoDBClient).Config.BaseURL,
		TerraformVersion: meta.(*config.MongoDBClient).Config.TerraformVersion,
	}

	clients, _ := cfg.NewClient(ctx)
	conn := clients.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]

	updateRequest := new(admin.AtlasOrganization)

	if d.HasChange("name") ||
		d.HasChange("skip_default_alerts_settings") {
		updateRequest.Name = d.Get("name").(string)
		updateRequest.SkipDefaultAlertsSettings = conversion.Pointer(d.Get("skip_default_alerts_settings").(bool))

		_, _, err := conn.OrganizationsApi.UpdateOrganization(ctx, orgID, updateRequest).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating Organization name: %s", err))
		}
	}

	if d.HasChange("api_access_list_required") ||
		d.HasChange("multi_factor_auth_required") ||
		d.HasChange("restrict_employee_access") ||
		d.HasChange("gen_ai_features_enabled") {
		if _, _, err := conn.OrganizationsApi.UpdateOrganizationSettings(ctx, orgID, newOrganizationSettings(d)).Execute(); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Organization settings: %s", err))
		}
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	cfg := config.Config{
		PublicKey:        d.Get("public_key").(string),
		PrivateKey:       d.Get("private_key").(string),
		BaseURL:          meta.(*config.MongoDBClient).Config.BaseURL,
		TerraformVersion: meta.(*config.MongoDBClient).Config.TerraformVersion,
	}

	clients, _ := cfg.NewClient(ctx)
	conn := clients.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]

	if _, _, err := conn.OrganizationsApi.DeleteOrganization(ctx, orgID).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting Organization: %s", err))
	}
	return nil
}

func newCreateOrganizationRequest(d *schema.ResourceData) *admin.CreateOrganizationRequest {
	// skip_default_alerts_settings defaults to `true` to prevent Atlas from automatically creating organization-level alerts not explicitly managed through Terraform.
	// Note that this deviates from the API default of `false` for this attribute.
	skipDefaultAlertsSettings := true

	if v, ok := d.GetOkExists("skip_default_alerts_settings"); ok {
		skipDefaultAlertsSettings = v.(bool)
	}

	createRequest := &admin.CreateOrganizationRequest{
		Name:                      d.Get("name").(string),
		OrgOwnerId:                conversion.Pointer(d.Get("org_owner_id").(string)),
		SkipDefaultAlertsSettings: conversion.Pointer(skipDefaultAlertsSettings),

		ApiKey: &admin.CreateAtlasOrganizationApiKey{
			Roles: conversion.ExpandStringList(d.Get("role_names").(*schema.Set).List()),
			Desc:  d.Get("description").(string),
		},
	}

	if federationSettingsID, ok := d.Get("federation_settings_id").(string); ok && federationSettingsID != "" {
		createRequest.FederationSettingsId = &federationSettingsID
	}

	return createRequest
}

func newOrganizationSettings(d *schema.ResourceData) *admin.OrganizationSettings {
	return &admin.OrganizationSettings{
		ApiAccessListRequired:   conversion.Pointer(d.Get("api_access_list_required").(bool)),
		MultiFactorAuthRequired: conversion.Pointer(d.Get("multi_factor_auth_required").(bool)),
		RestrictEmployeeAccess:  conversion.Pointer(d.Get("restrict_employee_access").(bool)),
		GenAIFeaturesEnabled:    conversion.Pointer(d.Get("gen_ai_features_enabled").(bool)),
	}
}

func ValidateAPIKeyIsOrgOwner(roles []string) error {
	for _, role := range roles {
		if role == constant.OrgOwner {
			return nil
		}
	}

	return fmt.Errorf("`role_names` for new API Key must have the ORG_OWNER role to use this resource")
}
