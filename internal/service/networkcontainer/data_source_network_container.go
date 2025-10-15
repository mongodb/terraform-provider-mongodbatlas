package networkcontainer

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
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
			"container_id": {
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
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	containerID := conversion.GetEncodedID(d.Get("container_id").(string), "container_id")

	container, resp, err := connV2.NetworkPeeringApi.GetGroupContainer(ctx, projectID, containerID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			return nil
		}
		return diag.FromErr(fmt.Errorf(ErrorContainerRead, containerID, err))
	}

	if err := d.Set("atlas_cidr_block", container.GetAtlasCidrBlock()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `atlas_cidr_block` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("provider_name", container.GetProviderName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `provider_name` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("region_name", container.GetRegionName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `region_name` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("region", container.GetRegion()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `region` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("azure_subscription_id", container.GetAzureSubscriptionId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `azure_subscription_id` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("provisioned", container.GetProvisioned()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `provisioned` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("gcp_project_id", container.GetGcpProjectId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `gcp_project_id` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("network_name", container.GetNetworkName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `network_name` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("vpc_id", container.GetVpcId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `vpc_id` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("vnet_name", container.GetVnetName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `vnet_name` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("regions", container.GetRegions()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `regions` for Network Container (%s): %s", d.Id(), err))
	}

	d.SetId(container.GetId())

	return nil
}
