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

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	errorPeersCreate = "error creating MongoDB Network Peering Connection: %s"
	errorPeersRead   = "error reading MongoDB Network Peering Connection (%s): %s"
	errorPeersDelete = "error deleting MongoDB Network Peering Connection (%s): %s"
	errorPeersUpdate = "error updating MongoDB Network Peering Connection (%s): %s"
)

func resourceMongoDBAtlasNetworkPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasNetworkPeeringCreate,
		Read:   resourceMongoDBAtlasNetworkPeeringRead,
		Update: resourceMongoDBAtlasNetworkPeeringUpdate,
		Delete: resourceMongoDBAtlasNetworkPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasNetworkPeeringImportState,
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
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "AZURE", "GCP"}, false),
				Default:      "AWS",
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
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasNetworkPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)

	//Get the required ones
	peerRequest := &matlas.Peer{
		ContainerID:  d.Get("container_id").(string),
		ProviderName: providerName,
	}

	if providerName == "AWS" {
		accepter, ok := d.GetOk("accepter_region_name")
		if !ok {
			return errors.New("`accepter_region_name` must be set when `provider_name` is `AWS`")
		}

		awsAccountID, ok := d.GetOk("aws_account_id")
		if !ok {
			return errors.New("`aws_account_id` must be set when `provider_name` is `AWS`")
		}

		rtCIDR, ok := d.GetOk("route_table_cidr_block")
		if !ok {
			return errors.New("`route_table_cidr_block` must be set when `provider_name` is `AWS`")
		}

		vpcID, ok := d.GetOk("vpc_id")
		if !ok {
			return errors.New("`vpc_id` must be set when `provider_name` is `AWS`")
		}

		peerRequest.AccepterRegionName = accepter.(string)
		peerRequest.AWSAccountId = awsAccountID.(string)
		peerRequest.RouteTableCIDRBlock = rtCIDR.(string)
		peerRequest.VpcID = vpcID.(string)

	}

	if providerName == "GCP" {
		gcpProjectID, ok := d.GetOk("gcp_project_id")
		if !ok {
			return errors.New("`gcp_project_id` must be set when `provider_name` is `GCP`")
		}

		networkName, ok := d.GetOk("network_name")
		if !ok {
			return errors.New("`network_name` must be set when `provider_name` is `GCP`")
		}

		peerRequest.GCPProjectID = gcpProjectID.(string)
		peerRequest.NetworkName = networkName.(string)
	}

	if providerName == "AZURE" {
		atlasCidrBlock, ok := d.GetOk("atlas_cidr_block")
		if !ok {
			return errors.New("`atlas_cidr_block` must be set when `provider_name` is `AZURE`")
		}

		azureDirectoryID, ok := d.GetOk("azure_directory_id")
		if !ok {
			return errors.New("`azure_directory_id` must be set when `provider_name` is `AZURE`")
		}

		azureSubscriptionID, ok := d.GetOk("azure_subscription_id")
		if !ok {
			return errors.New("`azure_subscription_id` must be set when `provider_name` is `AZURE`")
		}

		resourceGroupName, ok := d.GetOk("resource_group_name")
		if !ok {
			return errors.New("`resource_group_name` must be set when `provider_name` is `AZURE`")
		}

		vnetName, ok := d.GetOk("vnet_name")
		if !ok {
			return errors.New("`vnet_name` must be set when `provider_name` is `AZURE`")
		}

		peerRequest.AtlasCIDRBlock = atlasCidrBlock.(string)
		peerRequest.AzureDirectoryID = azureDirectoryID.(string)
		peerRequest.AzureSubscriptionId = azureSubscriptionID.(string)
		peerRequest.ResourceGroupName = resourceGroupName.(string)
		peerRequest.VNetName = vnetName.(string)

	}

	peer, _, err := conn.Peers.Create(context.Background(), projectID, peerRequest)
	if err != nil {
		return fmt.Errorf(errorPeersCreate, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"peer_id":    peer.ID,
	}))

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INITIATING", "FINALIZING", "ADDING_PEER", "WAITING_FOR_USER"},
		Target:     []string{"AVAILABLE", "PENDING_ACCEPTANCE"},
		Refresh:    resourceNetworkPeeringRefreshFunc(peer.ID, projectID, conn),
		Timeout:    1 * time.Hour,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorPeersCreate, err)
	}

	return resourceMongoDBAtlasNetworkPeeringRead(d, meta)
}

func resourceMongoDBAtlasNetworkPeeringRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	peerID := ids["peer_id"]

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
	if err := d.Set("peer_id", peer.ID); err != nil {
		return fmt.Errorf("error setting `peer_id` for Network Peering Connection (%s): %s", peerID, err)
	}
	return nil
}

func resourceMongoDBAtlasNetworkPeeringUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	peerID := ids["peer_id"]

	peer := new(matlas.Peer)

	if d.HasChange("accepter_region_name") {
		peer.AccepterRegionName = d.Get("accepter_region_name").(string)
	}

	if d.HasChange("aws_account_id") {
		peer.AWSAccountId = d.Get("aws_account_id").(string)
	}

	if d.HasChange("provider_name") {
		peer.ProviderName = d.Get("provider_name").(string)
	}

	if d.HasChange("route_table_cidr_block") {
		peer.RouteTableCIDRBlock = d.Get("route_table_cidr_block").(string)
	}

	if d.HasChange("vpc_id") {
		peer.VpcID = d.Get("vpc_id").(string)
	}

	if d.HasChange("atlas_cidr_block") {
		peer.AtlasCIDRBlock = d.Get("atlas_cidr_block").(string)
	}

	if d.HasChange("azure_directory_id") {
		peer.AzureDirectoryID = d.Get("azure_directory_id").(string)
	}

	if d.HasChange("azure_subscription_id") {
		peer.AzureSubscriptionId = d.Get("azure_subscription_id").(string)
	}

	if d.HasChange("resource_group_name") {
		peer.ResourceGroupName = d.Get("resource_group_name").(string)
	}

	if d.HasChange("vnet_name") {
		peer.VNetName = d.Get("vnet_name").(string)
	}

	if d.HasChange("gcp_project_id") {
		peer.GCPProjectID = d.Get("gcp_project_id").(string)
	}

	if d.HasChange("network_name") {
		peer.NetworkName = d.Get("network_name").(string)
	}

	// Has changes
	if !reflect.DeepEqual(peer, matlas.Peer{}) {
		_, _, err := conn.Peers.Update(context.Background(), projectID, peerID, peer)
		if err != nil {
			return fmt.Errorf(errorPeersUpdate, peerID, err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INITIATING", "FINALIZING", "ADDING_PEER", "WAITING_FOR_USER"},
		Target:     []string{"AVAILABLE", "PENDING_ACCEPTANCE"},
		Refresh:    resourceNetworkPeeringRefreshFunc(peerID, projectID, conn),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorPeersCreate, err)
	}

	return resourceMongoDBAtlasNetworkPeeringRead(d, meta)
}

func resourceMongoDBAtlasNetworkPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	peerID := ids["peer_id"]

	_, err := conn.Peers.Delete(context.Background(), projectID, peerID)

	if err != nil {
		return fmt.Errorf(errorPeersDelete, peerID, err)
	}

	log.Println("[INFO] Waiting for MongoDB Network Peering Connection to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"AVAILABLE", "INITIATING", "PENDING_ACCEPTANCE", "FINALIZING", "ADDING_PEER", "WAITING_FOR_USER", "TERMINATING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkPeeringRefreshFunc(peerID, projectID, conn),
		Timeout:    1 * time.Hour,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorPeersDelete, peerID, err)
	}
	return nil
}

func resourceMongoDBAtlasNetworkPeeringImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a peer, use the format {project_id}-{peer_id}")
	}

	projectID := parts[0]
	peerID := parts[1]

	peer, _, err := conn.Peers.Get(context.Background(), projectID, peerID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import peer %s in project %s, error: %s", peerID, projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"peer_id":    peer.ID,
	}))

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", peerID, err)
	}

	if err := d.Set("container_id", peer.ContainerID); err != nil {
		log.Printf("[WARN] Error setting container_id for (%s): %s", peerID, err)
	}

	//Check wich provider is using.
	provider := "AWS"
	if peer.VNetName != "" {
		provider = "AZURE"
	} else if peer.NetworkName != "" {
		provider = "GCP"
	}

	if err := d.Set("provider_name", provider); err != nil {
		log.Printf("[WARN] Error setting provider_name for (%s): %s", peerID, err)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceNetworkPeeringRefreshFunc(peerID, projectID string, client *matlas.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, resp, err := client.Peers.Get(context.Background(), projectID, peerID)
		if err != nil {
			if resp.StatusCode == 404 {
				return 42, "DELETED", nil
			}
			log.Printf("error reading MongoDB Network Peering Connection %s: %s", peerID, err)
			return nil, "", err
		}

		status := c.Status

		if len(c.StatusName) > 0 {
			status = c.StatusName
		}

		log.Printf("[DEBUG] status for MongoDB Network Peering Connection: %s: %s", peerID, status)

		return c, status, nil
	}
}
