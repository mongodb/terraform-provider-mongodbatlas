package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spf13/cast"

	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasPrivateEndpointLink() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasPrivateEndpointLinkRead,
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
	}
}

func dataSourceMongoDBAtlasPrivateEndpointLinkRead(d *schema.ResourceData, meta interface{}) error {
	//Get client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	privateLinkID := d.Get("private_link_id").(string)
	interfaceEndpointID := d.Get("interface_endpoint_id").(string)

	interfaceEndpoint, _, err := conn.PrivateEndpoints.GetOneInterfaceEndpoint(context.Background(), projectID, privateLinkID, interfaceEndpointID)
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

	d.SetId(encodeStateID(map[string]string{
		"project_id":            projectID,
		"private_link_id":       privateLinkID,
		"interface_endpoint_id": interfaceEndpoint.ID,
	}))

	return nil
}
