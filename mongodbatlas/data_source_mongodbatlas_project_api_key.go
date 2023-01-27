package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasProjectAPIKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasProjectAPIKeyRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
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
			"private_key": {
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

func dataSourceMongoDBAtlasProjectAPIKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	apiKeyID := d.Get("api_key_id").(string)
	projectAPIKeys, _, err := conn.ProjectAPIKeys.List(ctx, projectID, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}

	for _, val := range projectAPIKeys {
		if val.ID == apiKeyID {
			if err := d.Set("description", val.Desc); err != nil {
				return diag.FromErr(fmt.Errorf("error setting `description`: %s", err))
			}

			if err := d.Set("public_key", val.PublicKey); err != nil {
				return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
			}

			if err := d.Set("private_key", val.PrivateKey); err != nil {
				return diag.FromErr(fmt.Errorf("error setting `private_key`: %s", err))
			}

			if err := d.Set("role_names", flattenProjectAPIKeyRoles(projectID, val.Roles)); err != nil {
				return diag.FromErr(fmt.Errorf("error setting `roles`: %s", err))
			}
		}
	}

	d.SetId(resource.UniqueId())

	return nil
}
