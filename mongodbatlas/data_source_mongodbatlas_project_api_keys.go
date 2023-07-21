package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasProjectAPIKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasProjectAPIKeysRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_key_id": {
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
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasProjectAPIKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	projectID := d.Get("project_id").(string)

	apiKeys, _, err := conn.ProjectAPIKeys.List(ctx, projectID, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api keys information: %s", err))
	}

	if err := d.Set("results", flattenProjectAPIKeys(ctx, conn, projectID, apiKeys)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}
