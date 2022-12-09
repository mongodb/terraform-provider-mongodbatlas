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

func resourceMongoDBAtlasAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasAPIKeyCreate,
		ReadContext:   resourceMongoDBAtlasAPIKeyRead,
		UpdateContext: resourceMongoDBAtlasAPIKeyUpdate,
		DeleteContext: resourceMongoDBAtlasAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasAPIKeyImportState,
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

func resourceMongoDBAtlasAPIKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	orgID := d.Get("org_id").(string)
	createRequest := new(matlas.APIKeyInput)

	createRequest.Desc = d.Get("description").(string)

	createRequest.Roles = expandStringList(d.Get("role_names").(*schema.Set).List())

	apiKey, resp, err := conn.APIKeys.Create(ctx, orgID, createRequest)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("error create API key: %s", err))
	}

	if err := d.Set("private_key", apiKey.PrivateKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id":     orgID,
		"api_key_id": apiKey.ID,
	}))

	return resourceMongoDBAtlasAPIKeyRead(ctx, d, meta)
}

func resourceMongoDBAtlasAPIKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	apiKeyID := ids["api_key_id"]

	apiKey, _, err := conn.APIKeys.Get(ctx, orgID, apiKeyID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}

	if err := d.Set("description", apiKey.Desc); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `description`: %s", err))
	}

	if err := d.Set("public_key", apiKey.PublicKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	if err := d.Set("role_names", flattenOrgAPIKeyRoles(orgID, apiKey.Roles)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `roles`: %s", err))
	}

	return nil
}

func resourceMongoDBAtlasAPIKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	apiKeyID := ids["api_key_id"]

	updateRequest := new(matlas.APIKeyInput)

	if d.HasChange("description") || d.HasChange("role_names") {
		updateRequest.Desc = d.Get("description").(string)

		updateRequest.Roles = expandStringList(d.Get("role_names").(*schema.Set).List())

		_, _, err := conn.APIKeys.Update(ctx, orgID, apiKeyID, updateRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating API key: %s", err))
		}
	}

	return resourceMongoDBAtlasAPIKeyRead(ctx, d, meta)
}

func resourceMongoDBAtlasAPIKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	orgID := ids["org_id"]
	apiKeyID := ids["api_key_id"]

	_, err := conn.APIKeys.Delete(ctx, orgID, apiKeyID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting custom db role information: %s", err))
	}
	return nil
}

func resourceMongoDBAtlasAPIKeyImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a api key use the format {org_id}-{api_key_id}")
	}

	orgID := parts[0]
	apiKeyID := parts[1]

	r, _, err := conn.APIKeys.Get(ctx, orgID, apiKeyID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import api key %s in project %s, error: %s", orgID, apiKeyID, err)
	}

	if err := d.Set("description", r.Desc); err != nil {
		return nil, fmt.Errorf("error setting `description`: %s", err)
	}

	if err := d.Set("public_key", r.PublicKey); err != nil {
		return nil, fmt.Errorf("error setting `public_key`: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"org_id":     orgID,
		"api_key_id": r.ID,
	}))

	return []*schema.ResourceData{d}, nil
}

func flattenOrgAPIKeys(ctx context.Context, conn *matlas.Client, orgID string, apiKeys []matlas.APIKey) []map[string]interface{} {
	var results []map[string]interface{}

	if len(apiKeys) > 0 {
		results = make([]map[string]interface{}, len(apiKeys))
		for k, apiKey := range apiKeys {
			results[k] = map[string]interface{}{
				"api_key_id":  apiKey.ID,
				"description": apiKey.Desc,
				"public_key":  apiKey.PublicKey,
				"role_names":  flattenOrgAPIKeyRoles(orgID, apiKey.Roles),
			}
		}
	}
	return results
}

func flattenOrgAPIKeyRoles(orgID string, apiKeyRoles []matlas.AtlasRole) []string {
	if len(apiKeyRoles) == 0 {
		return nil
	}

	flattenedOrgRoles := []string{}

	for _, role := range apiKeyRoles {
		if strings.HasPrefix(role.RoleName, "ORG_") && role.OrgID == orgID {
			flattenedOrgRoles = append(flattenedOrgRoles, role.RoleName)
		}
	}

	return flattenedOrgRoles
}
