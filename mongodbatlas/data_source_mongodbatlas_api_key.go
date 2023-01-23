package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasAPIKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasAPIKeyRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasAPIKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	orgID := d.Get("org_id").(string)
	apiKeyID := d.Get("api_key_id").(string)
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

	d.SetId(resource.UniqueId())

	return nil
}
