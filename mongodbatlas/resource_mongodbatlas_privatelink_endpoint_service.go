package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorServiceEndpointAdd  = "error adding MongoDB Private Service Endpoint Connection(%s) to a Private Endpoint (%s): %s"
	errorServiceEndpointRead = "error reading MongoDB Private Service Endpoint Connection(%s): %s"
	errorEndpointDelete      = "error deleting MongoDB Private Service Endpoint Connection(%s): %s"
	errorEndpointSetting     = "error setting `%s` for MongoDB Private Service Endpoint Connection(%s): %s"
)

func resourceMongoDBAtlasPrivateEndpointServiceLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasPrivateEndpointServiceLinkCreate,
		Read:   resourceMongoDBAtlasPrivateEndpointServiceLinkRead,
		Delete: resourceMongoDBAtlasPrivateEndpointServiceLinkDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasPrivateEndpointServiceLinkImportState,
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
			"private_link_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interface_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_endpoint_connection_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_endpoint_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_endpoint_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_requested": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_connection_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azure_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongoDBAtlasPrivateEndpointServiceLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	privateLinkID := d.Get("private_link_id").(string)
	providerName := d.Get("provider_name").(string)
	endpointServiceID := d.Get("endpoint_service_id").(string)

	request := &matlas.InterfaceEndpointConnection{
		ID:                       endpointServiceID,
		PrivateEndpointIPAddress: d.Get("private_endpoint_ip_address").(string),
	}

	_, _, err := conn.PrivateEndpoints.AddOnePrivateEndpoint(context.Background(), projectID, providerName, privateLinkID, request)
	if err != nil {
		return fmt.Errorf(errorServiceEndpointAdd, providerName, privateLinkID, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NONE", "INITIATING", "PENDING_ACCEPTANCE", "PENDING", "DELETING"},
		Target:     []string{"AVAILABLE", "REJECTED", "DELETED", "FAILED"},
		Refresh:    resourceServiceEndpointRefreshFunc(conn, projectID, providerName, privateLinkID, endpointServiceID),
		Timeout:    1 * time.Hour,
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
	}
	// Wait, catching any errors
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(errorServiceEndpointAdd, endpointServiceID, privateLinkID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":          projectID,
		"private_link_id":     privateLinkID,
		"endpoint_service_id": endpointServiceID,
		"provider_name":       providerName,
	}))

	return resourceMongoDBAtlasPrivateEndpointServiceLinkRead(d, meta)
}

func resourceMongoDBAtlasPrivateEndpointServiceLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	endpointServiceID := ids["endpoint_service_id"]
	providerName := ids["provider_name"]

	privateEndpoint, _, err := conn.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, providerName, privateLinkID, endpointServiceID)
	if err != nil {
		return fmt.Errorf(errorServiceEndpointRead, endpointServiceID, err)
	}

	if err := d.Set("delete_requested", cast.ToBool(privateEndpoint.DeleteRequested)); err != nil {
		return fmt.Errorf(errorEndpointSetting, "delete_requested", endpointServiceID, err)
	}

	if err := d.Set("error_message", privateEndpoint.ErrorMessage); err != nil {
		return fmt.Errorf(errorEndpointSetting, "error_message", endpointServiceID, err)
	}

	if err := d.Set("aws_connection_status", privateEndpoint.AWSConnectionStatus); err != nil {
		return fmt.Errorf(errorEndpointSetting, "aws_connection_status", endpointServiceID, err)
	}

	if err := d.Set("azure_status", privateEndpoint.AzureStatus); err != nil {
		return fmt.Errorf(errorEndpointSetting, "azure_status", endpointServiceID, err)
	}

	if err := d.Set("interface_endpoint_id", privateEndpoint.InterfaceEndpointID); err != nil {
		return fmt.Errorf(errorEndpointSetting, "connection_status", endpointServiceID, err)
	}

	if err := d.Set("private_endpoint_connection_name", privateEndpoint.PrivateEndpointConnectionName); err != nil {
		return fmt.Errorf(errorEndpointSetting, "connection_status", endpointServiceID, err)
	}

	if err := d.Set("private_endpoint_ip_address", privateEndpoint.PrivateEndpointIPAddress); err != nil {
		return fmt.Errorf(errorEndpointSetting, "connection_status", endpointServiceID, err)
	}

	if err := d.Set("private_endpoint_resource_id", privateEndpoint.PrivateEndpointResourceID); err != nil {
		return fmt.Errorf(errorEndpointSetting, "connection_status", endpointServiceID, err)
	}

	if err := d.Set("endpoint_service_id", endpointServiceID); err != nil {
		return fmt.Errorf(errorEndpointSetting, "endpoint_service_id", endpointServiceID, err)
	}

	if err := d.Set("private_link_id", privateLinkID); err != nil {
		return fmt.Errorf(errorEndpointSetting, "private_link_id", endpointServiceID, err)
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointServiceLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*matlas.Client)

	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	privateLinkID := ids["private_link_id"]
	endpointServiceID := ids["endpoint_service_id"]
	providerName := ids["provider_name"]

	if endpointServiceID != "" {
		_, err := conn.PrivateEndpoints.DeleteOnePrivateEndpoint(context.Background(), projectID, providerName, privateLinkID, endpointServiceID)
		if err != nil {
			return fmt.Errorf(errorEndpointDelete, endpointServiceID, err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"NONE", "PENDING_ACCEPTANCE", "PENDING", "DELETING", "INITIATING"},
			Target:     []string{"REJECTED", "DELETED", "FAILED"},
			Refresh:    resourceServiceEndpointRefreshFunc(conn, projectID, providerName, privateLinkID, endpointServiceID),
			Timeout:    1 * time.Hour,
			MinTimeout: 5 * time.Second,
			Delay:      3 * time.Second,
		}

		// Wait, catching any errors
		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(errorEndpointDelete, endpointServiceID, err)
		}
	}

	return nil
}

func resourceMongoDBAtlasPrivateEndpointServiceLinkImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "--", 4)
	if len(parts) != 4 {
		return nil, errors.New("import format error: to import a MongoDB Private Endpoint, use the format {project_id}--{private_link_id}--{endpoint_service_id}--{provider_name}")
	}

	projectID := parts[0]
	privateLinkID := parts[1]
	endpointServiceID := parts[2]
	providerName := parts[3]

	_, _, err := conn.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, providerName, privateLinkID, endpointServiceID)
	if err != nil {
		return nil, fmt.Errorf(errorServiceEndpointRead, endpointServiceID, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointsSetting, "project_id", privateLinkID, err)
	}

	if err := d.Set("private_link_id", privateLinkID); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointsSetting, "private_link_id", privateLinkID, err)
	}

	if err := d.Set("endpoint_service_id", endpointServiceID); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointsSetting, "endpoint_service_id", privateLinkID, err)
	}

	if err := d.Set("provider_name", providerName); err != nil {
		return nil, fmt.Errorf(errorPrivateEndpointsSetting, "provider_name", privateLinkID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":          projectID,
		"private_link_id":     privateLinkID,
		"endpoint_service_id": endpointServiceID,
		"provider_name":       providerName,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceServiceEndpointRefreshFunc(client *matlas.Client, projectID, providerName, privateLinkID, endpointServiceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		i, resp, err := client.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, providerName, privateLinkID, endpointServiceID)
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return "", "DELETED", nil
			}

			return nil, "", err
		}

		if strings.EqualFold(providerName, "azure") {
			if i.AzureStatus != "AVAILABLE" {
				return "", i.AzureStatus, nil
			}
			return i, i.AzureStatus, nil
		}
		if i.AWSConnectionStatus != "AVAILABLE" {
			return "", i.AWSConnectionStatus, nil
		}

		return i, i.AWSConnectionStatus, nil
	}
}
