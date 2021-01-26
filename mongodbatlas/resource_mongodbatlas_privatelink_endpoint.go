package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorPrivateLinkEndpointsCreate  = "error creating MongoDB Private Endpoints Connection: %s"
	errorPrivateLinkEndpointsRead    = "error reading MongoDB Private Endpoints Connection(%s): %s"
	errorPrivateLinkEndpointsDelete  = "error deleting MongoDB Private Endpoints Connection(%s): %s"
	errorPrivateLinkEndpointsSetting = "error setting `%s` for MongoDB Private Endpoints Connection(%s): %s"
)

func resourceMongoDBAtlasPrivateLinkEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasPrivateLinkEndpointCreate,
		Read:   resourceMongoDBAtlasPrivateLinkEndpointRead,
		Delete: resourceMongoDBAtlasPrivateLinkEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasPrivateLinkEndpointImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "AZURE"}, false),
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_link_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_link_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_link_service_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasPrivateLinkEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)
	region := d.Get("region").(string)

	request := &matlas.PrivateEndpointConnection{
		ProviderName: providerName,
		Region:       region,
	}

	privateEndpointConn, _, err := conn.PrivateEndpoints.Create(context.Background(), projectID, request)
	if err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsCreate, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INITIATING", "DELETING"},
		Target:     []string{"WAITING_FOR_USER", "FAILED", "DELETED", "AVAILABLE"},
		Refresh:    resourcePrivateLinkEndpointRefreshFunc(conn, projectID, providerName, privateEndpointConn.ID),
		Timeout:    1 * time.Hour,
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsCreate, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"private_link_id": privateEndpointConn.ID,
		"project_id":      projectID,
		"provider_name":   providerName,
		"region":          region,
	}))

	return resourceMongoDBAtlasPrivateLinkEndpointRead(d, meta)
}

func resourceMongoDBAtlasPrivateLinkEndpointRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	providerName := ids["provider_name"]
	region := ids["region"]

	privateEndpoint, _, err := conn.PrivateEndpoints.Get(context.Background(), projectID, providerName, privateLinkID)
	if err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsRead, privateLinkID, err)
	}

	if err := d.Set("private_link_id", privateEndpoint.ID); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "private_link_id", privateLinkID, err)
	}

	if err := d.Set("endpoint_service_name", privateEndpoint.EndpointServiceName); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "endpoint_service_name", privateLinkID, err)
	}

	if err := d.Set("error_message", privateEndpoint.ErrorMessage); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "error_message", privateLinkID, err)
	}

	if err := d.Set("interface_endpoints", privateEndpoint.InterfaceEndpoints); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err)
	}

	if err := d.Set("private_endpoints", privateEndpoint.PrivateEndpoints); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err)
	}

	if err := d.Set("private_link_service_name", privateEndpoint.PrivateLinkServiceName); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err)
	}

	if err := d.Set("private_link_service_resource_id", privateEndpoint.PrivateLinkServiceResourceID); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err)
	}

	if err := d.Set("status", privateEndpoint.Status); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "status", privateLinkID, err)
	}

	if err := d.Set("provider_name", providerName); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "provider_name", privateLinkID, err)
	}

	if err := d.Set("region", region); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "region", privateLinkID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateLinkEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	privateLinkID := ids["private_link_id"]
	projectID := ids["project_id"]
	providerName := ids["provider_name"]

	resp, err := conn.PrivateEndpoints.Delete(context.Background(), projectID, providerName, privateLinkID)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			return nil
		}

		return fmt.Errorf(errorPrivateLinkEndpointsDelete, privateLinkID, err)
	}

	log.Println("[INFO] Waiting for MongoDB Private Endpoints Connection to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED", "FAILED"},
		Refresh:    resourcePrivateLinkEndpointRefreshFunc(conn, projectID, providerName, privateLinkID),
		Timeout:    10 * time.Minute,
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsDelete, privateLinkID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateLinkEndpointImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.Split(d.Id(), "-")
	if len(parts) != 6 && len(parts) != 4 {
		return nil, errors.New("import format error: to import a MongoDB Private Endpoint, use the format {project_id}-{private_link_id}-{provider_name}-{region}")
	}

	projectID := parts[0]
	privateLinkID := parts[1]
	providerName := parts[2]
	region := parts[3] // If region it's azure or Atlas format like US_EAST_1
	if len(parts) == 6 {
		region = fmt.Sprintf("%s-%s-%s", parts[3], parts[4], parts[5])
	}

	privateEndpoint, _, err := conn.PrivateEndpoints.Get(context.Background(), projectID, providerName, privateLinkID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import peer %s in project %s with cloud provider name %s, error: %s", privateLinkID, projectID, providerName, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf(errorPrivateLinkEndpointsSetting, "project_id", privateLinkID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"private_link_id": privateEndpoint.ID,
		"project_id":      projectID,
		"provider_name":   providerName,
		"region":          region,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourcePrivateLinkEndpointRefreshFunc(client *matlas.Client, projectID, providerName, privateLinkID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, resp, err := client.PrivateEndpoints.Get(context.Background(), projectID, providerName, privateLinkID)
		if err != nil {
			if resp.Response.StatusCode == 404 {
				return "", "DELETED", nil
			}

			return nil, "REJECTED", err
		}

		if p.Status != "WAITING_FOR_USER" {
			return "", p.Status, nil
		}

		return p, p.Status, nil
	}
}
