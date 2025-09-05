package apikey

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	orgID := d.Get("org_id").(string)
	apiKeyID := d.Get("api_key_id").(string)
	apiKey, _, err := connV2.ProgrammaticAPIKeysApi.GetOrgApiKey(ctx, orgID, apiKeyID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}

	if err := d.Set("description", apiKey.Desc); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `description`: %s", err))
	}

	if err := d.Set("public_key", apiKey.PublicKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	if err := d.Set("role_names", flattenOrgAPIKeyRoles(orgID, apiKey.GetRoles())); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `roles`: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}
