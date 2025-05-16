package networkpeering

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"

	"go.mongodb.org/atlas-sdk/v20250312003/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
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

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)

	peers, _, err := conn.NetworkPeeringApi.ListPeeringConnections(ctx, projectID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting network peering connections information: %s", err))
	}
	peersMap, err := flattenNetworkPeerings(ctx, conn.NetworkPeeringApi, peers.GetResults(), projectID)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("results", peersMap); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `result` for network peering connections: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenNetworkPeerings(ctx context.Context, conn admin.NetworkPeeringApi, peers []admin.BaseNetworkPeeringConnectionSettings, projectID string) ([]map[string]any, error) {
	var peersMap []map[string]any

	if len(peers) > 0 {
		peersMap = make([]map[string]any, len(peers))
		for i := range peers {
			p := peers[i]
			accepterRegionName, err := ensureAccepterRegionName(ctx, &p, conn, projectID)
			if err != nil {
				return nil, err
			}
			atlasCidrBlock, err := readAtlasCidrBlock(ctx, conn, projectID, p.GetContainerId())
			if err != nil {
				return nil, err
			}
			peersMap[i] = map[string]any{
				"peering_id":             p.GetId(),
				"container_id":           p.GetContainerId(),
				"accepter_region_name":   accepterRegionName,
				"aws_account_id":         p.GetAwsAccountId(),
				"provider_name":          getProviderNameByPeer(&p),
				"route_table_cidr_block": p.GetRouteTableCidrBlock(),
				"vpc_id":                 p.GetVpcId(),
				"connection_id":          p.GetConnectionId(),
				"error_state_name":       p.GetErrorStateName(),
				"status_name":            p.GetStatusName(),
				"atlas_cidr_block":       atlasCidrBlock,
				"azure_directory_id":     p.GetAzureDirectoryId(),
				"azure_subscription_id":  p.GetAzureSubscriptionId(),
				"resource_group_name":    p.GetResourceGroupName(),
				"vnet_name":              p.GetVnetName(),
				"error_state":            p.GetErrorState(),
				"status":                 p.GetStatus(),
				"gcp_project_id":         p.GetGcpProjectId(),
				"network_name":           p.GetNetworkName(),
				"error_message":          p.GetErrorMessage(),
			}
		}
	}

	return peersMap, nil
}

func getProviderNameByPeer(peer *admin.BaseNetworkPeeringConnectionSettings) string {
	provider := "AWS"
	if peer.GetVnetName() != "" {
		provider = "AZURE"
	} else if peer.GetNetworkName() != "" {
		provider = "GCP"
	}

	return provider
}
