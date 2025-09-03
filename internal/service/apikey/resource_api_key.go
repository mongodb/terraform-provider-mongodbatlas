package apikey

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
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
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	orgID := d.Get("org_id").(string)
	createRequest := &admin.CreateAtlasOrganizationApiKey{
		Desc:  d.Get("description").(string),
		Roles: conversion.ExpandStringList(d.Get("role_names").(*schema.Set).List()),
	}

	apiKey, resp, err := connV2.ProgrammaticAPIKeysApi.CreateApiKey(ctx, orgID, createRequest).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error create API key: %s", err))
	}

	if err := d.Set("private_key", apiKey.GetPrivateKey()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `private_key`: %s", err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id":     orgID,
		"api_key_id": apiKey.GetId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	apiKeyID := ids["api_key_id"]

	apiKey, resp, err := connV2.ProgrammaticAPIKeysApi.GetApiKey(ctx, orgID, apiKeyID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) || validate.StatusBadRequest(resp) {
			log.Printf("warning API key deleted will recreate: %s \n", err.Error())
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}

	if err := d.Set("api_key_id", apiKey.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `api_key_id`: %s", err))
	}

	if err := d.Set("description", apiKey.GetDesc()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `description`: %s", err))
	}

	if err := d.Set("public_key", apiKey.GetPublicKey()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	if err := d.Set("role_names", flattenOrgAPIKeyRoles(orgID, apiKey.GetRoles())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `roles`: %s", err))
	}

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	apiKeyID := ids["api_key_id"]

	if d.HasChange("description") || d.HasChange("role_names") {
		updateRequest := &admin.UpdateAtlasOrganizationApiKey{
			Desc: conversion.StringPtr(d.Get("description").(string)),
		}
		if roles := conversion.ExpandStringListFromSetSchema(d.Get("role_names").(*schema.Set)); roles != nil {
			updateRequest.Roles = &roles
		}
		_, _, err := connV2.ProgrammaticAPIKeysApi.UpdateApiKey(ctx, orgID, apiKeyID, updateRequest).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating API key: %s", err))
		}
	}
	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	orgID := ids["org_id"]
	apiKeyID := ids["api_key_id"]

	_, err := connV2.ProgrammaticAPIKeysApi.DeleteApiKey(ctx, orgID, apiKeyID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error API Key: %s", err))
	}
	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a api key use the format {org_id}-{api_key_id}")
	}

	orgID := parts[0]
	apiKeyID := parts[1]

	r, _, err := connV2.ProgrammaticAPIKeysApi.GetApiKey(ctx, orgID, apiKeyID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import api key %s in project %s, error: %s", orgID, apiKeyID, err)
	}

	if err := d.Set("org_id", orgID); err != nil {
		return nil, fmt.Errorf("error setting `org_id`: %s", err)
	}

	if err := d.Set("description", r.Desc); err != nil {
		return nil, fmt.Errorf("error setting `description`: %s", err)
	}

	if err := d.Set("public_key", r.PublicKey); err != nil {
		return nil, fmt.Errorf("error setting `public_key`: %s", err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"org_id":     orgID,
		"api_key_id": r.GetId(),
	}))

	return []*schema.ResourceData{d}, nil
}

func flattenOrgAPIKeyRoles(orgID string, apiKeyRoles []admin.CloudAccessRoleAssignment) []string {
	flattenedOrgRoles := make([]string, 0, len(apiKeyRoles))
	for _, role := range apiKeyRoles {
		if strings.HasPrefix(role.GetRoleName(), "ORG_") && role.GetOrgId() == orgID {
			flattenedOrgRoles = append(flattenedOrgRoles, role.GetRoleName())
		}
	}
	return flattenedOrgRoles
}
