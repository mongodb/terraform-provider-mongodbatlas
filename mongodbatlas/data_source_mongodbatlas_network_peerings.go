package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasNetworkPeerings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasNetworkPeeringsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peering_id": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasNetworkPeeringsRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	peers, _, err := conn.Peers.List(context.Background(), projectID, nil)
	if err != nil {
		return fmt.Errorf("error getting network peering connections information: %s", err)
	}

	if err := d.Set("results", flattenNetworkPeerings(peers)); err != nil {
		return fmt.Errorf("error setting `result` for network peering connections: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenNetworkPeerings(peers []matlas.Peer) []map[string]interface{} {
	var peersMap []map[string]interface{}

	if len(peers) > 0 {
		peersMap = make([]map[string]interface{}, len(peers))

		for k, peer := range peers {
			peersMap[k] = map[string]interface{}{
				"peering_id":             peer.ID,
				"container_id":           peer.ContainerID,
				"accepter_region_name":   peer.AccepterRegionName,
				"aws_account_id":         peer.AWSAccountID,
				"provider_name":          getProviderNameByPeer(peer),
				"route_table_cidr_block": peer.RouteTableCIDRBlock,
				"vpc_id":                 peer.VpcID,
				"connection_id":          peer.ConnectionID,
				"error_state_name":       peer.ErrorStateName,
				"status_name":            peer.StatusName,
				"atlas_cidr_block":       peer.AtlasCIDRBlock,
				"azure_directory_id":     peer.AzureDirectoryID,
				"azure_subscription_id":  peer.AzureSubscriptionID,
				"resource_group_name":    peer.ResourceGroupName,
				"vnet_name":              peer.VNetName,
				"error_state":            peer.ErrorState,
				"status":                 peer.Status,
				"gcp_project_id":         peer.GCPProjectID,
				"network_name":           peer.NetworkName,
				"error_message":          peer.ErrorMessage,
			}
		}
	}
	return peersMap
}

func getProviderNameByPeer(peer matlas.Peer) string {
	provider := "AWS"
	if peer.VNetName != "" {
		provider = "AZURE"
	} else if peer.NetworkName != "" {
		provider = "GCP"
	}
	return provider
}
