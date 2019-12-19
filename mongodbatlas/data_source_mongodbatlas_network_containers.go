package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasNetworkContainers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasNetworkContainersRead,
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
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasNetworkContainersRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	containers, _, err := conn.Containers.List(context.Background(), projectID, &matlas.ContainersListOptions{
		ProviderName: d.Get("provider_name").(string),
	})

	if err != nil {
		return fmt.Errorf("error getting network peering containers information: %s", err)
	}

	if err := d.Set("results", flattenNetworkContainers(containers)); err != nil {
		return fmt.Errorf("error setting `result` for network containers: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenNetworkContainers(containers []matlas.Container) []map[string]interface{} {
	var containersMap []map[string]interface{}

	if len(containers) > 0 {
		containersMap = make([]map[string]interface{}, len(containers))

		for k, container := range containers {
			containersMap[k] = map[string]interface{}{
				"id":                    container.ID,
				"atlas_cidr_block":      container.AtlasCIDRBlock,
				"provider_name":         container.ProviderName,
				"region_name":           container.RegionName,
				"region":                container.Region,
				"azure_subscription_id": container.AzureSubscriptionID,
				"provisioned":           container.Provisioned,
				"gcp_project_id":        container.GCPProjectID,
				"network_name":          container.NetworkName,
				"vpc_id":                container.VPCID,
				"vnet_name":             container.VNetName,
			}
		}
	}
	return containersMap
}
