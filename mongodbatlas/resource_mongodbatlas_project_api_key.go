package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasProjectAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasProjectAPIKeyCreate,
		ReadContext:   resourceMongoDBAtlasProjectAPIKeyRead,
		UpdateContext: resourceMongoDBAtlasProjectAPIKeyUpdate,
		DeleteContext: resourceMongoDBAtlasProjectAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasProjectAPIKeyImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
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

func resourceMongoDBAtlasProjectAPIKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	createRequest := new(matlas.APIKeyInput)

	createRequest.Desc = d.Get("description").(string)

	createRequest.Roles = expandStringList(d.Get("role_names").(*schema.Set).List())

	apiKey, resp, err := conn.ProjectAPIKeys.Create(ctx, projectID, createRequest)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error create API key: %s", err))
	}

	if err := d.Set("public_key", apiKey.PublicKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	if err := d.Set("private_key", apiKey.PrivateKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `private_key`: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"api_key_id": apiKey.ID,
	}))

	return resourceMongoDBAtlasProjectAPIKeyRead(ctx, d, meta)
}

func resourceMongoDBAtlasProjectAPIKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	apiKeyID := ids["api_key_id"]

	projectAPIKeys, _, err := conn.ProjectAPIKeys.List(ctx, projectID, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}
	for _, val := range projectAPIKeys {
		if val.ID == apiKeyID {
			if err := d.Set("api_key_id", val.ID); err != nil {
				return diag.FromErr(fmt.Errorf("error setting `api_key_id`: %s", err))
			}

			if err := d.Set("description", val.Desc); err != nil {
				return diag.FromErr(fmt.Errorf("error setting `description`: %s", err))
			}

			if err := d.Set("public_key", val.PublicKey); err != nil {
				return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
			}

			if err := d.Set("role_names", flattenProjectAPIKeyRoles(projectID, val.Roles)); err != nil {
				return diag.FromErr(fmt.Errorf("error setting `roles`: %s", err))
			}
		}
	}

	if err := d.Set("project_id", projectID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `project_id`: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"api_key_id": apiKeyID,
	}))

	return nil
}

func resourceMongoDBAtlasProjectAPIKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	apiKeyID := ids["api_key_id"]

	updateRequest := new(matlas.AssignAPIKey)

	if d.HasChange("role_names") {
		updateRequest.Roles = expandStringList(d.Get("role_names").(*schema.Set).List())

		_, err := conn.ProjectAPIKeys.Assign(ctx, projectID, apiKeyID, updateRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating API key: %s", err))
		}
	}

	return resourceMongoDBAtlasProjectAPIKeyRead(ctx, d, meta)
}

func resourceMongoDBAtlasProjectAPIKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	apiKeyID := ids["api_key_id"]

	_, err := conn.ProjectAPIKeys.Unassign(ctx, projectID, apiKeyID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting project api key: %s", err))
	}
	d.SetId("")
	return nil
}

func resourceMongoDBAtlasProjectAPIKeyImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a api key use the format {org_id}-{api_key_id}")
	}

	projectID := parts[0]
	apiKeyID := parts[1]

	projectAPIKeys, _, err := conn.ProjectAPIKeys.List(ctx, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't import api key %s in project %s, error: %s", projectID, apiKeyID, err)
	}
	for _, val := range projectAPIKeys {
		if val.ID == apiKeyID {
			if err := d.Set("description", val.Desc); err != nil {
				return nil, fmt.Errorf("error setting `description`: %s", err)
			}

			if err := d.Set("public_key", val.PublicKey); err != nil {
				return nil, fmt.Errorf("error setting `public_key`: %s", err)
			}

			d.SetId(encodeStateID(map[string]string{
				"project_id": projectID,
				"api_key_id": val.ID,
			}))
		}
	}
	return []*schema.ResourceData{d}, nil
}

func flattenProjectAPIKeys(ctx context.Context, conn *matlas.Client, projectID string, apiKeys []matlas.APIKey) []map[string]interface{} {
	var results []map[string]interface{}

	if len(apiKeys) > 0 {
		results = make([]map[string]interface{}, len(apiKeys))
		for k, apiKey := range apiKeys {
			results[k] = map[string]interface{}{
				"api_key_id":  apiKey.ID,
				"description": apiKey.Desc,
				"public_key":  apiKey.PublicKey,
				"private_key": apiKey.PrivateKey,
				"role_names":  flattenProjectAPIKeyRoles(projectID, apiKey.Roles),
			}
		}
	}
	return results
}

func flattenProjectAPIKeyRoles(projectID string, apiKeyRoles []matlas.AtlasRole) []string {
	if len(apiKeyRoles) == 0 {
		return nil
	}

	flattenedOrgRoles := []string{}

	for _, role := range apiKeyRoles {
		if strings.HasPrefix(role.RoleName, "GROUP_") && role.GroupID == projectID {
			flattenedOrgRoles = append(flattenedOrgRoles, role.RoleName)
		}
	}

	return flattenedOrgRoles
}
