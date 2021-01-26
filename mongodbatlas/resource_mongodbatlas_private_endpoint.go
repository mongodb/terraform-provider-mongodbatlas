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
	errorPrivateEndpointsCreate  = "error creating MongoDB Private Endpoints  Connection: %s"
	errorPrivateEndpointsRead    = "error reading MongoDB Private Endpoints  Connection(%s): %s"
	errorPrivateEndpointsDelete  = "error deleting MongoDB Private Endpoints  Connection(%s): %s"
	errorPrivateEndpointsSetting = "error setting `%s` for MongoDB Private Endpoints  Connection(%s): %s"
)

func resourceMongoDBAtlasPrivateEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasPrivateEndpointCreate,
		Read:   resourceMongoDBAtlasPrivateEndpointRead,
		Delete: resourceMongoDBAtlasPrivateEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasPrivateEndpointImportState,
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
				ValidateFunc: validation.StringInSlice([]string{"AWS"}, false),
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		DeprecationMessage: "this resource is deprecated, please transition as soon as possible to mongodbatlas_privatelink_endpoint",
	}
}

func resourceMongoDBAtlasPrivateEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	request := &matlas.PrivateEndpointConnectionDeprecated{
		ProviderName: d.Get("provider_name").(string),
		Region:       d.Get("region").(string),
	}

	privateEndpointConn, _, err := conn.PrivateEndpointsDeprecated.Create(context.Background(), projectID, request)
	if err != nil {
		return fmt.Errorf(errorPrivateEndpointsCreate, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INITIATING", "DELETING"},
		Target:     []string{"WAITING_FOR_USER", "FAILED", "DELETED"},
		Refresh:    resourcePrivateEndpointRefreshFunc(conn, projectID, privateEndpointConn.ID),
		Timeout:    1 * time.Hour,
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorPrivateEndpointsCreate, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"private_link_id": privateEndpointConn.ID,
		"project_id":      projectID,
	}))

	return resourceMongoDBAtlasPrivateEndpointRead(d, meta)
}

func resourceMongoDBAtlasPrivateEndpointRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	providerName := ids["provider_name"]
	region := ids["region"]

	privateEndpoint, _, err := conn.PrivateEndpointsDeprecated.Get(context.Background(), projectID, privateLinkID)
	if err != nil {
		return fmt.Errorf(errorPrivateEndpointsRead, privateLinkID, err)
	}

	if err := d.Set("private_link_id", privateEndpoint.ID); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "private_link_id", privateLinkID, err)
	}

	if err := d.Set("endpoint_service_name", privateEndpoint.EndpointServiceName); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "endpoint_service_name", privateLinkID, err)
	}

	if err := d.Set("error_message", privateEndpoint.ErrorMessage); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "error_message", privateLinkID, err)
	}

	if err := d.Set("interface_endpoints", privateEndpoint.InterfaceEndpoints); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "interface_endpoints", privateLinkID, err)
	}

	if err := d.Set("status", privateEndpoint.Status); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "status", privateLinkID, err)
	}

	if err := d.Set("provider_name", providerName); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "provider_name", privateLinkID, err)
	}

	if err := d.Set("region", region); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "provider_name", privateLinkID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	privateLinkID := ids["private_link_id"]
	projectID := ids["project_id"]

	resp, err := conn.PrivateEndpointsDeprecated.Delete(context.Background(), projectID, privateLinkID)
	if err != nil {
		if resp.Response.StatusCode == 404 {
			return nil
		}

		return fmt.Errorf(errorPrivateEndpointsDelete, privateLinkID, err)
	}

	log.Println("[INFO] Waiting for MongoDB Private Endpoints  Connection to be destroyed")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED", "FAILED"},
		Refresh:    resourcePrivateEndpointRefreshFunc(conn, projectID, privateLinkID),
		Timeout:    10 * time.Minute,
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorPrivateEndpointsDelete, privateLinkID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.Split(d.Id(), "-")
	if len(parts) != 6 {
		return nil, errors.New("import format error: to import a MongoDB Private Endpoint, use the format {project_id}-{private_link_id}-{provider_name}-{region} ")
	}

	projectID := parts[0]
	privateLinkID := parts[1]
	providerName := parts[2]
	region := fmt.Sprintf("%s-%s-%s", parts[3], parts[4], parts[5])

	privateEndpoint, _, err := conn.PrivateEndpointsDeprecated.Get(context.Background(), projectID, privateLinkID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import peer %s in project %s, error: %s", privateLinkID, projectID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf(errorPrivateEndpointsSetting, "project_id", privateLinkID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"private_link_id": privateEndpoint.ID,
		"project_id":      projectID,
		"provider_name":   providerName,
		"region":          region,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourcePrivateEndpointRefreshFunc(client *matlas.Client, projectID, privateLinkID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, resp, err := client.PrivateEndpointsDeprecated.Get(context.Background(), projectID, privateLinkID)
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
