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
	errorContainterCreate = "error creating MongoDB Network Peering Container: %s"
	errorContainerRead    = "error reading MongoDB Network Peering Container (%s): %s"
	errorContainerDelete  = "error deleting MongoDB Network Peering Container (%s): %s"
	errorContainerUpdate  = "error updating MongoDB Network Peering Container (%s): %s"
)

func resourceMongoDBAtlasNetworkContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasNetworkContainerCreate,
		Read:   resourceMongoDBAtlasNetworkContainerRead,
		Update: resourceMongoDBAtlasNetworkContainerUpdate,
		Delete: resourceMongoDBAtlasNetworkContainerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasNetworkContainerImportState,
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
				Default:      "AWS",
				ValidateFunc: validation.StringInSlice([]string{"AWS", "GCP", "AZURE"}, false),
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
		},
	}
}

func resourceMongoDBAtlasNetworkContainerCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)

	containerRequest := &matlas.Container{
		AtlasCIDRBlock: d.Get("atlas_cidr_block").(string),
		ProviderName:   providerName,
	}

	if providerName == "AWS" {
		region, err := valRegion(d.Get("region_name"))
		if err != nil {
			return fmt.Errorf("`region_name` must be set when `provider_name` is AWS")
		}
		containerRequest.RegionName = region
	}

	if providerName == "AZURE" {
		region, err := valRegion(d.Get("region"))
		if err != nil {
			return fmt.Errorf("`region` must be set when `provider_name` is AZURE")
		}
		containerRequest.Region = region
	}

	container, _, err := conn.Containers.Create(context.Background(), projectID, containerRequest)
	if err != nil {
		return fmt.Errorf(errorContainterCreate, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"container_id": container.ID,
	}))

	return resourceMongoDBAtlasNetworkContainerRead(d, meta)
}

func resourceMongoDBAtlasNetworkContainerRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	containerID := ids["container_id"]

	container, resp, err := conn.Containers.Get(context.Background(), projectID, containerID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {

			return nil
		}
		return fmt.Errorf(errorContainerRead, containerID, err)
	}

	if err = d.Set("region_name", container.RegionName); err != nil {
		return fmt.Errorf("error setting `region_name` for Network Container (%s): %s", containerID, err)
	}

	if err = d.Set("region", container.Region); err != nil {
		return fmt.Errorf("error setting `region` for Network Container (%s): %s", containerID, err)
	}

	if err = d.Set("azure_subscription_id", container.AzureSubscriptionID); err != nil {
		return fmt.Errorf("error setting `azure_subscription_id` for Network Container (%s): %s", containerID, err)
	}
	if err = d.Set("provisioned", container.Provisioned); err != nil {
		return fmt.Errorf("error setting `provisioned` for Network Container (%s): %s", containerID, err)
	}
	if err = d.Set("gcp_project_id", container.GCPProjectID); err != nil {
		return fmt.Errorf("error setting `gcp_project_id` for Network Container (%s): %s", containerID, err)
	}
	if err = d.Set("network_name", container.NetworkName); err != nil {
		return fmt.Errorf("error setting `network_name` for Network Container (%s): %s", containerID, err)
	}
	if err = d.Set("gcp_project_id", container.GCPProjectID); err != nil {
		return fmt.Errorf("error setting `gcp_project_id` for Network Container (%s): %s", containerID, err)
	}
	if err = d.Set("vpc_id", container.VPCID); err != nil {
		return fmt.Errorf("error setting `vpc_id` for Network Container (%s): %s", containerID, err)
	}
	if err = d.Set("vnet_name", container.VNetName); err != nil {
		return fmt.Errorf("error setting `vnet_name` for Network Container (%s): %s", containerID, err)
	}
	if err = d.Set("container_id", container.ID); err != nil {
		return fmt.Errorf("error setting `container_id` for Network Container (%s): %s", containerID, err)
	}
	return nil
}

func resourceMongoDBAtlasNetworkContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	containerID := ids["container_id"]

	container := new(matlas.Container)

	if d.HasChange("atlas_cidr_block") {
		container.AtlasCIDRBlock = d.Get("atlas_cidr_block").(string)
		container.ProviderName = d.Get("provider_name").(string)
	}

	if d.HasChange("provider_name") {
		container.ProviderName = d.Get("provider_name").(string)
	}

	if d.HasChange("region_name") {
		region, _ := valRegion(d.Get("region_name"))
		container.RegionName = region
	}

	if d.HasChange("region") {
		region, _ := valRegion(d.Get("region"))
		container.Region = region
	}

	// Has changes
	if !reflect.DeepEqual(container, matlas.Container{}) {
		_, _, err := conn.Containers.Update(context.Background(), projectID, containerID, container)
		if err != nil {
			return fmt.Errorf(errorContainerUpdate, containerID, err)
		}
	}

	return resourceMongoDBAtlasNetworkContainerRead(d, meta)
}

func resourceMongoDBAtlasNetworkContainerDelete(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"provisioned_container"},
		Target:     []string{"deleted"},
		Refresh:    resourceNetworkContainerRefreshFunc(d, conn),
		Timeout:    1 * time.Hour,
		MinTimeout: 10 * time.Second,
		Delay:      2 * time.Minute,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorContainerDelete, decodeStateID(d.Id())["container_id"], err)
	}

	return nil
}

func resourceMongoDBAtlasNetworkContainerImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a container, use the format {project_id}-{container_id}")
	}

	projectID := parts[0]
	containerID := parts[1]

	u, _, err := conn.Containers.Get(context.Background(), projectID, containerID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import container %s in project %s, error: %s", containerID, projectID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"container_id": u.ID,
	}))

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", containerID, err)
	}
	if err := d.Set("provider_name", u.ProviderName); err != nil {
		log.Printf("[WARN] Error setting provider_name for (%s): %s", containerID, err)
	}
	if err := d.Set("atlas_cidr_block", u.AtlasCIDRBlock); err != nil {
		log.Printf("[WARN] Error setting atlas_cidr_block for (%s): %s", containerID, err)
	}
	if err = d.Set("container_id", u.ID); err != nil {
		log.Printf("[WARN] Error setting container_id (%s): %s", containerID, err)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceNetworkContainerRefreshFunc(d *schema.ResourceData, client *matlas.Client) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ids := decodeStateID(d.Id())
		projectID := ids["project_id"]
		containerID := ids["container_id"]

		var err error
		container, res, err := client.Containers.Get(context.Background(), projectID, containerID)
		if err != nil {
			if res.StatusCode == 404 {
				return 42, "deleted", nil
			}
			return nil, "", err
		}
		if *container.Provisioned && err == nil {
			return nil, "provisioned_container", nil
		}

		_, err = client.Containers.Delete(context.Background(), projectID, containerID)
		if err != nil {
			return nil, "provisioned_container", nil
		}

		return 42, "deleted", nil
	}
}

func resourceListNetworkContainerRefreshFunc(d *schema.ResourceData, client *matlas.Client, originaList []matlas.Container) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ids := decodeStateID(d.Id())
		projectID := ids["project_id"]
		providerName := ids["provider_name"]
		options := &matlas.ContainersListOptions{
			ProviderName: providerName,
		}

		var err error
		containers, res, err := client.Containers.List(context.Background(), projectID, options)
		if err != nil {
			if res.StatusCode == 404 {
				return 42, "deleted", nil
			}
			return nil, "", err
		}
		if len(containers) == len(originaList) {
			return nil, "same", nil
		}

		return containers, "different", nil
	}
}
