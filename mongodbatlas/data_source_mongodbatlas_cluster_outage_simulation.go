package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasClusterOutageSimulation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasClusterOutageSimulationRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"outage_filters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"start_request_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"simulation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasClusterOutageSimulationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas

	projectID, projectIDOk := d.GetOk("project_id")
	clusterName, clusterNameOk := d.GetOk("cluster_name")

	if !(projectIDOk && clusterNameOk) {
		return diag.Errorf("project_id and cluster_name must be configured")
	}

	outageSimulation, _, err := conn.ClusterOutageSimulation.GetOutageSimulation(ctx, projectID.(string), clusterName.(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterOutageSimulationRead, projectID, clusterName, err))
	}

	err = convertOutageSimulationToSchema(outageSimulation, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID.(string),
		"cluster_name": clusterName.(string),
	}))

	return nil
}
