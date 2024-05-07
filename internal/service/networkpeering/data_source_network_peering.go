package networkpeering

import (
	"context"
	"fmt"
	"net/http"

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

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	peerID := conversion.GetEncodedID(d.Get("peering_id").(string), "peer_id")

	peer, resp, err := conn.Peers.Get(ctx, projectID, peerID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorPeersRead, peerID, err))
	}

	provider := "AWS"
	if peer.VNetName != "" {
		provider = "AZURE"
	} else if peer.NetworkName != "" {
		provider = "GCP"
	}

	if err := d.Set("provider_name", provider); err != nil {
		return diag.FromErr(fmt.Errorf("[WARN] Error setting provider_name for (%s): %s", d.Id(), err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":    projectID,
		"peer_id":       peer.ID,
		"provider_name": provider,
	}))

	accepterRegionName, err := ensureAccepterRegionName(ctx, peer, conn, projectID)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("peering_id", peerID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `peering_id` for Network Peering Connection (%s): %s", peerID, err))
	}
	if err := d.Set("container_id", peer.ContainerID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `container_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	return setCommonFields(d, peer, peerID, accepterRegionName)
}
