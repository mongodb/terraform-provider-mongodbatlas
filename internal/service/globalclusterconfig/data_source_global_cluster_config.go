package globalclusterconfig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"managed_namespaces": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collection": {
							Type:     schema.TypeString,
							Required: true,
						},
						"custom_shard_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"db": {
							Type:     schema.TypeString,
							Required: true,
						},
						"is_custom_shard_key_hashed": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_shard_key_unique": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"custom_zone_mapping_zone_id": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	notFound, err := readGlobalClusterConfig(ctx, meta, projectID, clusterName, d)
	if err != nil {
		return diag.FromErr(err)
	}
	if notFound {
		d.SetId("")
	} else {
		d.SetId(clusterName)
	}
	return nil
}
