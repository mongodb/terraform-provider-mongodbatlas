package clusteroutagesimulation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID, projectIDOk := d.GetOk("project_id")
	clusterName, clusterNameOk := d.GetOk("cluster_name")

	if !projectIDOk || !clusterNameOk {
		return diag.Errorf("project_id and cluster_name must be configured")
	}

	outageSimulation, _, err := connV2.ClusterOutageSimulationApi.GetOutageSimulation(ctx, projectID.(string), clusterName.(string)).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorClusterOutageSimulationRead, projectID, clusterName, err))
	}

	if err = convertOutageSimulationToSchema(outageSimulation, d); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID.(string),
		"cluster_name": clusterName.(string),
	}))

	return nil
}
