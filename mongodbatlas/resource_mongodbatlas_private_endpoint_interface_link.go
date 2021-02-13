package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorInterfaceEndpointAdd     = "error adding MongoDB Interface Endpoint Connection(%s) to a Private Endpoint (%s): %s"
	errorInterfaceEndpointRead    = "error reading MongoDB Interface Endpoint Connection(%s): %s"
	errorInterfaceEndpointDelete  = "error deleting MongoDB Interface Endpoints Connection(%s): %s"
	errorInterfaceEndpointSetting = "error setting `%s` for MongoDB Interface Endpoints Connection(%s): %s"
)

func resourceMongoDBAtlasPrivateEndpointInterfaceLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasPrivateEndpointInterfaceLinkCreate,
		Read:   resourceMongoDBAtlasPrivateEndpointInterfaceLinkRead,
		Delete: resourceMongoDBAtlasPrivateEndpointInterfaceLinkDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasPrivateEndpointInterfaceLinkImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_link_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interface_endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delete_requested": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		DeprecationMessage: "this resource is deprecated, please transition as soon as possible to mongodbatlas_privatelink_endpoint_service",
	}
}

func resourceMongoDBAtlasPrivateEndpointInterfaceLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	privateLinkID := d.Get("private_link_id").(string)
	interfaceEndpointID := d.Get("interface_endpoint_id").(string)

	interfaceEndpointConn, _, err := conn.PrivateEndpointsDeprecated.AddOneInterfaceEndpoint(context.Background(), projectID, privateLinkID, interfaceEndpointID)
	if err != nil {
		return fmt.Errorf(errorInterfaceEndpointAdd, interfaceEndpointID, privateLinkID, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NONE", "PENDING_ACCEPTANCE", "PENDING", "DELETING"},
		Target:     []string{"AVAILABLE", "REJECTED", "DELETED"},
		Refresh:    resourceInterfaceEndpointRefreshFunc(conn, projectID, privateLinkID, interfaceEndpointConn.ID),
		Timeout:    1 * time.Hour,
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorInterfaceEndpointAdd, interfaceEndpointConn.ID, privateLinkID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":            projectID,
		"private_link_id":       privateLinkID,
		"interface_endpoint_id": interfaceEndpointConn.ID,
	}))

	return resourceMongoDBAtlasPrivateEndpointInterfaceLinkRead(d, meta)
}

func resourceMongoDBAtlasPrivateEndpointInterfaceLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	interfaceEndpointID := ids["interface_endpoint_id"]

	interfaceEndpoint, _, err := conn.PrivateEndpointsDeprecated.GetOneInterfaceEndpoint(context.Background(), projectID, privateLinkID, interfaceEndpointID)
	if err != nil {
		return fmt.Errorf(errorInterfaceEndpointRead, interfaceEndpointID, err)
	}

	if err := d.Set("delete_requested", cast.ToBool(interfaceEndpoint.DeleteRequested)); err != nil {
		return fmt.Errorf(errorInterfaceEndpointSetting, "delete_requested", interfaceEndpointID, err)
	}

	if err := d.Set("error_message", interfaceEndpoint.ErrorMessage); err != nil {
		return fmt.Errorf(errorInterfaceEndpointSetting, "error_message", interfaceEndpointID, err)
	}

	if err := d.Set("connection_status", interfaceEndpoint.ConnectionStatus); err != nil {
		return fmt.Errorf(errorInterfaceEndpointSetting, "connection_status", interfaceEndpointID, err)
	}

	if err := d.Set("private_link_id", privateLinkID); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "private_link_id", privateLinkID, err)
	}

	if err := d.Set("interface_endpoint_id", interfaceEndpointID); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "interface_endpoint_id", privateLinkID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointInterfaceLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	interfaceEndpointID := ids["interface_endpoint_id"]

	if interfaceEndpointID != "" {
		_, err := conn.PrivateEndpointsDeprecated.DeleteOneInterfaceEndpoint(context.Background(), projectID, privateLinkID, interfaceEndpointID)
		if err != nil {
			return fmt.Errorf(errorInterfaceEndpointDelete, interfaceEndpointID, err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"NONE", "PENDING_ACCEPTANCE", "PENDING", "DELETING"},
			Target:     []string{"REJECTED", "DELETED"},
			Refresh:    resourceInterfaceEndpointRefreshFunc(conn, projectID, privateLinkID, interfaceEndpointID),
			Timeout:    1 * time.Hour,
			MinTimeout: 5 * time.Second,
			Delay:      3 * time.Second,
		}

		// Wait, catching any errors
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(errorInterfaceEndpointDelete, interfaceEndpointID, err)
		}
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointInterfaceLinkImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 4)
	if len(parts) != 4 {
		return nil, errors.New("import format error: to import a MongoDB Private Endpoint, use the format {project_id}-{private_link_id}-{interface_endpoint_id}")
	}

	projectID := parts[0]
	privateLinkID := parts[1]
	interfaceEndpointID := parts[2] + "-" + parts[3]

	_, _, err := conn.PrivateEndpointsDeprecated.GetOneInterfaceEndpoint(context.Background(), projectID, privateLinkID, interfaceEndpointID)
	if err != nil {
		return nil, fmt.Errorf(errorInterfaceEndpointRead, interfaceEndpointID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointsSetting, "project_id", privateLinkID, err)
	}

	if err := d.Set("private_link_id", privateLinkID); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointsSetting, "private_link_id", privateLinkID, err)
	}

	if err := d.Set("interface_endpoint_id", interfaceEndpointID); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointsSetting, "interface_endpoint_id", privateLinkID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":            projectID,
		"private_link_id":       privateLinkID,
		"interface_endpoint_id": interfaceEndpointID,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceInterfaceEndpointRefreshFunc(client *matlas.Client, projectID, privateLinkID, interfaceEndpointID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		i, resp, err := client.PrivateEndpointsDeprecated.GetOneInterfaceEndpoint(context.Background(), projectID, privateLinkID, interfaceEndpointID)
		if err != nil {
			if resp.Response.StatusCode == 404 {
				return "", "DELETED", nil
			}

			return nil, "FAILED", err
		}

		if i.ConnectionStatus != "AVAILABLE" {
			return "", i.ConnectionStatus, nil
		}

		return i, i.ConnectionStatus, nil
	}
}
