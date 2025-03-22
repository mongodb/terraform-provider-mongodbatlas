package apikey

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func ResourceProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectAssignmentCreate,
		ReadContext:   resourceProjectAssignmentRead,
		UpdateContext: resourceProjectAssignmentUpdate,
		DeleteContext: resourceProjectAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceProjectAssignmentImport,
		},
		Schema: map[string]*schema.Schema{
			"api_key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_names": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceProjectAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	apiKeyID := d.Get("api_key_id").(string)

	roles := conversion.ExpandStringList(d.Get("role_names").(*schema.Set).List())
	createRequest := &[]admin.UserAccessRoleAssignment{
		{
			Roles:  &roles,
			UserId: &apiKeyID,
		},
	}

	apiKeyAssignment, resp, err := connV2.ProgrammaticAPIKeysApi.AddProjectApiKey(ctx, projectID, apiKeyID, createRequest).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error create API key assignment: %s", err))
	}

	log.Printf("apiKeyAssignment: %+v", apiKeyAssignment)

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"api_key_id": apiKeyID,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	apiKeyID := ids["api_key_id"]

	apiKeys, resp, err := connV2.ProgrammaticAPIKeysApi.ListProjectApiKeys(ctx, projectID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) || validate.StatusBadRequest(resp) {
			log.Printf("warning API key deleted will recreate: %s \n", err.Error())
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}
	apiKeyUserDetails := apiKeys.GetResults()

	for _, apiKey := range apiKeyUserDetails {
		if apiKey.GetId() == apiKeyID {
			d.Set("role_names", flattenAPIKeyProjectRoles(projectID, apiKey.GetRoles()))
		}
	}

	return nil
}

func resourceProjectAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	apiKeyID := ids["api_key_id"]

	if d.HasChange("role_names") {
		roles := conversion.ExpandStringList(d.Get("role_names").(*schema.Set).List())
		updateRequest := &admin.UpdateAtlasProjectApiKey{
			Roles: &roles,
		}
		_, _, err := connV2.ProgrammaticAPIKeysApi.UpdateApiKeyRoles(ctx, projectID, apiKeyID, updateRequest).Execute()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating API key assignments: %s", err))
		}
	}
	return resourceRead(ctx, d, meta)
}

func resourceProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	apiKeyID := ids["api_key_id"]

	_, _, err := connV2.ProgrammaticAPIKeysApi.RemoveProjectApiKey(ctx, projectID, apiKeyID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error removing API Key project role assignments: %s", err))
	}
	return nil
}

func resourceProjectAssignmentImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a api key use the format {org_id}-{api_key_id}")
	}

	projectID := parts[0]
	apiKeyID := parts[1]

	apiKeys, _, err := connV2.ProgrammaticAPIKeysApi.ListProjectApiKeys(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import api key assignment for %s in project %s, error: %s", apiKeyID, projectID, err)
	}
	apiKeyUserDetails := apiKeys.GetResults()

	for _, apiKey := range apiKeyUserDetails {
		if apiKey.GetId() == apiKeyID {
			d.Set("role_names", flattenAPIKeyProjectRoles(projectID, apiKey.GetRoles()))
		}
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
		"api_key_id": apiKeyID,
	}))

	return []*schema.ResourceData{d}, nil
}

func flattenAPIKeyProjectRoles(projectID string, apiKeyRoles []admin.CloudAccessRoleAssignment) []string {
	flattenedProjectRoles := make([]string, 0, len(apiKeyRoles))
	for _, role := range apiKeyRoles {
		if !strings.HasPrefix(role.GetRoleName(), "ORG_") && role.GetGroupId() == projectID {
			flattenedProjectRoles = append(flattenedProjectRoles, role.GetRoleName())
		}
	}
	return flattenedProjectRoles
}
