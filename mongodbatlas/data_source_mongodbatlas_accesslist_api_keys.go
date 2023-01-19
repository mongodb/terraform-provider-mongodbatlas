package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasAccessListAPIKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasAccessListAPIKeysRead,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"items_per_page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_used": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_used_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasAccessListAPIKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	orgID := d.Get("org_id").(string)
	apiKeyID := d.Get("api_key_id").(string)

	accessListAPIKeys, _, err := conn.AccessListAPIKeys.List(ctx, orgID, apiKeyID, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting access list api keys information: %s", err))
	}

	if err := d.Set("results", flattenAccessListAPIKeys(ctx, conn, orgID, accessListAPIKeys.Results)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}
