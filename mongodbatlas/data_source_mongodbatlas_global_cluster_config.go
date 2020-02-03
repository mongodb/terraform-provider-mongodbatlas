package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasGlobalCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasGlobalClusterRead,
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

func dataSourceMongoDBAtlasGlobalClusterRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	globalCluster, resp, err := conn.GlobalClusters.Get(context.Background(), projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {

			return nil
		}
		return fmt.Errorf(errorGlobalClusterRead, clusterName, err)
	}

	if err := d.Set("managed_namespaces", flattenManagedNamespaces(globalCluster.ManagedNamespaces)); err != nil {
		return fmt.Errorf(errorGlobalClusterRead, clusterName, err)
	}

	if err := d.Set("custom_zone_mapping", globalCluster.CustomZoneMapping); err != nil {
		return fmt.Errorf(errorGlobalClusterRead, clusterName, err)
	}

	d.SetId(clusterName)

	return nil
}
