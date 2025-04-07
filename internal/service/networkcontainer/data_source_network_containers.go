package networkcontainer

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312001/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
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

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	// Returns all providers independently of provider
	containers, _, err := connV2.NetworkPeeringApi.ListPeeringContainers(ctx, projectID).Execute()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting network peering containers information: %s", err))
	}

	// Necessary to keep same behavior. Only have to return the containers of the specified provider. This was the behavior of old SDK
	containersOfSpecifiedProvider := filterContainersByProvider(containers.GetResults(), d.Get("provider_name").(string))

	if err := d.Set("results", flattenNetworkContainers(containersOfSpecifiedProvider)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for network containers: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func filterContainersByProvider(containers []admin.CloudProviderContainer, provider string) []admin.CloudProviderContainer {
	result := make([]admin.CloudProviderContainer, 0)
	for _, container := range containers {
		if container.GetProviderName() == provider {
			result = append(result, container)
		}
	}
	return result
}

func flattenNetworkContainers(containers []admin.CloudProviderContainer) []map[string]any {
	var containersMap []map[string]any

	if len(containers) > 0 {
		containersMap = make([]map[string]any, len(containers))

		for i := range containers {
			containersMap[i] = map[string]any{
				"id":                    containers[i].GetId(),
				"atlas_cidr_block":      containers[i].GetAtlasCidrBlock(),
				"provider_name":         containers[i].GetProviderName(),
				"region_name":           containers[i].GetRegionName(),
				"region":                containers[i].GetRegion(),
				"azure_subscription_id": containers[i].GetAzureSubscriptionId(),
				"provisioned":           containers[i].GetProvisioned(),
				"gcp_project_id":        containers[i].GetGcpProjectId(),
				"network_name":          containers[i].GetNetworkName(),
				"vpc_id":                containers[i].GetVpcId(),
				"vnet_name":             containers[i].GetVnetName(),
				"regions":               containers[i].GetRegions(),
			}
		}
	}

	return containersMap
}
