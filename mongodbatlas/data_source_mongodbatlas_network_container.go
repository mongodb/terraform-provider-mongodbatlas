package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasNetworkContainer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasNetworkContainerRead,
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

func dataSourceMongoDBAtlasNetworkContainerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	containerID := getEncodedID(d.Get("container_id").(string), "container_id")

	container, resp, err := conn.Containers.Get(ctx, projectID, containerID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorContainerRead, containerID, err))
	}

	if err := d.Set("atlas_cidr_block", container.AtlasCIDRBlock); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `atlas_cidr_block` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("provider_name", container.ProviderName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `provider_name` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("region_name", container.RegionName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `region_name` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("region", container.Region); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `region` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("azure_subscription_id", container.AzureSubscriptionID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `azure_subscription_id` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("provisioned", container.Provisioned); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `provisioned` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("gcp_project_id", container.GCPProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `gcp_project_id` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("network_name", container.NetworkName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `network_name` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("gcp_project_id", container.GCPProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `gcp_project_id` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("vpc_id", container.VPCID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `vpc_id` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("vnet_name", container.VNetName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `vnet_name` for Network Container (%s): %s", d.Id(), err))
	}

	if err = d.Set("regions", container.Regions); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `regions` for Network Container (%s): %s", d.Id(), err))
	}

	d.SetId(container.ID)

	return nil
}
