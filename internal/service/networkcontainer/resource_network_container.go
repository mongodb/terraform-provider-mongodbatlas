package networkcontainer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312008/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorContainterCreate = "error creating MongoDB Network Peering Container: %s"
	ErrorContainerRead    = "error reading MongoDB Network Peering Container (%s): %s"
	errorContainerDelete  = "error deleting MongoDB Network Peering Container (%s): %s"
	errorContainerUpdate  = "error updating MongoDB Network Peering Container (%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImport,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"atlas_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      constant.AWS,
				ValidateFunc: validation.StringInSlice([]string{constant.AWS, constant.GCP, constant.AZURE}, false),
			},
			"region_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
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
			"container_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)

	atlasCidrBlock := d.Get("atlas_cidr_block").(string)
	containerRequest := &admin.CloudProviderContainer{
		AtlasCidrBlock: &atlasCidrBlock,
		ProviderName:   &providerName,
	}

	if providerName == constant.AWS {
		region, err := conversion.ValRegion(d.Get("region_name"))
		if err != nil {
			return diag.FromErr(fmt.Errorf("`region_name` must be set when `provider_name` is AWS"))
		}
		containerRequest.RegionName = &region
	}

	if providerName == constant.AZURE {
		region, err := conversion.ValRegion(d.Get("region"))
		if err != nil {
			return diag.FromErr(fmt.Errorf("`region` must be set when `provider_name` is AZURE"))
		}
		containerRequest.Region = &region
	}

	if providerName == constant.GCP {
		regions, ok := d.GetOk("regions")
		if ok {
			regionsSlice := cast.ToStringSlice(regions)
			containerRequest.Regions = &regionsSlice
		}
	}

	container, _, err := connV2.NetworkPeeringApi.CreateGroupContainer(ctx, projectID, containerRequest).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorContainterCreate, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID,
		"container_id": container.GetId(),
	}))

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	containerID := ids["container_id"]

	container, resp, err := connV2.NetworkPeeringApi.GetGroupContainer(ctx, projectID, containerID).Execute()
	if err != nil {
		if validate.StatusNotFound(resp) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf(ErrorContainerRead, containerID, err))
	}

	if err = d.Set("region_name", container.GetRegionName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `region_name` for Network Container (%s): %s", containerID, err))
	}

	if err = d.Set("region", container.GetRegion()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `region` for Network Container (%s): %s", containerID, err))
	}

	if err = d.Set("azure_subscription_id", container.GetAzureSubscriptionId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `azure_subscription_id` for Network Container (%s): %s", containerID, err))
	}

	if err = d.Set("provisioned", container.GetProvisioned()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `provisioned` for Network Container (%s): %s", containerID, err))
	}

	if err = d.Set("gcp_project_id", container.GetGcpProjectId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `gcp_project_id` for Network Container (%s): %s", containerID, err))
	}

	if err = d.Set("network_name", container.GetNetworkName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `network_name` for Network Container (%s): %s", containerID, err))
	}

	if err = d.Set("vpc_id", container.GetVpcId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `vpc_id` for Network Container (%s): %s", containerID, err))
	}

	if err = d.Set("vnet_name", container.GetVnetName()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `vnet_name` for Network Container (%s): %s", containerID, err))
	}

	if err = d.Set("container_id", container.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `container_id` for Network Container (%s): %s", containerID, err))
	}

	if err = d.Set("regions", container.GetRegions()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `regions` for Network Container (%s): %s", containerID, err))
	}

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if !d.HasChange("provider_name") && !d.HasChange("atlas_cidr_block") && !d.HasChange("region_name") && !d.HasChange("region") && !d.HasChange("regions") {
		return resourceRead(ctx, d, meta)
	}

	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	containerID := ids["container_id"]

	providerName := d.Get("provider_name").(string)
	cidr := d.Get("atlas_cidr_block").(string)

	params := &admin.CloudProviderContainer{
		ProviderName:   conversion.StringPtr(providerName),
		AtlasCidrBlock: conversion.StringPtr(cidr),
	}

	switch providerName {
	case constant.AWS:
		regionName, err := conversion.ValRegion(d.Get("region_name"))
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorContainerUpdate, containerID, err))
		}
		params.SetRegionName(regionName)
	case constant.AZURE:
		region, err := conversion.ValRegion(d.Get("region"))
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorContainerUpdate, containerID, err))
		}
		params.SetRegion(region)
	case constant.GCP:
		if regionList, ok := d.GetOk("regions"); ok {
			if regions := cast.ToStringSlice(regionList); regions != nil {
				params.SetRegions(regions)
			}
		}
	}
	_, _, err := connV2.NetworkPeeringApi.UpdateGroupContainer(ctx, projectID, containerID, params).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorContainerUpdate, containerID, err))
	}
	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"provisioned_container"},
		Target:     []string{"deleted"},
		Refresh:    resourceRefreshFunc(ctx, d, connV2),
		Timeout:    1 * time.Hour,
		MinTimeout: 10 * time.Second,
		Delay:      2 * time.Minute,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorContainerDelete, conversion.DecodeStateID(d.Id())["container_id"], err))
	}

	return nil
}

func resourceImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a container, use the format {project_id}-{container_id}")
	}

	projectID := parts[0]
	containerID := parts[1]

	networkContainer, _, err := connV2.NetworkPeeringApi.GetGroupContainer(ctx, projectID, containerID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import container %s in project %s, error: %s", containerID, projectID, err)
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID,
		"container_id": networkContainer.GetId(),
	}))

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", containerID, err)
	}

	if err := d.Set("provider_name", networkContainer.GetProviderName()); err != nil {
		log.Printf("[WARN] Error setting provider_name for (%s): %s", containerID, err)
	}

	if err := d.Set("atlas_cidr_block", networkContainer.GetAtlasCidrBlock()); err != nil {
		log.Printf("[WARN] Error setting atlas_cidr_block for (%s): %s", containerID, err)
	}

	if err = d.Set("container_id", networkContainer.GetId()); err != nil {
		log.Printf("[WARN] Error setting container_id (%s): %s", containerID, err)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceRefreshFunc(ctx context.Context, d *schema.ResourceData, client *admin.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		ids := conversion.DecodeStateID(d.Id())
		projectID := ids["project_id"]
		containerID := ids["container_id"]

		var err error
		container, res, err := client.NetworkPeeringApi.GetGroupContainer(ctx, projectID, containerID).Execute()
		if err != nil {
			if validate.StatusNotFound(res) {
				return "", "deleted", nil
			}

			return nil, "", err
		}

		if *container.Provisioned {
			return nil, "provisioned_container", nil
		}

		// Atlas Delete is called inside refresh to retry when error: HTTP 409 Conflict (Error code: "CONTAINERS_IN_USE").
		_, err = client.NetworkPeeringApi.DeleteGroupContainer(ctx, projectID, containerID).Execute()
		if err != nil {
			return nil, "provisioned_container", nil
		}

		return "", "deleted", nil
	}
}
