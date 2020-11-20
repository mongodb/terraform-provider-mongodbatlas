package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasPrivateEndpointInterfaceLink() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasPrivateEndpointInterfaceLinkRead,
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
		DeprecationMessage: "this data source is deprecated, please transition as soon as possible to mongodbatlas_privatelink_endpoint_service",
	}
}

func dataSourceMongoDBAtlasPrivateEndpointInterfaceLinkRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	privateLinkID := d.Get("private_link_id").(string)
	interfaceEndpointID := d.Get("interface_endpoint_id").(string)

	interfaceEndpoint, _, err := conn.PrivateEndpointsDeprecated.GetOneInterfaceEndpoint(context.Background(), projectID, privateLinkID, interfaceEndpointID)
	if err != nil {
		return fmt.Errorf(errorServiceEndpointRead, interfaceEndpointID, err)
	}

	if err := d.Set("delete_requested", cast.ToBool(interfaceEndpoint.DeleteRequested)); err != nil {
		return fmt.Errorf(errorEndpointSetting, "delete_requested", interfaceEndpointID, err)
	}

	if err := d.Set("error_message", interfaceEndpoint.ErrorMessage); err != nil {
		return fmt.Errorf(errorEndpointSetting, "error_message", interfaceEndpointID, err)
	}

	if err := d.Set("connection_status", interfaceEndpoint.ConnectionStatus); err != nil {
		return fmt.Errorf(errorEndpointSetting, "connection_status", interfaceEndpointID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":            projectID,
		"private_link_id":       privateLinkID,
		"interface_endpoint_id": interfaceEndpoint.ID,
	}))

	return nil
}
