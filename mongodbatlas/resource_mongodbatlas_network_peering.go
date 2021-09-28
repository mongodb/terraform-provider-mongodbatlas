package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorPeersCreate = "error creating MongoDB Network Peering Connection: %s"
	errorPeersRead   = "error reading MongoDB Network Peering Connection (%s): %s"
	errorPeersDelete = "error deleting MongoDB Network Peering Connection (%s): %s"
	errorPeersUpdate = "error updating MongoDB Network Peering Connection (%s): %s"
)

func resourceMongoDBAtlasNetworkPeering() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasNetworkPeeringCreate,
		ReadContext:   resourceMongoDBAtlasNetworkPeeringRead,
		UpdateContext: resourceMongoDBAtlasNetworkPeeringUpdate,
		DeleteContext: resourceMongoDBAtlasNetworkPeeringDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasNetworkPeeringImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accepter_region_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aws_account_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"peer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "AZURE", "GCP"}, false),
			},
			"route_table_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"azure_directory_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"azure_subscription_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vnet_name": {
				Type:     schema.TypeString,
				Optional: true,
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
				Optional: true,
				Computed: true,
			},
			"network_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"atlas_gcp_project_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"atlas_vpc_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasNetworkPeeringCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)

	// Get the required ones
	peerRequest := &matlas.Peer{
		ContainerID:  getEncodedID(d.Get("container_id").(string), "container_id"),
		ProviderName: providerName,
	}

	if providerName == "AWS" {
		region, err := valRegion(d.Get("accepter_region_name"), "network_peering")
		if err != nil {
			return diag.FromErr(errors.New("`accepter_region_name` must be set when `provider_name` is `AWS`"))
		}

		awsAccountID, ok := d.GetOk("aws_account_id")
		if !ok {
			return diag.FromErr(errors.New("`aws_account_id` must be set when `provider_name` is `AWS`"))
		}

		rtCIDR, ok := d.GetOk("route_table_cidr_block")
		if !ok {
			return diag.FromErr(errors.New("`route_table_cidr_block` must be set when `provider_name` is `AWS`"))
		}

		vpcID, ok := d.GetOk("vpc_id")
		if !ok {
			return diag.FromErr(errors.New("`vpc_id` must be set when `provider_name` is `AWS`"))
		}

		peerRequest.AccepterRegionName = region
		peerRequest.AWSAccountID = awsAccountID.(string)
		peerRequest.RouteTableCIDRBlock = rtCIDR.(string)
		peerRequest.VpcID = vpcID.(string)
	}

	if providerName == "GCP" {
		gcpProjectID, ok := d.GetOk("gcp_project_id")
		if !ok {
			return diag.FromErr(errors.New("`gcp_project_id` must be set when `provider_name` is `GCP`"))
		}

		networkName, ok := d.GetOk("network_name")
		if !ok {
			return diag.FromErr(errors.New("`network_name` must be set when `provider_name` is `GCP`"))
		}

		peerRequest.GCPProjectID = gcpProjectID.(string)
		peerRequest.NetworkName = networkName.(string)
	}

	if providerName == "AZURE" {
		azureDirectoryID, ok := d.GetOk("azure_directory_id")
		if !ok {
			return diag.FromErr(errors.New("`azure_directory_id` must be set when `provider_name` is `AZURE`"))
		}

		azureSubscriptionID, ok := d.GetOk("azure_subscription_id")
		if !ok {
			return diag.FromErr(errors.New("`azure_subscription_id` must be set when `provider_name` is `AZURE`"))
		}

		resourceGroupName, ok := d.GetOk("resource_group_name")
		if !ok {
			return diag.FromErr(errors.New("`resource_group_name` must be set when `provider_name` is `AZURE`"))
		}

		vnetName, ok := d.GetOk("vnet_name")
		if !ok {
			return diag.FromErr(errors.New("`vnet_name` must be set when `provider_name` is `AZURE`"))
		}

		peerRequest.AzureDirectoryID = azureDirectoryID.(string)
		peerRequest.AzureSubscriptionID = azureSubscriptionID.(string)
		peerRequest.ResourceGroupName = resourceGroupName.(string)
		peerRequest.VNetName = vnetName.(string)
	}

	peer, _, err := conn.Peers.Create(ctx, projectID, peerRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPeersCreate, err))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INITIATING", "FINALIZING", "ADDING_PEER", "WAITING_FOR_USER"},
		Target:     []string{"AVAILABLE", "PENDING_ACCEPTANCE"},
		Refresh:    resourceNetworkPeeringRefreshFunc(ctx, peer.ID, projectID, peerRequest.ContainerID, conn),
		Timeout:    1 * time.Hour,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPeersCreate, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":    projectID,
		"peer_id":       peer.ID,
		"provider_name": providerName,
	}))

	return resourceMongoDBAtlasNetworkPeeringRead(ctx, d, meta)
}

func resourceMongoDBAtlasNetworkPeeringRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	peerID := ids["peer_id"]
	providerName := ids["provider_name"]

	peer, resp, err := conn.Peers.Get(ctx, projectID, peerID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorPeersRead, peerID, err))
	}

	/* This fix the bug https://github.com/mongodb/terraform-provider-mongodbatlas/issues/53
	 * If the region name of the peering connection resource is the same as the container resource,
	 * the API returns it as a null value, so this causes the issue mentioned.
	 */
	var acepterRegionName string
	if peer.AccepterRegionName != "" {
		acepterRegionName = peer.AccepterRegionName
	} else {
		container, _, err := conn.Containers.Get(ctx, projectID, peer.ContainerID)
		if err != nil {
			return diag.FromErr(err)
		}
		acepterRegionName = strings.ToLower(strings.ReplaceAll(container.RegionName, "_", "-"))
	}

	if err := d.Set("accepter_region_name", acepterRegionName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `accepter_region_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("aws_account_id", peer.AWSAccountID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `aws_account_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("route_table_cidr_block", peer.RouteTableCIDRBlock); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `route_table_cidr_block` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("vpc_id", peer.VpcID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `vpc_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("connection_id", peer.ConnectionID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `connection_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("error_state_name", peer.ErrorStateName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `error_state_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("atlas_id", peer.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `atlas_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("status_name", peer.StatusName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `status_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("azure_directory_id", peer.AzureDirectoryID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `azure_directory_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("azure_subscription_id", peer.AzureSubscriptionID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `azure_subscription_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("resource_group_name", peer.ResourceGroupName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `resource_group_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("vnet_name", peer.VNetName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `vnet_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("error_state", peer.ErrorState); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `error_state` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("status", peer.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `status` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("gcp_project_id", peer.GCPProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `gcp_project_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("network_name", peer.NetworkName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `network_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("error_message", peer.ErrorMessage); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `error_message` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("peer_id", peer.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `peer_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	// If provider name is GCP we need to get the parameters to configure the the reciprocal connection
	//  between Mongo and Google
	container := &matlas.Container{}

	if strings.EqualFold(providerName, "GCP") {
		container, _, err = conn.Containers.Get(ctx, projectID, peer.ContainerID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("atlas_gcp_project_id", container.GCPProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `atlas_gcp_project_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("atlas_vpc_name", container.NetworkName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `atlas_vpc_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	return nil
}

func resourceMongoDBAtlasNetworkPeeringUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	peerID := ids["peer_id"]

	// All the request to update the peer require the ProviderName and ContainerID attribute.
	peer := &matlas.Peer{
		ProviderName: ids["provider_name"],
		ContainerID:  getEncodedID(d.Get("container_id").(string), "container_id"),
	}

	// Depending of the Provider name the request will be set
	switch peer.ProviderName {
	case "GCP":
		peer.GCPProjectID = d.Get("gcp_project_id").(string)
		peer.NetworkName = d.Get("network_name").(string)
	case "AZURE":
		if d.HasChange("azure_directory_id") {
			peer.AzureDirectoryID = d.Get("azure_directory_id").(string)
		}

		if d.HasChange("azure_subscription_id") {
			peer.AzureSubscriptionID = d.Get("azure_subscription_id").(string)
		}

		if d.HasChange("resource_group_name") {
			peer.ResourceGroupName = d.Get("resource_group_name").(string)
		}

		if d.HasChange("vnet_name") {
			peer.VNetName = d.Get("vnet_name").(string)
		}
	default: // AWS by default
		region, _ := valRegion(d.Get("accepter_region_name"), "network_peering")
		peer.AccepterRegionName = region
		peer.AWSAccountID = d.Get("aws_account_id").(string)
		peer.RouteTableCIDRBlock = d.Get("route_table_cidr_block").(string)
		peer.VpcID = d.Get("vpc_id").(string)
	}

	// Has changes
	if !reflect.DeepEqual(peer, matlas.Peer{}) {
		_, _, err := conn.Peers.Update(ctx, projectID, peerID, peer)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorPeersUpdate, peerID, err))
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INITIATING", "FINALIZING", "ADDING_PEER", "WAITING_FOR_USER"},
		Target:     []string{"AVAILABLE", "PENDING_ACCEPTANCE"},
		Refresh:    resourceNetworkPeeringRefreshFunc(ctx, peerID, projectID, "", conn),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPeersCreate, err))
	}

	return resourceMongoDBAtlasNetworkPeeringRead(ctx, d, meta)
}

func resourceMongoDBAtlasNetworkPeeringDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	peerID := ids["peer_id"]

	_, err := conn.Peers.Delete(ctx, projectID, peerID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPeersDelete, peerID, err))
	}

	log.Println("[INFO] Waiting for MongoDB Network Peering Connection to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"AVAILABLE", "INITIATING", "PENDING_ACCEPTANCE", "FINALIZING", "ADDING_PEER", "WAITING_FOR_USER", "TERMINATING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkPeeringRefreshFunc(ctx, peerID, projectID, "", conn),
		Timeout:    1 * time.Hour,
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second, // Wait 10 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPeersDelete, peerID, err))
	}

	return nil
}

func resourceMongoDBAtlasNetworkPeeringImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "-", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a peer, use the format {project_id}-{peer_id}-{provider_name}")
	}

	projectID := parts[0]
	peerID := parts[1]
	providerName := parts[2]

	peer, _, err := conn.Peers.Get(ctx, projectID, peerID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import peer %s in project %s, error: %s", peerID, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", peerID, err)
	}

	if err := d.Set("container_id", peer.ContainerID); err != nil {
		log.Printf("[WARN] Error setting container_id for (%s): %s", peerID, err)
	}

	if err := d.Set("provider_name", providerName); err != nil {
		log.Printf("[WARN] Error setting provider_name for (%s): %s", peerID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":    projectID,
		"peer_id":       peer.ID,
		"provider_name": providerName,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceNetworkPeeringRefreshFunc(ctx context.Context, peerID, projectID, containerID string, client *matlas.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, resp, err := client.Peers.Get(ctx, projectID, peerID)
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return "", "DELETED", nil
			}

			log.Printf("error reading MongoDB Network Peering Connection %s: %s", peerID, err)

			return nil, "", err
		}

		status := c.Status

		if len(c.StatusName) > 0 {
			status = c.StatusName
		}

		log.Printf("[DEBUG] status for MongoDB Network Peering Connection: %s: %s", peerID, status)

		/* We need to get the provisioned status from Mongo container that contains the peering connection
		 * to validate if it has changed to true. This means that the reciprocal connection in Mongo side
		 * is right, and the Mongo parameters used on the Google side to configure the reciprocal connection
		 * are now available. */
		if status == "WAITING_FOR_USER" {
			container, _, err := client.Containers.Get(ctx, projectID, containerID)

			if err != nil {
				return nil, "", fmt.Errorf(errorContainerRead, containerID, err)
			}

			if *container.Provisioned {
				return container, "PENDING_ACCEPTANCE", nil
			}
		}

		return c, status, nil
	}
}
