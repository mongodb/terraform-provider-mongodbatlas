package organization

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115003/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasOrganizationCreate,
		ReadContext:   resourceMongoDBAtlasOrganizationRead,
		UpdateContext: resourceMongoDBAtlasOrganizationUpdate,
		DeleteContext: resourceMongoDBAtlasOrganizationDelete,
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
		},
	}
}

func resourceMongoDBAtlasOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	organization, resp, err := conn.OrganizationsApi.CreateOrganization(ctx, newCreateOrganizationRequest(d)).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error create Organization: %s", err))
	}

	orgID := *organization.Organization.Id

	// update settings using new keys for this created org
	cfg := config.Config{
		PublicKey:  *organization.ApiKey.PublicKey,
		PrivateKey: *organization.ApiKey.PrivateKey,
		BaseURL:    meta.(*config.MongoDBClient).Config.BaseURL,
	}

	clients, _ := cfg.NewClient(ctx)
	conn = clients.(*config.MongoDBClient).AtlasV2

	_, _, err = conn.OrganizationsApi.UpdateOrganizationSettings(ctx, orgID, newOrganizationSettings(d)).Execute()
	if err != nil {
		if _, _, err := conn.OrganizationsApi.DeleteOrganization(ctx, orgID).Execute(); err != nil {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("unable to delete organization because of create failure. An error occurred when updating Organization settings: %s", err))
		}

		return diag.FromErr(fmt.Errorf("error create Organization: %s", err))
	}

	if err := d.Set("private_key", organization.ApiKey.PrivateKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `private_key`: %s", err))
	}

	if err := d.Set("public_key", organization.ApiKey.PublicKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	if err := d.Set("org_id", *organization.Organization.Id); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `org_id`: %s", err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id": *organization.Organization.Id,
	}))

	return resourceMongoDBAtlasOrganizationRead(ctx, d, meta)
}

func resourceMongoDBAtlasOrganizationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	cfg := config.Config{
		PublicKey:  d.Get("public_key").(string),
		PrivateKey: d.Get("private_key").(string),
		BaseURL:    meta.(*config.MongoDBClient).Config.BaseURL,
	}

	clients, _ := cfg.NewClient(ctx)
	conn := clients.(*config.MongoDBClient).AtlasV2

	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]

	// organization, resp, err := conn.Organizations.Get(ctx, orgID)
	organization, resp, err := conn.OrganizationsApi.GetOrganization(ctx, orgID).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("warning Organization deleted will recreate: %s \n", err.Error())
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading organization information: %s", err))
	}

	if err := d.Set("name", organization.Name); err != nil {
		return diag.Errorf("error setting `name` for organization (%s): %s", *organization.Id, err)
	}

	settings, _, err := conn.OrganizationsApi.GetOrganizationSettings(ctx, orgID).Execute()

	if err := d.Set("api_access_list_required", settings.ApiAccessListRequired); err != nil {
		return diag.Errorf("error setting `api_access_list_required` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("multi_factor_auth_required", settings.MultiFactorAuthRequired); err != nil {
		return diag.Errorf("error setting `multi_factor_auth_required` for organization (%s): %s", orgID, err)
	}
	if err := d.Set("restrict_employee_access", settings.RestrictEmployeeAccess); err != nil {
		return diag.Errorf("error setting `restrict_employee_access` for organization (%s): %s", orgID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id": *organization.Id,
	}))
	return nil
}

func resourceMongoDBAtlasOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	cfg := config.Config{
		PublicKey:  d.Get("public_key").(string),
		PrivateKey: d.Get("private_key").(string),
		BaseURL:    meta.(*config.MongoDBClient).Config.BaseURL,
	}

	clients, _ := cfg.NewClient(ctx)
	conn := clients.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]

	updateRequest := new(admin.AtlasOrganization)
	if d.HasChange("name") {
		updateRequest.Name = d.Get("name").(string)
		_, _, err := conn.OrganizationsApi.RenameOrganization(ctx, orgID, updateRequest).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating Organization name: %s", err))
		}
	}

	if d.HasChange("api_access_list_required") || d.HasChange("multi_factor_auth_required") || d.HasChange("restrict_employee_access") {
		if _, _, err := conn.OrganizationsApi.UpdateOrganizationSettings(ctx, orgID, newOrganizationSettings(d)).Execute(); err != nil {
			return diag.FromErr(fmt.Errorf("error updating Organization settings: %s", err))
		}
	}

	return resourceMongoDBAtlasOrganizationRead(ctx, d, meta)
}

func resourceMongoDBAtlasOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	cfg := config.Config{
		PublicKey:  d.Get("public_key").(string),
		PrivateKey: d.Get("private_key").(string),
		BaseURL:    meta.(*config.MongoDBClient).Config.BaseURL,
	}

	clients, _ := cfg.NewClient(ctx)
	conn := clients.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]

	if _, _, err := conn.OrganizationsApi.DeleteOrganization(ctx, orgID).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("error Organization: %s", err))
	}
	return nil
}

// func newCreateOrganizationRequest(d *schema.ResourceData) *matlas.CreateOrganizationRequest {
// 	createRequest := &matlas.CreateOrganizationRequest{
// 		Name:       d.Get("name").(string),
// 		OrgOwnerID: pointy.String(d.Get("org_owner_id").(string)),
// 		APIKey: &matlas.APIKeyInput{
// 			Roles: conversion.ExpandStringList(d.Get("role_names").(*schema.Set).List()),
// 			Desc:  d.Get("description").(string),
// 		},
// 	}

// 	if federationSettingsID, ok := d.Get("federation_settings_id").(string); ok && federationSettingsID != "" {
// 		createRequest.FederationSettingsID = &federationSettingsID
// 	}

// 	return createRequest
// }

func newCreateOrganizationRequest(d *schema.ResourceData) *admin.CreateOrganizationRequest {
	createRequest := &admin.CreateOrganizationRequest{
		Name:       d.Get("name").(string),
		OrgOwnerId: pointy.String(d.Get("org_owner_id").(string)),

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
		ApiAccessListRequired:   pointy.Bool(d.Get("api_access_list_required").(bool)),
		MultiFactorAuthRequired: pointy.Bool(d.Get("multi_factor_auth_required").(bool)),
		RestrictEmployeeAccess:  pointy.Bool(d.Get("restrict_employee_access").(bool)),
	}
}
