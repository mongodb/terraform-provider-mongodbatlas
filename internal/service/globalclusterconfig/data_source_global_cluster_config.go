package globalclusterconfig

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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
			"custom_zone_mapping": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV220240530 := meta.(*config.MongoDBClient).AtlasV220240530
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	globalCluster, resp, err := connV220240530.GlobalClustersApi.GetManagedNamespace(ctx, projectID, clusterName).Execute()
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorGlobalClusterRead, clusterName, err))
	}

	if err := d.Set("managed_namespaces", flattenManagedNamespaces(globalCluster.GetManagedNamespaces())); err != nil {
		return diag.FromErr(fmt.Errorf(errorGlobalClusterRead, clusterName, err))
	}

	if err := d.Set("custom_zone_mapping", globalCluster.GetCustomZoneMapping()); err != nil {
		return diag.FromErr(fmt.Errorf(errorGlobalClusterRead, clusterName, err))
	}

	d.SetId(clusterName)

	return nil
}
