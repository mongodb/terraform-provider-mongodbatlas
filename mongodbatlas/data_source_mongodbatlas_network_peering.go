package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasNetworkPeering() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasNetworkPeeringRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"peering_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"container_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"accepter_region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_table_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_state_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"atlas_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"atlas_cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azure_directory_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azure_subscription_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vnet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
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
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasNetworkPeeringRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	peerID := d.Get("peering_id").(string)

	peer, resp, err := conn.Peers.Get(context.Background(), projectID, peerID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf(errorPeersRead, peerID, err)
	}

	//Workaround until fix.
	if peer.AccepterRegionName != "" {
		if err := d.Set("accepter_region_name", peer.AccepterRegionName); err != nil {
			return fmt.Errorf("error setting `accepter_region_name` for Network Peering Connection (%s): %s", peerID, err)
		}
	}

	if err := d.Set("aws_account_id", peer.AWSAccountId); err != nil {
		return fmt.Errorf("error setting `aws_account_id` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("container_id", peer.ContainerID); err != nil {
		return fmt.Errorf("error setting `container_id` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("route_table_cidr_block", peer.RouteTableCIDRBlock); err != nil {
		return fmt.Errorf("error setting `route_table_cidr_block` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("vpc_id", peer.VpcID); err != nil {
		return fmt.Errorf("error setting `vpc_id` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("connection_id", peer.ConnectionID); err != nil {
		return fmt.Errorf("error setting `connection_id` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("error_state_name", peer.ErrorStateName); err != nil {
		return fmt.Errorf("error setting `error_state_name` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("atlas_id", peer.ID); err != nil {
		return fmt.Errorf("error setting `atlas_id` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("status_name", peer.StatusName); err != nil {
		return fmt.Errorf("error setting `status_name` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("atlas_cidr_block", peer.AtlasCIDRBlock); err != nil {
		return fmt.Errorf("error setting `atlas_cidr_block` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("azure_directory_id", peer.AzureDirectoryID); err != nil {
		return fmt.Errorf("error setting `azure_directory_id` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("azure_subscription_id", peer.AzureSubscriptionId); err != nil {
		return fmt.Errorf("error setting `azure_subscription_id` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("resource_group_name", peer.ResourceGroupName); err != nil {
		return fmt.Errorf("error setting `resource_group_name` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("vnet_name", peer.VNetName); err != nil {
		return fmt.Errorf("error setting `vnet_name` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("error_state", peer.ErrorState); err != nil {
		return fmt.Errorf("error setting `error_state` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("status", peer.Status); err != nil {
		return fmt.Errorf("error setting `status` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("gcp_project_id", peer.GCPProjectID); err != nil {
		return fmt.Errorf("error setting `gcp_project_id` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("network_name", peer.NetworkName); err != nil {
		return fmt.Errorf("error setting `network_name` for Network Peering Connection (%s): %s", peerID, err)
	}

	if err := d.Set("error_message", peer.ErrorMessage); err != nil {
		return fmt.Errorf("error setting `error_message` for Network Peering Connection (%s): %s", peerID, err)
	}

	provider := "AWS"
	if peer.VNetName != "" {
		provider = "AZURE"
	} else if peer.NetworkName != "" {
		provider = "GCP"
	}

	if err := d.Set("provider_name", provider); err != nil {
		return fmt.Errorf("[WARN] Error setting provider_name for (%s): %s", d.Id(), err)
	}

	d.SetId(peer.ID)

	return nil
}
