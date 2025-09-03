package networkpeering

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/networkcontainer"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

const (
	errorPeersCreate = "error creating MongoDB Network Peering Connection: %s"
	errorPeersRead   = "error reading MongoDB Network Peering Connection (%s): %s"
	errorPeersDelete = "error deleting MongoDB Network Peering Connection (%s): %s"
	errorPeersUpdate = "error updating MongoDB Network Peering Connection (%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportState,
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
				ForceNew: true,
			},
			"azure_subscription_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vnet_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
				ForceNew: true,
			},
			"network_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)

	peerRequest := &admin.BaseNetworkPeeringConnectionSettings{
		ContainerId:  conversion.GetEncodedID(d.Get("container_id").(string), "container_id"),
		ProviderName: conversion.StringPtr(providerName),
	}

	if providerName == "AWS" {
		region, err := conversion.ValRegion(d.Get("accepter_region_name"), "network_peering")
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

		peerRequest.SetAccepterRegionName(region)
		peerRequest.SetAwsAccountId(awsAccountID.(string))
		peerRequest.SetRouteTableCidrBlock(rtCIDR.(string))
		peerRequest.SetVpcId(vpcID.(string))
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

		peerRequest.SetGcpProjectId(gcpProjectID.(string))
		peerRequest.SetNetworkName(networkName.(string))
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

		peerRequest.SetAzureDirectoryId(azureDirectoryID.(string))
		peerRequest.SetAzureSubscriptionId(azureSubscriptionID.(string))
		peerRequest.SetResourceGroupName(resourceGroupName.(string))
		peerRequest.SetVnetName(vnetName.(string))
	}

	peer, _, err := conn.NetworkPeeringApi.CreateGroupPeer(ctx, projectID, peerRequest).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPeersCreate, err))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"INITIATING", "FINALIZING", "ADDING_PEER", "WAITING_FOR_USER"},
		Target:     []string{"FAILED", "AVAILABLE", "PENDING_ACCEPTANCE"},
		Refresh:    resourceRefreshFunc(ctx, peer.GetId(), projectID, peerRequest.GetContainerId(), conn.NetworkPeeringApi),
		Timeout:    1 * time.Hour,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPeersCreate, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":    projectID,
		"peer_id":       peer.GetId(),
		"provider_name": providerName,
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	peerID := ids["peer_id"]
	providerName := d.Get("provider_name").(string)

	peer, resp, err := conn.NetworkPeeringApi.GetGroupPeer(ctx, projectID, peerID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorPeersRead, peerID, err))
	}

	container, err := getContainer(ctx, conn.NetworkPeeringApi, projectID, peer.GetContainerId())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("atlas_cidr_block", container.GetAtlasCidrBlock()); err != nil {
		return diag.Errorf("error setting `atlas_cidr_block` for Network Peering Connection (%s): %s", peerID, err)
	}

	// If provider name is GCP we need to get the parameters to configure the the reciprocal connection between Mongo and Google
	if providerName == "GCP" {
		if err := d.Set("atlas_gcp_project_id", container.GetGcpProjectId()); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `atlas_gcp_project_id` for Network Peering Connection (%s): %s", peerID, err))
		}

		if err := d.Set("atlas_vpc_name", container.GetNetworkName()); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `atlas_vpc_name` for Network Peering Connection (%s): %s", peerID, err))
		}
	}

	if err := d.Set("peer_id", peerID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `peer_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	accepterRegionName, err := ensureAccepterRegionName(ctx, peer, conn.NetworkPeeringApi, projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	diags := setCommonFields(d, peer, peerID, accepterRegionName)

	failedStatusDiag := errorIfFailedStatusIsPresent(peer)

	return append(diags, failedStatusDiag...)
}

func errorIfFailedStatusIsPresent(peer *admin.BaseNetworkPeeringConnectionSettings) diag.Diagnostics {
	// for Azure/GCP "status" and "errorState" is returned, for AWS "statusName" and "errorStateName" :-/
	if peer.GetStatus() == "FAILED" || peer.GetStatusName() == "FAILED" {
		errorState := peer.GetErrorState()
		if peer.GetErrorStateName() != "" {
			errorState = peer.GetErrorStateName()
		}
		return diag.FromErr(fmt.Errorf("peer networking is in a failed state: %s", errorState))
	}
	return nil
}

func setCommonFields(d *schema.ResourceData, peer *admin.BaseNetworkPeeringConnectionSettings, peerID, accepterRegionName string) diag.Diagnostics {
	if err := d.Set("accepter_region_name", accepterRegionName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `accepter_region_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("aws_account_id", peer.GetAwsAccountId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `aws_account_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("route_table_cidr_block", peer.GetRouteTableCidrBlock()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `route_table_cidr_block` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("vpc_id", peer.GetVpcId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `vpc_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("connection_id", peer.GetConnectionId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `connection_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("error_state_name", peer.GetErrorStateName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `error_state_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("atlas_id", peer.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `atlas_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("status_name", peer.GetStatusName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `status_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("azure_directory_id", peer.GetAzureDirectoryId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `azure_directory_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("azure_subscription_id", peer.GetAzureSubscriptionId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `azure_subscription_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("resource_group_name", peer.GetResourceGroupName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `resource_group_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("vnet_name", peer.GetVnetName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `vnet_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("error_state", peer.GetErrorState()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `error_state` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("status", peer.GetStatus()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `status` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("gcp_project_id", peer.GetGcpProjectId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `gcp_project_id` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("network_name", peer.GetNetworkName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `network_name` for Network Peering Connection (%s): %s", peerID, err))
	}

	if err := d.Set("error_message", peer.GetErrorMessage()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `error_message` for Network Peering Connection (%s): %s", peerID, err))
	}

	return nil
}

func ensureAccepterRegionName(ctx context.Context, peer *admin.BaseNetworkPeeringConnectionSettings, conn admin.NetworkPeeringApi, projectID string) (string, error) {
	/* This fix the bug https://github.com/mongodb/terraform-provider-mongodbatlas/issues/53
	 * If the region name of the peering connection resource is the same as the container resource,
	 * the API returns it as a null value, so this causes the issue mentioned.
	 */
	var acepterRegionName string
	if peer.GetAccepterRegionName() != "" {
		acepterRegionName = peer.GetAccepterRegionName()
	} else {
		container, _, err := conn.GetGroupContainer(ctx, projectID, peer.GetContainerId()).Execute()
		if err != nil {
			return "", err
		}
		// network_peering resource must use region of format eu-west-2 while network_container uses EU_WEST_2.
		acepterRegionName = conversion.MongoDBRegionToAWSRegion(container.GetRegionName())
	}
	return acepterRegionName, nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	peerID := ids["peer_id"]

	peer := &admin.BaseNetworkPeeringConnectionSettings{
		ProviderName: conversion.StringPtr(ids["provider_name"]),
		ContainerId:  conversion.GetEncodedID(d.Get("container_id").(string), "container_id"),
	}

	if peer.GetProviderName() == "AWS" {
		region, _ := conversion.ValRegion(d.Get("accepter_region_name"), "network_peering")
		peer.SetAccepterRegionName(region)
		peer.SetAwsAccountId(d.Get("aws_account_id").(string))
		peer.SetRouteTableCidrBlock(d.Get("route_table_cidr_block").(string))
		peer.SetVpcId(d.Get("vpc_id").(string))
	}
	peerConn, resp, getErr := conn.NetworkPeeringApi.GetGroupPeer(ctx, projectID, peerID).Execute()
	if getErr != nil {
		if validate.StatusNotFound(resp) {
			return nil
		}
	}
	fmt.Print(peerConn.GetStatus())

	_, _, err := conn.NetworkPeeringApi.UpdateGroupPeer(ctx, projectID, peerID, peer).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPeersUpdate, peerID, err))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"INITIATING", "FINALIZING", "ADDING_PEER", "WAITING_FOR_USER"},
		Target:     []string{"FAILED", "AVAILABLE", "PENDING_ACCEPTANCE"},
		Refresh:    resourceRefreshFunc(ctx, peerID, projectID, "", conn.NetworkPeeringApi),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPeersCreate, err))
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	peerID := ids["peer_id"]

	_, _, err := conn.NetworkPeeringApi.DeleteGroupPeer(ctx, projectID, peerID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPeersDelete, peerID, err))
	}

	log.Println("[INFO] Waiting for MongoDB Network Peering Connection to be destroyed")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"AVAILABLE", "INITIATING", "PENDING_ACCEPTANCE", "FINALIZING", "ADDING_PEER", "WAITING_FOR_USER", "TERMINATING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceRefreshFunc(ctx, peerID, projectID, "", conn.NetworkPeeringApi),
		Timeout:    1 * time.Hour,
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second, // Wait 10 secs before starting
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPeersDelete, peerID, err))
	}

	return nil
}

func resourceImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "-", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a peer, use the format {project_id}-{peer_id}-{provider_name}")
	}

	projectID := parts[0]
	peerID := parts[1]
	providerName := parts[2]

	peer, _, err := conn.NetworkPeeringApi.GetGroupPeer(ctx, projectID, peerID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import peer %s in project %s, error: %s", peerID, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", peerID, err)
	}

	if err := d.Set("container_id", peer.GetContainerId()); err != nil {
		log.Printf("[WARN] Error setting container_id for (%s): %s", peerID, err)
	}

	if err := d.Set("provider_name", providerName); err != nil {
		log.Printf("[WARN] Error setting provider_name for (%s): %s", peerID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":    projectID,
		"peer_id":       peer.GetId(),
		"provider_name": providerName,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceRefreshFunc(ctx context.Context, peerID, projectID, containerID string, api admin.NetworkPeeringApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		c, resp, err := api.GetGroupPeer(ctx, projectID, peerID).Execute()
		if err != nil {
			if validate.StatusNotFound(resp) {
				return "", "DELETED", nil
			}

			log.Printf("error reading MongoDB Network Peering Connection %s: %s", peerID, err)

			return nil, "", err
		}

		status := c.GetStatus()

		if c.GetStatusName() != "" {
			status = c.GetStatusName()
		}

		log.Printf("[DEBUG] status for MongoDB Network Peering Connection: %s: %s", peerID, status)

		/* We need to get the provisioned status from Mongo container that contains the peering connection
		 * to validate if it has changed to true. This means that the reciprocal connection in Mongo side
		 * is right, and the Mongo parameters used on the Google side to configure the reciprocal connection
		 * are now available. */
		if status == "WAITING_FOR_USER" {
			container, _, err := api.GetGroupContainer(ctx, projectID, containerID).Execute()

			if err != nil {
				return nil, "", fmt.Errorf(networkcontainer.ErrorContainerRead, containerID, err)
			}

			if container.GetProvisioned() {
				return container, "PENDING_ACCEPTANCE", nil
			}
		}

		return c, status, nil
	}
}
