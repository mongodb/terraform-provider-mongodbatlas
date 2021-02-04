package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasPrivateEndpointServiceLink() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasPrivateEndpointServiceLinkRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func dataSourceMongoDBAtlasPrivateEndpointServiceLinkRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	privateLinkID := d.Get("private_link_id").(string)
	endpointServiceID := d.Get("endpoint_service_id").(string)
	providerName := d.Get("provider_name").(string)

	serviceEndpoint, _, err := conn.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, providerName, privateLinkID, endpointServiceID)
	if err != nil {
		return fmt.Errorf(errorServiceEndpointRead, endpointServiceID, err)
	}

	if err := d.Set("delete_requested", cast.ToBool(serviceEndpoint.DeleteRequested)); err != nil {
		return fmt.Errorf(errorEndpointSetting, "delete_requested", endpointServiceID, err)
	}

	if err := d.Set("error_message", serviceEndpoint.ErrorMessage); err != nil {
		return fmt.Errorf(errorEndpointSetting, "error_message", endpointServiceID, err)
	}

	if err := d.Set("aws_connection_status", serviceEndpoint.AWSConnectionStatus); err != nil {
		return fmt.Errorf(errorEndpointSetting, "aws_connection_status", endpointServiceID, err)
	}

	if err := d.Set("azure_status", serviceEndpoint.AzureStatus); err != nil {
		return fmt.Errorf(errorEndpointSetting, "azure_status", endpointServiceID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":          projectID,
		"private_link_id":     privateLinkID,
		"endpoint_service_id": endpointServiceID,
		"provider_name":       providerName,
	}))

	return nil
}
