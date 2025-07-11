package projectapikey

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportState,
		},
		Schema: map[string]*schema.Schema{
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
			"project_assignment": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_names": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := validateUniqueProjectIDs(d); err != nil {
		return diag.FromErr(err)
	}
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	val, ok := d.GetOk("project_assignment")
	if !ok {
		return diag.FromErr(errors.New("project_assignment not found"))
	}
	assignments := expandProjectAssignments(val.(*schema.Set))
	projectIDs := make([]string, 0, len(assignments))
	for projectID := range assignments {
		projectIDs = append(projectIDs, projectID)
	}

	req := &admin.CreateAtlasProjectApiKey{
		Desc:  d.Get("description").(string),
		Roles: assignments[projectIDs[0]],
	}
	ret, _, err := connV2.ProgrammaticAPIKeysApi.CreateProjectApiKey(ctx, projectIDs[0], req).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	apiKeyID := ret.GetId()
	d.SetId(conversion.EncodeStateID(map[string]string{
		"api_key_id": apiKeyID,
	}))
	if err := d.Set("public_key", ret.GetPublicKey()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}
	if err := d.Set("private_key", ret.GetPrivateKey()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `private_key`: %s", err))
	}

	for _, projectID := range projectIDs[1:] {
		roles := assignments[projectID]
		req := &[]admin.UserAccessRoleAssignment{{Roles: &roles}}
		if _, err := connV2.ProgrammaticAPIKeysApi.AddProjectApiKey(ctx, projectID, apiKeyID, req).Execute(); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	apiKeyID := ids["api_key_id"]

	details, _, err := getKeyDetails(ctx, connV2, apiKeyID)
	if err != nil {
		return diag.FromErr(err)
	}
	if details == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("api_key_id", details.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `api_key_id`: %s", err))
	}

	if err := d.Set("description", details.GetDesc()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `description`: %s", err))
	}

	if err := d.Set("public_key", details.GetPublicKey()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	if err := d.Set("project_assignment", flattenProjectAssignments(details.GetRoles())); err != nil {
		return diag.Errorf("error setting `project_assignment` : %s", err)
	}
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := validateUniqueProjectIDs(d); err != nil {
		return diag.FromErr(err)
	}

	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	apiKeyID := ids["api_key_id"]

	details, orgID, err := getKeyDetails(ctx, connV2, apiKeyID)
	if err != nil {
		return diag.FromErr(err)
	}
	if details == nil {
		return diag.Errorf("error updating project api_key (%s): not found", apiKeyID)
	}

	if d.HasChange("project_assignment") {
		add, remove, update := getAssignmentChanges(d)

		for projectID := range remove {
			_, err := connV2.ProgrammaticAPIKeysApi.RemoveProjectApiKey(ctx, projectID, apiKeyID).Execute()
			if err != nil {
				if admin.IsErrorCode(err, "GROUP_NOT_FOUND") {
					continue // allows removing assignment for a project that has been deleted
				}
				return diag.Errorf("error removing project_api_key(%s) from project(%s): %s", apiKeyID, projectID, err)
			}
		}

		for projectID, roles := range add {
			req := &[]admin.UserAccessRoleAssignment{{Roles: &roles}}
			_, err := connV2.ProgrammaticAPIKeysApi.AddProjectApiKey(ctx, projectID, apiKeyID, req).Execute()
			if err != nil {
				return diag.Errorf("error adding project_api_key(%s) to project(%s): %s", apiKeyID, projectID, err)
			}
		}

		for projectID, roles := range update {
			req := &admin.UpdateAtlasProjectApiKey{Roles: &roles}
			_, _, err := connV2.ProgrammaticAPIKeysApi.UpdateApiKeyRoles(ctx, projectID, apiKeyID, req).Execute()
			if err != nil {
				return diag.Errorf("error changing project_api_key(%s) in project(%s): %s", apiKeyID, projectID, err)
			}
		}
	}

	if d.HasChange("description") {
		req := &admin.UpdateAtlasOrganizationApiKey{Desc: conversion.StringPtr(d.Get("description").(string))}
		if _, _, err := connV2.ProgrammaticAPIKeysApi.UpdateApiKey(ctx, orgID, apiKeyID, req).Execute(); err != nil {
			return diag.Errorf("error updating description in api key(%s): %s", apiKeyID, err)
		}
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	apiKeyID := ids["api_key_id"]
	details, orgID, err := getKeyDetails(ctx, connV2, apiKeyID)
	if err != nil {
		return diag.FromErr(err)
	}
	if details != nil && orgID != "" {
		if _, err = connV2.ProgrammaticAPIKeysApi.DeleteApiKey(ctx, orgID, apiKeyID).Execute(); err != nil {
			return diag.FromErr(fmt.Errorf("error deleting project key (%s): %s", apiKeyID, err))
		}
	}
	d.SetId("")
	return nil
}

func resourceImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a api key use the format {project_id}-{api_key_id}")
	}

	// projectID is not needed for import any more, but kept to maintain import format and avoid breaking changes
	apiKeyID := parts[1]

	d.SetId(conversion.EncodeStateID(map[string]string{
		"api_key_id": apiKeyID,
	}))
	return []*schema.ResourceData{d}, nil
}
