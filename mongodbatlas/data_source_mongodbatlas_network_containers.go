package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasNetworkContainers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasNetworkContainersRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"atlas_cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_subscription_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provisioned": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"gcp_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vnet_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"regions": {
							Type:     schema.TypeList,
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

func dataSourceMongoDBAtlasNetworkContainersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	containers, _, err := conn.Containers.List(ctx, projectID, &matlas.ContainersListOptions{
		ProviderName: d.Get("provider_name").(string),
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting network peering containers information: %s", err))
	}

	if err := d.Set("results", flattenNetworkContainers(containers)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for network containers: %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenNetworkContainers(containers []matlas.Container) []map[string]interface{} {
	var containersMap []map[string]interface{}

	if len(containers) > 0 {
		containersMap = make([]map[string]interface{}, len(containers))

		for i := range containers {
			containersMap[i] = map[string]interface{}{
				"id":                    containers[i].ID,
				"atlas_cidr_block":      containers[i].AtlasCIDRBlock,
				"provider_name":         containers[i].ProviderName,
				"region_name":           containers[i].RegionName,
				"region":                containers[i].Region,
				"azure_subscription_id": containers[i].AzureSubscriptionID,
				"provisioned":           containers[i].Provisioned,
				"gcp_project_id":        containers[i].GCPProjectID,
				"network_name":          containers[i].NetworkName,
				"vpc_id":                containers[i].VPCID,
				"vnet_name":             containers[i].VNetName,
				"regions":               containers[i].Regions,
			}
		}
	}

	return containersMap
}
